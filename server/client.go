package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 25 * time.Second
	sendBuffer = 64
)

// Client is one WebSocket connection.
//
// Concurrency rules (so the hub never freezes and gorilla's "one writer, one
// reader" rule holds):
//   - readPump is the only reader, writePump is the only writer.
//   - everyone else talks to the client through enqueue(), which never blocks.
//   - c.send is never closed (closing a channel others may still send on
//     panics); c.done signals shutdown instead.
//   - the fields below "hub-owned" are touched only by the hub goroutine.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	done chan struct{}
	once sync.Once

	// hub-owned
	name        string // display spelling
	key         string // lowercased, used for identity
	registered  bool
	openRooms   map[string]bool
	summarizing bool
}

func newClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:       h,
		conn:      conn,
		send:      make(chan []byte, sendBuffer),
		done:      make(chan struct{}),
		openRooms: map[string]bool{},
	}
}

// enqueue queues a frame without ever blocking. It reports false if the
// client's buffer is full (a slow client), so the caller can drop it instead
// of letting one slow reader freeze everyone.
func (c *Client) enqueue(b []byte) bool {
	select {
	case <-c.done:
		return false
	default:
	}
	select {
	case c.send <- b:
		return true
	default:
		return false
	}
}

func (c *Client) kill() {
	c.once.Do(func() {
		close(c.done)
		c.conn.Close()
	})
}

// readPump runs in the HTTP handler goroutine of the upgrade request.
func (c *Client) readPump() {
	defer c.hub.post(leave{c})

	c.conn.SetReadLimit(maxFrameBytes)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		mt, data, err := c.conn.ReadMessage()
		if err != nil {
			return // closed, timed out, or frame too large
		}
		c.conn.SetReadDeadline(time.Now().Add(pongWait))

		if mt != websocket.TextMessage {
			c.enqueue(encode("error", errorPayload{codeBadRequest, "text frames only"}, ""))
			continue
		}
		var env Envelope
		if err := json.Unmarshal(data, &env); err != nil || env.Type == "" {
			c.enqueue(encode("error", errorPayload{codeBadRequest, "malformed JSON envelope"}, ""))
			continue
		}
		c.hub.post(inbound{c, env})
	}
}

// writePump is the only goroutine that writes to the connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer c.conn.Close()

	for {
		select {
		case b := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				c.kill()
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.kill()
				return
			}
		case <-c.done:
			return
		}
	}
}
