// Phase 1 (P2 — งานย่อย 1.1 & 1.2): TCP Echo Server ที่รองรับหลาย Client พร้อมกัน
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9001", "TCP listen address")
	flag.Parse()

	// 1. Listen: เปิดเบอร์โทรศัพท์หน้าร้านรอให้ลูกค้าติดต่อเข้ามา
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("P2 TCP echo server listening on %s", ln.Addr())
	serveEcho(ln)
}

func serveEcho(ln net.Listener) {
	for {
		// 2. Accept: บล็อก (หยุดรอ) จนกว่าจะมีลูกค้าเชื่อมต่อเข้ามา
		// เมื่อเชื่อมสำเร็จ จะคืนค่า conn (ท่อสื่อสารเฉพาะของลูกค้าคนนั้น)
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		// งานย่อย 1.2: เติมคำว่า `go` เพื่อเสกพนักงานประจำโต๊ะ (Goroutine)
		// ไปดูแลลูกค้าคนนี้แยกต่างหาก พนักงานต้อนรับจะได้กลับไปรอ Accept คนถัดไปทันที
		go handleEchoClient(conn)
	}
}

func handleEchoClient(conn net.Conn) {
	// 5. Close: ปิดท่อเสมอเมื่อจบฟังก์ชัน (เมื่อลูกค้าวางสาย)
	defer conn.Close()

	// เนื่องจาก TCP เป็น byte stream (สายน้ำข้อมูล) เราจึงใช้การอ่านทีละบรรทัดจบด้วย '\n'
	reader := bufio.NewReader(conn)
	for {
		// 3. Read: รออ่านข้อความ 1 บรรทัดจากลูกค้า
		// ถ้าลูกค้าตัดการเชื่อมต่อ จะได้ err (เช่น io.EOF) แล้วหยุดลูปทันที
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			// 4. Write: ส่งข้อความเดิมสะท้อนกลับไปหาลูกค้าคนเดิม
			if _, werr := conn.Write(line); werr != nil {
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
			}
			return
		}
	}
}
