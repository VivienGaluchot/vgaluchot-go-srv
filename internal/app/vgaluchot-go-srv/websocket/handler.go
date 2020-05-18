package websocket

import (
	"net/http"
)

// InstallHandlers installs websocket handler
func InstallHandlers(mux *http.ServeMux) {
	hub := newHub()
	go hub.run()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}
