## What this PR does
Adds the whole first working version of the chat app, split into 13 commits by area so each owner can review their own part: raw-TCP and WebSocket exercises, the Go chat server (hub, DMs, groups, invites, stickers, unread, AI summary), the React client, a multi-client browser test, and the docs (protocol, decisions, architecture, running guide, demo script, next steps).

It is bigger than the "under ~200 lines" rule on purpose: it is the baseline everything else builds on. Please review it **by commit**, not as one diff.

## Whose module / which requirement
- Owner: P4 wrote it (with an AI tool), including code that belongs to P1, P2 and P3. See README Section 0.7 for exactly what was written on whose behalf.
- Touches: R2 to R11, F1 to F4, docs, tests
- Phase: 1 to 7 (code); 8 not done; 9 partly

## How to review it (by commit)
| Commit area | Owner who should review it | Read these |
|---|---|---|
| `feat(phase1)`, `feat(phase2)` | everyone (try yours first) | `server/phase1/`, `server/phase2/wsecho` |
| `feat(server): hub, per-client pumps ...` | P1 (hub, pumps), P2 (identity, DM, stickers), P3 (groups, invites) | `server/hub.go`, `client.go`, `protocol.go`, `hub_test.go` |
| `feat(server): AI summary` | P4 | `server/summary.go` |
| `feat(client): WebSocket connection ... (P1 scope)` | P1 | `client-web/src/socket/`, `components/shared/` |
| `feat(client): name entry ... (P2 scope)` | P2 | `features/users-dm/`, `public/stickers/` |
| `feat(client): group panel ... (P3 scope)` | P3 | `features/groups/` |
| `feat(client): unread list ... (P4 scope)` | P4 | `features/ai/` |
| `test(e2e)` | P3 | `e2e/` |
| `docs: ...` | all | `docs/`, README Section 0 |

## What I learned (author)
<!-- P4: fill this in after reading docs/ARCHITECTURE.md. Until you can write it in your own words, keep this PR as a draft. -->

## Explain-back (author answers, reviewer checks)
- [ ] The author can explain the server without looking at the code.
- Q1: Walk through a browser connecting until it disconnects.
- Q2: How does the server guarantee unique names with simultaneous registrations?
- Reviewer's questions:

## Checks
- [x] `cd server; go test ./...` passes (hub tests drive real WebSocket clients: R3 to R11, stickers, invites, unread, summary failure and slow-summary cases)
- [x] `cd client-web; npm run build` passes
- [x] 35-check browser test in `e2e/` passes against a running server, including third-tab-sees-nothing checks
- [ ] Tried on 2 physical computers: **not done** (Phase 8)
- [ ] `docs/PROTOCOL.md` acknowledged by the team: **not yet** (it is a proposal, v0.1)
- [x] Learning-log entry added
- [x] README Section 0 updated

## Not verified (please do not assume these work)
- The live Claude call (no API key was available; only a fake summarizer is tested).
- A real phone scanning the QR code; OS-level notifications (browsers need HTTPS or localhost plus permission).
- `go test -race` (no 64-bit C compiler on the first machine).

## Safety
- [x] Membership and privacy are enforced in `server/hub.go`, not in the UI
- [x] No secrets in the diff; `.env` is git-ignored; the API key is read from the server environment only
- [x] Chat goes over the WebSocket only; HTTP serves static files and the upgrade
- [x] User text is rendered as text

## Who wrote it
- [ ] I wrote this myself
- [x] Written with an AI tool; includes code that belongs to P1, P2 and P3. Each owner must read and understand their part before it counts as done.

## Reviewer
Everyone. P1 reviews P2's parts, P2 reviews P3's, P3 reviews P4's, P4 reviews P1's (README Section 3).
