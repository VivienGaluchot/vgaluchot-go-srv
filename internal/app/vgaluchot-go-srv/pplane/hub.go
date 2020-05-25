package pplane

import "log"

// AdmincastMessage is a message to send to admin
type AdmincastMessage struct {
	src     *Client
	payload []byte
}

// Router direct messages from client to client.
type Router struct {
	// Registered admins.
	admins map[*Client]bool
	// Register an admin client.
	registerAdmin chan *Client
	// Unregister a client.
	unregister chan *Client
	// Payload to send to the admin client.
	admincast chan AdmincastMessage
}

func newRouter() *Router {
	return &Router{
		admins:        make(map[*Client]bool),
		registerAdmin: make(chan *Client),
		unregister:    make(chan *Client),
		admincast:     make(chan AdmincastMessage),
	}
}

func (h *Router) run() {
	for {
		select {
		case client := <-h.registerAdmin:
			log.Printf("Router : admin %s registered\n", client.conn.RemoteAddr())
			h.admins[client] = true
		case client := <-h.unregister:
			log.Printf("Router : admin %s unregistered\n", client.conn.RemoteAddr())
			if _, ok := h.admins[client]; ok {
				delete(h.admins, client)
				close(client.send)
			}
		case message := <-h.admincast:
			log.Printf("Router : admincast %s\n", message.payload)
			for client := range h.admins {
				if client != message.src {
					select {
					case client.send <- message.payload:
					default:
						close(client.send)
						delete(h.admins, client)
					}
				}
			}
		}
	}
}
