// Command server is the chat backend: one HTTP server that serves the React
// build (static files) and the /ws WebSocket endpoint. Chat traffic only ever
// travels over /ws (R2).
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address (all interfaces by default so other computers can connect)")
	static := flag.String("static", "../client-web/dist", "directory with the built React app")
	dev := flag.Bool("dev", false, "allow WebSocket connections from any origin (Vite dev server only)")
	flag.Parse()

	summarizer := NewAnthropicSummarizer()
	if summarizer == nil {
		log.Println("ANTHROPIC_API_KEY not set: AI summary will answer with an error")
	}

	hub := NewHub(summarizer, *dev)
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.ServeWS)
	mux.Handle("/", staticHandler(*static))

	printURLs(*addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// staticHandler serves the built frontend. Unknown paths without a file
// extension fall back to index.html so the single-page app handles them
// (for example the /?invite=<token> links inside QR codes).
func staticHandler(dir string) http.Handler {
	if _, err := os.Stat(filepath.Join(dir, "index.html")); err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "frontend not built: run `npm run build` in client-web (or use the Vite dev server)", http.StatusNotFound)
		})
	}
	files := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && filepath.Ext(r.URL.Path) == "" {
			r.URL.Path = "/"
		}
		files.ServeHTTP(w, r)
	})
}

func printURLs(addr string) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Println("listening on", addr)
		return
	}
	fmt.Println("chat server listening. Open one of:")
	fmt.Printf("  http://localhost:%s\n", port)
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
			continue
		}
		if strings.HasPrefix(ipnet.IP.String(), "169.254.") {
			continue
		}
		fmt.Printf("  http://%s:%s   (other computers on the same network)\n", ipnet.IP, port)
	}
}
