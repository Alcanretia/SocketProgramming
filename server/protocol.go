package main

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

// Limits from docs/PROTOCOL.md.
const (
	maxNameLen      = 20
	maxGroupNameLen = 30
	maxTextLen      = 1000
	maxFrameBytes   = 8192
)

// Envelope is every frame the client sends: {"type","payload","reqId"}.
type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
	ReqID   string          `json:"reqId,omitempty"`
}

// outFrame is every frame the server sends.
type outFrame struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
	ReqID   string `json:"reqId,omitempty"`
}

func encode(typ string, payload any, reqID string) []byte {
	if payload == nil {
		payload = struct{}{}
	}
	b, _ := json.Marshal(outFrame{Type: typ, Payload: payload, ReqID: reqID})
	return b
}

// Error codes (PROTOCOL.md 3.8).
const (
	codeBadRequest    = "bad_request"
	codeNotRegistered = "not_registered"
	codeNoSuchUser    = "no_such_user"
	codeNoSuchGroup   = "no_such_group"
	codeNotMember     = "not_member"
	codeAlreadyMember = "already_member"
	codeBadSticker    = "bad_sticker"
	codeBadInvite     = "bad_invite"
	codeRateLimited   = "rate_limited"
)

// Inbound payloads.
type registerReq struct {
	Name string `json:"name"`
}

type dmReq struct {
	To        string `json:"to"`
	Kind      string `json:"kind"`
	Text      string `json:"text"`
	StickerID string `json:"stickerId"`
}

type createGroupReq struct {
	Name string `json:"name"`
}

type joinGroupReq struct {
	GroupID string `json:"groupId"`
	Invite  string `json:"invite"`
}

type groupMsgReq struct {
	GroupID   string `json:"groupId"`
	Kind      string `json:"kind"`
	Text      string `json:"text"`
	StickerID string `json:"stickerId"`
}

type groupIDReq struct {
	GroupID string `json:"groupId"`
}

type tokenReq struct {
	Token string `json:"token"`
}

type roomReq struct {
	Room string `json:"room"`
}

// Outbound payloads.
type registerOK struct {
	Name string `json:"name"`
}

type registerErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type userList struct {
	Users []string `json:"users"`
}

// Message is a delivered DM or group message.
type Message struct {
	Room      string `json:"room"`
	From      string `json:"from"`
	To        string `json:"to,omitempty"`
	GroupID   string `json:"groupId,omitempty"`
	Kind      string `json:"kind"`
	Text      string `json:"text,omitempty"`
	StickerID string `json:"stickerId,omitempty"`
	Ts        int64  `json:"ts"`
	MsgID     string `json:"msgId"`
}

type GroupInfo struct {
	GroupID string   `json:"groupId"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type groupList struct {
	Groups []GroupInfo `json:"groups"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// validName trims and checks a user or group name.
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

// validContent checks a text-or-sticker body. The server, not the client,
// decides what is acceptable (unknown sticker ids are rejected here).
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
