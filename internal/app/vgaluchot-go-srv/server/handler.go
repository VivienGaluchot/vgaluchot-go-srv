package server

import (
	"log"
	"net/http"
)

// Serve the web application forever on the given port.
// Not expected to terminate.
func Serve(port string, serveStatics bool) {
	installHandlers(http.DefaultServeMux, serveStatics)
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}

func installHandlers(mux *http.ServeMux, serveStatics bool) {
	views := []struct {
		url     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{"/", indexView},
		{"/index", indexView},
		{"/index.html", indexView},
		{"/contact", contactView},
		{"/contact.html", contactView},
	}
	for _, it := range views {
		handler := makeGetExactURLHandler(it.url, it.handler)
		mux.HandleFunc(it.url, handler)
	}

	if serveStatics {
		staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
		mux.Handle("/static/", staticHandler)
	}
}

// create a http handler with the following properties
// - if the request url does not march the given url return error 404
// - if the request method is not GET return error 404
// - else call the next handler
func makeGetExactURLHandler(url string, next func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != url {
			http.NotFound(w, r)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
