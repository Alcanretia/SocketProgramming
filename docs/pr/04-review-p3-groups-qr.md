## Heads-up for P3
While building the app, P4 also wrote the P3 part with an AI tool: create / list / join groups with member lists, group message fan-out, invite tokens, the group panel, the QR code and the join-confirmation dialog (`server/hub.go`, `client-web/src/features/groups/`), plus the multi-client browser test (`e2e/`). This PR is P3's review of that code, two bug fixes, and the follow-up work. Notes file: `docs/reviews/p3-groups-qr.md`.

## What this PR does
- **Client fix 1: stale QR after reconnect.** Invite tokens cached in `inviteCreated` survived a dropped connection. After a server restart, group IDs restart at `g_1`, so "Invite (QR)" could show a dead link or a link to a different group. The cache is now cleared on disconnect, together with users, groups and unread.
- **Client fix 2: stuck join dialog.** If the group was deleted between "Join group X?" and pressing Join, the server answers `no_such_group`; the dialog ignored it and every Join press only added an error toast. That error now closes a pending invite, the same way `bad_invite` does.
- Read the group and invite handlers (`onCreateGroup` to `onInviteInfo`, and the group part of `onLeave`) and `client-web/src/features/groups/`. No server change was needed.
- Answered the five review questions in `docs/reviews/p3-groups-qr.md`.
- Tested the QR invite with a real phone on the same Wi-Fi via the LAN address: **passed** (2026-11-02).
- Added 12 browser checks to `e2e/e2e.mjs` for invite edge cases (table below).
- Wrote the group and invite rules into `docs/DECISIONS.md` (D7), added a learning-log entry, updated README Section 0.

New e2e cases (section "QR invite trap cases (P3)"):

| Case | What it proves |
|---|---|
| Hana opens the Study link and presses **Not now** | Declining does not add her (R10); the token is cleared from the URL; the same link worked for Dave first, so invites are reusable (D7) |
| Eve opens `/?invite=deadbeef` | A forged token gets a `bad_invite` toast, no dialog, and is removed from the URL |
| Frank creates "Temp", makes an invite, then leaves; Gina opens the link | The empty group disappears for everyone; the stale invite gets a "no longer exists" toast, no dialog |

## Whose module / which requirement
- Owner: P3
- Touches: R5, R8, R9, R10, R11, F2, the multi-client test script
- Files changed: `client-web/src/socket/store.tsx` (2 fixes), `e2e/e2e.mjs`, `docs/reviews/p3-groups-qr.md`, `docs/DECISIONS.md`, `docs/LEARNING_LOG.md`, `README.md`, this file

## What I learned (author)
- **Fan-out:** `onGroupMsg` checks `g.has(c.key)` first, then loops over `g.Members` and delivers to every member who is connected (looked up in `h.byKey`), including the sender. A non-member gets `not_member` back and nobody else sees the attempt.
- **Scanning is not joining (R10):** opening the link only sends `invite_info`, which is read-only. `join_group` is a separate message sent when the user presses Join, and `onJoinGroup` only ever adds the sender, so nobody can add someone else.
- **Why a token, not the group ID:** group IDs are reused after a server restart, so a `?group=g_1` link could silently lead to the wrong group; a stale token just gets `bad_invite`. The token is not access control, because R9/R10 make groups open.
- **Disconnect:** `onLeave` removes the user from every group and broadcasts the new `group_list`; a group left with 0 members is deleted with its invite tokens. A page refresh counts as a disconnect.
- Two groups can never get the same ID because only the hub goroutine touches `groupSeq`, one event at a time.

## Explain-back (author answers, reviewer checks)
- [x] I can explain every line of my part.
- Q1: How does your server decide who receives a group message? (`onGroupMsg`, loop over `g.Members`)
- Q2: Why isn't scanning the QR the same as being added by the creator (R10)?
- Reviewer's questions (P2 reviews this PR):
  1. If a non-member sends a group message, what does the server do, and does anyone else see the attempt?
  2. What happens to membership when a member disconnects, and to a group when everyone has left? Where is that decided?

## How to test the fixes
1. **Stale QR:** open the app, create a group, press "Invite (QR)". Stop and restart the server, wait for the page to reconnect, register, create a group again, press "Invite (QR)". Expected: a new token is created (the link differs from before) and scanning it opens "Join group X?".
2. **Stuck dialog:** user A creates "Temp" and opens its QR. In another browser, open the link and stop at "Join group Temp?". Close A's tab (Temp is deleted). Press Join. Expected: an error toast, the dialog closes, `?invite=` is removed from the URL.

## Checks
- [x] `cd server; go test ./...` passes (13 tests; `go test -race ./...` clean). No server code changed in the fix commit.
- [x] `cd client-web; npm run build` passes after the fixes
- [x] 3-tab test: A creates, B joins, C does not; C sees nothing and cannot send
- [x] **Scanned the QR with a real phone** on the same Wi-Fi (LAN address, not `localhost`): passed, 2026-11-02
- [x] `e2e/` run and extended with my own cases: **47/47 pass**. Joining twice is covered by `TestGroupsMembershipAndPrivacy` in `server/hub_test.go`.
- [x] Re-run `e2e/` after the client fixes: all pass
- [x] Group rules written in `docs/DECISIONS.md` (D7). `docs/PROTOCOL.md` 3.3 already states them. The invite rules (reusable, no expiry) still need a line in 3.5; left to P2 as the `PROTOCOL.md` keeper.
- [x] Learning-log entry added; README Section 0 updated

## Safety
- [x] Membership is enforced on the server (`g.has(c.key)` in `onGroupMsg`, `onCreateInvite` and `roomAccess`; only `onJoinGroup` adds members, and only the sender)
- [x] Invite tokens are random and unguessable (`crypto/rand`, 128 bits); a token for another group is rejected (`bad_invite`)
- [x] No secrets in the diff (`.env` is ignored and not tracked)

## Notes for the reviewer
- The "You are already a member / Open group" branch in `InviteJoinDialog` cannot be reached from the UI: opening a link starts a new connection, which is in no groups. Harmless; left in.
- Disconnecting, including a page refresh, drops all group memberships. **Do not refresh during the demo recording.**

## Who wrote it
- [x] The original group and invite code, the QR UI and `e2e/` were written by P4 with an AI tool; I have read them and can explain them.
- [x] The two client fixes were found and drafted with an AI assistant (Claude); I reviewed them and can explain both.
- [x] The new e2e cases were drafted with an AI assistant; I reviewed and ran them and can explain each check.
- [x] The review answers, D7 decisions and learning-log entry are mine; some wording was drafted with AI help.

## Reviewer
P2
