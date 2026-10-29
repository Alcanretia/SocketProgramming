## What this PR does
P3's review and learning changes for groups and QR invites: <!-- e.g. my Phase 1 exercise, group rules I decided, real phone scan result, new e2e cases, learning-log entry -->

## Whose module / which requirement
- Owner: P3
- Touches: R5, R8, R9, R10, R11, F2, the multi-client test script
- Files: `onCreateGroup` to `onInviteInfo` in `server/hub.go`, `client-web/src/features/groups/`, `e2e/`

## What I learned (author)
<!-- Your own words. Cover: how group fan-out works, why scanning is not joining, what happens on disconnect. -->

## Explain-back (author answers, reviewer checks)
- [ ] I can explain every line of my part.
- Q1: How does your server decide who receives a group message? (`onGroupMsg`, loop over `g.Members`)
- Q2: Why isn't scanning the QR the same as being added by the creator (R10)?
- Reviewer's questions (P2 reviews this PR):
  1. If a non-member sends a group message, what does the server do, and does anyone else see the attempt?
  2. What happens to membership when a member disconnects, and to a group when everyone has left? Where is that decided?

## Checks
- [ ] `cd server; go test ./...` passes
- [ ] `cd client-web; npm run build` passes
- [ ] 3-tab test: A creates, B joins, C does not; C sees nothing and cannot send
- [ ] **Scanned the QR with a real phone** on the same Wi-Fi (app opened via the LAN address, not `localhost`); result written down
- [ ] `e2e/` run and extended with my own cases (join twice, forged or reused invite, group deleted before a scan)
- [ ] Group rules I chose are written in `docs/PROTOCOL.md` section 3.3 and `docs/DECISIONS.md` (D7)
- [ ] Learning-log entry added
- [ ] README Section 0 updated

## Safety
- [ ] Membership is enforced on the server
- [ ] Invite tokens are random and unguessable; a token for another group is rejected
- [ ] No secrets in the diff

## Who wrote it
- [ ] I wrote my changes myself
- [x] The original group and invite code, the QR UI and `e2e/` were written by P4 with an AI tool; I have read them and can explain them

## Reviewer
P2
