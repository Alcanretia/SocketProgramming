# P2 review and learning notes: users, DMs, stickers

Branch `review/p2-users-dm`. Reviewer: P1. Fill this in with your own words; commit it to this branch.

> The code in your area, and the first draft of `docs/PROTOCOL.md` (which you keep), were written by P4 with an AI tool before you started. Read it, change what you disagree with, and answer below. It is not done until you can explain it.

## Read in this order
1. Do your own Phase 1 first.
2. `docs/PROTOCOL.md` completely. Then `server/protocol.go`.
3. `server/hub.go`: `onRegister`, `onLeave` (name freed), `dmRoom`, `deliver`, `onDM`; `server/stickers.go`, `validContent`.
4. `server/hub_test.go`: `TestRegisterAndUserList`, `TestSimultaneousRegistrationOnlyOneWins`, `TestDMPrivacyAndSpoofing`, `TestStickers`.
5. `client-web/src/features/users-dm/`.

## Answer in your own words
1. How does the server guarantee unique names even with two simultaneous registrations?

   Answer:
   The server uses a single owner goroutine (`Hub.Run` in `server/hub.go`) that reads events sequentially from the `h.events` channel, rather than allowing multiple client `readPump` goroutines to touch `h.byKey` directly. When two clients send `register` at the exact same instant, each client's `readPump` posts an `inbound{c, env}` event into `h.events`. The hub goroutine processes them one at a time in `onRegister`: it normalizes the name with `key := strings.ToLower(name)` and checks `if _, taken := h.byKey[key]; taken`. Because no other goroutine can read or write `h.byKey`, this check-then-set (`h.byKey[key] = c`) cannot interleave or race. The first event processed claims `h.byKey[key]` and gets `register_ok`; the second immediately sees `taken == true` and receives `register_error` with `code: "taken"` (verified by `TestSimultaneousRegistrationOnlyOneWins`).

2. Show the exact line where a message is sent to only one person.

   Answer:
   At the end of `onDM` in `server/hub.go` (lines 433–434):
   ```go
   h.deliver(c, c, "dm", m)
   h.deliver(to, c, "dm", m)
   ```
   Inside `h.deliver(to, from, typ, m)` (`server/hub.go` line 359), the server calls:
   ```go
   h.send(to, typ, m, "")
   ```
   which calls `to.enqueue(encode(typ, payload, reqID))` (`server/hub.go` line 138) to push the encoded JSON frame into **only** that single recipient's buffered channel (`to.send`). `onDM` calls `deliver` twice: once for the sender `c` (so the sender renders the server-confirmed copy) and once for the single target client `to := h.byKey[toKey]`. No broadcast loop is used, so no third client ever receives the frame.

3. Can a client send a DM that claims to be from someone else? Where is that stopped?

   Answer:
   No. First, `dmReq` in `server/protocol.go` (lines 57–62) only decodes `to`, `kind`, `text`, and `stickerId` — any extra `"from"` field sent in the client's JSON payload is ignored during `json.Unmarshal`. Second, in `onDM` (`server/hub.go` lines 427–431), the server constructs the outgoing `Message` and explicitly stamps `From: c.name` using the registered name bound to the sender's `*Client` connection (`c`):
   ```go
   m := Message{
       Room: dmRoom(c.key, to.key), From: c.name, To: to.name,
       Kind: req.Kind, Text: text, StickerID: req.StickerID,
       Ts: time.Now().UnixMilli(), MsgID: h.newMsgID(),
   }
   ```
   Because `c.name` is owned by the server and set only upon a successful `onRegister`, a client can never spoof another user's identity (verified by `TestDMPrivacyAndSpoofing`).

4. Is the room between A and B the same as between B and A? How is it identified?

   Answer:
   Yes, they are the exact same room. It is identified by `dmRoom(a, b string) string` in `server/hub.go` (lines 315–320) and mirrored by `dmRoomId(a, b)` in `client-web/src/socket/types.ts` (lines 43–46):
   ```go
   func dmRoom(a, b string) string {
       if a > b {
           a, b = b, a
       }
       return "dm:" + a + "|" + b
   }
   ```
   Both user keys (`c.key` and `to.key`) are already lowercased (`strings.ToLower`). `dmRoom` compares them lexicographically and swaps them if `a > b` so the two names are always ordered alphabetically, then formats the ID as `"dm:<smaller>|<larger>"` (for example, both Alice→Bob and Bob→Alice produce `"dm:alice|bob"`).

5. What exactly travels over the WebSocket when I send a sticker?

   Answer:
   Only a small JSON text frame containing the sticker's string identifier (`stickerId`) travels over the WebSocket — never the SVG image bytes.
   - **Client → Server:**
     `{"type":"dm","payload":{"to":"bob","kind":"sticker","stickerId":"smile"}}` (or `"type":"group_msg"` with `"groupId"`).
   - **Server validation:** `validContent` (`server/protocol.go` lines 159–162) checks `isSticker(stickerID)` against the allowed list in `server/stickers.go` (`smile`, `laugh`, `love`, `cry`, `angry`, `surprised`, `sleepy`, `thumbs_up`). Unknown IDs (such as `"../../etc/passwd"`) are rejected with `error` (`code: "bad_sticker"`).
   - **Server → Client (delivered):**
     `{"type":"dm","payload":{"room":"dm:alice|bob","from":"Alice","to":"Bob","kind":"sticker","stickerId":"smile","ts":1730460000000,"msgId":"m_1"}}`
   - **Rendering:** `StickerView` in `client-web/src/features/users-dm/stickers.tsx` renders `<img src="/stickers/<id>.svg" />`, which loads the static SVG asset over HTTP (satisfying R2, since chat messages travel only over the WebSocket while static assets use HTTP).

## `docs/PROTOCOL.md`
- [x] I read all of it
- Open questions in section 5 (history? unread retention? mentions? AI bot user? rate limiting?), my answers:
  - **Message history:** No server-side message history replay on join. Clients only see messages received while online; the server keeps up to 200 unread messages per room per user (`maxUnreadKept = 200`) strictly to feed `summarize_request`.
  - **Unread retention on disconnect:** Dropping unread state on disconnect (`delete(h.unread, c.key)` in `onLeave`) is the right choice because there is no password authentication; if unread messages survived disconnect, someone else registering the freed name could read or summarize the previous user's unread messages.
  - **Mentions (`@name`):** Out of scope for v0.1; F4 notifications already alert on any incoming DM or group message when the room is not actively open/visible.
  - **AI bot user vs. summary only:** Summary only (`summarize_request` → `summary_result` / `summary_error`). It fulfills F1 cleanly without injecting a fake user into `user_list`.
  - **`rate_limited`:** Keep as implemented (`c.summarizing` prevents a user from running more than one concurrent `summarize_request`).
- Things I want to change (doc first, then code and `client-web/src/socket/types.ts`):
  - No structural wire changes needed; v0.1 cleanly satisfies R2–R11 and F1–F4. Updated `docs/PROTOCOL.md` Section 5 with the answers above and ticked P2 acknowledgement in Section 6.
- [x] I agree with it and ticked my acknowledgement box

## Things I tried
- [x] Tricky names: unicode, emoji, 500 characters, only spaces (added explicit 500-char, control-char, and Thai+emoji test cases in `TestRegisterAndUserList` in `server/hub_test.go`)
- [x] Third-tab privacy test (verified in `TestDMPrivacyAndSpoofing`: Carol's third connection receives nothing when Alice and Bob exchange DMs)
- [x] Sticker artwork looks fine / replaced with my own (reviewed all 8 SVGs in `client-web/public/stickers/` and verified sync with `server/stickers.go` via `TestStickerAssetsExist`)
