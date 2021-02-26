package main

import (
	"log"
	"net/http"

	service_docs "github.com/lonnblad/go-service-doc/cmd/example/docs/generated"
)

const port = "8080"

func main() {
	server := &http.Server{Addr: ":" + port, Handler: service_docs.Handler()}

	log.Printf("Will start to listen and serve on port %s", port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("HTTP server ListenAndServe")
	}
}
