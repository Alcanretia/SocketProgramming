// Phase 1 helper: tiny TCP client.
// A goroutine prints whatever the server sends; the main loop sends what you type.
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Receiver: runs alongside the typing loop, so messages forwarded
	// by the hub show up even while you are not typing.
	go func() {
		sc := bufio.NewScanner(conn)
		for sc.Scan() {
			fmt.Println(sc.Text())
		}
		fmt.Println("server closed")
		os.Exit(0)
	}()

	// Sender: each line you type goes to the server.
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		conn.Write([]byte(stdin.Text() + "\n"))
	}
}
