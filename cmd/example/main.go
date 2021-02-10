package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	service_docs "github.com/lonnblad/go-service-doc/cmd/example/docs/generated"
)

const port = "8080"

func main() {
	mux := mux.NewRouter()

	mux.PathPrefix("/bars").Handler(service_docs.Handler())

	server := &http.Server{Addr: ":" + port, Handler: mux}

	log.Printf("Will start to listen and serve on port %s", port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("HTTP server ListenAndServe")
	}
}
