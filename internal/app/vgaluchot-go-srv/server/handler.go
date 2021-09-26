package server

import (
	"log"
	"net/http"
)

// Serve the web application forever on the given port.
// Not expected to terminate.
func Serve(addr string, serveStatics bool) {
	installHandlers(http.DefaultServeMux, serveStatics)
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}

func installHandlers(mux *http.ServeMux, serveStatics bool) {
	if serveStatics {
		staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
		mux.Handle("/static/", staticHandler)
	}

	views := []struct {
		url     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{"/", indexView},
		{"/index.html", indexView},
		{"/index", indexView},
		{"/portfolio", portfolioView},
		{"/portfolio.html", portfolioView},
	}
	for _, it := range views {
		handler := makeGetExactURLHandler(it.url, it.handler)
		mux.HandleFunc(it.url, handler)
	}
}

// create a http handler with the following properties
// - if the request url does not march the given url return error 404
// - if the request method is not GET return error 404
// - else call the next handler
func makeGetExactURLHandler(url string, next func(http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != url {
			w.WriteHeader(404)
			error404View(w, r)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
