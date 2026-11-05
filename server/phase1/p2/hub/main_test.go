package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

// TestMiniHubForwarding ทดสอบว่าเมื่อ A ส่งข้อความ B และ C จะได้รับ แต่ A จะไม่ได้รับเสียงสะท้อนของตัวเอง
func TestMiniHubForwarding(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	h := newMiniHub()
	go h.serve(ln)

	dial := func() net.Conn {
		c, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		_ = c.SetDeadline(time.Now().Add(2 * time.Second))
		return c
	}

	cA := dial()
	defer cA.Close()
	cB := dial()
	defer cB.Close()
	cC := dial()
	defer cC.Close()

	// รอเสี้ยววินาทีให้ทั้ง 3 คนถูกจดลงสมุดรายชื่อเรียบร้อย
	time.Sleep(30 * time.Millisecond)

	msg := "hello from A\n"
	if _, err := cA.Write([]byte(msg)); err != nil {
		t.Fatal(err)
	}

	for name, conn := range map[string]net.Conn{"B": cB, "C": cC} {
		got, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Fatalf("client %s read error: %v", name, err)
		}
		if got != msg {
			t.Fatalf("client %s got %q, want %q", name, got, msg)
		}
	}
}
