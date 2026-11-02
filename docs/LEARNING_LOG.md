# LEARNING_LOG.md

3 lines per session: what I built / what confused me / what clicked.

## 2026-10-29: P4, Phase 1 (with an AI mentor)
- Built: ran a TCP echo server, then a version that serves 2 clients with `go handle(conn)`, then a mini chat hub that forwards a line to everyone else (`server/phase1/`), from mentor-provided worked examples.
- Confused: what `Listen` and `Accept` each return, and why `go` is needed.
- Clicked: `Accept` and `Read` both *block*; one goroutine per connection stops one quiet client from freezing the rest. (Still owed: explain-back on the hub's mutex and `go` removal.)

## 2026-10-29: P4, the rest was generated, not learned
- The Phase 2 to 7 code (`server/*.go`, `client-web/`) was produced by an AI helper in one sitting. It has tests, but **P4 has not yet studied it**. Next entry should be written after reading `docs/ARCHITECTURE.md` and answering its explain-backs in your own words.

## 2026-11-2: P3, Phase 2-7
Built: reviewed groups and QR, testing QR invitation with phone (passed), review additional e2e cases that are written by AI
Confused: why the token exists if anyone can join, or why localhost breaks the QR.
Clicked: the single hub goroutine makes the check-then-act safe, and invite_info vs join_group is how R10 is kept, the effect of user disconnecting from the server. 
