# P2 review and learning notes: users, DMs, stickers

Branch `review/p2-users-dm`. Reviewer: P1. Fill this in with your own words; commit it to this branch.

> The code in your area, and the first draft of `docs/PROTOCOL.md` (which you keep), were written by P4 with an AI tool before you started. Read it, change what you disagree with, and answer below. It is not done until you can explain it.

## Read in this order
1. Do your own Phase 1 first.
2. `docs/PROTOCOL.md` completely. Then `server/protocol.go`.
3. `server/hub.go`: `onRegister`, `onLeave` (name freed), `dmRoom`, `deliver`, `onDM`; `server/stickers.go`, `validContent`.
4. `server/hub_test.go`: `TestRegisterAndUserList`, `TestSimultaneousRegistrationOnlyOneWins`, `TestDMPrivacyAndSpoofing`, `TestStickers`.
5. `client-web/src/features/users-dm/`.

## Answer in your own words
1. How does the server guarantee unique names even with two simultaneous registrations?

   Answer:

2. Show the exact line where a message is sent to only one person.

   Answer:

3. Can a client send a DM that claims to be from someone else? Where is that stopped?

   Answer:

4. Is the room between A and B the same as between B and A? How is it identified?

   Answer:

5. What exactly travels over the WebSocket when I send a sticker?

   Answer:

## `docs/PROTOCOL.md`
- [ ] I read all of it
- Open questions in section 5 (history? unread retention? mentions? AI bot user? rate limiting?), my answers:
- Things I want to change (doc first, then code and `client-web/src/socket/types.ts`):
- [ ] I agree with it and ticked my acknowledgement box

## Things I tried
- [ ] Tricky names: unicode, emoji, 500 characters, only spaces
- [ ] Third-tab privacy test
- [ ] Sticker artwork looks fine / replaced with my own
