package pplane

import (
	"net/http"
)

// InstallHandlers installs websocket handler
func InstallHandlers(mux *http.ServeMux) {
	router := newRouter()
	go router.run()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(router, w, r)
	})
}
