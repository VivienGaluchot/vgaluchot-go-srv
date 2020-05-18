package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/websocket"
)

func templateDir() string {
	return filepath.Join(os.Getenv("WEB_DIR"), "template")
}

var (
	templates = template.Must(
		template.ParseFiles(
			filepath.Join(templateDir(), "index.html"),
			filepath.Join(templateDir(), "chat.html")),
	)
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
	mux.HandleFunc("/", makeTemplateHandler("/", "index.html"))
	mux.HandleFunc("/chat", makeTemplateHandler("/chat", "chat.html"))

	if serveStatics {
		staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
		mux.Handle("/static/", staticHandler)
	}

	websocket.InstallHandlers(mux)
}

func makeTemplateHandler(url string, templateName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != url {
			http.NotFound(w, r)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type indexData struct {
			Static      string
			RequestTime string
		}
		data := indexData{
			Static:      "/static",
			RequestTime: time.Now().Format(time.RFC822),
		}
		if err := templates.ExecuteTemplate(w, templateName, data); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
