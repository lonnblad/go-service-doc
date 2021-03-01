[![Release](https://img.shields.io/github/v/release/lonnblad/go-service-doc)](https://github.com/lonnblad/go-service-doc/releases/latest)
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Flonnblad%2Fgo-service-doc%2Fbadge%3Fref%3Dmain&style=flat)](https://actions-badge.atrox.dev/lonnblad/go-service-doc/goto?ref=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/lonnblad/go-service-doc)](https://goreportcard.com/report/github.com/lonnblad/go-service-doc)
[![Coverage Status](https://coveralls.io/repos/github/lonnblad/go-service-doc/badge.svg?branch=main)](https://coveralls.io/github/lonnblad/go-service-doc?branch=main)

# go-service-doc

This a tool to generate basic Service Documentation web pages based on standard Markdown files.

It will convert the Markdown files to HTML pages, generate a menu based on `#` and `##` elements and add CSS similar to the CSS used by github to display Markdown files. HTML pages will have the kebab-case version of the filename as the URL.

Apart from standard Markdown syntax support, it features support for embedding svg files.

It currently has support for generating standard HTML files and a `go` handler.

Here you can find a [deployed example](https://lonnblad.github.io/go-service-doc).

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

> go-service-doc -s bars.md -d docs/src -o docs/generated -p /go-service-doc

Example code:

```go
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
```

## Features

- HTML Page Generator

  It will convert the Markdown files to HTML pages and add CSS similar to the CSS used by github to display Markdown files. The URL for the generated HTML page will be the kebab-case version of the filename excluding the extension, i.e. `monkey_bar.md` will be `/<base_path>/monkey-bar`.

- Side Menu Generator

  The Side Menu is generated based on the Markdown Header Elements: `#` and `##`. It will only generate entries for the headers that have a defined Header ID, like: `{#header_id}`.

- Search Engine

  The Side Menu features a Search field that can be used to search in all generated pages. The search engine will index content based on Markdown Headers.

- Embedding Images

  All SVG, PNG and ICO files found in `<src_dir>/static` will be embedded in the generated go-handler and can be referenced through `<base_path>/static/<file_name>`.

- Favicon

  If a file called `favicon.ico` is found in `<src_dir>/static`, it will be used as the sites favicon.
