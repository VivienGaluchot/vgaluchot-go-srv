package pplane

import "log"

// DistributionMessage is a payload to distribute to all clients except the sender
type DistributionMessage struct {
	src     *Client
	payload []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Payload to broadcast to all clients.
	broadcast chan []byte

	// Payload to distribute to all clients except the sender.
	distribute chan DistributionMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		distribute: make(chan DistributionMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("HUB : %s registered\n", client.conn.RemoteAddr())
			h.clients[client] = true
		case client := <-h.unregister:
			log.Printf("HUB : %s unregistered\n", client.conn.RemoteAddr())
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.distribute:
			for client := range h.clients {
				if message.src != client {
					select {
					case client.send <- message.payload:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
