# LEARNING_LOG.md

3 lines per session: what I built / what confused me / what clicked.

## 2026-10-29: P4, Phase 1 (with an AI mentor)
- Built: ran a TCP echo server, then a version that serves 2 clients with `go handle(conn)`, then a mini chat hub that forwards a line to everyone else (`server/phase1/`), from mentor-provided worked examples.
- Confused: what `Listen` and `Accept` each return, and why `go` is needed.
- Clicked: `Accept` and `Read` both *block*; one goroutine per connection stops one quiet client from freezing the rest. (Still owed: explain-back on the hub's mutex and `go` removal.)

## 2026-10-29: P4, the rest was generated, not learned
- The Phase 2 to 7 code (`server/*.go`, `client-web/`) was produced by an AI helper in one sitting. It has tests, but **P4 has not yet studied it**. Next entry should be written after reading `docs/ARCHITECTURE.md` and answering its explain-backs in your own words.

## 2026-11-05: P1 (Connection Core Review & Hardening)
- Built/Read: Studied `main.go`, `client.go`, and `hub.go` line by line until I could trace the full lifecycle. Ran the break-it checks (10 tabs, abrupt disconnects, malformed JSON, and oversized frames) and verified F4 desktop notifications on `localhost`.
- Confused: Why browser OS notifications work fine on `localhost` but silently fail when accessed from other devices over LAN.
- Clicked: The Web Notification API mandates a Secure Context; browsers treat `localhost` as secure, but plain `http://<LAN-IP>` is considered insecure, cleanly falling back to our in-app toasts, sound, and tab title badge. Also clicked why the hub uses a single event loop and dedicated client write pumps: to eliminate concurrent socket writes and ensure slow clients never block the server.
