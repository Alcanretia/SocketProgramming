# LEARNING_LOG.md

3 lines per session: what I built / what confused me / what clicked.

## 2026-10-29: P4, Phase 1 (with an AI mentor)
- Built: ran a TCP echo server, then a version that serves 2 clients with `go handle(conn)`, then a mini chat hub that forwards a line to everyone else (`server/phase1/`), from mentor-provided worked examples.
- Confused: what `Listen` and `Accept` each return, and why `go` is needed.
- Clicked: `Accept` and `Read` both *block*; one goroutine per connection stops one quiet client from freezing the rest. (Still owed: explain-back on the hub's mutex and `go` removal.)

## 2026-10-29: P4, the rest was generated, not learned
- The Phase 2 to 7 code (`server/*.go`, `client-web/`) was produced by an AI helper in one sitting. It has tests, but **P4 has not yet studied it**. Next entry should be written after reading `docs/ARCHITECTURE.md` and answering its explain-backs in your own words.

## 2026-11-01: P2, Phase 1 + Phases 4, 5, 7 (F3) & PROTOCOL.md review
- Built: reviewed `onRegister`, `onLeave`, `dmRoom`, `deliver`, `onDM`, `validContent`, `stickers.go`, and `client-web/src/features/users-dm/`; answered open questions in `docs/PROTOCOL.md` and extended `TestRegisterAndUserList` with unicode/emoji, 500-char, and control-char name cases.
- Confused: how `onRegister` prevents duplicate names without a `sync.Mutex` when 10 clients register at the exact same millisecond, and why both `alice→bob` and `bob→alice` map to the same room on both client and server.
- Clicked: all events pass through the single `Hub.Run` goroutine over `h.events`, making the `byKey` check-then-set strictly sequential and race-free; `dmRoom(a, b)` sorts lowercased keys alphabetically (`dm:alice|bob`) and `onDM` stamps `From: c.name` on the server before calling `h.deliver` for only the sender and receiver.
