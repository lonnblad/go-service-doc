package example

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	service_docs "github.com/lonnblad/go-service-doc/example/docs/service"
)

const port = "8080"

func MainExample() {
	mux := mux.NewRouter()

	mux.PathPrefix("/docs/service").Handler(service_docs.Handler())
	mux.PathPrefix("/docs/service/").Handler(service_docs.Handler())

	server := &http.Server{Addr: ":" + port, Handler: mux}

	log.Printf("Will start to listen and serve on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("HTTP server ListenAndServe")
	}
}
