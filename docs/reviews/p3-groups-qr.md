# P3 review and learning notes: groups and QR invites

Branch `review/p3-groups-qr`. Reviewer: P2. Fill this in with your own words; commit it to this branch.

> The code in your area, and the multi-client test `e2e/` that was yours to write, were written by P4 with an AI tool before you started. Read it, change what you disagree with, and answer below. It is not done until you can explain it.

## Read in this order
1. Do your own Phase 1 first.
2. `server/hub.go`: `Group`, `onCreateGroup`, `onJoinGroup`, `onGroupMsg`, `onCreateInvite`, `onInviteInfo`, and the group part of `onLeave`.
3. `server/hub_test.go`: `TestGroupsMembershipAndPrivacy`, `TestInviteFlow`.
4. `client-web/src/features/groups/` (`GroupPanel`, `InviteQr`, `InviteJoinDialog`).
5. `e2e/e2e.mjs` and `e2e/README.md`.

## Answer in your own words
1. How does the server decide who receives a group message?

   Answer:

2. If a non-member sends a group message, what does the server do, and does anyone else see it?

   Answer:

3. Why isn't scanning the QR the same as being added by the creator (R10)?

   Answer:

4. What does the QR encode, and why a random token instead of the group id?

   Answer:

5. What happens to membership when a member disconnects, and to a group when everyone has left?

   Answer:

## Rules I decide (my proposals in the code are only defaults)
- Group names unique? Empty allowed?
- Disconnect leaves all groups, or keep membership?
- Empty group deleted?
- Invite expiry or revocation?
Write the decisions in `docs/PROTOCOL.md` 3.3 and `docs/DECISIONS.md` (D7), then change the code.

## Things I tried
- [ ] 3-tab test: A creates, B joins, C does not
- [ ] **Real phone scan** on the same Wi-Fi (app opened via the LAN address): result:
- [ ] Ran `e2e/` and added my own cases (join twice, forged or reused invite, group deleted before scan)
