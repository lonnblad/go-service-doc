[![Build Status](https://travis-ci.org/lonnblad/go-service-doc.svg?branch=master)](https://travis-ci.org/lonnblad/go-service-doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/lonnblad/go-service-doc)](https://goreportcard.com/report/github.com/lonnblad/go-service-doc)

# go-service-doc
This a tool to generate basic Service Documentation web pages based on standard Markdown files.

It will convert the Markdown files to HTML pages, generate a menu based on `#` and `##` elements and adds CSS similar to the CSS used by github.

Apart from standard Markdown syntax support, it features support for embedding svg files.

It currently has support for generating standard HTML files and a `go` handler.

## Usage

### Install
> go get -u github.com/lonnblad/go-service-doc

### Run
> go-service-doc

#### Flags
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

## Features
- Side Menu Generator
  - The Side Menu is generated based on the Markdown Header Elements: `#` and `##`. It will only generate entries for the headers that have a defined Header ID, like: `{#header_id}`.

- Embedding SVG files
  - All SVG files found in `<src_dir>/static` will be embedded in the generated go-handler and can be referenced through `<base_path>/static/<file_name>`.