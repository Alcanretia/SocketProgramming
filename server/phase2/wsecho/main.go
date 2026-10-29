// NEEDS TEAM REVIEW: Phase 2 worked example (WebSocket echo server).
// Same idea as phase1/echo, but the TCP connection first starts life as an
// HTTP request and is then "upgraded" to a WebSocket.
package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed index.html
var indexHTML []byte

// Upgrader performs the HTTP -> WebSocket handshake.
// By default it rejects cross-origin browser requests (CheckOrigin).
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", handleWS)

	log.Println("listening on :8080 (open http://localhost:8080)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// serveIndex serves a static test page. HTTP is only used for static files
// and the upgrade, never for chat messages (R2).
func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

// handleWS runs in its own goroutine per request (net/http does that for us),
// so it is the "table waiter" for exactly one WebSocket client.
func handleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade replies "101 Switching Protocols" and hands us the raw connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}
	defer conn.Close()
	log.Println("ws connected:", conn.RemoteAddr())

	for {
		// ReadMessage returns one whole WebSocket message (frames are
		// reassembled), so unlike raw TCP we do not split lines ourselves.
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws gone:", conn.RemoteAddr(), err)
			return
		}
		if err := conn.WriteMessage(mt, msg); err != nil {
			log.Println("write failed:", err)
			return
		}
	}
}
