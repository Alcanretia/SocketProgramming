// Phase 1 (P2): TCP Client สำหรับทดสอบพิมพ์คุยผ่าน Terminal
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9001", "server TCP address")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Printf("connected to %s (type lines and press Enter)", *addr)

	// Goroutine แยกสำหรับรอรับข้อความจาก Server มาแสดงบนหน้าจอตลอดเวลา
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
	}()

	// อ่านสิ่งที่ผู้ใช้พิมพ์จากคีย์บอร์ดส่งขึ้นไปหา Server
	_, _ = io.Copy(conn, os.Stdin)
}
