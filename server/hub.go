package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxUnreadKept  = 200 // newest unread messages kept per room, for summaries
	summaryTimeout = 25 * time.Second
)

// Events posted to the hub goroutine.
type (
	inbound struct {
		c   *Client
		env Envelope
	}
	leave       struct{ c *Client }
	summaryDone struct {
		c     *Client
		room  string
		text  string
		count int
		err   error
	}
)

type Group struct {
	ID      string
	Name    string
	Seq     int
	Members []string // lowercased keys, in join order
}

func (g *Group) has(key string) bool {
	for _, m := range g.Members {
		if m == key {
			return true
		}
	}
	return false
}

type unreadBox struct {
	Count int
	Msgs  []Message
}

// Hub owns ALL shared state. Only the goroutine running Run touches these
// fields, so there are no locks: other goroutines talk to the hub by posting
// events on h.events (design D4: one owner goroutine + channels).
type Hub struct {
	events chan any
	quit   chan struct{}

	byKey    map[string]*Client
	groups   map[string]*Group
	invites  map[string]string // token -> groupID
	unread   map[string]map[string]*unreadBox
	groupSeq int
	msgSeq   int

	summarizer Summarizer
	upgrader   websocket.Upgrader
}

func NewHub(summarizer Summarizer, allowAnyOrigin bool) *Hub {
	h := &Hub{
		events:     make(chan any, 256),
		quit:       make(chan struct{}),
		byKey:      map[string]*Client{},
		groups:     map[string]*Group{},
		invites:    map[string]string{},
		unread:     map[string]map[string]*unreadBox{},
		summarizer: summarizer,
	}
	h.upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	if allowAnyOrigin {
		h.upgrader.CheckOrigin = func(*http.Request) bool { return true }
	}
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.quit:
			return
		case ev := <-h.events:
			switch e := ev.(type) {
			case inbound:
				h.handle(e.c, e.env)
			case leave:
				h.onLeave(e.c)
			case summaryDone:
				h.onSummaryDone(e)
			}
		}
	}
}

func (h *Hub) Stop() { close(h.quit) }

func (h *Hub) post(ev any) {
	select {
	case h.events <- ev:
	case <-h.quit:
	}
}

// ServeWS upgrades an HTTP request to a WebSocket and runs the client.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := newClient(h, conn)
	go c.writePump()
	c.readPump()
}

// ---- sending helpers -------------------------------------------------------

func (h *Hub) send(c *Client, typ string, payload any, reqID string) {
	if !c.enqueue(encode(typ, payload, reqID)) {
		c.kill()
	}
}

func (h *Hub) fail(c *Client, reqID, code, msg string) {
	h.send(c, "error", errorPayload{code, msg}, reqID)
}

func (h *Hub) broadcast(typ string, payload any) {
	frame := encode(typ, payload, "")
	for _, c := range h.byKey {
		if !c.enqueue(frame) {
			c.kill()
		}
	}
}

func (h *Hub) groupInfos() []GroupInfo {
	gs := make([]*Group, 0, len(h.groups))
	for _, g := range h.groups {
		gs = append(gs, g)
	}
	sort.Slice(gs, func(i, j int) bool { return gs[i].Seq < gs[j].Seq })
	out := make([]GroupInfo, 0, len(gs))
	for _, g := range gs {
		members := make([]string, 0, len(g.Members))
		for _, k := range g.Members {
			if u := h.byKey[k]; u != nil {
				members = append(members, u.name)
			}
		}
		out = append(out, GroupInfo{GroupID: g.ID, Name: g.Name, Members: members})
	}
	return out
}

func (h *Hub) broadcastGroups() { h.broadcast("group_list", groupList{Groups: h.groupInfos()}) }

func (h *Hub) sendGroups(c *Client) { h.send(c, "group_list", groupList{Groups: h.groupInfos()}, "") }

func (h *Hub) sendUnread(c *Client) {
	counts := map[string]int{}
	for room, box := range h.unread[c.key] {
		counts[room] = box.Count
	}
	h.send(c, "unread", map[string]any{"counts": counts}, "")
}

// ---- dispatch --------------------------------------------------------------

func decode[T any](h *Hub, c *Client, env Envelope) (T, bool) {
	var v T
	raw := env.Payload
	if len(raw) == 0 {
		raw = []byte("{}")
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		h.fail(c, env.ReqID, codeBadRequest, "malformed payload")
		return v, false
	}
	return v, true
}

func (h *Hub) handle(c *Client, env Envelope) {
	if !c.registered && env.Type != "register" {
		h.fail(c, env.ReqID, codeNotRegistered, "register a name first")
		return
	}
	switch env.Type {
	case "register":
		h.onRegister(c, env)
	case "dm":
		h.onDM(c, env)
	case "create_group":
		h.onCreateGroup(c, env)
	case "join_group":
		h.onJoinGroup(c, env)
	case "group_msg":
		h.onGroupMsg(c, env)
	case "create_invite":
		h.onCreateInvite(c, env)
	case "invite_info":
		h.onInviteInfo(c, env)
	case "open_room":
		h.onOpenRoom(c, env, true)
	case "close_room":
		h.onOpenRoom(c, env, false)
	case "summarize_request":
		h.onSummarize(c, env)
	default:
		h.fail(c, env.ReqID, codeBadRequest, "unknown type")
	}
}

// ---- disconnect cleanup ----------------------------------------------------
// P2's functions (validName, broadcastUsers, onRegister, dmRoom, onDM)
// will be implemented in server/users_dm.go during Stages 3 & 4.

func (h *Hub) onLeave(c *Client) {
	c.kill()
	if !c.registered {
		return
	}
	c.registered = false
	if h.byKey[c.key] == c {
		delete(h.byKey, c.key)
	}
	for id, g := range h.groups {
		kept := g.Members[:0]
		for _, m := range g.Members {
			if m != c.key {
				kept = append(kept, m)
			}
		}
		g.Members = kept
		if len(g.Members) == 0 {
			h.deleteGroup(id)
		}
	}
	delete(h.unread, c.key)
	h.broadcastUsers()
	h.broadcastGroups()
}

func (h *Hub) deleteGroup(id string) {
	delete(h.groups, id)
	for tok, gid := range h.invites {
		if gid == id {
			delete(h.invites, tok)
		}
	}
}

// ---- rooms, delivery, unread ----------------------------------------------

func (h *Hub) roomAccess(c *Client, room string) (label, code string) {
	switch {
	case strings.HasPrefix(room, "dm:"):
		parts := strings.SplitN(strings.TrimPrefix(room, "dm:"), "|", 2)
		if len(parts) != 2 || (parts[0] != c.key && parts[1] != c.key) {
			return "", codeNotMember
		}
		other := parts[0]
		if other == c.key {
			other = parts[1]
		}
		if u := h.byKey[other]; u != nil {
			other = u.name
		}
		return "DM with " + other, ""
	case strings.HasPrefix(room, "group:"):
		g := h.groups[strings.TrimPrefix(room, "group:")]
		if g == nil {
			return "", codeNoSuchGroup
		}
		if !g.has(c.key) {
			return "", codeNotMember
		}
		return "Group " + g.Name, ""
	}
	return "", codeBadRequest
}

func (h *Hub) newMsgID() string {
	h.msgSeq++
	return fmt.Sprintf("m_%d", h.msgSeq)
}

func (h *Hub) deliver(to, from *Client, typ string, m Message) {
	h.send(to, typ, m, "")
	if to == from || to.openRooms[m.Room] {
		return
	}
	boxes := h.unread[to.key]
	if boxes == nil {
		boxes = map[string]*unreadBox{}
		h.unread[to.key] = boxes
	}
	box := boxes[m.Room]
	if box == nil {
		box = &unreadBox{}
		boxes[m.Room] = box
	}
	box.Count++
	box.Msgs = append(box.Msgs, m)
	if len(box.Msgs) > maxUnreadKept {
		box.Msgs = box.Msgs[len(box.Msgs)-maxUnreadKept:]
	}
	h.sendUnread(to)
}

func (h *Hub) onOpenRoom(c *Client, env Envelope, open bool) {
	req, ok := decode[roomReq](h, c, env)
	if !ok {
		return
	}
	if _, code := h.roomAccess(c, req.Room); code != "" {
		h.fail(c, env.ReqID, code, "you cannot access that room")
		return
	}
	if !open {
		delete(c.openRooms, req.Room)
		return
	}
	c.openRooms[req.Room] = true
	if boxes := h.unread[c.key]; boxes != nil {
		if _, had := boxes[req.Room]; had {
			delete(boxes, req.Room)
			h.sendUnread(c)
		}
	}
}

// ---- groups (R8 to R11) ----------------------------------------------------

func (h *Hub) onCreateGroup(c *Client, env Envelope) {
	req, ok := decode[createGroupReq](h, c, env)
	if !ok {
		return
	}
	name, valid := validName(req.Name, maxGroupNameLen)
	if !valid {
		h.fail(c, env.ReqID, codeBadRequest, "group name must be 1 to 30 printable characters")
		return
	}
	h.groupSeq++
	g := &Group{ID: fmt.Sprintf("g_%d", h.groupSeq), Name: name, Seq: h.groupSeq, Members: []string{c.key}}
	h.groups[g.ID] = g
	h.send(c, "create_group_ok", map[string]string{"groupId": g.ID}, env.ReqID)
	h.broadcastGroups()
}

func (h *Hub) onJoinGroup(c *Client, env Envelope) {
	req, ok := decode[joinGroupReq](h, c, env)
	if !ok {
		return
	}
	g := h.groups[req.GroupID]
	if g == nil {
		h.fail(c, env.ReqID, codeNoSuchGroup, "no such group")
		return
	}
	if req.Invite != "" && h.invites[req.Invite] != g.ID {
		h.fail(c, env.ReqID, codeBadInvite, "invite is invalid for this group")
		return
	}
	if g.has(c.key) {
		h.fail(c, env.ReqID, codeAlreadyMember, "you are already a member")
		return
	}
	g.Members = append(g.Members, c.key)
	h.send(c, "join_group_ok", map[string]string{"groupId": g.ID}, env.ReqID)
	h.broadcastGroups()
}

func (h *Hub) onGroupMsg(c *Client, env Envelope) {
	req, ok := decode[groupMsgReq](h, c, env)
	if !ok {
		return
	}
	g := h.groups[req.GroupID]
	if g == nil {
		h.fail(c, env.ReqID, codeNoSuchGroup, "no such group")
		return
	}
	if !g.has(c.key) {
		h.fail(c, env.ReqID, codeNotMember, "join the group first")
		return
	}
	text, code, msg := validContent(req.Kind, req.Text, req.StickerID)
	if code != "" {
		h.fail(c, env.ReqID, code, msg)
		return
	}
	m := Message{
		Room: "group:" + g.ID, From: c.name, GroupID: g.ID,
		Kind: req.Kind, Text: text, StickerID: req.StickerID,
		Ts: time.Now().UnixMilli(), MsgID: h.newMsgID(),
	}
	for _, key := range g.Members {
		if member := h.byKey[key]; member != nil {
			h.deliver(member, c, "group_msg", m)
		}
	}
}

// ---- invites / QR ----------------------------------------------------------

func (h *Hub) onCreateInvite(c *Client, env Envelope) {
	req, ok := decode[groupIDReq](h, c, env)
	if !ok {
		return
	}
	g := h.groups[req.GroupID]
	if g == nil {
		h.fail(c, env.ReqID, codeNoSuchGroup, "no such group")
		return
	}
	if !g.has(c.key) {
		h.fail(c, env.ReqID, codeNotMember, "only members can create invites")
		return
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		h.fail(c, env.ReqID, codeBadRequest, "could not create invite")
		return
	}
	token := hex.EncodeToString(buf)
	h.invites[token] = g.ID
	h.send(c, "invite_created", map[string]string{"groupId": g.ID, "token": token}, env.ReqID)
}

func (h *Hub) onInviteInfo(c *Client, env Envelope) {
	req, ok := decode[tokenReq](h, c, env)
	if !ok {
		return
	}
	g := h.groups[h.invites[req.Token]]
	if g == nil {
		h.fail(c, env.ReqID, codeBadInvite, "invite is invalid or the group no longer exists")
		return
	}
	h.send(c, "invite_info_ok", map[string]any{
		"groupId": g.ID, "groupName": g.Name, "memberCount": len(g.Members),
		"alreadyMember": g.has(c.key),
	}, env.ReqID)
}

// ---- AI summary ------------------------------------------------------------

func (h *Hub) onSummarize(c *Client, env Envelope) {
	req, ok := decode[roomReq](h, c, env)
	if !ok {
		return
	}
	label, code := h.roomAccess(c, req.Room)
	if code != "" {
		h.summaryErr(c, req.Room, "not_allowed", "you cannot access that room", env.ReqID)
		return
	}
	box := h.unread[c.key][req.Room]
	if box == nil || len(box.Msgs) == 0 {
		h.summaryErr(c, req.Room, "nothing_unread", "no unread messages in this room", env.ReqID)
		return
	}
	if h.summarizer == nil {
		h.summaryErr(c, req.Room, "llm_error", errNotConfigured.Error(), env.ReqID)
		return
	}
	if c.summarizing {
		h.fail(c, env.ReqID, codeRateLimited, "a summary is already in progress")
		return
	}
	c.summarizing = true

	lines := make([]string, 0, len(box.Msgs))
	for _, m := range box.Msgs {
		if m.Kind == "sticker" {
			lines = append(lines, fmt.Sprintf("%s: [sticker: %s]", m.From, m.StickerID))
		} else {
			lines = append(lines, fmt.Sprintf("%s: %s", m.From, m.Text))
		}
	}
	count := box.Count
	summarizer := h.summarizer
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), summaryTimeout)
		defer cancel()
		text, err := summarizer(ctx, label, lines)
		h.post(summaryDone{c: c, room: req.Room, text: text, count: count, err: err})
	}()
}

func (h *Hub) summaryErr(c *Client, room, code, msg, reqID string) {
	h.send(c, "summary_error", map[string]string{"room": room, "code": code, "message": msg}, reqID)
}

func (h *Hub) onSummaryDone(e summaryDone) {
	e.c.summarizing = false
	if !e.c.registered {
		return
	}
	if e.err != nil {
		log.Println("summary failed:", e.err) // detail stays on the server; the client gets a generic message
		code, msg := "llm_error", "the AI service failed, try again later"
		switch {
		case errors.Is(e.err, context.DeadlineExceeded):
			code, msg = "timeout", "the AI service took too long"
		case errors.Is(e.err, errNotConfigured):
			msg = e.err.Error()
		}
		h.summaryErr(e.c, e.room, code, msg, "")
		return
	}
	h.send(e.c, "summary_result", map[string]any{
		"room": e.room, "text": e.text, "messageCount": e.count,
	}, "")
}
