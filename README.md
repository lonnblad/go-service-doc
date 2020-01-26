[![Build Status](https://travis-ci.org/lonnblad/go-service-doc.svg?branch=master)](https://travis-ci.org/lonnblad/go-service-doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/lonnblad/go-service-doc)](https://goreportcard.com/report/github.com/lonnblad/go-service-doc)

# go-service-doc
This a tool to generate Service Documentation based on standard Markdown files for a `go` application.

It will convert the Markdown files to HTML pages, generate a menu based on `#` and `##` elements and adds CSS similar to the CSS used by github.

## Usage

### Install
> go get -u github.com/lonnblad/go-service-doc/cmd/go-service-doc

### Run

> go-service-doc

### Flags
- **-s**

    > The filename of the Markdown file to use for the base path, defaults to `service.md`.

- **-d**

    > The Directory where the markdown files are located, defaults to `docs`.

- **-o**

    > The Directory where to write the generated files, defaults to `docs`.

- **-p**

    > Base path to add for the generated documentation, defaults to `/docs`.

### Example
You can find this example with the markdown source files and the generated output in [cmd/example](cmd/example).

To generate the output, the following is executed from [cmd/example](cmd/example).

> go-service-doc -d docs/src -o docs/generated -p /docs/service

Example code:
```go
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

	mux.PathPrefix("/docs/service").Handler(service_docs.Handler())

	server := &http.Server{Addr: ":" + port, Handler: mux}

	log.Printf("Will start to listen and serve on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("HTTP server ListenAndServe")
	}
}
```