# P2 Review & Explain-Back Notes: Users, DMs, Stickers & `PROTOCOL.md`

Owner: **P2**
Reviewer: **P1** (Phase 4) / **P3** (Phase 5)
Branch: `teach/p2-coaching`

---

## 1. Files Owned & Implemented by P2

- **Phase 1 Fundamentals:**
  - `server/phase1/echo/main.go` & `main_test.go` — Raw TCP echo server supporting multiple concurrent clients via goroutines (`go handleEchoClient(conn)`).
  - `server/phase1/hub/main.go` & `main_test.go` — Raw TCP mini chat hub forwarding messages to all other clients using `sync.Mutex`.
  - `server/phase1/client/main.go` — Terminal TCP client.
- **Phase 2 Protocol Contract:**
  - `docs/PROTOCOL.md` (Keeper) & `docs/DECISIONS.md` — Answered Section 5 open questions and recorded P2 acknowledgement.
- **Phase 4, 5 & 7 (F3) Backend:**
  - `server/users_dm.go` — `validName`, `normalizeKey`, `broadcastUsers`, `onRegister`, `dmRoom`, `onDM`.
  - `server/stickers.go` — `stickerPack`, `isSticker`, `validContent`.
  - `server/users_dm_test.go` — Unit tests for tricky names (Thai/Emoji, whitespace-only, 500 chars, control chars), DM room symmetry, and sticker validation.
- **Frontend (`client-web/src/features/users-dm/`):**
  - `NameEntry.tsx` — Name registration screen and error display.
  - `UserList.tsx` — Live online user list (including `(you)` with self-DM disabled) and DM room switcher.
  - `stickers.tsx` — 8-sticker picker (`StickerPicker`) and safe sticker renderer (`StickerView`).

---

## 2. Explain-Back Answers (in P2's Own Words)

### Phase 1: *"Why does a chat server need one goroutine per connection (or something equivalent)?"*
- คำสั่ง `Accept()` และ `Read()` เป็นการทำงานแบบ **Blocking (หยุดรอ)**
- ถ้ามีแค่ลูปเดียวโดยไม่แยก Goroutine (`go handleEchoClient(conn)`) เมื่อลูกค้าคนที่ 1 เชื่อมต่อเข้ามาแล้วนั่งเงียบไม่ส่งข้อมูล โปรแกรมจะค้างรออยู่ที่ `Read()` ของคนที่ 1 ทำให้กลับไปรับสาย (`Accept()`) ลูกค้าคนที่ 2 ไม่ได้เลย การแยก Goroutine เหมือนการมีพนักงานประจำแต่ละโต๊ะแยกจากพนักงานต้อนรับหน้าร้าน

### Phase 2: *"What is the difference between an HTTP request and a WebSocket connection? Where does one turn into the other?"*
- **HTTP** เป็นการสื่อสารแบบถาม-ตอบครั้งเดียวจบ (Server พูดหา Client เองไม่ได้ถ้า Client ไม่ถามก่อน) แต่ระบบแชทต้องการส่งข้อความต่อเนื่องทันที
- **WebSocket** เปิดท่อค้างไว้แบบสองทาง (Full-Duplex) ทำให้ Server ดัน (Push) ข้อความหาผู้รับได้ทันที
- จุดเปลี่ยนเกิดขึ้นตอนเริ่มต้นที่ Browser ส่ง HTTP Request พร้อม Header `Upgrade: websocket` แล้ว Server ตอบกลับด้วยรหัส **`101 Switching Protocols`**

### Phase 4 (R3, R4): *"How does the server guarantee uniqueness even with two simultaneous registrations?"*
- ใน `server/hub.go` คำขอจากทุกท่อจะต้องส่งเข้าคิว **`h.events` (Go Channel)** ซึ่งทำหน้าที่เหมือนประตูหมุนที่รับเข้าคิวทีละคน
- และมี Goroutine ของ `Hub.Run` เพียงตัวเดียวที่เป็นคนหยิบคำขอจากคิวมาเรียก [`onRegister`](../../server/users_dm.go) ทีละรายการจนจบ (เช็ก `h.byKey[key]` ด้วยกุญแจตัวพิมพ์เล็ก `normalizeKey` แล้วบันทึกชื่อทันที) ดังนั้นคนที่เข้ามาเป็นคิวที่ 2 จึงเห็นชื่อถูกใช้ไปแล้วเสมอโดยไม่เกิด Race Condition แม้ไม่ใช้ Mutex

### Phase 5 (R5, R7): *"Show me the exact line where a message is sent to only one person, and how spoofing is prevented."*
- **บรรทัดที่ส่งแบบเจาะจง (Unicast):** ใน [`onDM`](../../server/users_dm.go) บรรทัด:
  ```go
  h.deliver(c, c, "dm", m)   // ส่งยืนยันกลับไปที่ท่อของผู้ส่งเอง
  h.deliver(to, c, "dm", m)  // ส่งไปที่ท่อของผู้รับปลายทาง (to) เพียงคนเดียว
  ```
  ข้อความจะถูกส่งเข้าเฉพาะท่อของผู้ส่งและผู้รับเท่านั้น ไม่มีการส่งเข้าท่อของคนที่ 3
- **การป้องกันการปลอมชื่อ (Anti-Spoofing):** ผู้ส่งส่งมาแค่ `to` และเนื้อหาข้อความ ส่วนชื่อผู้ส่ง (`From: c.name`) **Server เป็นคนประทับตราเอง** จากท่อเชื่อมต่อที่ลงทะเบียนไว้แล้ว

### Phase 7 (F3 — Stickers): *"What exactly travels over the WebSocket when I send a sticker?"*
- ส่งเฉพาะข้อความ JSON ที่มีรหัสสติกเกอร์สั้นๆ เช่น `{"kind":"sticker","stickerId":"smile"}` วิ่งผ่านท่อ WebSocket เพื่อประหยัดพื้นที่และส่งได้รวดเร็ว
- ส่วนไฟล์ภาพ `/stickers/smile.svg` หน้าเว็บจะไปเปิดจาก Static Asset เอง โดย Server จะตรวจสอบรหัสด้วย `isSticker` และ `validContent` เสมอเพื่อไม่ให้ส่งรหัสปลอมได้
