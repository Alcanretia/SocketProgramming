package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// ---- test harness: real WebSocket clients against a real hub --------------

type frame struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	ReqID   string          `json:"reqId"`
}

func (f frame) into(t *testing.T, v any) {
	t.Helper()
	if err := json.Unmarshal(f.Payload, v); err != nil {
		t.Fatalf("decode %s payload: %v", f.Type, err)
	}
}

func (f frame) str(t *testing.T, key string) string {
	t.Helper()
	var m map[string]any
	f.into(t, &m)
	s, _ := m[key].(string)
	return s
}

type testClient struct {
	t    *testing.T
	conn *websocket.Conn
	in   chan frame
}

func startHub(t *testing.T, s Summarizer) string {
	t.Helper()
	h := NewHub(s, false)
	go h.Run()
	srv := httptest.NewServer(http.HandlerFunc(h.ServeWS))
	t.Cleanup(func() { srv.Close(); h.Stop() })
	return "ws" + strings.TrimPrefix(srv.URL, "http")
}

func dial(t *testing.T, url string) *testClient {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	c := &testClient{t: t, conn: conn, in: make(chan frame, 512)}
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				close(c.in)
				return
			}
			var f frame
			json.Unmarshal(data, &f)
			c.in <- f
		}
	}()
	t.Cleanup(func() { conn.Close() })
	return c
}

func (c *testClient) send(typ string, payload any) {
	c.t.Helper()
	b, _ := json.Marshal(map[string]any{"type": typ, "payload": payload})
	if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
		c.t.Fatal("write:", err)
	}
}

// expect returns the next frame of the given type (skipping others).
func (c *testClient) expect(typ string) frame {
	c.t.Helper()
	timeout := time.After(3 * time.Second)
	for {
		select {
		case f, ok := <-c.in:
			if !ok {
				c.t.Fatalf("connection closed while waiting for %q", typ)
			}
			if f.Type == typ {
				return f
			}
		case <-timeout:
			c.t.Fatalf("timed out waiting for %q", typ)
		}
	}
}

// expectNone fails if a frame of the given type arrives within d.
func (c *testClient) expectNone(typ string, d time.Duration) {
	c.t.Helper()
	timeout := time.After(d)
	for {
		select {
		case f, ok := <-c.in:
			if !ok {
				return
			}
			if f.Type == typ {
				c.t.Fatalf("unexpected %q frame: %s", typ, f.Payload)
			}
		case <-timeout:
			return
		}
	}
}

func (c *testClient) register(name string) {
	c.t.Helper()
	c.send("register", map[string]any{"name": name})
	c.expect("register_ok")
}

func (c *testClient) errCode() string {
	c.t.Helper()
	return c.expect("error").str(c.t, "code")
}

const quiet = 300 * time.Millisecond

// ---- R3 + R4: names and user list ------------------------------------------

func TestRegisterAndUserList(t *testing.T) {
	url := startHub(t, nil)
	alice, bob := dial(t, url), dial(t, url)

	// Before registering, anything but `register` is refused.
	alice.send("dm", map[string]any{"to": "x", "kind": "text", "text": "hi"})
	if got := alice.errCode(); got != "not_registered" {
		t.Fatalf("want not_registered, got %s", got)
	}

	alice.register("Alice")
	// R4: the user list includes the client's own name.
	var ul struct{ Users []string }
	alice.expect("user_list").into(t, &ul)
	if len(ul.Users) != 1 || ul.Users[0] != "Alice" {
		t.Fatalf("want [Alice], got %v", ul.Users)
	}

	// R3: names are unique, case-insensitively.
	bob.send("register", map[string]any{"name": "aLiCe"})
	if got := bob.expect("register_error").str(t, "code"); got != "taken" {
		t.Fatalf("want taken, got %s", got)
	}
	bob.send("register", map[string]any{"name": "   "})
	if got := bob.expect("register_error").str(t, "code"); got != "invalid" {
		t.Fatalf("want invalid, got %s", got)
	}
	bob.send("register", map[string]any{"name": strings.Repeat("x", 21)})
	if got := bob.expect("register_error").str(t, "code"); got != "invalid" {
		t.Fatalf("want invalid for long name, got %s", got)
	}

	bob.register("Bob")
	// Live update on join (both see both names).
	for _, c := range []*testClient{alice, bob} {
		deadline := time.After(2 * time.Second)
		for {
			var u struct{ Users []string }
			select {
			case f := <-c.in:
				if f.Type != "user_list" {
					continue
				}
				f.into(t, &u)
				if len(u.Users) == 2 {
					goto ok
				}
			case <-deadline:
				t.Fatal("never saw both names")
			}
		}
	ok:
	}

	// Live update on leave, and the name is freed.
	bob.conn.Close()
	for {
		var u struct{ Users []string }
		alice.expect("user_list").into(t, &u)
		if len(u.Users) == 1 {
			break
		}
	}
	again := dial(t, url)
	again.register("bob")
}

func TestSimultaneousRegistrationOnlyOneWins(t *testing.T) {
	url := startHub(t, nil)
	const n = 10
	results := make(chan string, n)
	for i := 0; i < n; i++ {
		c := dial(t, url)
		go func() {
			c.send("register", map[string]any{"name": "Same"})
			deadline := time.After(3 * time.Second)
			for {
				select {
				case f := <-c.in:
					if f.Type == "register_ok" || f.Type == "register_error" {
						results <- f.Type
						return
					}
				case <-deadline:
					results <- "timeout"
					return
				}
			}
		}()
	}
	ok := 0
	for i := 0; i < n; i++ {
		switch <-results {
		case "register_ok":
			ok++
		case "timeout":
			t.Fatal("a client got no answer")
		}
	}
	if ok != 1 {
		t.Fatalf("want exactly 1 winner, got %d", ok)
	}
}

func TestMalformedInputKeepsConnectionOpen(t *testing.T) {
	url := startHub(t, nil)
	a := dial(t, url)
	a.conn.WriteMessage(websocket.TextMessage, []byte("{not json"))
	if got := a.errCode(); got != "bad_request" {
		t.Fatalf("want bad_request, got %s", got)
	}
	a.send("nonsense", nil)
	if got := a.errCode(); got != "not_registered" {
		t.Fatalf("want not_registered, got %s", got)
	}
	a.register("Alice") // still usable
	a.send("nonsense", nil)
	if got := a.errCode(); got != "bad_request" {
		t.Fatalf("want bad_request, got %s", got)
	}
}

func TestOversizedFrameDropsConnection(t *testing.T) {
	url := startHub(t, nil)
	a := dial(t, url)
	a.register("Alice")
	a.conn.WriteMessage(websocket.TextMessage, []byte(strings.Repeat("x", maxFrameBytes+100)))
	select {
	case _, ok := <-drain(a):
		if ok {
			t.Fatal("expected the server to close the connection")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("connection was not closed")
	}
}

func drain(c *testClient) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for range c.in {
		}
		close(done)
	}()
	return done
}

// ---- R5 + R7: direct messages ----------------------------------------------

func TestDMPrivacyAndSpoofing(t *testing.T) {
	url := startHub(t, nil)
	alice, bob, carol := dial(t, url), dial(t, url), dial(t, url)
	alice.register("Alice")
	bob.register("Bob")
	carol.register("Carol")

	// A -> B reaches exactly A and B.
	alice.send("dm", map[string]any{"to": "bob", "kind": "text", "text": "secret"})
	for _, c := range []*testClient{alice, bob} {
		f := c.expect("dm")
		var m Message
		f.into(t, &m)
		if m.Text != "secret" || m.From != "Alice" || m.To != "Bob" || m.Room != "dm:alice|bob" {
			t.Fatalf("unexpected DM: %+v", m)
		}
	}
	// R7: the third tab sees nothing.
	carol.expectNone("dm", quiet)

	// Same room whichever direction it goes.
	bob.send("dm", map[string]any{"to": "ALICE", "kind": "text", "text": "reply"})
	var m Message
	alice.expect("dm").into(t, &m)
	if m.Room != "dm:alice|bob" {
		t.Fatalf("room should be symmetric, got %s", m.Room)
	}
	bob.expect("dm") // the sender also gets the server's copy of its own message

	// A client cannot choose the sender: extra "from" is ignored.
	alice.send("dm", map[string]any{"to": "bob", "kind": "text", "text": "x", "from": "Carol"})
	bob.expect("dm").into(t, &m)
	if m.From != "Alice" {
		t.Fatalf("sender must be stamped by the server, got %s", m.From)
	}

	// Errors are graceful.
	alice.send("dm", map[string]any{"to": "nobody", "kind": "text", "text": "hi"})
	if got := alice.errCode(); got != "no_such_user" {
		t.Fatalf("want no_such_user, got %s", got)
	}
	alice.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "hi"})
	if got := alice.errCode(); got != "bad_request" {
		t.Fatalf("want bad_request for self DM, got %s", got)
	}
	alice.send("dm", map[string]any{"to": "bob", "kind": "text", "text": "   "})
	if got := alice.errCode(); got != "bad_request" {
		t.Fatalf("want bad_request for empty text, got %s", got)
	}
}

// ---- R8 to R11: groups -----------------------------------------------------

func createGroup(t *testing.T, c *testClient, name string) string {
	t.Helper()
	c.send("create_group", map[string]any{"name": name})
	return c.expect("create_group_ok").str(t, "groupId")
}

func lastGroup(t *testing.T, c *testClient, id string, wantMembers int) GroupInfo {
	t.Helper()
	deadline := time.After(3 * time.Second)
	for {
		select {
		case f := <-c.in:
			if f.Type != "group_list" {
				continue
			}
			var gl groupList
			f.into(t, &gl)
			for _, g := range gl.Groups {
				if g.GroupID == id && len(g.Members) == wantMembers {
					return g
				}
			}
		case <-deadline:
			t.Fatalf("never saw group %s with %d members", id, wantMembers)
		}
	}
}

func TestGroupsMembershipAndPrivacy(t *testing.T) {
	url := startHub(t, nil)
	alice, bob, carol := dial(t, url), dial(t, url), dial(t, url)
	alice.register("Alice")
	bob.register("Bob")
	carol.register("Carol")

	gid := createGroup(t, alice, "Study")
	// R8: creator is the only member. R9: everyone sees the list with members.
	for _, c := range []*testClient{alice, bob, carol} {
		g := lastGroup(t, c, gid, 1)
		if g.Name != "Study" || g.Members[0] != "Alice" {
			t.Fatalf("unexpected group: %+v", g)
		}
	}

	// R10: nobody is added automatically; Bob has to join himself.
	bob.send("join_group", map[string]any{"groupId": gid})
	bob.expect("join_group_ok")
	g := lastGroup(t, carol, gid, 2)
	if g.Members[1] != "Bob" {
		t.Fatalf("want Bob second, got %v", g.Members)
	}
	bob.send("join_group", map[string]any{"groupId": gid})
	if got := bob.errCode(); got != "already_member" {
		t.Fatalf("want already_member, got %s", got)
	}

	// R11: members get it, the non-member (Carol) does not.
	alice.send("group_msg", map[string]any{"groupId": gid, "kind": "text", "text": "hello team"})
	for _, c := range []*testClient{alice, bob} {
		var m Message
		c.expect("group_msg").into(t, &m)
		if m.Text != "hello team" || m.From != "Alice" || m.Room != "group:"+gid {
			t.Fatalf("unexpected group msg: %+v", m)
		}
	}
	carol.expectNone("group_msg", quiet)

	// A non-member cannot send, and members never see the attempt.
	carol.send("group_msg", map[string]any{"groupId": gid, "kind": "text", "text": "let me in"})
	if got := carol.errCode(); got != "not_member" {
		t.Fatalf("want not_member, got %s", got)
	}
	alice.expectNone("group_msg", quiet)
	bob.expectNone("group_msg", quiet)

	carol.send("group_msg", map[string]any{"groupId": "g_999", "kind": "text", "text": "x"})
	if got := carol.errCode(); got != "no_such_group" {
		t.Fatalf("want no_such_group, got %s", got)
	}

	// Leaving (disconnect) removes the member; an empty group disappears.
	bob.conn.Close()
	lastGroup(t, alice, gid, 1)
	alice.conn.Close()
	deadline := time.After(3 * time.Second)
	for {
		var gl groupList
		select {
		case f := <-carol.in:
			if f.Type != "group_list" {
				continue
			}
			f.into(t, &gl)
			if len(gl.Groups) == 0 {
				return
			}
		case <-deadline:
			t.Fatal("empty group was not removed")
		}
	}
}

// ---- stickers --------------------------------------------------------------

func TestStickers(t *testing.T) {
	url := startHub(t, nil)
	alice, bob, carol := dial(t, url), dial(t, url), dial(t, url)
	alice.register("Alice")
	bob.register("Bob")
	carol.register("Carol")

	alice.send("dm", map[string]any{"to": "bob", "kind": "sticker", "stickerId": "smile"})
	var m Message
	bob.expect("dm").into(t, &m)
	if m.Kind != "sticker" || m.StickerID != "smile" || m.Text != "" {
		t.Fatalf("unexpected sticker msg: %+v", m)
	}
	carol.expectNone("dm", quiet)

	alice.send("dm", map[string]any{"to": "bob", "kind": "sticker", "stickerId": "../../etc/passwd"})
	if got := alice.errCode(); got != "bad_sticker" {
		t.Fatalf("want bad_sticker, got %s", got)
	}
	bob.expectNone("dm", quiet)

	gid := createGroup(t, alice, "G")
	alice.send("group_msg", map[string]any{"groupId": gid, "kind": "sticker", "stickerId": "nope"})
	if got := alice.errCode(); got != "bad_sticker" {
		t.Fatalf("want bad_sticker, got %s", got)
	}
}

// Every sticker id the server accepts must have an image in the frontend.
func TestStickerAssetsExist(t *testing.T) {
	dir := filepath.Join("..", "client-web", "public", "stickers")
	if _, err := os.Stat(dir); err != nil {
		t.Skip("client-web/public/stickers not present")
	}
	for _, id := range stickerPack {
		if _, err := os.Stat(filepath.Join(dir, id+".svg")); err != nil {
			t.Errorf("missing sticker image for %q: %v", id, err)
		}
	}
}

// ---- invites / QR ------------------------------------------------------------

func TestInviteFlow(t *testing.T) {
	url := startHub(t, nil)
	alice, bob, carol := dial(t, url), dial(t, url), dial(t, url)
	alice.register("Alice")
	bob.register("Bob")
	carol.register("Carol")
	gid := createGroup(t, alice, "Invite me")
	other := createGroup(t, alice, "Other")

	// Only members can create invites.
	bob.send("create_invite", map[string]any{"groupId": gid})
	if got := bob.errCode(); got != "not_member" {
		t.Fatalf("want not_member, got %s", got)
	}
	alice.send("create_invite", map[string]any{"groupId": gid})
	token := alice.expect("invite_created").str(t, "token")
	if len(token) < 32 {
		t.Fatalf("token looks guessable: %q", token)
	}

	// Looking at an invite changes nothing (R10): Bob is still not a member.
	bob.send("invite_info", map[string]any{"token": token})
	var info struct {
		GroupID, GroupName string
		MemberCount        int
	}
	bob.expect("invite_info_ok").into(t, &info)
	if info.GroupID != gid || info.GroupName != "Invite me" || info.MemberCount != 1 {
		t.Fatalf("unexpected invite info: %+v", info)
	}
	bob.send("group_msg", map[string]any{"groupId": gid, "kind": "text", "text": "x"})
	if got := bob.errCode(); got != "not_member" {
		t.Fatalf("scanning must not add anyone, got %s", got)
	}

	// A forged or mismatched token is rejected.
	bob.send("invite_info", map[string]any{"token": "deadbeef"})
	if got := bob.errCode(); got != "bad_invite" {
		t.Fatalf("want bad_invite, got %s", got)
	}
	carol.send("join_group", map[string]any{"groupId": other, "invite": token})
	if got := carol.errCode(); got != "bad_invite" {
		t.Fatalf("token for another group must fail, got %s", got)
	}

	// Joining is the user's own explicit action.
	bob.send("join_group", map[string]any{"groupId": gid, "invite": token})
	bob.expect("join_group_ok")
	lastGroup(t, carol, gid, 2)
}

// ---- unread + AI summary -----------------------------------------------------

func unreadCount(t *testing.T, c *testClient, room string, want int) {
	t.Helper()
	deadline := time.After(3 * time.Second)
	for {
		select {
		case f := <-c.in:
			if f.Type != "unread" {
				continue
			}
			var u struct{ Counts map[string]int }
			f.into(t, &u)
			if u.Counts[room] == want {
				return
			}
		case <-deadline:
			t.Fatalf("never saw unread[%s]=%d", room, want)
		}
	}
}

func TestUnreadAndSummary(t *testing.T) {
	type call struct {
		label string
		lines []string
	}
	calls := make(chan call, 1)
	fake := func(ctx context.Context, label string, lines []string) (string, error) {
		calls <- call{label, lines}
		return "SUMMARY", nil
	}
	url := startHub(t, fake)
	alice, bob, carol := dial(t, url), dial(t, url), dial(t, url)
	alice.register("Alice")
	bob.register("Bob")
	carol.register("Carol")
	room := "dm:alice|bob"
	bob.expect("unread") // the initial (empty) unread state sent after register_ok

	bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "one"})
	bob.send("dm", map[string]any{"to": "alice", "kind": "sticker", "stickerId": "love"})
	bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "three"})
	unreadCount(t, alice, room, 3)
	bob.expectNone("unread", quiet) // own messages are never unread for the sender

	// Carol cannot ask for a summary of a room she is not in, and gets nothing of it.
	carol.send("summarize_request", map[string]any{"room": room})
	if got := carol.expect("summary_error").str(t, "code"); got != "not_allowed" {
		t.Fatalf("want not_allowed, got %s", got)
	}
	carol.send("open_room", map[string]any{"room": room})
	if got := carol.errCode(); got != "not_member" {
		t.Fatalf("want not_member, got %s", got)
	}

	// Alice summarizes her unread messages: only those, stickers as text.
	alice.send("summarize_request", map[string]any{"room": room})
	res := alice.expect("summary_result")
	if res.str(t, "text") != "SUMMARY" {
		t.Fatalf("unexpected result: %s", res.Payload)
	}
	got := <-calls
	want := []string{"Bob: one", "Bob: [sticker: love]", "Bob: three"}
	if strings.Join(got.lines, "|") != strings.Join(want, "|") {
		t.Fatalf("LLM input = %v, want %v", got.lines, want)
	}
	if got.label != "DM with Bob" {
		t.Fatalf("label = %q", got.label)
	}

	// Opening the room clears unread; a message while it is open is not unread.
	alice.send("open_room", map[string]any{"room": room})
	unreadCount(t, alice, room, 0)
	bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "live"})
	alice.expect("dm")
	alice.expectNone("unread", quiet)
	alice.send("summarize_request", map[string]any{"room": room})
	if got := alice.expect("summary_error").str(t, "code"); got != "nothing_unread" {
		t.Fatalf("want nothing_unread, got %s", got)
	}

	// After closing the room, new messages count again.
	alice.send("close_room", map[string]any{"room": room})
	bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "later"})
	unreadCount(t, alice, room, 1)
}

func TestSummaryFailuresDegradeGracefully(t *testing.T) {
	for _, tc := range []struct {
		name string
		fn   Summarizer
		code string
	}{
		{"not configured", nil, "llm_error"},
		{"api error", func(context.Context, string, []string) (string, error) { return "", errors.New("boom") }, "llm_error"},
		{"timeout", func(context.Context, string, []string) (string, error) { return "", context.DeadlineExceeded }, "timeout"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			url := startHub(t, tc.fn)
			alice, bob := dial(t, url), dial(t, url)
			alice.register("Alice")
			bob.register("Bob")
			bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "hi"})
			unreadCount(t, alice, "dm:alice|bob", 1)
			alice.send("summarize_request", map[string]any{"room": "dm:alice|bob"})
			f := alice.expect("summary_error")
			if got := f.str(t, "code"); got != tc.code {
				t.Fatalf("want %s, got %s (%s)", tc.code, got, f.Payload)
			}
			// The connection still works afterwards.
			bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "still alive"})
			alice.expect("dm")
		})
	}
}

func TestSlowSummaryDoesNotBlockChat(t *testing.T) {
	release := make(chan struct{})
	slow := func(ctx context.Context, _ string, _ []string) (string, error) {
		select {
		case <-release:
		case <-ctx.Done():
		}
		return "late", nil
	}
	url := startHub(t, slow)
	alice, bob := dial(t, url), dial(t, url)
	alice.register("Alice")
	bob.register("Bob")
	bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "hi"})
	unreadCount(t, alice, "dm:alice|bob", 1)

	alice.send("summarize_request", map[string]any{"room": "dm:alice|bob"})
	alice.send("summarize_request", map[string]any{"room": "dm:alice|bob"})
	if got := alice.errCode(); got != "rate_limited" {
		t.Fatalf("second request in flight should be rate_limited, got %s", got)
	}
	// While the "LLM" is stuck, chat keeps flowing.
	bob.send("dm", map[string]any{"to": "alice", "kind": "text", "text": "chat continues"})
	alice.expect("dm")
	close(release)
	alice.expect("summary_result")
}
