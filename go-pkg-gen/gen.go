package gen

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"github.com/stroem/go-service-doc/core"
)

type Gen struct {
	pages          core.Pages
	staticFiles    core.Files
	indexDocuments []core.IndexDocument
	searchPage     string
	css            string
	basePath       string
}

func New() *Gen {
	return &Gen{}
}

func (g *Gen) WithPages(pages core.Pages) *Gen {
	for _, page := range pages {
		page.HTML = strings.ReplaceAll(page.HTML, "`", "` + \"`\" + `")
		g.pages = append(g.pages, page)
		g.indexDocuments = append(g.indexDocuments, page.IndexDocuments...)
	}
	return g
}

func (g *Gen) WithStaticFiles(files core.Files) *Gen {
	g.staticFiles = files
	return g
}

func (g *Gen) WithSearchPage(page string) *Gen {
	g.searchPage = page
	return g
}

func (g *Gen) WithCSS(css string) *Gen {
	g.css = css
	return g
}

func (g *Gen) WithBasePath(basePath string) *Gen {
	g.basePath = basePath
	return g
}

func (g *Gen) Build() (_ []byte, err error) {
	templateInfo := struct {
		Timestamp      time.Time
		Pages          core.Pages
		StaticFiles    core.Files
		IndexDocuments []core.IndexDocument
		CSS            string
		BasePath       string
		SearchPage     string
	}{
		Timestamp:      time.Now(),
		Pages:          g.pages,
		StaticFiles:    g.staticFiles,
		IndexDocuments: g.indexDocuments,
		CSS:            g.css,
		BasePath:       g.basePath,
		SearchPage:     g.searchPage,
	}

	generator, err := template.New("go_pkg").Parse(packageTemplate)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse package template")
		return
	}

	buffer := &bytes.Buffer{}
	err = generator.Execute(buffer, templateInfo)
	if err != nil {
		err = errors.Wrapf(err, "failed to execute generator")
		return
	}

	return buffer.Bytes(), nil
}

const packageTemplate = `// This file was generated by lonnblad/go-service-doc at
// {{ .Timestamp }}
package docs

import (
	"net/http"
	"strings"

	"github.com/blevesearch/bleve"
)

const contentType = "Content-Type"
const mimeHTML = "text/html"
const mimeCSS = "text/css"

func Handler() http.Handler {
	index, _ := createSearchIndex()

	mux := http.NewServeMux()
	mux.HandleFunc("{{.BasePath}}/markdown.css", cssHandler)
	mux.HandleFunc("{{.BasePath}}/search", searchHandler(index))

{{- range .Pages}}
	mux.HandleFunc("{{.WebPath}}", {{.Name}}PageHandler)
{{- end}}

{{- range .StaticFiles}}
	mux.HandleFunc("{{.Href}}", {{.Name}}StaticFileHandler)
{{- end}}

	return mux
}

func cssHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentType, mimeCSS)

	const content = ` + "`{{.CSS}}`" + `

	// nolint: errcheck
	w.Write([]byte(content))
}

{{- range .Pages}}
func {{.Name}}PageHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentType, mimeHTML)

	const content = ` + "`{{.HTML}}`" + `

	// nolint: errcheck
	w.Write([]byte(content))
}
{{end}}

{{- range .StaticFiles}}
func {{.Name}}StaticFileHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentType, "{{.ContentType}}")

	const content = ` + "`{{.Content}}`" + `

	// nolint: errcheck
	w.Write([]byte(content))
}
{{end}}
func searchHandler(searchIndex bleve.Index) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		queryString := req.URL.Query().Get("q")

		disQuery := bleve.NewDisjunctionQuery()

		for _, q := range strings.Split(queryString, " ") {
			contentFuzzyQuery := bleve.NewFuzzyQuery(q)
			contentFuzzyQuery.FieldVal = "Content"
			contextFuzzyQuery := bleve.NewFuzzyQuery(q)
			contextFuzzyQuery.FieldVal = "Context"

			contentMatchQuery := bleve.NewMatchQuery(q)
			contentMatchQuery.FieldVal = "Content"
			contextMatchQuery := bleve.NewMatchQuery(q)
			contextMatchQuery.FieldVal = "Context"

			disQuery.Disjuncts = append(disQuery.Disjuncts,
				contentFuzzyQuery,
				contextFuzzyQuery,
				contentMatchQuery,
				contextMatchQuery,
			)
		}

		searchRequest := bleve.NewSearchRequest(disQuery)
		searchRequest.Fields = []string{"Context", "HTML", "Link"}

		// nolint: errcheck
		searchResult, _ := searchIndex.Search(searchRequest)

		var result = make([]document, len(searchResult.Hits))
		for idx, hit := range searchResult.Hits {
			cs, ok := hit.Fields["Context"].([]interface{})
			if ok {
				for _, c := range cs {
					result[idx].Context = append(result[idx].Context, c.(string))
				}
			} else {
				result[idx].Context = []string{hit.Fields["Context"].(string)}
			}

			result[idx].HTML = hit.Fields["HTML"].(string)
			result[idx].Link = hit.Fields["Link"].(string)
		}

		w.Header().Set(contentType, mimeHTML)

		// nolint: errcheck
		w.Write(createSearchPage(queryString, result))
	}
}

func createSearchIndex() (searchIndex bleve.Index, err error) {
	indexMapping := bleve.NewIndexMapping()
	if searchIndex, err = bleve.NewMemOnly(indexMapping); err != nil {
 		return
	}

	var doc document
{{- range .IndexDocuments}}
	doc = document{
		Link:    "{{.Link}}",
		Context: []string{ {{- range $index, $element := .Context}} ` + "`{{$element}}`" + `, {{- end}} },
		Content: []string{ {{- range $index, $element := .Content}} ` + "`{{$element}}`" + `, {{- end}} },
		HTML: ` + "`{{.HTML}}`" + `,
	}

	if err = searchIndex.Index("{{.Link}}", doc); err != nil {
		return
	}
{{end}}
	return
}

func createSearchPage(queryString string, searchResult []document) []byte {
	var result = "<div><h1>Search Result for (" + queryString + ")</h1>"

	for _, doc := range searchResult {
		title := strings.Join(doc.Context, " > ")

		result += ` + "`" +
	`<div class=search-result-card onclick="location.href='` + "`" + ` + doc.Link + ` + "`" + `';">` +
	"<h2>` + title + `</h2><div class=search-result-content>` + doc.HTML + `</div></div>`" + `
	}

	result += "</div>"

	page := strings.ReplaceAll(searchPage, "<query_string>", queryString)
	page = strings.ReplaceAll(page, "<search_result>", result)

	return []byte(page)
}

const searchPage = ` + "`" + `{{.SearchPage}}` + "`" + `

type document struct {
	Link    string
	Context []string
	Content []string
	HTML    string
}
`
