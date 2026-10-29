# Running the project

## Requirements
- Go 1.26 (`server/go.mod` pins `go 1.26.5`; an older Go will try to download that toolchain) and Node 20+ with npm.

## Build and run (one server, one port)

```powershell
# 1. build the React app once (and again after frontend changes)
cd client-web
npm install
npm run build

# 2. run the Go server (serves the app and /ws on :8080)
cd ../server
go build -o chat.exe .
./chat.exe
```

The server prints the addresses to open, for example `http://localhost:8080` and `http://192.168.1.20:8080`.

Flags: `-addr :8080`, `-static ../client-web/dist`, `-dev` (see below).

## AI summary (special feature F1)

```powershell
$env:ANTHROPIC_API_KEY = "sk-ant-..."     # never commit this; .env is git-ignored
$env:LLM_MODEL = "claude-haiku-5-5"       # optional, default claude-opus-5-5
./chat.exe
```

Without a key everything else works and the summary button answers "not available".

## Frontend development with hot reload

```powershell
cd server;     go run . -dev          # -dev accepts the Vite origin
cd client-web; npm run dev            # http://localhost:5173
```

## Tests

```powershell
cd server; go test ./...              # hub behaviour with real WebSocket clients (R3 to R11, stickers, invites, summary)
cd client-web; npm run build          # type-checks the frontend
```

`go test -race` needs a 64-bit C compiler for cgo; it did not work on the first machine, so run it on a machine where it does.

## Two or more computers (R1) and phones (QR)

1. Put every machine on the same Wi-Fi/LAN.
2. On the server machine run the server and note the LAN address it prints (not `localhost`).
3. Allow the program through Windows Firewall when asked (private networks). If no prompt appeared: *Windows Security → Firewall → Allow an app*, or allow inbound TCP 8080.
   Running `go run .` creates a new temporary exe each time, so Windows may ask again. Prefer `go build` and running the same `chat.exe`.
4. On the other computer open `http://<server-LAN-IP>:8080`. Check with `ping <ip>` first if it does not load.
5. For the QR code, open the app on the server machine **via the LAN address** (not `localhost`) before pressing "Invite (QR)", otherwise the QR contains a link that a phone cannot reach (the UI warns about this).
6. Desktop notifications need HTTPS or `localhost`; over plain `http://<LAN-IP>` the app falls back to in-app toasts, sound and the tab-title badge. See ARCHITECTURE.md.

## Repo map

| Path | What |
|---|---|
| `server/*.go` | the chat server (`hub.go` is the core) |
| `server/phase1`, `server/phase2` | learning exercises: raw TCP echo/hub, WebSocket echo |
| `client-web/src/socket` | WebSocket connection, store, notifications (P1) |
| `client-web/src/components/shared` | `ChatRoom` (P1) |
| `client-web/src/features/users-dm` | names, user list, stickers (P2) |
| `client-web/src/features/groups` | groups and QR invite (P3) |
| `client-web/src/features/ai` | unread list and AI summary (P4) |
| `docs/` | PROTOCOL, DECISIONS, ARCHITECTURE, LEARNING_LOG, DEMO_SCRIPT |
