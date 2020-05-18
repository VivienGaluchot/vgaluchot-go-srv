package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var (
	indexTmpl = template.Must(
		template.ParseFiles(filepath.Join(os.Getenv("WEB_DIR"), "template", "index.html")),
	)
)

// Serve the web application forever on the given port.
// Not expected to terminate.
func Serve(port string, serveStatics bool) {
	http.HandleFunc("/", indexHandler)

	if serveStatics {
		// Serve static files out of the static directory.
		// By configuring a static handler in app.yaml, App Engine serves all the
		// static content itself. As a result, the following two lines are in
		// effect for development only.
		static := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
		http.Handle("/static/", static)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	type indexData struct {
		Style       string
		RequestTime string
	}
	data := indexData{
		Style:       "/static/style.css",
		RequestTime: time.Now().Format(time.RFC822),
	}
	if err := indexTmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
