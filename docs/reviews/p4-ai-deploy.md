# P4 review and learning notes: AI, deployment, demo

Branch `learn/p4-explain-back`. Reviewer: P3. Fill this in with your own words; commit it to this branch.

> Even this part was written with an AI tool. It is not done until you can explain it and the live checks below have been run.

## Owed from Phase 1 and 2
1. Why does a chat server need one goroutine per connection? What happens if you remove the `go` in front of `handle(conn)`?

   Answer:

2. Why does `server/phase1/hub` need `mu.Lock()` around `clients`? Give an example of what could go wrong without it.

   Answer:

3. Draw the lifecycle of a TCP connection (Listen, Accept, Read, Write, Close). Which side does each?

   Answer:

4. What is the difference between an HTTP request and a WebSocket connection? Where does one turn into the other?

   Answer:

## AI summary
1. Trace one summary from the trigger to the screen. Which protocol is involved at each hop?

   Answer:

2. If the LLM takes 5 seconds, what must not block? Which line makes sure of it?

   Answer:

3. Which messages are sent to the third party, and how does the code guarantee only messages this user may see?

   Answer:

## Live checks (never run yet)
- [ ] Real summary with `ANTHROPIC_API_KEY` (model, cost, result):
- [ ] Phase 8: server on one laptop; two other physical computers and a phone connected (firewall and IP steps):
- [ ] `docs/DEMO_SCRIPT.md` rehearsed with a timer (time taken):
- [ ] "Implementation Explanation" slide made

## Things I found or want to change
-
