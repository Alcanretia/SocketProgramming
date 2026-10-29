# NEXT_STEPS: what is still left, per person

The app is built and passes automated tests on one computer, but it was **written by P4 with an AI helper, not by the owners below**. Nothing is "done" for a person until they have read their part, answered the explain-backs in their own words, and had it reviewed by someone who did not write it (README Section 8).

Suggested branches (create your own from `feature/app-p4-initial` or from `main` after it is merged):
`review/p1-connection-core`, `review/p2-users-dm`, `review/p3-groups-qr`, `learn/p4-explain-back`.
Commit your learning-log entry, your fixes and your review notes on your branch, and open a small PR.

## Everyone
- [ ] Read `docs/PROTOCOL.md` and `docs/DECISIONS.md`; tick your column; argue about anything you disagree with and change the code to match.
- [ ] Do the Phase 1 exercises (TCP echo, 2 clients, forwarding) **yourself before** reading `server/phase1/`.
- [ ] Read `docs/ARCHITECTURE.md` and write one `LEARNING_LOG.md` entry (built / confused / clicked).
- [ ] Mock instructor Q&A: take turns pointing at random lines ("what does this do, what if it fails, why this choice?").
- [ ] Make the "Implementation Explanation" slide (architecture diagram from ARCHITECTURE.md, protocol table from README Section 7, message flow for one DM, one group message, one summary).
- [ ] Agree what the course allows regarding AI-written code and say so to the instructor if required.

## P1: Connection core
- [ ] Phase 1 + Phase 2 yourself (WebSocket echo from scratch, find `101 Switching Protocols` in DevTools → Network → WS and keep a screenshot for the slide).
- [ ] Read `server/hub.go`, `server/client.go`, `client-web/src/socket/`. Explain-back: *connect to disconnect walk-through*; *why one writer goroutine*; *why `enqueue` never blocks*; *why `send` is never closed*.
- [ ] Break it (Phase 9): kill a client mid-message, 10 tabs, a huge message, malformed JSON. Write down what you saw.
- [ ] F4: try OS notifications on `localhost` (permission prompt, tab hidden). Decide what to say about plain-HTTP LAN and whether to try HTTPS. Decide if mentions (`@name`) are in scope.
- [ ] Review P2's part.

## P2: Users, DMs, stickers
- [ ] Phase 1 yourself, then read `onRegister`, `onDM`, `validContent`, `stickers.go`.
- [ ] Own `docs/PROTOCOL.md`: answer its open questions (history? unread retention? mentions?) and mark it acknowledged only when you agree. Any change goes into the doc first, then the code and `client-web/src/socket/types.ts`.
- [ ] Try to break names: unicode, emoji, 500 characters, spaces only, same name at the same instant (a test exists; understand why it passes).
- [ ] Explain-back: *how uniqueness holds under simultaneous registration*; *show the line where a DM goes to only one person*; *what travels over the socket for a sticker*.
- [ ] Review the sticker artwork (replace with your own if you want) and keep the id list in `server/stickers.go` and `client-web/src/features/users-dm/stickers.tsx` in sync (`TestStickerAssetsExist` checks the files).
- [ ] Review P3's part.

## P3: Groups and QR
- [ ] Phase 1 yourself, then read the group and invite handlers in `server/hub.go` (`onCreateGroup` to `onInviteInfo`).
- [ ] Decide and document the group rules you want (disconnect leaves all groups? empty group deleted? unique group names? invite expiry or revocation?).
- [ ] Scan the QR with a **real phone** on the same Wi-Fi (open the app via the LAN address, not `localhost`). Fix anything that breaks.
- [ ] Own `e2e/`: run it, read it, extend it with your own trap cases (join twice, forged or reused invite, group deleted before a scan).
- [ ] Explain-back: *how the server decides who receives a group message*; *why scanning a QR is not the same as being added (R10)*.
- [ ] Review P4's part.

## P4: AI, deployment, slides, demo
- [ ] Explain-backs still owed from Phase 1 (`go` removal, why the mutex in `phase1/hub`, TCP lifecycle) and Phase 2 (HTTP vs WebSocket).
- [ ] Get an `ANTHROPIC_API_KEY` and run the **real** summary once (`docs/RUNNING.md`). It has never run live. Check the cost and decide on `LLM_MODEL`.
- [ ] Read `server/summary.go` and `onSummarize`. Explain-back: *trace one summary from trigger to screen, naming the protocol at each hop*; *what must not block*; *what is sent to the third party*.
- [ ] Phase 8 (R1): run the server on one laptop and connect from two other physical computers and a phone. Fix firewall/IP problems. Optional: cloud deploy with `wss://`.
- [ ] Rehearse `docs/DEMO_SCRIPT.md` with a timer.
- [ ] Review P1's part.

## Known gaps (not anyone's fault, but must be known)
- Live Claude call never executed (no key). Real phone scan and OS notifications never executed (headless test only).
- No message history; state is lost on disconnect; no authentication.
- `go test -race` could not run on the first machine (needs a 64-bit C compiler).
- There is no CI. Run `cd server && go test ./...` and `cd client-web && npm run build` before every PR.
