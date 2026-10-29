## What this PR does
P2's review and learning changes for identity, DMs and stickers: <!-- e.g. my Phase 1 exercise, answers to PROTOCOL.md open questions, fixes, learning-log entry -->

## Whose module / which requirement
- Owner: P2
- Touches: R3, R4, R5, R7, F3, `docs/PROTOCOL.md` (P2 is the keeper)
- Files: `onRegister`, `onDM`, `validContent`, `server/stickers.go`, `server/protocol.go`, `client-web/src/features/users-dm/`, `client-web/public/stickers/`

## What I learned (author)
<!-- Your own words. Cover: how uniqueness holds, how a DM reaches exactly two people, what a sticker message contains. -->

## Explain-back (author answers, reviewer checks)
- [ ] I can explain every line of my part.
- Q1: How does the server guarantee uniqueness even with two simultaneous registrations? (`TestSimultaneousRegistrationOnlyOneWins`)
- Q2: Show me the exact line where a message is sent to only one person. (the two `deliver` calls at the end of `onDM`)
- Q3: What exactly travels over the WebSocket when I send a sticker?
- Reviewer's questions (P1 reviews this PR):
  1. Can a client send a DM that claims to be from someone else? Where is that stopped?
  2. Is the room between A and B the same as between B and A? How is it identified?

## Checks
- [ ] `cd server; go test ./...` passes
- [ ] `cd client-web; npm run build` passes
- [ ] Third-tab privacy test done (a third tab must never see an A and B message)
- [ ] Tried tricky names (unicode, emoji, 500 characters, only spaces, same name at the same instant)
- [ ] **`docs/PROTOCOL.md`**: open questions in section 5 answered, and the acknowledgement box ticked **only if** I agree. Any change went into the doc first, then the code and `client-web/src/socket/types.ts`
- [ ] Learning-log entry added
- [ ] README Section 0 updated

## Safety
- [ ] Sender is stamped by the server (spoofing test)
- [ ] Unknown sticker ids are rejected by the server
- [ ] No secrets in the diff

## Who wrote it
- [ ] I wrote my changes myself
- [x] The original code and the `PROTOCOL.md` draft were written by P4 with an AI tool; I have read them and can explain them

## Reviewer
P1
