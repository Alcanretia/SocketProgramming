# R1 slide text

Two slides for the "Implementation Explanation" deck. Source facts: `docs/ARCHITECTURE.md`, README Section 7, `docs/PROTOCOL.md`.

## Slide 1: R1 – System Architecture Design

**Picture:** insert `docs/slides/r1-architecture.png` (2760x840) into the "System Architecture Diagram" box. `r1-architecture.svg` is the editable source.

**Explanation (bullets for the slide):**
- Three physical computers: Computer A runs the Go server; Computer B and Computer C each run a browser client (React app).
- Clients and server exchange chat messages only over a WebSocket (a TCP socket); every message is a JSON text frame (R2).
- HTTP is used only to download the web app and the sticker images, and to start the WebSocket (HTTP Upgrade).
- Inside the server, every connected client has two goroutines: `readPump` reads frames from the socket and passes them to the hub, and `writePump` is the only writer of that socket and sends what the hub queued.
- The hub goroutine owns all shared state and is the authority: it decides who receives each private or group message.
- Only the AI summary leaves the local network: the server calls the Google Gemini API over HTTPS.

**Say it out loud (ภาษาไทย):**
- เครื่อง A รันเซิร์ฟเวอร์ Go เครื่อง B กับ C เปิดเว็บแชตผ่านเบราว์เซอร์ ทั้งสามเครื่องอยู่วง LAN เดียวกัน และ client เปิด `http://<IP เครื่อง A>:8080`
- ข้อความแชตทุกอย่างวิ่งผ่าน WebSocket อย่างเดียว ส่วน HTTP ใช้โหลดหน้าเว็บกับรูปสติกเกอร์ และใช้เริ่มต้นการเชื่อมต่อ WebSocket
- ในเซิร์ฟเวอร์ client แต่ละคนมี 2 ลูป: `readPump` อ่านเฟรมจากซ็อกเก็ตแล้วส่งให้ Hub และ `writePump` เป็นตัวเดียวที่เขียนลงซ็อกเก็ต (ส่งสิ่งที่ Hub ใส่คิวไว้ รวมถึง ping) ส่วน Hub เป็น goroutine เดียวที่ถือข้อมูลทั้งหมด จึงไม่ต้องใช้ mutex
- สิทธิ์ว่าใครเห็นข้อความไหน เซิร์ฟเวอร์เป็นคนตัดสิน ไม่ใช่หน้าเว็บ
- มีแค่ฟีเจอร์สรุปด้วย AI ที่เซิร์ฟเวอร์ออกไปคุยกับผู้ให้บริการ AI บนอินเทอร์เน็ตผ่าน HTTPS

## Slide 2: R1 – Explain the computer network protocols used in the project

Suggested layout: one table with three columns that answer the slide's three bullets. "What" is the first column.

| What (protocol) | Why selected | How used in this project |
|---|---|---|
| **Our chat protocol** (JSON messages) | Server and clients must agree on one message format so both read every event the same way | Message types such as `register`, `dm`, `create_group`, `join_group`, `group_msg`, `summarize_request` (`docs/PROTOCOL.md`). The server writes the sender's name itself, so it cannot be faked |
| **WebSocket** (RFC 6455) | A connection that stays open in both directions, so the server can push a message without polling. Browsers support it; they cannot open raw TCP | One connection per browser tab carries all chat traffic (R2) |
| **HTTP/1.1** | The WebSocket handshake is an HTTP request with `Upgrade: websocket`; the web app also has to be downloaded from somewhere | Serves the React app and sticker images. The handshake is answered with `101 Switching Protocols`. No chat message travels over HTTP |
| **TCP** | Reliable, ordered delivery: a chat must not lose or reorder messages. WebSocket and HTTP both run on it | Underneath every connection, server port 8080 |
| **IPv4** | Addressing and routing between the computers | Clients connect to `<Computer A IP>:8080` |
| **HTTPS** (HTTP over TLS) | Encrypts and authenticates the call to the AI provider and protects the API key in transit | Server to AI provider for the AI summary only |
| **DNS** | Turns the AI provider's host name into an IP address | Used by the server to reach the AI provider. Not needed between our computers because we type the IP |
| **DHCP, ARP, ICMP** (LAN support) | They make the local network work: DHCP gives each laptop an IP, ARP maps an IP to a MAC address, ICMP (`ping`) checks connectivity | Used implicitly on the Wi-Fi/LAN; we use `ping` to check that Computer A is reachable |
| **Wi-Fi / Ethernet** (link layer) | Physical delivery between the computers | The network the demo runs on |

**If the instructor asks "why not X?"**
- **UDP:** no delivery or ordering guarantee; we would have to rebuild reliability ourselves.
- **Socket.IO:** hides the socket logic the course wants us to understand.
- **HTTP polling / Server-Sent Events:** polling wastes requests and adds delay; SSE only goes server to client, chat needs both ways.
- **Raw TCP from the browser:** browsers do not allow it. We used raw TCP only for the Phase 1 exercise in `server/phase1/`.

**Say it out loud (ภาษาไทย):** ไล่จากล่างขึ้นบน: Wi-Fi/LAN → IP → TCP → HTTP (ใช้โหลดเว็บและเริ่ม handshake) → WebSocket (ใช้ส่งข้อความแชต) → รูปแบบข้อความ JSON ของเราเอง และแยกเล่าเส้นทางไป AI ต่างหากว่าใช้ DNS + HTTPS
