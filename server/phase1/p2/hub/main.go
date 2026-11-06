// Phase 1 (P2 — งานย่อย 1.3): Mini Chat Hub (กระจายข้อความจาก A ไปหาทุกคนยกเว้น A)
package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"sync"
)

type miniHub struct {
	mu      sync.Mutex        // กุญแจล็อกสมุดรายชื่อ ป้องกันการอ่าน-เขียนพร้อมกัน (Race Condition)
	clients map[net.Conn]bool // สมุดรายชื่อส่วนกลาง เก็บท่อเชื่อมต่อของทุกคนที่ออนไลน์
}

func newMiniHub() *miniHub {
	return &miniHub{clients: make(map[net.Conn]bool)}
}

func (h *miniHub) add(c net.Conn) {
	h.mu.Lock() // ล็อกเฉพาะจังหวะที่จดชื่อลงสมุด
	h.clients[c] = true
	h.mu.Unlock() // ปลดล็อกทันทีเมื่อจดเสร็จ
}

func (h *miniHub) remove(c net.Conn) {
	h.mu.Lock() // ล็อกเฉพาะจังหวะที่ขีดฆ่าชื่อออกจากสมุดเมื่อวางสาย
	delete(h.clients, c)
	h.mu.Unlock()
	_ = c.Close()
}

// forward ส่งต่อข้อความจาก sender ไปให้ทุกคนในสมุด ยกเว้นตัว sender เอง
func (h *miniHub) forward(sender net.Conn, line []byte) {
	// เทคนิคลดการหน่วงเวลา: ล็อกเฉพาะตอนก๊อปปี้รายชื่อผู้รับออกมาใส่รายการชั่วคราว (targets)
	// แล้วปลดล็อกทันที ก่อนจะสั่ง Write ลงท่อจริง
	h.mu.Lock()
	targets := make([]net.Conn, 0, len(h.clients))
	for c := range h.clients {
		if c != sender {
			targets = append(targets, c)
		}
	}
	h.mu.Unlock()

	for _, c := range targets {
		_, _ = c.Write(line)
	}
}

func (h *miniHub) serve(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		h.add(conn)
		go h.handleClient(conn)
	}
}

func (h *miniHub) handleClient(conn net.Conn) {
	defer h.remove(conn) // เมื่อหลุดลูป (ลูกค้าวางสาย) ให้ลบชื่อออกจากสมุดและปิดท่อ

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			h.forward(conn, line)
		}
		if err != nil {
			return
		}
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1:9002", "TCP hub listen address")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("P2 mini chat hub listening on %s", ln.Addr())
	newMiniHub().serve(ln)
}
