package main

import (
	"strings"
	"unicode/utf8"
)

// ============================================================================
// P2 Stickers & Content Validation Module — ด่านที่ 5 (Phase 7: F3)
// ============================================================================

// stickerPack คือรายชื่อรหัสสติกเกอร์มาตรฐานทั้ง 8 ตัวในระบบ
// รหัสทุกตัวจะตรงกับไฟล์รูปภาพใน client-web/public/stickers/<id>.svg
// เวลาส่งสติกเกอร์ จะส่งเฉพาะรหัสข้อความ (id) วิ่งผ่าน WebSocket เท่านั้น (ไม่ส่งไฟล์รูป)
var stickerPack = []string{
	"smile", "laugh", "love", "cry", "angry", "surprised", "sleepy", "thumbs_up",
}

// stickerSet แปลงรายการ stickerPack เป็นแผนที่ (map) เพื่อให้ค้นหาว่ามีรหัสนี้จริงหรือไม่ได้อย่างรวดเร็ว O(1)
var stickerSet = func() map[string]bool {
	m := make(map[string]bool, len(stickerPack))
	for _, id := range stickerPack {
		m[id] = true
	}
	return m
}()

// isSticker (งานย่อย 5.2): ตรวจสอบว่ารหัสสติกเกอร์ที่ส่งมาอยู่ในชุดที่ระบบอนุญาตหรือไม่
func isSticker(id string) bool {
	return stickerSet[id]
}

// validContent (งานย่อย 5.2): ตรวจสอบความถูกต้องของเนื้อหาข้อความแชทฝั่ง Server (Never Trust the Client)
// รองรับ 2 รูปแบบ (kind):
//   - "text": ต้องไม่ว่างเปล่าหลังตัดช่องว่าง และยาวไม่เกิน 1,000 ตัวอักษร (maxTextLen)
//   - "sticker": ต้องเป็นรหัสสติกเกอร์ที่มีอยู่ใน stickerPack เท่านั้น มิฉะนั้นคืนรหัส error "bad_sticker"
func validContent(kind, text, stickerID string) (cleanText, code, msg string) {
	switch kind {
	case "text":
		t := strings.TrimSpace(text)
		if t == "" || utf8.RuneCountInString(t) > maxTextLen {
			return "", codeBadRequest, "text must be 1 to 1000 characters"
		}
		return t, "", ""
	case "sticker":
		if !isSticker(stickerID) {
			return "", codeBadSticker, "unknown sticker id"
		}
		return "", "", ""
	default:
		return "", codeBadRequest, "kind must be text or sticker"
	}
}
