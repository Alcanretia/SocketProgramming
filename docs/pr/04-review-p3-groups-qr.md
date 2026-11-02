## Heads-up for P3
While building the app I (P4) also wrote **your** part with an AI tool, without asking you first: create / list / join groups with member lists, group message fan-out, invite tokens, the group panel, the QR code and the join-confirmation dialog (`server/hub.go`, `client-web/src/features/groups/`). I also wrote the multi-client browser test (`e2e/`) that was yours to write. Sorry; it was not meant to take your part away. It is **not done for you** until you have read it, can explain it, scanned the QR with a real phone, and P2 has reviewed your explanation. The group rules (what happens when a member disconnects, whether empty groups disappear, invite expiry) are my proposals; you decide. Your notes file for this PR: `docs/reviews/p3-groups-qr.md`.

## What this PR does
P3's review of the groups and QR-invite code that P4 wrote, plus the follow-up work. **No server or client code changed.**
- Read the group and invite handlers (`onCreateGroup` to `onInviteInfo`, and the group part of `onLeave`) and `client-web/src/features/groups/`.
- Tested the QR invite with a real phone on the same Wi-Fi, using the LAN address: **passed** (2026-11-02).
- Added 12 browser checks to `e2e/e2e.mjs` for invite edge cases (table below).
- Wrote the group and invite rules into `docs/DECISIONS.md` (D7).
- Added a learning-log entry and updated README Section 0 (F2 phone scan, P3 status, two D7 items ticked).

New e2e cases (section "QR invite trap cases (P3)"):

| Case | What it proves |
|---|---|
| Hana opens the Study link and presses **Not now** | Declining does not add her (R10); the token is cleared from the URL; the same link worked for Dave first, so invites are reusable (D7) |
| Eve opens `/?invite=deadbeef` | A forged token gets a `bad_invite` toast, no dialog, and is removed from the URL |
| Frank creates "Temp", makes an invite, then leaves; Gina opens the link | The empty group disappears for everyone; the stale invite gets a "no longer exists" toast, no dialog |

## Whose module / which requirement
- Owner: P3
- Touches: R5, R8, R9, R10, R11, F2, the multi-client test script
- Files changed: `e2e/e2e.mjs`, `docs/DECISIONS.md`, `docs/LEARNING_LOG.md`, `README.md`, this file

## What I learned (author)
Group fan-out is a `for` loop over `g.Members` in `onGroupMsg`, so only members receive a group message. A non-member gets `not_member` from the server, whatever the UI shows.

Scanning a QR is not joining. The page only sends `invite_info`, which is read-only; `join_group` is a separate request the user sends by pressing Join (R10).

When a user disconnects, `onLeave` removes them from every group. If they were the last member, the group and its invite tokens are deleted.

Two groups can never get the same ID because only the hub goroutine touches `groupSeq`, and it handles one event at a time.

## Explain-back (author answers, reviewer checks)
- [x] I can explain every line of my part.
- Q1: How does your server decide who receives a group message? (`onGroupMsg`, loop over `g.Members`)
- Q2: Why isn't scanning the QR the same as being added by the creator (R10)?
- Reviewer's questions (P2 reviews this PR):
  1. If a non-member sends a group message, what does the server do, and does anyone else see the attempt?
  2. What happens to membership when a member disconnects, and to a group when everyone has left? Where is that decided?

## Checks
- [x] `cd server; go test ./...` passes (13 tests; also `go test -race ./...` is clean, the first time anyone has run it)
- [x] `cd client-web; npm run build` passes
- [x] 3-tab test: A creates, B joins, C does not; C sees nothing and cannot send
- [x] **Scanned the QR with a real phone** on the same Wi-Fi (app opened via the LAN address, not `localhost`): passed, 2026-11-02
- [x] `e2e/` run and extended with my own cases (declined invite, forged token, reused link, group deleted before a scan): **47/47 pass**. Joining twice is already covered by `TestGroupsMembershipAndPrivacy` in `server/hub_test.go`.
- [x] Group rules written in `docs/DECISIONS.md` (D7). `docs/PROTOCOL.md` section 3.3 already states them, so no change there. The invite rules (reusable, no expiry) still need a line in section 3.5; left to P2 as the `PROTOCOL.md` keeper.
- [x] Learning-log entry added
- [x] README Section 0 updated

## Safety
- [x] Membership is enforced on the server (`g.has(c.key)` in `onGroupMsg`, `onCreateInvite` and `roomAccess`; only `onJoinGroup` adds members, and only the sender)
- [x] Invite tokens are random and unguessable (`crypto/rand`, 128 bits); a token for another group is rejected (`bad_invite`)
- [x] No secrets in the diff (`.env` is ignored and not tracked)

## Notes for the reviewer
- The invite token is **not access control**. R9 and R10 make every group open, so anyone can join without one. The token keeps the QR link tied to one group and leaves room for expiry or revocation later.
- The "You are already a member / Open group" branch in `InviteJoinDialog` cannot be reached from the UI. Opening a link starts a new connection, and a new connection is in no groups. It is harmless; I left it.
- Disconnecting, including a page refresh, drops all group memberships. **Do not refresh during the demo recording.**

## Who wrote it
- [x] The original group and invite code, the QR UI and `e2e/` were written by P4 with an AI tool; I have read them and can explain them.
- [x] The new e2e cases were drafted with an AI assistant; I reviewed and ran them and can explain each check.
- [x] The D7 decisions and the learning-log entry are mine; the D7 wording was drafted with AI help.

## Reviewer
P2
