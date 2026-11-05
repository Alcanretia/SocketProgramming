package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

// TestEchoConcurrent ทดสอบว่าเมื่อ Client 1 เชื่อมต่อเข้ามาแล้วเงียบ
// Client 2 ยังสามารถเชื่อมต่อและได้รับข้อความ Echo กลับมาได้ทันที (เพราะมี Goroutine แยก)
func TestEchoConcurrent(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go serveEcho(ln)

	// ลูกค้าคนที่ 1 เชื่อมต่อเข้ามาแล้วเงียบ ไม่ส่งอะไรเลย
	c1, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Close()

	// ลูกค้าคนที่ 2 เชื่อมต่อเข้ามาทีหลัง และส่งข้อความ
	c2, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()
	_ = c2.SetDeadline(time.Now().Add(2 * time.Second))

	msg := "hello from client 2\n"
	if _, err := c2.Write([]byte(msg)); err != nil {
		t.Fatal(err)
	}
	got, err := bufio.NewReader(c2).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if got != msg {
		t.Fatalf("got %q, want %q", got, msg)
	}
}
