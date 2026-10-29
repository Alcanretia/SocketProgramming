# P1 review and learning notes: connection core

Branch `review/p1-connection-core`. Reviewer: P4. Fill this in with your own words; commit it to this branch.

> The code in your area was written by P4 with an AI tool before you started. Read it, change what you disagree with, and answer below. It is not done until you can explain it.

## Read in this order
1. `server/phase1/` is optional; do your own Phase 1 and Phase 2 first (TCP echo, 2 clients, forwarding, WebSocket echo + DevTools handshake).
2. `server/client.go`: `Client`, `enqueue`, `kill`, `readPump`, `writePump`.
3. `server/hub.go`: `Hub`, `Run`, `post`, `ServeWS`, `send`, `broadcast`, `onLeave`.
4. `client-web/src/socket/connection.ts`, `store.tsx`, `notifications.ts`; `components/shared/ChatRoom.tsx`.
5. `docs/ARCHITECTURE.md` sections 2 and 3.

## Answer in your own words
1. Walk through what happens from the moment a browser connects until it disconnects.

   Answer:

2. Why does each client need its own write loop and outgoing channel? What happens to a slow client?

   Answer:

3. Why is `Client.send` never closed? What would go wrong if it were?

   Answer:

4. Which parts of notifications are the server's job and which are the browser's? What happens on plain `http://<LAN-IP>`?

   Answer:

## Things I found or want to change
-

## Things I tried (write what you saw)
- [ ] My own Phase 1 and Phase 2 programs
- [ ] 10 tabs, a tab killed mid-message, a huge message, malformed JSON
- [ ] OS notification on `localhost` with the tab hidden

## Decision D4 (hub design: owner goroutine vs mutex)
Confirm or change, and why:
