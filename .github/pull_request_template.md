## What this PR does
<!-- 1 to 3 sentences. Say WHY, not just what. Keep PRs small (under ~200 lines of real change). -->

## Whose module / which requirement
- Owner: P1 / P2 / P3 / P4
- Touches: R? / F? / docs / tests
- Phase: ?

## What I learned (author)
<!-- In your own words. If you cannot write this, the PR is not ready (README Section 2, rule 2). -->

## Explain-back (author answers, reviewer checks)
- [ ] I can explain every changed line of server or socket code without looking at it.
- Question 1 (from the phase in README Section 5):
- Question 2:
- Reviewer's own question 1:
- Reviewer's own question 2:

## Checks
- [ ] `cd server; go test ./...` passes
- [ ] `cd client-web; npm run build` passes
- [ ] Tried with 3 browser tabs, including a "should not see it" tab (`e2e/` or by hand)
- [ ] Tried on 2 physical computers (required at least once per phase from Phase 4; tick "n/a" for docs-only PRs)
- [ ] `docs/PROTOCOL.md` updated first if a message type or field changed, and the team acknowledged it
- [ ] Learning-log entry added to `docs/LEARNING_LOG.md`
- [ ] README Section 0 (tracker) updated

## Safety (README Section 1.6)
- [ ] Privacy or membership rules are enforced on the **server**, not only hidden in the UI
- [ ] No secrets, `.env` files or API keys in the diff
- [ ] Chat traffic goes over the WebSocket only (HTTP only for static files and the upgrade)
- [ ] User text is rendered as text (no `dangerouslySetInnerHTML` with user content)

## Who wrote it
- [ ] I wrote this myself
- [ ] Parts were written with an AI tool or by a teammate (say which files, and that I have read and understood them)

## Reviewer
<!-- Someone who did NOT write this code. Review pairing round 1: P1 reviews P2, P2 reviews P3, P3 reviews P4, P4 reviews P1. -->
Reviewer:
