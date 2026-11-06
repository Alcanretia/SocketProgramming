package main

import (
	"strings"
	"testing"
)

// ============================================================================
// P2 Unit Tests (ด่านที่ 7): ทดสอบเคสพิเศษของการตั้งชื่อ, รหัสห้อง DM, และสติกเกอร์
// ============================================================================

// TestTrickyNames ทดสอบฟังก์ชัน validName กับเคสพิเศษต่างๆ (ภาษาไทย, อิโมจิ, เว้นวรรคล้วน, ยาว 500 ตัว, อักขระควบคุม)
func TestTrickyNames(t *testing.T) {
	cases := []struct {
		name      string
		raw       string
		wantClean string
		wantValid bool
	}{
		{
			name:      "Thai and Emoji within 20 runes",
			raw:       "  สมชาย 😊  ",
			wantClean: "สมชาย 😊",
			wantValid: true,
		},
		{
			name:      "Spaces only",
			raw:       "     ",
			wantClean: "",
			wantValid: false,
		},
		{
			name:      "500 characters (too long)",
			raw:       strings.Repeat("ก", 500),
			wantClean: "",
			wantValid: false,
		},
		{
			name:      "Control character newline inside name",
			raw:       "Alice\nBob",
			wantClean: "",
			wantValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clean, ok := validName(tc.raw, maxNameLen)
			if ok != tc.wantValid || clean != tc.wantClean {
				t.Fatalf("validName(%q) = (%q, %v), want (%q, %v)",
					tc.raw, clean, ok, tc.wantClean, tc.wantValid)
			}
		})
	}
}

// TestDMRoomSymmetry ทดสอบว่าไม่ว่าใครจะทักก่อน รหัสห้อง DM จะออกมาเหมือนกันเสมอ (R5)
func TestDMRoomSymmetry(t *testing.T) {
	r1 := dmRoom(normalizeKey("Alice"), normalizeKey("Bob"))
	r2 := dmRoom(normalizeKey("  bOb "), normalizeKey("ALICE"))
	if r1 != "dm:alice|bob" || r1 != r2 {
		t.Fatalf("dmRoom mismatch: r1=%q, r2=%q", r1, r2)
	}
}

// TestValidContentAndStickers ทดสอบการดักสติกเกอร์ปลอมและข้อความว่างฝั่ง Server (F3)
func TestValidContentAndStickers(t *testing.T) {
	// สติกเกอร์ที่มีอยู่จริงในระบบ 8 ตัวต้องผ่าน
	if _, code, _ := validContent("sticker", "", "smile"); code != "" {
		t.Fatalf("expected valid sticker 'smile', got error code %q", code)
	}
	// สติกเกอร์ปลอมต้องถูกปฏิเสธด้วยรหัส bad_sticker
	if _, code, _ := validContent("sticker", "", "hacked_dragon"); code != codeBadSticker {
		t.Fatalf("expected %q for unknown sticker, got %q", codeBadSticker, code)
	}
	// ข้อความที่มีแต่เว้นวรรคต้องถูกปฏิเสธด้วยรหัส bad_request
	if _, code, _ := validContent("text", "   ", ""); code != codeBadRequest {
		t.Fatalf("expected %q for empty text, got %q", codeBadRequest, code)
	}
}
