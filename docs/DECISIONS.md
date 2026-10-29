# DECISIONS.md

Each decision gets: what we chose, the alternatives, and **why**.
Status: `PROPOSED` (made by P4 with an AI helper, not yet agreed by the team) or `DECIDED`.
Nothing here is `DECIDED` until P1 to P3 have read it and said so in the table at the bottom.

## D1. Go WebSocket library: `github.com/gorilla/websocket` v1.5.3
- **Status:** PROPOSED (implemented)
- **Alternatives:** `github.com/coder/websocket` (formerly nhooyr), `golang.org/x/net/websocket` (legacy).
- **Why:** the most widely documented Go WebSocket library, so tutorials and course material most likely match. It exposes the upgrade step explicitly (`Upgrader.Upgrade`), which makes the HTTP → WebSocket handshake easy to explain. `coder/websocket` is a fine alternative if the team prefers a context-based API.
- **Caveat:** a gorilla `*Conn` allows only one concurrent writer and one concurrent reader, which shaped D4.

## D2. Frontend language: TypeScript
- **Status:** PROPOSED (implemented)
- **Alternative:** plain JavaScript.
- **Why:** `client-web/src/socket/types.ts` mirrors `PROTOCOL.md`, so a field typo is a compile error instead of a runtime mystery. Cost: a few more concepts to read.

## D3. Build tool: Vite (React 19)
- **Status:** PROPOSED (implemented)
- **Why:** fast dev server, proxies `/ws` to the Go server. In production the Go server serves the built `client-web/dist`, so there is one origin and one port.

## D4. Hub concurrency: one owner goroutine + channels (no mutex)
- **Status:** PROPOSED (implemented in `server/hub.go`)
- **Alternative:** a `sync.Mutex` around shared maps (this is what `server/phase1/hub` does).
- **Why:** all shared state (names, groups, invites, unread) is touched only by the goroutine running `Hub.Run`. Other goroutines send it events on a channel. "Check the name is free, then take it" is therefore atomic without any lock, which is what makes R3 safe under simultaneous registrations (tested). Per client there are exactly two goroutines (`readPump`, `writePump`), because gorilla allows one writer.
- **Cost:** every handler must be short and must never block, otherwise it stalls everyone. So sending is `enqueue` into a 64-frame buffer, never a blocking write, and a full buffer disconnects that client. The LLM call runs in its own goroutine and reports back as an event.
- **Rejected:** the mutex version, because holding a lock while writing to many sockets lets one slow client freeze all the others (the weakness left in `phase1/hub` on purpose).

## D5. Unread tracking: the server owns it
- **Status:** PROPOSED (implemented). **P1 and P4 were supposed to agree on this first.**
- **Alternative:** the client counts unread locally.
- **Why:** the AI summary needs the actual unread messages, and the server is the only place that knows what each user was allowed to receive (R7/R11). A client-side counter would also reset on refresh. Counts are dropped on disconnect (see PROTOCOL.md 3.6).

## D6. One HTTP server on `:8080` for static files and `/ws`
- **Status:** PROPOSED (implemented)
- **Why:** a phone that scans a QR code reaches the same origin as the app. HTTP carries only static files (the React build, sticker images) and the WebSocket upgrade; chat never uses it (R2). Listens on all interfaces so other computers can connect. The default `Origin` check of the upgrader is kept (blocks cross-site WebSocket hijacking); `-dev` turns it off for the Vite dev server only.
- The Phase 1 raw-TCP exercises use `:9000` and are throwaway.

## D7. Group lifecycle, sticker representation, invite tokens
- **Status:** PROPOSED in `PROTOCOL.md` (3.3, 3.4, 3.5); implemented.
- Disconnect removes a user from all groups; an empty group is deleted with its invites. A sticker is `kind: "sticker"` + `stickerId` (id only, images are static files). Invite = 128-bit random token, link `/?invite=<token>`, scanning never joins by itself.

## D8. LLM provider and model
- **Status:** PROPOSED (implemented in `server/summary.go`)
- Claude via the official Anthropic Go SDK, called from the server over HTTPS. Model defaults to `claude-opus-5-5`, overridable with `LLM_MODEL` (for a cheaper demo, e.g. `claude-haiku-5-5`). The key is read from `ANTHROPIC_API_KEY` in the server environment only. Without a key the app still works and summaries answer with a friendly error.
- **Privacy:** only the requesting user's own unread messages of one room are sent, and the UI tells the user that a third party processes them.
- **Not verified:** no API key was available when this was written, so the live call has never been run; only the plumbing is tested with a fake summarizer.

## D9. Sticker artwork
- **Status:** PROPOSED (implemented)
- Eight original SVG faces made for this project (no LINE or third-party assets), in `client-web/public/stickers/`. Ids live in `server/stickers.go`; a Go test checks every id has a file.

## Team acknowledgement

| Decision | P1 | P2 | P3 | P4 |
|---|---|---|---|---|
| D1 to D9 (read all of them) | ⬜ | ⬜ | ⬜ | ⬜ |
