package main

import (
	"log"
	"os"

	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/conf"
	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/server"
)

func main() {
	log.Printf("Startup in GAE configuration on version '%s'\n", conf.Version)

	// By configuring a static handler in app.yaml, App Engine serves all the
	// static content itself.
	serveStatic := false

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	server.Serve(":"+port, serveStatic)
	log.Fatalln("server.Serve terminated")
}
