package main

import (
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// ============================================================================
// P2 Backend Module — ด่านที่ 3: Identity & Live User List (R3, R4)
// ============================================================================

// validName (งานย่อย 3.1): ตัดช่องว่างหน้า-หลัง และตรวจสอบความยาว (1..max ตัวอักษร)
// รวมถึงห้ามมีอักขระควบคุม (control characters เช่น \n, \t, \x00)
func validName(raw string, max int) (string, bool) {
	s := strings.TrimSpace(raw)
	n := utf8.RuneCountInString(s)
	if n < 1 || n > max {
		return "", false
	}
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return "", false
		}
	}
	return s, true
}

// normalizeKey (งานย่อย 3.1): แปลงชื่อเป็นตัวพิมพ์เล็กทั้งหมดเพื่อใช้เป็นกุญแจเช็กชื่อซ้ำ (Case-Insensitive)
func normalizeKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// broadcastUsers (งานย่อย 3.3): ส่งรายชื่อคนออนไลน์ทั้งหมด (รวมถึงชื่อของผู้รับเอง ตามเกณฑ์ R4)
// เรียงตามตัวอักษรไปให้ลูกค้าทุกคนที่ลงทะเบียนแล้ว
func (h *Hub) broadcastUsers() {
	names := make([]string, 0, len(h.byKey))
	for _, c := range h.byKey {
		names = append(names, c.name)
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j])
	})
	h.broadcast("user_list", userList{Users: names})
}

// onRegister (งานย่อย 3.2): จัดการคำขอลงทะเบียนชื่อผู้ใช้ (R3, R4)
// ทำงานอยู่ใน Goroutine เดียวของ Hub.Run เสมอ จึงไม่มีปัญหาแย่งตั้งชื่อพร้อมกัน (Race Condition)
func (h *Hub) onRegister(c *Client, env Envelope) {
	req, ok := decode[registerReq](h, c, env)
	if !ok {
		return
	}
	if c.registered {
		h.fail(c, env.ReqID, codeBadRequest, "already registered")
		return
	}

	// 1. ตัดช่องว่างหน้า-หลัง และเช็กความถูกต้องของชื่อ (1-20 ตัวอักษร)
	name, valid := validName(req.Name, maxNameLen)
	if !valid {
		h.send(c, "register_error", registerErr{"invalid", "name must be 1 to 20 printable characters"}, env.ReqID)
		return
	}

	// 2. แปลงเป็นตัวพิมพ์เล็กเพื่อเช็กว่ามีชื่อนี้อยู่ในสมุด h.byKey แล้วหรือยัง
	key := normalizeKey(name)
	if _, taken := h.byKey[key]; taken {
		h.send(c, "register_error", registerErr{"taken", "that name is already in use"}, env.ReqID)
		return
	}

	// 3. บันทึกทั้งชื่อสำหรับแสดงผล (c.name) และกุญแจตัวพิมพ์เล็ก (c.key) ลงในสมุด
	c.name, c.key, c.registered = name, key, true
	h.byKey[key] = c

	// 4. ตอบยืนยันให้เจ้าตัว และกระจายรายชื่อคนออนไลน์ชุดใหม่ให้ทุกคน (R4)
	h.send(c, "register_ok", registerOK{Name: name}, env.ReqID)
	h.broadcastUsers()
	h.sendGroups(c)
	h.sendUnread(c)
}

// ============================================================================
// P2 Backend Module — ด่านที่ 4: Direct Messages & Rooms (R5, R7)
// ============================================================================

// dmRoom (งานย่อย 4.1): สร้างรหัสห้องแชทส่วนตัวระหว่างผู้ใช้ 2 คน (R5)
// เรียงลำดับชื่อตัวพิมพ์เล็กตามตัวอักษร (a < b) เพื่อให้ Alice->Bob และ Bob->Alice ได้ห้อง "dm:alice|bob" เหมือนกันเสมอ
func dmRoom(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return "dm:" + a + "|" + b
}

// onDM (งานย่อย 4.2): จัดการการส่งข้อความส่วนตัวแบบ 1-ต่อ-1 (R5, R7)
func (h *Hub) onDM(c *Client, env Envelope) {
	req, ok := decode[dmReq](h, c, env)
	if !ok {
		return
	}

	// 1. ตรวจสอบความถูกต้องของเนื้อหา (ข้อความธรรมดา หรือ สติกเกอร์)
	text, code, msg := validContent(req.Kind, req.Text, req.StickerID)
	if code != "" {
		h.fail(c, env.ReqID, code, msg)
		return
	}

	// 2. ค้นหาท่อของผู้รับจากสมุดรายชื่อ h.byKey ด้วยกุญแจตัวพิมพ์เล็ก
	toKey := normalizeKey(req.To)
	to := h.byKey[toKey]
	if to == nil {
		h.fail(c, env.ReqID, codeNoSuchUser, "that user is not online")
		return
	}
	if to == c {
		h.fail(c, env.ReqID, codeBadRequest, "you cannot message yourself")
		return
	}

	// 3. สร้างซองข้อความโดย Server เป็นผู้ประทับชื่อผู้ส่ง (From: c.name) เองเสมอ ป้องกันการปลอมตัว
	m := Message{
		Room:      dmRoom(c.key, to.key),
		From:      c.name,
		To:        to.name,
		Kind:      req.Kind,
		Text:      text,
		StickerID: req.StickerID,
		Ts:        time.Now().UnixMilli(),
		MsgID:     h.newMsgID(),
	}

	// 4. ส่งแบบเจาะจง (Unicast) ให้แค่ 2 คนเท่านั้น คือ ผู้ส่ง (เพื่อยืนยันว่าส่งถึงจริง) และ ผู้รับ (R7)
	h.deliver(c, c, "dm", m)
	h.deliver(to, c, "dm", m)
}
