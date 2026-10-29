# SocketProgramming

Team Workflow + Agent Starter: Socket Chat Project (Go + React)

A learning-first plan for a team of 4, **written so each teammate can hand this file to their own AI agent**. The goal is not only to pass the rubric but to **be able to explain every line of socket code to the instructor**.

---

## 0. Progress Tracker (update this every session)

> **Last updated:** 2026-10-29
> **Current phase:** the whole app (Phases 2 to 7) is **built and auto-tested by P4 with an AI helper**, ahead of the team. Learning/review status is behind the code: nobody has done the explain-backs, and P1 to P3 have not read it yet.
> **Next up (needs people, in this order; per-person checklists are in `docs/NEXT_STEPS.md`):**
> 1. Team reads and acknowledges `docs/PROTOCOL.md` (DRAFT v0.1) and `docs/DECISIONS.md` (D1 to D9). Change them if you disagree, then change the code.
> 2. P1 to P3 do the Phase 1 exercises themselves **before** reading `server/phase1/`, then every member reads `docs/ARCHITECTURE.md` and answers the explain-backs.
> 3. Provide an `ANTHROPIC_API_KEY` and try the real AI summary once (never run live; only a fake summarizer is tested).
> 4. Phase 8: run on 2+ physical computers and a phone (`docs/RUNNING.md`), including the QR scan and the firewall step.
> 5. Rehearse with `docs/DEMO_SCRIPT.md`; split reviews so every PR-sized piece gets a non-author reviewer.
>
> **Blockers (only people can do these):** API key for F1 live test; physical computers for R1; teammates' acknowledgement and review; instructor-style Q&A rehearsal.
> **Open questions:** see `docs/PROTOCOL.md` section 5 (history, unread retention, mentions, AI bot user).

**Legend:** ⬜ not started · 🟡 in progress · ✅ done

### 0.1 Phase overview

| # | Phase | Lead | Status | Finished on | Notes |
|---|---|---|---|---|---|
| 1 | Socket fundamentals (TCP echo, 2 clients, forwarding) | All | ⬜ | | |
| 2 | WebSocket basics + protocol contract (`PROTOCOL.md`) | P1 leads, all | 🟡 | | Draft by P4: WS echo server + tested handshake, `PROTOCOL.md` DRAFT v0, `DECISIONS.md`. Needs team ack. Frontend agents start after this |
| 3 | Connection core: client and hub | P1 | 🟡 | | Code done: `server/hub.go`, `client.go` (one owner goroutine, read/write pumps). Needs P1 to read it and pass the explain-back |
| 4 | Identity and user list (R3, R4) | P2 | 🟡 | | Code + tests done (case-insensitive, simultaneous registration). Needs P2 review |
| 5 | Direct messages and rooms (R5, R7) | P2 | 🟡 | | Code + tests + 3-tab browser test done. Needs P2 review |
| 6 | Groups (R8 to R11) | P3 | 🟡 | | Code + tests + 3-tab browser test done. Needs P3 review |
| 7 | Special features (Part 3) | one owner each | 🟡 | | F2, F3, F4 built and tested. F1 built; live Claude call **not** tested (no API key) |
| 8 | Multi-computer test + R1 | P1 + P4 | ⬜ | | Blocked: needs 2+ physical computers (see `docs/RUNNING.md`) |
| 9 | Hardening and rehearsal | All | 🟡 | | Partly: malformed JSON, oversize frame, simultaneous names, slow summary, privacy tests exist. The "trap hunt" and mock Q&A need people |

### 0.2 Rubric checklist (tick only when verified with the 3-tab test)

> `[x]` below means: implemented and verified by the automated multi-client tests on ONE computer (`server/hub_test.go` and `e2e/`). It does **not** mean the team has reviewed it or can explain it yet (see Section 8, Definition of done).

**Part 1**
- [ ] **R1:** 2+ clients on different physical computers; each protocol explained (cloud = bonus). Protocol explanations are written (Section 7, `docs/ARCHITECTURE.md`); the multi-computer run is not done
- [x] **R2:** chat messages travel only over sockets (HTTP serves static files and the upgrade only)

**Part 2**
- [x] **R3:** unique client names
- [x] **R4:** every client sees all connected names, including their own
- [x] **R5:** separate chat room for each DM and each group
- [x] **R6:** chat box + message window per room
- [x] **R7:** DM visible only to sender and receiver
- [x] **R8:** create a group; creator is the only initial member
- [x] **R9:** list of all groups with their member lists
- [x] **R10:** users join by themselves, never auto-added
- [x] **R11:** group messages visible only to members

**Part 3 (1 point each)**
- [ ] **F1:** AI text summary on unread chats (P4). Built; the failure path is tested; **the live Claude call has not been run (needs an API key)**
- [x] **F2:** QR code group invitation (P3). Scan flow tested via the invite link; a real phone camera scan is not tested
- [x] **F3:** default stickers (P2)
- [x] **F4:** notification on received messages (P1). In-app toast/sound/title badge tested; OS notifications need HTTPS or localhost and permission (not tested headless)

### 0.3 Per-person status

| Person | Role | Current task | Status |
|---|---|---|---|
| P1 | Connection Core | Do Phase 1 yourself, then read `server/hub.go`, `client.go`, `client-web/src/socket/`; notifications (F4) | ⬜ |
| P2 | Users, DMs, Stickers | Do Phase 1 yourself, then read `onRegister`, `onDM`, stickers; own `PROTOCOL.md` | ⬜ |
| P3 | Groups + QR Invite | Do Phase 1 yourself, then read the group and invite handlers; own `e2e/` | ⬜ |
| P4 | AI + Deployment | Phase 1: task 1 ✅; tasks 2 and 3 built via worked examples, explain-back pending | 🟡 |

### 0.4 Setup checklist

- [x] Git repo created
- [x] Repo layout created: `/server`, `/client-web`, `/docs`
- [x] Go module created (`server/go.mod`, `socketchat/server`)
- [x] `.gitignore` created (with `.env`) **before the first real commit**
- [x] `docs/PROTOCOL.md` created (P2 keeps it)
- [x] `docs/DECISIONS.md` created
- [x] `docs/LEARNING_LOG.md` created
- [ ] Each member has been assigned P1 to P4

### 0.5 Decisions to record in `docs/DECISIONS.md`

- [ ] Go WebSocket library (🟡 proposed: gorilla/websocket v1.5.3, awaiting team)
- [ ] TypeScript or JavaScript for the frontend
- [ ] Vite (or other build tool)
- [ ] Hub concurrency design (channel-owning goroutine vs mutex)
- [ ] Unread-tracking design and owner (shared by F1 and F4; P1 and P4 agree first)
- [ ] Sticker representation (new message type or field)
- [ ] Invite token design (QR payload)
- [ ] Group lifecycle rules (member disconnects, empty groups, duplicate join)

### 0.6 Session log

| Date | Session | What we did | Next |
|---|---|---|---|
| 2026-10-29 | 0 | Repo created, workflow written into README, folders + `.gitignore` + docs stubs scaffolded (Go 1.26.5 installed) | Start Phase 1 |
| 2026-10-29 | 2 | **P4 (with an AI helper) built beyond P4's own scope:** the whole chat server (`server/*.go`: hub, DM, groups, invites, stickers, unread, AI summary), the whole React app (`client-web/`, all four feature folders), Go tests, an e2e browser test (`e2e/`), and docs (`ARCHITECTURE`, `RUNNING`, `DEMO_SCRIPT`). All tests pass on one computer. This bypasses the "each person builds and explains their own module" rule, so P1 to P3 must read, review and explain their parts before they count as done | Team review and acknowledgement; API key for the live AI test; multi-computer run |
| 2026-10-29 | 1 | P4: Phase 1 worked examples (`server/phase1/echo`, `hub`, `client`), tested with 2 to 3 clients. P4 then drafted Phase 2 ahead: `server/phase2/wsecho` (test passes), `PROTOCOL.md` DRAFT v0, `DECISIONS.md` | Team ack of protocol + decisions; P1 to P3 do Phase 1; P4 explain-backs |

### 0.7 Work done on teammates' behalf (read this, P1 to P3)

P4 built the whole first version with an AI tool, so **code in your area already exists and you did not write it**. It is on branch `feature/app-p4-initial` (13 commits, split by area). Nothing below counts as "done" for you until you have read it, can explain it, and a non-author reviewed it. You may rewrite any of it; if you change a message type or field, change `docs/PROTOCOL.md` first.

| Owner | What was written or changed for them | Where | What the owner should do |
|---|---|---|---|
| **P1** Connection core | Hub (one owner goroutine + event channel), per-client read/write pumps, slow-client drop, ping/pong, size limit, WebSocket upgrade; browser connection with reconnect and re-register, state store, notification service with in-app fallback, shared `ChatRoom`, `NotificationSetting`, `roomInfo` | `server/hub.go`, `server/client.go`, `server/main.go`, `client-web/src/socket/`, `client-web/src/components/shared/` | Do Phase 1 and 2 yourself, read these files, pass the Phase 3 explain-back, test OS notifications |
| **P2** Users, DMs, stickers | Name validation and uniqueness, `user_list` broadcast, DM routing and room ids, sticker validation and the 8 original SVG stickers, name screen, user list, sticker picker. **Also wrote `docs/PROTOCOL.md` (you are the keeper)** | `onRegister`, `onDM`, `validContent` in `server/hub.go` and `server/protocol.go`, `server/stickers.go`, `client-web/src/features/users-dm/`, `client-web/public/stickers/`, `docs/PROTOCOL.md` | Read, answer the open questions in `PROTOCOL.md` section 5, acknowledge or change it, pass the Phase 4 and 5 explain-backs |
| **P3** Groups, QR | Create / list / join groups with member lists, group fan-out, disconnect rules, invite tokens and `invite_info`, group panel, QR code, join-confirmation dialog. **Also wrote `e2e/` (your "multi-client test script")** | `onCreateGroup` to `onInviteInfo` in `server/hub.go`, `client-web/src/features/groups/`, `e2e/` | Read, decide the group rules, scan the QR with a real phone, extend `e2e/`, pass the Phase 6 explain-back |
| **P4** (own scope) | AI summary (Anthropic SDK call, timeout, graceful errors), unread list and summary card, running guide, demo script, architecture doc | `server/summary.go`, `client-web/src/features/ai/`, `docs/RUNNING.md`, `docs/DEMO_SCRIPT.md`, `docs/ARCHITECTURE.md` | Live-test with an API key, multi-computer run, slides, explain-backs |
| **Shared** | Tests for R3 to R11, stickers, invites, unread and summary failures; `.gitignore`; Go module; decisions D1 to D9; per-person checklist; PR templates | `server/hub_test.go`, `.gitignore`, `server/go.mod`, `docs/DECISIONS.md`, `docs/NEXT_STEPS.md`, `.github/`, `docs/pr/` | Everyone: acknowledge decisions D1 to D9 |

**Decisions made for the team that you may overturn:** D1 to D9 in `docs/DECISIONS.md`. Two of them were supposed to be agreed first: D4 (hub concurrency design, P1 owns it) and D5 (unread tracking, P1 and P4 agree).

**Edits made to this workflow file itself** (the rules in Sections 1.1 to 1.6, 1.8, 2, 3, 4, 7, 8 and 9 were not changed):
- Added Section 0 (this tracker) and this subsection.
- "How to use this file": the agent opener now says `README.md`.
- Section 1.7 (session routine): agents now read Section 0 at the start and update it at the end.
- Phase sections: tasks, "Done when" and explain-back lines became checkboxes; each phase got a `Status:` line; Phase 1 got a per-person table; the Phase 2 protocol worksheet got a ✔ column; Section 6 got a "Done" column.

**What this changes about the plan:** the review pairing in Section 3 is no longer independent for the code that P4 wrote (for example P4 can no longer be the independent reviewer of P3's groups). Every phase from 3 to 7 is 🟡 instead of ✅ for that reason.

---

## How to use this file

1. This README is the team workflow file. Optionally also copy it to the repo root as `AGENTS.md` / `CLAUDE.md` (or into `/docs/`) so agents load it automatically.
2. Each teammate starts their agent with this opener (edit the bracketed parts):

```text
Read README.md fully before replying. I am [P1 / P2 / P3 / P4].
Follow Section 1 (Agent Instructions) strictly. Tell me which Phase we are in,
what my current task is, and start with the first guiding question.
```

3. Section 0 is the live status board; everything below Section 1 is shared context for the agent *and* the human team.

---

## 1. Agent Instructions (the agent must follow these)

You are the personal **mentor and frontend engineer** for ONE member of a 4-person university team building a socket-programming chat application (Computer Network term project). Your job is to make that person *understand* socket programming while the team ships a working app.

### 1.1 Identify the person
If the user hasn't said which role they are (P1 to P4, Section 3), ask once. Everything you do is scoped to their role and phase. Reply in the language the user writes in.

### 1.2 Two modes

| Mode | Applies to | Your behavior |
|---|---|---|
| **MENTOR** | The Go server, anything touching sockets, the protocol, concurrency, deployment, networking | **Guide, don't solve.** Follow the hint ladder below. |
| **BUILDER** | The React frontend, styling, build tooling, UI components, QR/sticker/notification UI | **Implement it fully**, following the team's `PROTOCOL.md`. Then explain the socket-related parts of what you wrote (how the browser `WebSocket` is used, how events are dispatched). |

If a task mixes both (e.g. "add a sticker feature"), build the UI part, and mentor the server part.

### 1.3 The hint ladder (MENTOR mode)
Go one level at a time. Move up only when the user has **tried and shown you their attempt**.

1. **Concept:** explain the idea with an analogy; no code.
2. **Guiding question:** use the questions in the phase, or ask your own.
3. **Structure:** pseudocode, a diagram, or a list of the steps and data structures.
4. **Isolated snippet:** a tiny example of *one* technique, on an unrelated toy problem (e.g. show a channel loop with a counter, not the chat hub).
5. **Code review:** when the user pastes their code, point to the *symptom and location* first ("look at what happens when two goroutines reach line X"), and let them find the fix.

**Never** hand over a full working implementation of a server feature at the start. If the user asks for "just the code":
- Remind them of the team rule, offer the next ladder level.
- If they have genuinely tried and still insist, you may give a worked example of **one function** with line-by-line explanation, label it `NEEDS TEAM REVIEW`, and follow with an explain-back quiz.

### 1.4 Explain-back before "done"
After each task, ask the person 2 to 3 questions (the phase has seeds). Don't mark the task done until they answer in their own words. Correct gaps kindly and specifically.

### 1.5 The protocol is the contract
- `/docs/PROTOCOL.md` is the single source of truth for message types and fields.
- Never invent or change an event, field, or semantics on your own. If a change is needed, tell the user to raise it with the team, update `PROTOCOL.md`, and get it acknowledged before coding.
- The frontend uses the plain browser `WebSocket` API (no Socket.IO). The server uses a Go WebSocket library (decide in Phase 2).

### 1.6 Quality and safety rules
- **Server is the authority.** Privacy and membership checks (R7, R10, R11) must happen on the server, not by hiding things in the UI. Tell the user if they're trusting the client.
- **Secrets:** LLM API keys live only in server environment variables; `.env` is in `.gitignore` before the first commit; never put keys in React code.
- **R2:** chat traffic (including AI replies, stickers, summaries) goes over the socket. REST/HTTP is allowed only for serving static files (React app, sticker images) and the WebSocket upgrade.
- Escape/sanitize message text in the UI (no `dangerouslySetInnerHTML` with user content).
- Keep PRs small, and say in the PR description what the author learned and what the reviewer should ask.

### 1.7 Session routine
- **Start:** read Section 0, state the current phase, the person's task, and what they should be able to explain by the end.
- **End:** write a 3-line entry for `LEARNING_LOG.md` (what I built, what confused me, what clicked), list open questions to raise with the team, and **update the checkboxes and "Last updated" in Section 0**.

### 1.8 Frontend ownership (BUILDER mode: avoid merge conflicts)

```
client-web/src/
  socket/              P1  connection manager, reconnect, event dispatcher, notification service
  components/shared/   P1  ChatRoom component (chat box + message window, R6), layout (changes need a PR)
  features/users-dm/   P2  name entry, user list, DM rooms, sticker picker + sticker rendering
  features/groups/     P3  group list/create/join, QR invite (show + join-by-link flow)
  features/ai/         P4  AI bot room, unread-summary UI
```

Only edit your own folder. If you need a change in someone else's, open an issue or ask them.

---

## 2. Ground rules (for humans)

1. **Learning beats speed.** If you can't explain a piece of backend code out loud, it isn't done.
2. **No copy-paste you don't understand.** Agents and tutorials are for *explanations*. Before merging, rewrite the idea in your own words.
3. **Explain-back before merge.** Every PR is reviewed by someone who did *not* write it. The author explains the code; the reviewer asks 2 questions. If the author can't answer, the PR waits.
4. **Rotate review.** Everyone owns one module but must read and review all the others.
5. **Keep a learning log.** One shared `docs/LEARNING_LOG.md`. This becomes raw material for your "Implementation Explanation" slide.
6. **Record decisions.** One shared `docs/DECISIONS.md` (Go WebSocket library, TypeScript or not, Vite, hub concurrency design, unread tracking design, etc.), each with *why*.

### Repo layout

```
/server        Go backend (separate process, own main file)
/client-web    React frontend
/docs          ARCHITECTURE, PROTOCOL.md, DECISIONS.md, LEARNING_LOG.md, slides
```

### Git flow

- `main` is always runnable. Work on `feature/<module>-<name>` branches.
- Small PRs (under ~200 lines). One non-author reviewer minimum.
- Commit messages say *why*, not just *what*.

---

## 3. Roles

| Person | Owns (backend) | Owns (frontend slice) | Rubric / special point |
|---|---|---|---|
| **P1: Connection Core** | WebSocket upgrade, per-client read/write loops, connection lifecycle, hub skeleton | `socket/`, shared `ChatRoom`, **notifications** | R2, R6, **Special: Notifications**; also R1 protocol explanation + architecture diagram |
| **P2: Users, DMs, Stickers** | Unique names, user list broadcast, DM routing, DM rooms, sticker messages | `features/users-dm/` | R3, R4, R5, R7, **Special: Stickers**; also keeper of `PROTOCOL.md` |
| **P3: Groups + QR Invite** | Create / list / join groups, member lists, group fan-out, invite handling | `features/groups/` | R5, R8, R9, R10, R11, **Special: QR invitation**; also multi-client test script |
| **P4: AI + Deployment** | LLM API integration, unread summary, (optional) AI bot user, multi-computer + optional cloud | `features/ai/` | **Special: AI summary (counts as the required AI feature)**; also deployment (R1 bonus), slides, demo script |

**Review pairing (round 1):** P1 reviews P2, P2 reviews P3, P3 reviews P4, P4 reviews P1. Shift by one each phase so everyone reads everything.

---

## 4. Rubric map

| Req | What the instructor will want to see |
|---|---|
| R1 | 2+ clients on different physical computers; explain each protocol (what, why, how); cloud is a bonus |
| R2 | Chat messages travel **only** over sockets |
| R3 | Unique client names |
| R4 | Every client sees all connected names, including their own |
| R5 | A separate chat room for each DM and each group |
| R6 | Chat box + message window per room |
| R7 | DM visible only to sender and receiver |
| R8 | Create a group; the creator is the only initial member |
| R9 | List of all groups with their member lists |
| R10 | Users join by themselves, never auto-added |
| R11 | Group messages visible only to members |
| Part 3 | 1 point per special feature: **AI summary, QR invitation, stickers, notifications** |

> Part 1 gates Part 2, and Part 2 gates Part 3. **Finish in order.** Prototype special-feature UIs early if you like, but don't integrate them until R3 to R11 work.

---

## 5. Phases

Each phase lists: **Goal**, **Concepts**, **Guiding questions**, **Hints** (directions, not answers), **Done when**, and **Explain-back**. A session is about 2 to 3 hours. Tick the checkboxes as you go, and mirror the result in Section 0.

---

### Phase 1: Socket fundamentals (all 4 together, 1 to 2 sessions)

**Status:** ⬜ not started

**Goal:** Everyone runs and modifies a raw socket program before touching WebSockets.

**Concepts:** IP address + port, TCP, listening vs connecting, client/server roles, blocking I/O, goroutines.

**Tasks**
- [ ] 1. Each person writes a tiny **TCP echo server** and **client** using Go's `net` package. Don't share code until everyone has tried. (P4 ✅ via worked example; P1 to P3 pending. Note: P4's example code exists in `server/phase1/`, so P1 to P3 should try before looking.)
- [ ] 2. Modify it so the server handles **two clients at the same time**.
- [ ] 3. Modify it again so a message from client A is **forwarded to client B**. This is a mini chat hub.

| Member | Task 1 | Task 2 | Task 3 |
|---|---|---|---|
| P1 | ⬜ | ⬜ | ⬜ |
| P2 | ⬜ | ⬜ | ⬜ |
| P3 | ⬜ | ⬜ | ⬜ |
| P4 | ✅ | 🟡 | 🟡 |

> P4 notes: tasks use mentor-provided worked examples marked `NEEDS TEAM REVIEW`: `server/phase1/echo` (tasks 1 and 2, tested with 2 simultaneous clients), `server/phase1/hub` (task 3, tested with 3 clients: sender not echoed, others receive), and `server/phase1/client`. Task 1 explain-back passed. Task 2 explain-back ("what if `go` is removed?") was explained by mentor; re-ask. Task 3 explain-back pending.

**Guiding questions**
- What does the server do between `Listen` and `Accept`? What does `Accept` return?
- If client 1 sends nothing for a minute, can client 2 still be served? Why or why not?
- TCP is a *stream*, not messages. If a client sends "hello" and "world" quickly, can the server read "helloworld" in one read? What does that mean for chat?
- What happens when a client disconnects? How does the server find out?

**Hints:** look up `net.Listen`, `Accept`, `Read`, `Write`, `Close`, and `go`. Test with two terminals and `nc localhost <port>`.

**Done when**
- [ ] Everyone can draw the lifecycle of a TCP connection (listen, accept, read, write, close) from memory.

**Explain-back:** *"Why does a chat server need one goroutine per connection (or something equivalent)?"*
- [ ] P1 · [ ] P2 · [ ] P3 · [ ] P4

---

### Phase 2: WebSocket basics and the protocol contract (all 4, 1 to 2 sessions)

**Status:** ⬜ not started

**Goal:** Move from raw TCP to WebSocket and agree on the message format. **The frontend agents build against this, so it must be complete.**

**Concepts:** HTTP upgrade handshake, why browsers can't open raw TCP, WebSocket frames, text vs binary, full-duplex.

**Tasks**
- [ ] 1. **P1 leads:** build a WebSocket **echo server** in Go. Connect from the browser console with `new WebSocket("ws://localhost:<port>/ws")`. (🟡 reference version exists at `server/phase2/wsecho` with a passing Go test; P1 should still build and explain their own.)
- [ ] 2. In DevTools, inspect the handshake: find the `Upgrade` header and the `101 Switching Protocols` response. (🟡 the Go test proves 101 + `Upgrade: websocket`; the DevTools look is still to do.)
- [ ] 3. **Everyone:** design the protocol with one JSON envelope for all messages, and record it in `/docs/PROTOCOL.md` (P2 keeps it). (🟡 DRAFT v0 written, awaiting team review.)
- [ ] 4. **Record team decisions** in `DECISIONS.md` (Go WS library, TypeScript or JS, Vite, unread-tracking owner). (🟡 D1 library proposed: gorilla/websocket; rest open.)

**Protocol worksheet** (fill in; no `?` may remain). Track each row in `PROTOCOL.md`; tick here once it is final. *Every row below now has a proposed answer in `PROTOCOL.md` DRAFT v0, but none is final until the team acknowledges it, so the ✔ column stays ⬜.*

| ✔ | Direction | Event name | Fields | Who receives it |
|---|---|---|---|---|
| ⬜ | client → server | `register` | ? | ? |
| ⬜ | server → client | `register_ok` / `register_error` | ? | ? |
| ⬜ | server → clients | `user_list` | ? | ? |
| ⬜ | client → server | `dm` | ? | ? |
| ⬜ | server → client | `dm` (delivered) | ? | ? |
| ⬜ | client → server | `create_group` | ? | ? |
| ⬜ | server → clients | `group_list` | ? | ? |
| ⬜ | client → server | `join_group` | ? | ? |
| ⬜ | client → server | `group_msg` | ? | ? |
| ⬜ | server → client | `group_msg` (delivered) | ? | ? |
| ⬜ | both | sticker message (new type, or a field on existing messages?) | ? | ? |
| ⬜ | both | invite / join-by-invite flow | ? | ? |
| ⬜ | both | unread state (who tracks it, how is it reported?) | ? | ? |
| ⬜ | client → server | `summarize_request` | ? | ? |
| ⬜ | server → client | `summary_result` / `summary_error` | ? | ? |
| ⬜ | server → client | `error` | ? | ? |

**Guiding questions**
- What fields must *every* message have?
- Who stamps the sender name: the client or the server? (Think about spoofing.)
- Which events are *requests*, and which are *broadcasts*?
- How will the frontend know which chat room a message belongs to?
- Notifications and AI summary both need "unread". **Who owns that state, the server or the client? What happens after a refresh or reconnect?**

**Hints:** start with `{ "type": "...", "payload": { ... } }` and add fields only when you can justify them.

**Done when**
- [ ] Echo works from the browser
- [ ] Protocol table is complete (no `?` left)
- [ ] Everyone agrees
- [ ] Decisions are logged

**Explain-back:** *"What is the difference between an HTTP request and a WebSocket connection? Where does one turn into the other?"*
- [ ] P1 · [ ] P2 · [ ] P3 · [ ] P4

---

### Phase 3: Connection core: client and hub (P1 leads, others pair-review)

**Status:** 🟡 code exists (written by P4 + AI, not by P1); review and explain-back pending

**Goal:** A server that accepts many WebSocket clients and can send a message to a specific one.

**Concepts:** Hub pattern, per-connection read and write loops, channels, shared state and concurrency, clean disconnects.

**Guiding questions**
- What is one `Client` in your program? What does it need to remember?
- Who may touch the map of connected clients? What if two goroutines modify it at once?
- Why might each client need its **own** write loop and outgoing channel?
- What should happen, in order, when a connection breaks?
- How do you stop one slow client from freezing everyone else?

**Hints**
- Two common designs: (a) one goroutine owns shared state and others talk to it via channels, or (b) a mutex guards shared state. Pick one and be ready to explain what the other costs.
- Read about why concurrent writes to one WebSocket connection are unsafe.
- Test: connect 3 tabs, close one abruptly, see if the server notices.

**Done when**
- [ ] 3 tabs connect
- [ ] Messages from one are logged by the server
- [ ] Closing a tab removes its client cleanly (no leaks or panics)
- [ ] Reviewed by a non-author (P4 reviews P1)

**Explain-back:** *"Walk me through what happens from the moment a browser connects until it disconnects."*
- [ ] Answered by P1 · [ ] Others can also answer

---

### Phase 4: Identity and user list (P2 leads) → R3, R4

**Status:** 🟡 code exists (written by P4 + AI, not by P2); review and explain-back pending

**Goal:** Unique names and a live list of connected users.

**Concepts:** Application-level handshake on top of the socket, state transitions (connected but not yet registered), broadcasting.

**Guiding questions**
- Between "socket connected" and "registered with a name", what should the server accept or reject?
- Is `Kasidis` the same as `kasidis`? Is an empty or 500-character name okay?
- What if two users try to register the same name at the same instant?
- When is the user list sent: only on join, or whenever anyone joins or leaves? To whom?
- What happens to a name when its owner disconnects?

**Hints:** think of registration as a request/response pair; a broadcast is "loop over clients and queue a message to each".

**Done when**
- [ ] Names can't collide (R3)
- [ ] Every client sees every name including its own (R4)
- [ ] The list updates live on join and leave
- [ ] Reviewed by a non-author (P1 reviews P2)

**Explain-back:** *"How does the server guarantee uniqueness even with two simultaneous registrations?"*
- [ ] Answered by P2 · [ ] Others can also answer

---

### Phase 5: Direct messages and rooms (P2 leads, P3 reviews) → R5, R7

**Status:** 🟡 code exists (written by P4 + AI, not by P2); review and explain-back pending

**Goal:** Private one-to-one chat that only two people can see.

**Concepts:** Routing (unicast vs broadcast), room identity, server-side authorization.

**Guiding questions**
- Given a DM from A to B, how does the server find B's connection?
- What if B doesn't exist or just disconnected?
- Is the room between A and B the same as between B and A? How do you identify it?
- Can a client send a message that *claims* to be from someone else? Where do you stop it?
- Does the sender get their own message echoed back from the server, or does the client just display it? Trade-offs?

**Hints:** prove privacy with a third tab (C) that must never receive an A↔B message; make this a repeatable test.

**Done when**
- [ ] A↔B messages appear only for A and B (R7)
- [ ] Each DM has its own room in the UI (R5, R6)
- [ ] Errors for missing users are handled gracefully
- [ ] Third-tab privacy test is repeatable
- [ ] Reviewed by a non-author (P3 reviews P2)

**Explain-back:** *"Show me the exact line where a message is sent to only one person."*
- [ ] Answered by P2 · [ ] Others can also answer

---

### Phase 6: Groups (P3 leads, P4 reviews) → R8, R9, R10, R11

**Status:** 🟡 code exists (written by P4 + AI, not by P3, so P4 can no longer be the independent reviewer); review and explain-back pending

**Goal:** Create, list, join, and chat in groups with correct visibility.

**Concepts:** Many-to-many membership, fan-out, group state, consistency between server state and what clients display.

**Guiding questions**
- What does a `Group` store? How do you make "creator starts as the only member" true?
- Must group names be unique? Can they be empty?
- When someone creates or joins, who is notified?
- If a non-member sends a group message, what does the *server* do?
- What happens to membership when a member disconnects? What happens to a group when everyone has left? Decide and document.
- What if a user joins the same group twice?

**Hints:** enforce membership on the server. Test with 3 tabs: A creates, B joins, C does not. C must see nothing. R9 needs member lists inside the group list, so think about when it goes stale.

**Done when**
- [ ] R8: creator is the only initial member
- [ ] R9: group list shows member lists, kept fresh
- [ ] R10: joining is always user-initiated
- [ ] R11: a non-member provably cannot read or send (3-tab test)
- [ ] Reviewed by a non-author (P4 reviews P3)

**Explain-back:** *"How does your server decide who receives a group message?"*
- [ ] Answered by P3 · [ ] Others can also answer

---

### Phase 7: Special features (Part 3), one owner each, all review

> Start only after R3 to R11 pass. All four are worth 1 point each. Unread tracking is shared by F1 and F4: **P1 and P4 agree on it first** (see Phase 2 decision).

**Status:** 🟡 F1 to F4 all built (not by their listed owners); F1 live Claude call untested; review and explain-backs pending

#### F1: AI text summary triggered on unread chats (P4)

**Goal:** When a user has unread messages in a chat, they can get (or are offered) an AI summary of what they missed.

**Concepts:** unread tracking, calling an external API from the server, non-blocking work, prompt design, failure handling, privacy.

**Guiding questions**
- What exactly does "unread" mean: room not open? tab not focused? Who decides, and where is the counter stored?
- What triggers the summary: a threshold (N unread), a button, or on opening the room? What are the cost and annoyance trade-offs?
- If the LLM takes 5 seconds, what must **not** block?
- Which messages are you allowed to send to the LLM? (Think R7/R11: only what *this* user may see.) Should users be told that messages go to a third party?
- How is a sticker represented in the text sent to the model?
- How does the summary reach the user: over which path (R2)? What if the API fails or times out?
- How much history do you send? What if there are 5,000 unread messages?

**Hints:** model the summarizer as a function from "list of visible messages" to "short text". Add a timeout and a friendly error. The key lives only in a server env var.

**Done when**
- [ ] A user with unread messages gets a correct summary
- [ ] Failures degrade gracefully (timeout, API error)
- [ ] No secret is in the repo
- [ ] Only messages the user may see are sent to the LLM

**Explain-back:** *"Trace one summary from the trigger to the screen. Which protocols are involved at each hop?"*
- [ ] Answered by P4

#### F2: QR code group invitation (P3)

**Goal:** A group member shows a QR code; another person scans it and ends up on a join screen for that group.

**Concepts:** URLs as the payload of a QR code, invite tokens vs guessable IDs, join-by-link flow, consent (R10).

**Guiding questions**
- What does the QR encode? A group ID you could guess, or something harder? What can an attacker do with each?
- Scanning must not **auto-add** the user (R10). What does the user do after scanning?
- What if the scanner hasn't registered a name yet? What if the group no longer exists? What if they're already a member?
- Which address goes in the QR on a LAN demo? Can a phone on the same Wi-Fi reach it? What changes when you deploy to the cloud?
- Should invites expire or be revocable? (Optional.)

**Hints:** the QR is just a link, so the "scan" is the phone camera app opening a URL. The join itself still goes through the socket. In-browser camera scanning needs a secure context, so prefer the phone's camera app for the demo.

**Done when**
- [ ] Scanning from a second device opens the app
- [ ] Registration completes if needed
- [ ] A "Join group X?" step is shown (no auto-add)
- [ ] Joining works via the socket
- [ ] Bad cases handled (group gone, already a member, forged token)

**Explain-back:** *"Why isn't scanning the QR the same as being added by the creator (R10)?"*
- [ ] Answered by P3

#### F3: Default stickers, LINE-style (P2)

**Goal:** Users can pick a sticker from a built-in pack and send it in DMs and groups.

**Concepts:** message types, payload size, static assets vs socket messages, input validation.

**Guiding questions**
- Do you send the image bytes over the socket or just an identifier? What are the trade-offs?
- How does the server reject an unknown sticker ID?
- Is a sticker a new message type or a field on a normal message? What does that mean for history and for the AI summary?
- Where are the images served from? Does that violate R2? (Think: what is the *chat message* and what is just an asset?)
- **Licensing:** can you use LINE's real stickers? Make your own, or use CC0 / open-licensed assets, and credit them.

**Hints:** sticker pack = a small list of IDs and image files shipped with the frontend. Picker UI sits next to the chat box in `ChatRoom`.

**Done when**
- [ ] A sticker sent in a DM or group is visible only to the correct people
- [ ] It renders properly
- [ ] Invalid IDs are rejected by the server
- [ ] Sticker assets are open-licensed or self-made, with credit

**Explain-back:** *"What exactly travels over the WebSocket when I send a sticker?"*
- [ ] Answered by P2

#### F4: Notification when messages are received (P1)

**Goal:** Users are alerted to new messages when they aren't looking at that chat.

**Concepts:** browser Notification API, permissions, tab visibility, unread badges, secure contexts.

**Guiding questions**
- Which events notify (DM, group, mention)? Which don't (your own messages, the room you're viewing)?
- When may you ask for notification permission, and what is the fallback if the user says no?
- Which parts are your server's job, and which are purely the browser's?
- **Secure context:** the Notification API generally needs HTTPS or `localhost`. Your 2-computer demo uses `http://<LAN-IP>`. What happens, and what are your options (HTTPS dev cert, tunnel, cloud deploy, or an in-app fallback)?
- Does this work when the tab is closed? (Hint: that needs Web Push + a service worker; decide if it's out of scope.)

**Hints:** plan an in-app fallback (toast, title badge, sound, unread counter on the room list) so the feature never silently fails. Coordinate the unread counter with F1.

**Done when**
- [ ] A message to a user who isn't viewing that room produces a visible alert on the demo setup
- [ ] A working in-app fallback exists
- [ ] Own messages and the open room do not notify

**Explain-back:** *"Which layer of the system detects that a message arrived, and which layer shows the notification?"*
- [ ] Answered by P1

---

### Phase 8: Multi-computer test and R1 (P1 + P4 lead, all participate) → R1

**Status:** ⬜ not started (blocked: needs physical computers; steps in `docs/RUNNING.md`)

**Goal:** Prove it works across real machines and explain the protocols (Section 7).

**Tasks**
- [ ] 1. Run the server on one laptop; connect from two other physical computers.
- [ ] 2. Walk through every protocol in Section 7: what it is, why, and how this project uses it.
- [ ] 3. *(Bonus)* Deploy to the cloud with a secure WebSocket (`wss://`).

**Guiding questions**
- Why does `localhost` work on your machine but fail on your friend's?
- What does it mean to listen on `127.0.0.1` vs all interfaces?
- Could an OS firewall block your port?
- Why would an HTTPS page refuse to open a plain `ws://` connection?
- Which hosting models can keep long-lived connections open?

**Done when**
- [ ] 2 browsers on 2 physical computers chat privately
- [ ] ...and in a group
- [ ] All four special features have been demonstrated

**Explain-back:** *"What would break if we put this behind a service that only supports short request/response calls?"*
- [ ] P1 · [ ] P2 · [ ] P3 · [ ] P4

---

### Phase 9: Hardening and rehearsal (all 4, 1 to 2 sessions)

**Status:** ⬜ not started

- [ ] 1. **Break it on purpose:** kill a client mid-message, send malformed JSON, send a huge message, open 10 tabs.
- [ ] 2. **Trap hunt:** each person tries to trigger one of these on a teammate's module:
  - [ ] unsafe concurrent map access
  - [ ] concurrent writes to one connection
  - [ ] goroutine leaks on disconnect
  - [ ] message to a user who just left
  - [ ] spoofed sender name
  - [ ] joining a group twice
  - [ ] invalid sticker ID
  - [ ] reused or forged invite
- [ ] 3. **Mock Q&A:** take turns acting as the instructor, pointing at random lines: *"What does this do? What if it fails? Why this choice?"*
- [ ] 4. **Slide:** "Implementation Explanation" (architecture diagram, protocol table from Section 7, message flow for one DM, one group message, one summary).
- [ ] 5. **Demo script:** timed run-through covering every requirement in order.

**Done when**
- [ ] Every member can answer questions about **any** part of the server.

---

## 6. Timeline (adjust to your deadline)

| Session | Focus | Done |
|---|---|---|
| 1 | Phase 1: TCP echo and forwarding | ⬜ |
| 2 | Phase 1 wrap-up + Phase 2: WebSocket echo, start protocol | ⬜ |
| 3 | Phase 2: finalize protocol and decisions (frontend agents start after this) | ⬜ |
| 4 | Phase 3: hub and connection lifecycle | ⬜ |
| 5 | Phase 4: names + user list | ⬜ |
| 6 | Phase 5: direct messages | ⬜ |
| 7 | Phase 6: groups | ⬜ |
| 8 | Phase 6 finish + F1 to F4 design on paper | ⬜ |
| 9 | Phase 7: special features (parallel, one each) | ⬜ |
| 10 | Phase 7 finish + integration + cross-review | ⬜ |
| 11 | Phase 8: multi-computer test, R1 write-up | ⬜ |
| 12 | Phase 9: hardening + mock Q&A | ⬜ |
| 13 | Final rehearsal and buffer | ⬜ |

Parallelism: while P1 finishes the hub, P2 and P3 design data structures on paper, and P4 experiments with the LLM API in a standalone script. Frontend agents can build their slices against a mock socket once `PROTOCOL.md` is final.

---

## 7. Protocols we use, and why

### 7.1 Protocol stack for this project

| Layer | Protocol | Why we use it | How it appears in our project |
|---|---|---|---|
| Application (our own) | **Chat protocol: JSON messages over WebSocket** (`PROTOCOL.md`) | The rubric requires socket-based chat; we need an agreed message format so server and clients interpret events the same way | `register`, `dm`, `group_msg`, sticker, invite, summary events, etc. |
| Application | **WebSocket (RFC 6455)** | Persistent, full-duplex channel; the server can push messages without polling; supported natively by browsers (unlike raw TCP) | All chat traffic between the React client and the Go server (R2) |
| Application | **HTTP/1.1** | WebSocket starts as an HTTP request with `Upgrade: websocket`; also serves the React app and sticker images; the QR code encodes an HTTP(S) link | Handshake (`101 Switching Protocols`), static assets, invite links. **Not** used for chat messages |
| Application | **HTTPS (HTTP over TLS)** | Encrypts and authenticates calls from our server to the LLM provider's API; protects the API key in transit | Server → LLM API for the AI summary / bot |
| Security | **TLS** (`wss://`, `https://`) | Confidentiality and integrity; needed by browsers for secure-context features; recommended for cloud deployment | Cloud bonus, and a way to make notifications work on a non-localhost demo |
| Application | **DNS** | Resolves hostnames (LLM API host, cloud domain) to IP addresses; unnecessary on a LAN if you use the IP directly | Server resolving the LLM API; clients resolving the cloud domain |
| Transport | **TCP** | Reliable, ordered, connection-oriented delivery; chat can't tolerate lost or reordered messages; WebSocket and HTTP both run on it | Underlies every connection above; the Phase 1 raw-socket exercise |
| Network | **IP (IPv4)** | Addressing and routing between the computers (LAN IPs or a cloud IP) | Clients connect to `<server-IP>:<port>` |
| Network / LAN support | **ARP, DHCP, ICMP** | ARP maps IP to MAC on the LAN; DHCP assigns each laptop's IP; ICMP (`ping`) is how we verify connectivity | Used implicitly on the Wi-Fi/LAN during the multi-computer test |
| Link | **Wi-Fi (802.11) / Ethernet** | Physical and link-layer delivery between the machines | The demo network |

**Note:** DNS usually uses UDP for lookups, so UDP is used *indirectly*, but not for chat.

### 7.2 Things that sound like protocols but aren't

| Item | What it really is |
|---|---|
| JSON | A data format that our chat protocol uses to encode messages |
| Notification API | A browser API; our server isn't involved in displaying notifications |
| QR code | An encoding of a URL (we use it to carry an HTTP(S) link) |
| LLM provider API | An HTTPS API we call from the server |

### 7.3 Alternatives we rejected (be ready to defend these)

| Alternative | Why not |
|---|---|
| Raw TCP sockets from the browser | Browsers can't open raw TCP sockets; we use it only for the Phase 1 learning exercise |
| UDP for chat | No delivery or ordering guarantees; we'd have to rebuild reliability ourselves |
| Socket.IO | Hides the socket logic that the course wants us to understand and explain |
| REST polling / long polling | Wasteful and higher latency; violates the spirit of R2 |
| Server-Sent Events | One-way (server → client) only; chat needs both directions |
| MQTT / gRPC | Extra infrastructure and concepts beyond the project's scope |

---

## 8. Definition of done (per feature)

- [ ] Works with 3 browser tabs, including a "should not see it" tab
- [ ] Works on 2 physical computers (at least once per phase from Phase 4)
- [ ] Reviewed by a non-author who asked questions
- [ ] `PROTOCOL.md` updated if the message format changed
- [ ] Learning-log entry written
- [ ] Author can explain it without looking at the code

---

## 9. Study resources

- Go: `net` package docs, your chosen WebSocket library's docs, "Go by Example: goroutines, channels, mutexes"
- WebSockets: MDN "WebSocket API" and "Writing WebSocket servers"
- Browser: MDN "Notifications API" and "Secure Contexts"
- Course materials in the project brief (Python sockets, Go WebSocket tutorial)
- Browser DevTools → Network → WS tab (your best debugging tool)
