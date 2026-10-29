## What this PR does
P4's learning, live-test and deployment changes: <!-- e.g. explain-back answers in the learning log, result of the first live AI summary, multi-computer test notes, slides, demo rehearsal fixes -->

## Whose module / which requirement
- Owner: P4
- Touches: R1, F1, Phase 1, 2, 8 and 9
- Files: `server/summary.go`, `server/phase1/`, `client-web/src/features/ai/`, `docs/RUNNING.md`, `docs/DEMO_SCRIPT.md`, slides

## What I learned (author)
<!-- Your own words. Cover: Accept vs conn, why `go`, why the mutex in phase1/hub, HTTP vs WebSocket, how one summary travels. -->

## Explain-back (author answers, reviewer checks)
- [ ] I can explain `summary.go` and `onSummarize` without looking.
- Q1: Why does a chat server need one goroutine per connection? What happens if you remove the `go` in front of `handle(conn)`?
- Q2: Why does `phase1/hub` need `mu.Lock()` around `clients`? Give an example of what could go wrong without it.
- Q3: Trace one summary from the trigger to the screen. Which protocol is involved at each hop? What must not block? What is sent to the third party?
- Q4: What is the difference between an HTTP request and a WebSocket connection? Where does one turn into the other?
- Reviewer's questions (P3 reviews this PR):
  1. What happens if the LLM takes 5 seconds, or fails?
  2. Which messages are you allowed to send to the LLM, and how does the code make sure of it?

## Checks
- [ ] `cd server; go test ./...` passes
- [ ] **Real AI summary run once** with `ANTHROPIC_API_KEY` (result and cost noted; model chosen via `LLM_MODEL`)
- [ ] Phase 8: server on one laptop, two other physical computers and a phone connected; firewall and IP steps worked
- [ ] `docs/DEMO_SCRIPT.md` rehearsed with a timer
- [ ] Learning-log entry added
- [ ] README Section 0 updated

## Safety
- [ ] The API key is only in a server environment variable, never in the repo or the browser
- [ ] Only the requesting user's own unread messages are sent to the LLM
- [ ] No secrets in the diff

## Who wrote it
- [ ] I wrote my changes myself
- [x] The original `summary.go` and `features/ai/` were written with an AI tool; I have read them and can explain them

## Reviewer
P3
