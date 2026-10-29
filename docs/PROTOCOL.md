# PROTOCOL.md

Single source of truth for message types and fields. Keeper: P2.

> **STATUS: DRAFT v0.1, proposed by P4 (with an AI mentor). NOT yet acknowledged by the team.**
> v0.1 matches what is actually implemented in `server/` (Go) and `client-web/src/socket/` (TypeScript);
> `server/hub_test.go` tests it. Because the team has not agreed yet, treat every rule here as a proposal:
> P1 to P3 must read it, change what they disagree with (and then the code), and tick the
> acknowledgement list at the bottom.

## 1. Transport and envelope

- One WebSocket per client at `ws://<host>:8080/ws` (`wss://` when deployed behind TLS).
- Every frame is one **text** message containing one JSON object (the envelope):

```json
{ "type": "dm", "payload": { } }
```

| Field | Required | Meaning |
|---|---|---|
| `type` | always | Event name (tables below) |
| `payload` | always | Object with the event's fields (`{}` if none) |
| `reqId` | optional, client → server | Client-chosen string. If present, the server echoes it in the direct reply (`*_ok`, `*_error`, `error`) so the UI can match answers to requests |

Rules
- **The server stamps `from` / `sender`.** Clients never send their own name in `dm` / `group_msg`; the server uses the name bound to that connection. (Prevents spoofing, R7.)
- A connection can send only `register` until it has registered. Anything else gets `error` with `code: "not_registered"`.
- Names are compared **case-insensitively** and trimmed; allowed length 1 to 20 chars. The server stores the original spelling for display. A name is freed when its owner disconnects.
- Unknown `type` or malformed JSON: reply `error` (`code: "bad_request"`), keep the connection open.
- Max frame size 8 KB. A larger frame makes the server close the connection (the WebSocket library enforces the read limit and gives no chance to reply first). Non-text frames get `error` (`bad_request`).
- After a dropped connection the client reconnects and re-sends `register` with the same name. The server forgot the old connection, so the name is free again unless someone else took it (then `register_error: taken`).
- Liveness: the server pings every 25 s and closes a connection that has not answered within 60 s (browsers answer pings automatically).
- A client whose outgoing buffer is full (64 frames, a slow or stuck reader) is disconnected instead of making everyone else wait.

## 2. Room identity (how the UI knows where a message belongs)

Every chat message carries a `room` string, built by the **server**:

| Room | Format | Example |
|---|---|---|
| DM | `dm:<nameA>\|<nameB>` with both names lowercased and sorted alphabetically (A↔B and B↔A are the same room) | `dm:alice\|bob` |
| Group | `group:<groupId>` (groupId is a server-generated id, not the name) | `group:g_17` |

## 3. Messages

### 3.1 Identity (R3, R4)

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| client → server | `register` | `name: string` | server |
| server → client | `register_ok` | `name: string` (canonical display name) | only the registering client |
| server → client | `register_error` | `code: "taken" \| "invalid"`, `message: string` | only the registering client |
| server → clients | `user_list` | `users: string[]` (every registered name, including the receiver's own) | **all** registered clients, sent after every join and leave; also sent to a new client right after `register_ok` |

### 3.2 Direct messages (R5, R7)

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| client → server | `dm` | `to: string`, `kind: "text" \| "sticker"`, `text?: string`, `stickerId?: string` | server |
| server → client | `dm` (delivered) | `room`, `from`, `to`, `kind`, `text?`, `stickerId?`, `ts: number` (unix ms, server clock), `msgId: string` | **only** the sender and the receiver (the sender receives it as confirmation, so both sides render the server's version) |
| server → client | `error` | `code: "no_such_user"` etc. | only the sender |

### 3.3 Groups (R8 to R11)

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| client → server | `create_group` | `name: string` (1 to 30 chars) | server |
| server → client | `create_group_ok` | `groupId: string` | only the creator |
| server → clients | `group_list` | `groups: [{ groupId, name, members: string[] }]` (all groups, with member lists, R9) | **all** registered clients, after any create / join / leave / disconnect; also to a new client after `register_ok` |
| client → server | `join_group` | `groupId: string`, `invite?: string` (invite token, see 3.5) | server |
| server → client | `join_group_ok` | `groupId` | only the joiner |
| client → server | `group_msg` | `groupId: string`, `kind: "text" \| "sticker"`, `text?`, `stickerId?` | server |
| server → client | `group_msg` (delivered) | `room`, `groupId`, `from`, `kind`, `text?`, `stickerId?`, `ts`, `msgId` | **only current members** of the group (including the sender) |
| server → client | `error` | `code: "not_member" \| "no_such_group" \| "already_member"` | only the requester |

Decisions baked into this draft (change them if the team disagrees)
- The creator is the **only** member at creation (R8). Nobody is ever added by someone else (R10): `join_group` is always sent by the person joining.
- Group names need not be unique; `groupId` is the identity. Duplicate `join_group` returns `already_member` (harmless, no second entry).
- A disconnecting user is **removed from all groups** (they must join again). A group with zero members is deleted.
- A `group_msg` from a non-member is rejected with `not_member`; nobody else sees it (R11).

### 3.4 Stickers (Special: stickers)

- A sticker is a normal `dm` / `group_msg` with `kind: "sticker"` and `stickerId` instead of `text`. **Only the id travels over the socket**, never image bytes.
- The server validates `stickerId` against its built-in pack list; an unknown id gets `error` with `code: "bad_sticker"`.
- Images are static files shipped with the frontend (HTTP static assets, not chat messages).
- For the AI summary a sticker is rendered as the text `[sticker: <stickerId>]`.

### 3.5 Invites / QR (Special: QR)

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| client → server | `create_invite` | `groupId` | server (requester must be a member, else `not_member`) |
| server → client | `invite_created` | `groupId`, `token: string` (random, unguessable) | only the requester |
| client → server | `invite_info` | `token` | server |
| server → client | `invite_info_ok` | `groupId`, `groupName`, `memberCount`, `alreadyMember: boolean` | only the requester (used to show "Join group X?") |
| client → server | `join_group` | `groupId`, `invite: token` | server (token must match the group) |

- The QR encodes `http(s)://<host>/?invite=<token>`. Opening it does **not** join anything: the page registers a name if needed, sends `invite_info`, shows "Join group X?", and only sends `join_group` when the user clicks Join (R10).
- Invalid or unknown tokens return `error` with `code: "bad_invite"`.
- Plain `join_group` without a token is still allowed (join from the group list). Whether invites should also expire is optional and undecided.

### 3.6 Unread state (shared by AI summary and notifications)

**Proposed owner: the server**, so it survives refresh and reconnect.

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| client → server | `open_room` | `room: string` (the room the user is now looking at) | server |
| client → server | `close_room` | `room` (tab hidden or room switched away) | server |
| server → client | `unread` | `counts: { [room: string]: number }` | only that user; sent when a count changes and right after `register_ok` |

- The server counts a delivered message as unread for a recipient if that recipient does not currently have the room open. `open_room` resets that room's count to 0.
- The client uses `unread` to draw badges and decide when to show a notification (the Notification API itself is purely browser-side).
- A user's unread state is dropped when they disconnect (proposal; the simplest rule, and nobody can inherit another person's unread by taking their freed name). After a reconnect the counts start at 0. There is no message history, so a refresh loses the local message window too.
- For summaries the server keeps the newest 200 unread messages per room per user; the count itself keeps growing past 200.
- A message never counts as unread for its own sender, nor for a recipient who has that room open (`open_room`) in a visible tab. The client sends `close_room` when the tab is hidden or the user switches away.
- `open_room` / `close_room` for a room the user may not access return `error` (`not_member` / `no_such_group`).

### 3.7 AI summary (Special: AI summary)

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| client → server | `summarize_request` | `room: string` | server |
| server → client | `summary_result` | `room`, `text: string`, `messageCount: number` | only the requester |
| server → client | `summary_error` | `room`, `code: "timeout" \| "llm_error" \| "nothing_unread" \| "not_allowed"`, `message` | only the requester |

- The server builds the input **only from messages the requester is allowed to see** (member of the group, party to the DM) and only the unread ones, capped to the most recent N (N to be decided; a first guess is 200).
- The LLM call runs in its own goroutine with a 25 s timeout; it must not block the read loop or the hub. Only one summary per user can be in flight: a second request gets `error` with `code: "rate_limited"`.
- If `ANTHROPIC_API_KEY` is not set on the server, the answer is `summary_error` with `code: "llm_error"` (the rest of the app works without AI).
- Sticker messages appear in the LLM input as `[sticker: <id>]`; message text is wrapped as untrusted data in the prompt.
- The API key is a server environment variable only.

### 3.8 Generic error

| Direction | `type` | `payload` fields | Who receives it |
|---|---|---|---|
| server → client | `error` | `code: string`, `message: string`, `reqId?` | only the client that caused it |

Known codes: `bad_request`, `not_registered`, `no_such_user`, `no_such_group`, `not_member`, `already_member`, `bad_sticker`, `bad_invite`, `rate_limited`.

## 4. Message flow examples

DM from alice to bob (carol must see nothing):
```
alice → server : {"type":"dm","payload":{"to":"bob","kind":"text","text":"hi"}}
server → alice : {"type":"dm","payload":{"room":"dm:alice|bob","from":"alice","to":"bob","kind":"text","text":"hi","ts":...,"msgId":"m_1"}}
server → bob   : (the same delivered message)
server → carol : nothing
```

## 5. Open questions for the team

- [ ] Message history: is history stored on the server and replayed on join, or does a client only see messages received while online? (This draft assumes **no history**; the AI summary then covers only messages the server still holds as unread.)
- [ ] Unread counts are dropped on disconnect in v0.1. Is that acceptable, or should they survive a short reconnect?
- [ ] Mentions (`@name`) as a notification trigger: in or out of scope?
- [ ] Should the AI also be a chat bot user, or summary only?
- [ ] `rate_limited`: do we implement it?

## 6. Team acknowledgement (required before coding against this)

- [ ] P1 acknowledged
- [ ] P2 acknowledged (keeper)
- [ ] P3 acknowledged
- [ ] P4 acknowledged
