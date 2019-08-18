package gen

import (
	"bytes"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"github.com/lonnblad/go-service-doc/core"
)

type Gen struct {
	api   string
	pages core.Pages
	doc   string
}

func New() *Gen {
	return &Gen{}
}
func (g *Gen) WithAPITitle(apiTitle string) *Gen {
	g.api = apiTitle
	return g
}

func (g *Gen) WithPages(pages core.Pages) *Gen {
	g.pages = pages
	return g
}

func (g *Gen) WithDocument(doc string) *Gen {
	g.doc = doc
	return g
}

func (g *Gen) Build() (_ []byte, err error) {
	templateInfo := struct {
		Timestamp time.Time
		API       string
		Pages     core.Pages
		Doc       string
	}{
		Timestamp: time.Now(),
		API:       g.api,
		Pages:     g.pages,
		Doc:       g.doc,
	}

	generator, err := template.New("html_page").Parse(htmlPageTemplate)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse HTML page template")
		return
	}

	buffer := &bytes.Buffer{}
	if err = generator.Execute(buffer, templateInfo); err != nil {
		err = errors.Wrapf(err, "failed to execute generator")
		return
	}

	return buffer.Bytes(), nil
}

func GetMarkdownCSS() []byte {
	return []byte(markdownCSS)
}

const htmlPageTemplate = `<!DOCTYPE html>
<html lang=en>
<head>
  <title>{{.API}}</title>
  <meta name='generator' content='github.com/lonnblad/go-service-doc'>
  <link rel="stylesheet" href="/docs/service/markdown.css">
</head>
<body class="markdown-body">
  <div class="flex-container">
	<div class="menu-container">
      <ul>
      {{range .Pages}}
        {{range .Headers}}
        <li><a href="{{.Link}}">{{.Title}}</a>
        {{if not .Headers}}</li>{{else}}<ul>{{end}}
            {{range .Headers}}
            <li><a href="{{.Link}}">{{.Title}}</a></li>
            {{end}}
        {{if .Headers}}
          </ul>
        </li>  
        {{end}}
        {{end}}
      {{end}}
      </ul>
    </div>
    <div class="doc-container">
      {{.Doc}}
    </div>
  </div>
</body>
</html>`
