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
- Group names are not unique; groups are identified by ID (g_N).
- Invites are reusable and don’t expire, but die with the group.
- The token is a link convenience and consistency check, not access control, because R9 and R10 make groups open.
- There is no auto-rejoin after a reconnect. The alternative considered was client-side rejoin; it was rejected because it touches P1’s store.tsx and only helps with connection drops, not refreshes.

## D8. LLM provider and model
- **Status:** PROPOSED (implemented in `server/summary.go` and `server/gemini.go`)
- Two interchangeable providers behind the one `Summarizer` function type, chosen at startup by `NewSummarizerFromEnv` (`LLM_PROVIDER` forces one; otherwise whichever key exists, Anthropic first).
  - **Gemini** (added because it has a free tier and the team did not want to pay for an Anthropic API key): plain `net/http` POST to `generateContent` over HTTPS, key in the `x-goog-api-key` header, key from `GEMINI_API_KEY`, model from `GEMINI_MODEL` (default `gemini-3.5-flash-lite`). Google calls `generateContent` "legacy but supported" and now recommends its Interactions API; we kept `generateContent` because its response format is stable and well known. Revisit if it is removed.
  - **Claude** via the official Anthropic Go SDK. Model defaults to `claude-opus-5-5`, overridable with `LLM_MODEL` (cheaper: `claude-haiku-5-5`). Key from `ANTHROPIC_API_KEY`. This needs a paid API key; a claude.ai subscription does not work for it.
- Keys are read from the server environment only. Without a key the app still works and summaries answer with a friendly error. Provider errors are logged on the server and never sent to the browser.
- **Privacy:** only the requesting user's own unread messages of one room are sent, and the UI tells the user that a third party processes them.
- **Not verified:** the live calls have never been run against the real services. The Gemini request/response handling is tested only against a fake HTTP server written from Google's documentation (`server/gemini_test.go`), and the Claude path only through a fake summarizer. Whether the default Gemini model is free for a given key must be checked in AI Studio.

## D9. Sticker artwork
- **Status:** PROPOSED (implemented)
- Eight original SVG faces made for this project (no LINE or third-party assets), in `client-web/public/stickers/`. Ids live in `server/stickers.go`; a Go test checks every id has a file.

## Team acknowledgement

| Decision | P1 | P2 | P3 | P4 |
|---|---|---|---|---|
| D1 to D9 (read all of them) | ⬜ | ⬜ | ⬜ | ⬜ |
