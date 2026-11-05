## Heads-up for P2
While building the app I (P4) also wrote **your** part with an AI tool, without asking you first: name validation and uniqueness, the user list, DM routing, sticker validation and the eight sticker SVGs (`server/hub.go`, `server/protocol.go`, `server/stickers.go`, `client-web/src/features/users-dm/`). I also wrote the first draft of `docs/PROTOCOL.md`, which you are supposed to keep. Sorry; it was not meant to take your part away. It is **not done for you** until you have read it, can explain it, and P1 has reviewed your explanation. The protocol is a proposal until you acknowledge it: change anything you disagree with (doc first, then code). Your notes file for this PR: `docs/reviews/p2-users-dm.md`.

## What this PR does
P2's implementation, review, and learning changes for identity, DMs and stickers:
- Completed Phase 1 raw TCP exercises (`server/phase1/echo`, `server/phase1/hub`, `server/phase1/client` + unit tests).
- Modularized and documented P2's backend handlers (`validName`, `normalizeKey`, `broadcastUsers`, `onRegister`, `dmRoom`, `onDM` in `server/users_dm.go`, and `stickerPack`, `isSticker`, `validContent` in `server/stickers.go`).
- Implemented and documented P2's frontend slice (`client-web/src/features/users-dm/NameEntry.tsx`, `UserList.tsx`, `stickers.tsx`).
- Added unit tests for tricky names (Thai/Emoji, spaces-only, 500 chars, control chars), DM room symmetry, and sticker validation in `server/users_dm_test.go`.
- Answered all open questions in `docs/PROTOCOL.md` Section 5 and ticked P2's acknowledgement in `docs/PROTOCOL.md` and `docs/DECISIONS.md`.
- Completed all explain-back notes in `docs/reviews/p2-users-dm.md` and logged the session in `docs/LEARNING_LOG.md`.

## Whose module / which requirement
- Owner: P2
- Touches: R3, R4, R5, R7, F3, `docs/PROTOCOL.md` (P2 is the keeper)
- Files: `server/users_dm.go`, `server/users_dm_test.go`, `server/stickers.go`, `server/phase1/`, `client-web/src/features/users-dm/`, `docs/PROTOCOL.md`, `docs/reviews/p2-users-dm.md`

## What I learned (author)
- **How uniqueness holds (`onRegister`):** All shared state (`h.byKey`) is owned by the single `Hub.Run` goroutine and accessed sequentially through the `h.events` channel. Even if multiple clients send `register` simultaneously, `onRegister` runs one at a time, so the case-insensitive check-then-set (`key := normalizeKey(name)`; `if _, taken := h.byKey[key]; taken`) cannot interleave or race without needing a mutex.
- **How a DM reaches exactly two people (`onDM`):** The server stamps `From: c.name` from the sender's connection (preventing spoofing), computes the canonical room ID via `dmRoom(c.key, to.key)` (`"dm:<smaller>|<larger>"`), and calls `h.deliver(c, c, "dm", m)` and `h.deliver(to, c, "dm", m)`. Only the sender's and recipient's pipes receive the message.
- **What a sticker message contains:** Only the string `stickerId` (with `kind: "sticker"`) travels over the WebSocket (`{"type":"dm","payload":{"to":"bob","kind":"sticker","stickerId":"smile"}}`). The server validates `stickerId` against `stickerPack` in `server/stickers.go` (`validContent`), and the browser loads `/stickers/<id>.svg` as a static asset over HTTP (preserving R2).

## Explain-back (author answers, reviewer checks)
- [x] I can explain every line of my part.
- Q1: How does the server guarantee uniqueness even with two simultaneous registrations? (`TestSimultaneousRegistrationOnlyOneWins`) — Answered in `docs/reviews/p2-users-dm.md`.
- Q2: Show me the exact line where a message is sent to only one person. (the two `deliver` calls at the end of `onDM`) — Answered in `docs/reviews/p2-users-dm.md` (`server/users_dm.go` lines 139–140).
- Q3: What exactly travels over the WebSocket when I send a sticker? — Answered in `docs/reviews/p2-users-dm.md`.
- Reviewer's questions (P1 reviews this PR):
  1. Can a client send a DM that claims to be from someone else? Where is that stopped? — Answered in `docs/reviews/p2-users-dm.md` (`dmReq` ignores `"from"`, and `onDM` stamps `From: c.name` in `server/users_dm.go`).
  2. Is the room between A and B the same as between B and A? How is it identified? — Answered in `docs/reviews/p2-users-dm.md` (`dmRoom(a, b)` sorts lowercased keys alphabetically into `"dm:<a>|<b>"`).

## Checks
- [x] `cd server; go test ./...` passes
- [x] `cd client-web; npm run build` passes
- [x] Third-tab privacy test done (a third tab must never see an A and B message)
- [x] Tried tricky names (unicode, emoji, 500 characters, only spaces, same name at the same instant)
- [x] **`docs/PROTOCOL.md`**: open questions in section 5 answered, and the acknowledgement box ticked **only if** I agree. Any change went into the doc first, then the code and `client-web/src/socket/types.ts`
- [x] Learning-log entry added
- [x] README Section 0 updated

## Safety
- [x] Sender is stamped by the server (spoofing test)
- [x] Unknown sticker ids are rejected by the server
- [x] No secrets in the diff

## Who wrote it
- [x] I wrote my changes myself (guided step-by-step via coaching session)
- [x] The original code and the `PROTOCOL.md` draft were written by P4 with an AI tool; I have read them and can explain them

## Reviewer
P1
