// NEEDS TEAM REVIEW: Phase 1 worked example (TCP echo server, many clients at once).
package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen: open the shop. The OS now reserves port 9000 for us.
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("listening on :9000")

	// The door waiter: loop forever, only ever accepting.
	for {
		// Accept blocks until a client connects, then returns that client's own conn.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		// `go` hires a new table waiter (goroutine) for this conn,
		// and the door waiter immediately loops back to Accept.
		go handle(conn)
	}
}

// handle is the table waiter: it serves exactly one client.
func handle(conn net.Conn) {
	defer conn.Close() // runs when handle returns, i.e. when this client leaves
	fmt.Println("client connected:", conn.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		// Read blocks until this client sends something. That is fine now:
		// it only blocks this goroutine, not the door waiter or other tables.
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client gone:", conn.RemoteAddr(), err)
			return
		}
		conn.Write(buf[:n])
	}
}
