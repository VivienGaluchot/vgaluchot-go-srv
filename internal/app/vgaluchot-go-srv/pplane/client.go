package pplane

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// MsgState is typing the message state values, the state is seen by the client
type MsgState int

const (
	// MsgStateLocal : message not sent to server yet
	MsgStateLocal MsgState = 0
	// MsgStateSentToSrv : message has been sent to server
	MsgStateSentToSrv MsgState = 1
	// MsgStateSentToPair : message has been sent to pair
	MsgStateSentToPair MsgState = 2
	// MsgStateReadByPair : message has been read by pair
	MsgStateReadByPair MsgState = 3
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func handleIncommingMessage(c *Client, raw []byte) {
	log.Printf("CLIENT : %s -> '%s'\n", c.conn.RemoteAddr(), raw)

	// decode message
	type IncommingMessage struct {
		UID       string `json:"uid"`
		Data      string `json:"data"`
		Timestamp string `json:"timestamp"`
		Counter   int64  `json:"counter"`
	}
	var message IncommingMessage
	if err := json.Unmarshal(raw, &message); err != nil {
		log.Printf("unmarshal error with input '%s'\n", raw)
		return
	}

	// send server ack
	type StateUpdateMessage struct {
		UID     string   `json:"uid"`
		Counter int64    `json:"counter"`
		State   MsgState `json:"state"`
	}
	ackMessage := StateUpdateMessage{UID: message.UID, Counter: message.Counter, State: MsgStateSentToSrv}
	var ackBytes, err = json.Marshal(ackMessage)
	if err != nil {
		log.Fatal(err)
	}
	c.send <- ackBytes

	// broadcast
	if c.hub != nil {
		c.hub.broadcast <- raw
	}
}

// readPump pumps messages from the websocket connection to the handleIncommingMessage.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	log.Printf("CLIENT : %s connected\n", c.conn.RemoteAddr())
	defer func() {
		if c.hub != nil {
			c.hub.unregister <- c
		}
		c.conn.Close()
		log.Printf("CLIENT : %s disconnected\n", c.conn.RemoteAddr())
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		handleIncommingMessage(c, rawMessage)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			log.Printf("CLIENT : %s <- '%s'\n", c.conn.RemoteAddr(), message)
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: nil, conn: conn, send: make(chan []byte, 256)}
	// client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
