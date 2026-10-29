// NEEDS TEAM REVIEW: Phase 1 worked example (mini chat hub over raw TCP).
// A line sent by one client is forwarded to every OTHER connected client.
package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// clients is shared by every handle() goroutine, so all access goes through mu.
// (Two goroutines touching one map at the same time is a data race.)
var (
	mu      sync.Mutex
	clients = map[net.Conn]bool{}
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("hub listening on :9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	name := conn.RemoteAddr().String()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()
	fmt.Println("joined:", name)

	// On any exit path: forget this client first, then close.
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
		fmt.Println("left:", name)
	}()

	// TCP is a stream, not messages. Scanner splits the stream into lines,
	// so "hello" and "world" sent quickly never merge into "helloworld".
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		broadcast(conn, name+": "+sc.Text()+"\n")
	}
}

// broadcast sends msg to everyone except the sender.
func broadcast(from net.Conn, msg string) {
	mu.Lock()
	defer mu.Unlock()
	for c := range clients {
		if c != from {
			c.Write([]byte(msg))
		}
	}
}
