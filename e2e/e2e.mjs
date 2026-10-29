import { chromium } from "playwright-core";

const BASE = process.env.BASE_URL ?? "http://localhost:8080";
const SHOTS = new URL("./shots/", import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1");
import fs from "node:fs";
fs.mkdirSync(SHOTS, { recursive: true });

const browser = await chromium.launch({
  executablePath: process.env.BROWSER_PATH ?? "C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe",
  headless: true,
});

let failures = 0;
const check = (ok, what) => {
  console.log(`${ok ? "PASS" : "FAIL"}  ${what}`);
  if (!ok) failures++;
};
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const shot = (page, name) => page.screenshot({ path: `${SHOTS}${name}.png` });

async function join(name, url = BASE) {
  const ctx = await browser.newContext({ viewport: { width: 1000, height: 700 } });
  const page = await ctx.newPage();
  page.on("pageerror", (e) => console.log(`  [pageerror ${name}]`, e.message));
  await page.goto(url);
  await page.getByLabel("Your name").fill(name);
  await page.getByRole("button", { name: "Join" }).click();
  await page.getByText("Signed in as").waitFor({ timeout: 5000 });
  return page;
}
const visible = async (loc, ms = 4000) => {
  try {
    await loc.first().waitFor({ state: "visible", timeout: ms });
    return true;
  } catch {
    return false;
  }
};
const absent = async (loc, ms = 800) => {
  await sleep(ms);
  return (await loc.count()) === 0;
};

try {
  // ---- R3 / R4 --------------------------------------------------------
  const alice = await join("Alice");
  const bob = await join("Bob");
  const carol = await join("Carol");

  const dup = await browser.newContext().then((c) => c.newPage());
  await dup.goto(BASE);
  await dup.getByLabel("Your name").fill("aLICE");
  await dup.getByRole("button", { name: "Join" }).click();
  check(await visible(dup.getByText("already in use")), "R3 duplicate name (different case) is rejected in the UI");

  const online = alice.locator("aside").getByRole("button", { name: /Bob/ });
  check(await visible(online), "R4 Alice sees Bob in the online list");
  check(await visible(alice.locator("aside").getByText("(you)")), "R4 Alice sees her own name too");

  // ---- DM (R5, R7) + notification fallback ---------------------------------
  await alice.locator("aside").getByRole("button", { name: /Bob/ }).click();
  await alice.getByLabel("Message").fill("hello bob <b>bold?</b>");
  await alice.getByRole("button", { name: "Send" }).click();
  check(await visible(alice.locator(".bubble", { hasText: "hello bob" })), "Alice sees her own DM in the room");
  check(
    await visible(bob.locator(".toast.message", { hasText: "Alice: hello bob" })),
    "F4 Bob gets an in-app toast for a message in a room he is not viewing",
  );
  check(
    await visible(bob.locator("aside .badge", { hasText: "1" })),
    "Bob sees a server-counted unread badge",
  );
  check((await bob.title()).startsWith("(1)"), "F4 tab title shows the unread count");
  check(await absent(carol.getByText("hello bob")), "R7 Carol (third tab) sees nothing of the DM");
  check(await absent(carol.locator(".toast")), "R7 Carol gets no toast either");
  check(
    (await alice.locator(".bubble b").count()) === 0,
    "message text is escaped (no HTML injected)",
  );

  await bob.locator("aside").getByRole("button", { name: /Alice/ }).first().click();
  check(await visible(bob.locator(".bubble", { hasText: "hello bob" })), "Bob opens the DM and sees the message");
  check(await absent(bob.locator("aside .badge")), "unread clears after Bob opens the room");

  // ---- stickers (F3) -----------------------------------------------------------
  await bob.getByLabel("Stickers").click();
  await bob.getByTitle("love").click();
  check(await visible(bob.locator("img.sticker[alt='sticker: love']")), "Bob sees the sticker he sent");
  check(await visible(alice.locator("img.sticker[alt='sticker: love']")), "Alice receives and renders the sticker");
  check(await absent(carol.locator("img.sticker")), "Carol does not receive the sticker");
  await shot(alice, "dm-sticker");

  // ---- groups (R8-R11) -----------------------------------------------------------
  await alice.getByPlaceholder("New group name").fill("Study");
  await alice.getByRole("button", { name: "Create" }).click();
  const groupRow = (page) => page.locator("li.group-item", { hasText: "Study" });
  check(await visible(groupRow(bob)), "R9 Bob sees the new group in the list");
  check(await visible(groupRow(carol)), "R9 Carol sees the new group in the list");
  check(
    (await groupRow(bob).locator(".members").innerText()).trim() === "Alice",
    "R8 creator is the only initial member",
  );
  await groupRow(bob).getByRole("button", { name: "Join" }).click();
  check(await visible(bob.locator(".chat-header", { hasText: "Study" })), "R10 Bob joined by himself and the room opened");
  check(
    await visible(groupRow(alice).locator(".members", { hasText: "Alice, Bob" })),
    "R9 member list updates live for Alice",
  );

  await alice.getByLabel("Message").fill("group hello");
  await alice.getByRole("button", { name: "Send" }).click();
  check(await visible(bob.locator(".bubble", { hasText: "group hello" })), "R11 member Bob receives the group message");
  check(await absent(carol.getByText("group hello")), "R11 non-member Carol does not see it");
  check(await visible(groupRow(carol).getByRole("button", { name: "Join" })), "Carol still has to opt in (Join button)");

  // ---- QR invite (F2) --------------------------------------------------------------
  await groupRow(alice).getByRole("button", { name: "Invite (QR)" }).click();
  check(await visible(alice.locator(".qr img[alt*='QR']")), "F2 Alice sees a QR code image");
  const link = (await alice.locator(".qr code").innerText()).trim();
  check(/\/\?invite=[0-9a-f]{32}$/.test(link), `F2 QR link has an unguessable token (${link.slice(0, 40)}...)`);
  await shot(alice, "qr");

  const dave = await browser.newContext({ viewport: { width: 400, height: 760 } }).then((c) => c.newPage());
  dave.on("pageerror", (e) => console.log("  [pageerror dave]", e.message));
  await dave.goto(link);
  check(await visible(dave.getByText("You were invited to a group")), "scanning the link opens the app on the name screen");
  await dave.getByLabel("Your name").fill("Dave");
  await dave.getByRole("button", { name: "Join" }).click();
  check(await visible(dave.getByRole("dialog").getByText('Join group "Study"?')), 'F2 Dave is asked "Join group Study?"');
  check(
    await absent(groupRow(alice).locator(".members", { hasText: "Dave" }), 600),
    "R10 scanning alone did NOT add Dave to the group",
  );
  await shot(dave, "invite-dialog");
  await dave.getByRole("dialog").getByRole("button", { name: "Join" }).click();
  check(await visible(dave.locator(".chat-header", { hasText: "Study" })), "Dave is in the group after pressing Join");
  await sleep(300);
  check(!dave.url().includes("invite="), "invite token is removed from the address bar after joining");
  check(
    await visible(groupRow(alice).locator(".members", { hasText: "Dave" })),
    "group member list now includes Dave",
  );

  // ---- AI summary (F1): no API key here, so it must fail gracefully -----------
  await carol.locator("aside").getByRole("button", { name: /Alice/ }).first().waitFor().catch(() => {});
  await alice.locator("aside").getByRole("button", { name: /Carol/ }).click();
  await alice.getByLabel("Message").fill("are you coming tonight?");
  await alice.getByRole("button", { name: "Send" }).click();
  check(await visible(carol.locator("aside .unread-item", { hasText: "Alice" })), "F1 Carol sees an Unread entry");
  await carol.getByRole("button", { name: /Summarize/ }).click();
  check(
    await visible(carol.locator(".summary-card .error")),
    "F1 summary without an API key shows a friendly error instead of hanging",
  );
  await shot(carol, "summary-error");

  // chat still works after the failure
  await alice.getByLabel("Message").fill("still there?");
  await alice.getByRole("button", { name: "Send" }).click();
  check(await visible(carol.locator(".toast.message", { hasText: "still there?" })), "chat keeps working after an AI failure");

  // ---- disconnect / reconnect -----------------------------------------------------------
  await bob.context().close();
  check(
    await visible(alice.locator("aside").getByText("Online (3)")),
    "user list shrinks (Alice, Carol, Dave) when Bob's tab closes",
  );
} catch (e) {
  console.log("SCRIPT ERROR:", e.message);
  failures++;
} finally {
  await browser.close();
}
console.log(failures === 0 ? "\nALL PASSED" : `\n${failures} FAILED`);
process.exit(failures === 0 ? 0 : 1);
