package main

import (
	"log"

	"github.com/VivienGaluchot/vgaluchot-go-srv/internal/app/vgaluchot-go-srv/server"
)

func main() {
	log.Println("Startup in DEV configuration")

	port := "8000"
	serveStatic := true

	server.Serve(port, serveStatic)
	log.Fatalln("server.Serve terminated")
}
