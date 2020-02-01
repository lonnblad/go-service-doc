package golang

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/lonnblad/go-service-doc/core"
	go_gen "github.com/lonnblad/go-service-doc/go-pkg-gen"
	html_gen "github.com/lonnblad/go-service-doc/html-gen"
)

type GoExporter struct {
	pages       core.Pages
	staticFiles core.Files
	basepath    string
	searchPage  string
	outputDir   string
	err         error
}

func NewExporter() *GoExporter {
	return &GoExporter{}
}

func (goex *GoExporter) WithPages(pages core.Pages) *GoExporter {
	goex.pages = pages
	return goex
}

func (goex *GoExporter) WithStaticFiles(staticFiles core.Files) *GoExporter {
	goex.staticFiles = staticFiles
	return goex
}

func (goex *GoExporter) WithBasepath(basepath string) *GoExporter {
	goex.basepath = basepath
	return goex
}

func (goex *GoExporter) WithSearchPage(searchPage string) *GoExporter {
	goex.searchPage = searchPage
	return goex
}

func (goex *GoExporter) WithOutputDir(outputDir string) *GoExporter {
	goex.outputDir = outputDir
	return goex
}

func (goex *GoExporter) Error() error {
	return goex.err
}

func (goex *GoExporter) Run() {
	zap.L().Info("building go pkg")
	css := html_gen.GetMarkdownCSS()

	fileContent, err := go_gen.New().
		WithPages(goex.pages).
		WithStaticFiles(goex.staticFiles).
		WithCSS(string(css)).
		WithBasePath(goex.basepath).
		WithSearchPage(goex.searchPage).
		Build()

	if err != nil {
		goex.err = errors.Wrap(err, "go_gen.Build failed")
		return
	}

	filepath := goex.outputDir + "/" + "docs.go"

	zap.L().Info("exporting go pkg")
	if err := os.MkdirAll(goex.outputDir, os.ModePerm); err != nil {
		goex.err = errors.Wrap(err, "os.MkdirAll failed")
		return
	}

	if err := ioutil.WriteFile(filepath, fileContent, 0644); err != nil {
		goex.err = errors.Wrap(err, "ioutil.WriteFile failed")
		return
	}
}
