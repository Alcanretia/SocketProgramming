# DEMO_SCRIPT (draft, ~10 minutes)

Setup: server on laptop A (run from a built `chat.exe`, LAN address noted, firewall allowed). Laptop B and C (or a phone) on the same Wi-Fi. `ANTHROPIC_API_KEY` set on A. Three browser windows visible: Alice (A), Bob (B), Carol (C). Open the app on A via the LAN address, not `localhost`.

| # | Show | Say / point at | Rubric |
|---|---|---|---|
| 1 | Three machines, three different physical computers | "Server on A, clients on B and C over the LAN." Show the printed addresses. | R1 |
| 2 | DevTools → Network → WS on one client | `Upgrade: websocket`, `101 Switching Protocols`, frames in both directions. "Chat travels only here; HTTP just served the files." | R1, R2 |
| 3 | Register Alice, Bob, Carol; try "alice" again | Unique, case-insensitive. Everyone sees all names, including their own. | R3, R4 |
| 4 | Alice → Bob DM; Carol's screen stays empty | Separate room per DM, chat box + message window. "One `deliver` to sender, one to receiver." | R5, R6, R7 |
| 5 | Alice creates "Study" | Only Alice is a member. Bob and Carol see the group and its member list. | R8, R9 |
| 6 | Bob joins himself; Carol does not | Nobody is auto-added. Alice posts: Bob sees it, Carol does not. | R10, R11 |
| 7 | Sticker in a DM and in the group | "Only `stickerId` crosses the socket." Show the frame in DevTools. | F3 |
| 8 | Alice shows the QR; phone scans it | Phone opens the app, picks a name, sees "Join group Study?", presses Join. Scanning alone did not add anyone. | F2, R10 |
| 9 | Carol is on another room; Alice messages her | Toast + sound + tab badge; unread count in the list. Mention the secure-context limit and the fallback. | F4 |
| 10 | Carol presses "Summarize" on the unread chat | Summary appears. Trace: browser → WebSocket → server → HTTPS → Claude → back. Mention third-party disclosure. | F1 |
| 11 | Stop the API key / unplug network, ask again | Friendly error, chat still works. | F1 |
| 12 | Close Bob's laptop lid / kill a tab | User list updates, group member list updates, no crash. | robustness |

Backup plan if Wi-Fi fails: phone hotspot as the LAN. If the API is down: show step 11 and the test with the fake summarizer (`go test -run TestUnreadAndSummary -v`).

Before the demo, run `go test ./...` and `npm run build` once.
