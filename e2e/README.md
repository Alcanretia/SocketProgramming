# e2e: multi-client browser test

Drives real browser windows (Alice, Bob, Carol, Dave) against a running server and checks R3 to R11 and the special features end to end, including "the third tab must see nothing".

```powershell
# terminal 1: build the frontend and start the server (see docs/RUNNING.md)
# terminal 2:
cd e2e
npm install
$env:BROWSER_PATH = "C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe"   # or Chrome
npm test
```

`BASE_URL` (default `http://localhost:8080`) can point at another machine. Screenshots go to `e2e/shots/`.
Run it against a server **without** `ANTHROPIC_API_KEY` (one check expects the friendly "AI not available" error).
