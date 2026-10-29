## What this PR does
P1's review and learning changes for the connection core: <!-- e.g. my Phase 1 and Phase 2 exercises, fixes found while reading hub.go, learning-log entry -->

## Whose module / which requirement
- Owner: P1
- Touches: R2, R6, F4, Phase 1 to 3
- Files: `server/hub.go`, `server/client.go`, `client-web/src/socket/`, `client-web/src/components/shared/`

## What I learned (author)
<!-- Your own words. Cover: what one Client is, who may touch shared state, why one writer goroutine, why enqueue never blocks. -->

## Explain-back (author answers, reviewer checks)
- [ ] I can explain every line of `hub.go` and `client.go`.
- Q1: Walk me through what happens from the moment a browser connects until it disconnects. (`ServeWS`, `readPump`, `onRegister`, `onLeave`)
- Q2: Why does each client need its own write loop and outgoing channel? What happens to a slow client?
- Reviewer's questions (P4 reviews this PR):
  1. Why is `Client.send` never closed? What would go wrong if it were?
  2. What would you change if you used a mutex instead of the owner goroutine, and what does it cost?

## Checks
- [ ] `cd server; go test ./...` passes
- [ ] `cd client-web; npm run build` passes
- [ ] Tried 10 tabs, killed a tab mid-message, sent a huge message and malformed JSON (write down what you saw)
- [ ] Tried OS notifications on `localhost` (permission prompt, tab hidden) and noted the plain-HTTP LAN limitation
- [ ] Learning-log entry added to `docs/LEARNING_LOG.md`
- [ ] README Section 0 updated

## Safety
- [ ] Nothing here trusts the client for privacy
- [ ] No secrets in the diff

## Who wrote it
- [ ] I wrote my changes myself
- [x] The original `hub.go` / `client.go` / `socket/` were written by P4 with an AI tool; I have read them and can explain them

## Reviewer
P4
