package main

import (
	"log"

	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/conf"
	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/server"
)

func main() {
	log.Printf("Startup in DEV configuration on version '%s'\n", conf.Version)

	port := "8000"
	serveStatic := true

	server.Serve(port, serveStatic)
	log.Fatalln("server.Serve terminated")
}
