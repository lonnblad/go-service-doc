package simple

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/stroem/go-service-doc/core"
	html_gen "github.com/stroem/go-service-doc/html-gen"
)

type SimpleExporter struct {
	sourceDir   string
	outputDir   string
	pages       core.Pages
	staticFiles core.Files
	err         error
}

func NewExporter() *SimpleExporter {
	return &SimpleExporter{}
}

func (se *SimpleExporter) WithSourceDir(sourceDir string) *SimpleExporter {
	se.sourceDir = sourceDir
	return se
}

func (se *SimpleExporter) WithOutputDir(outputDir string) *SimpleExporter {
	se.outputDir = outputDir
	return se
}

func (se *SimpleExporter) WithPages(pages core.Pages) *SimpleExporter {
	se.pages = pages
	return se
}

func (se *SimpleExporter) WithStaticFiles(staticFiles core.Files) *SimpleExporter {
	se.staticFiles = staticFiles
	return se
}

func (se *SimpleExporter) Error() error {
	return se.err
}

func (se *SimpleExporter) Run() {
	zap.L().Info("exporting simple files")

	if err := os.MkdirAll(se.outputDir+"/static", os.ModePerm); err != nil {
		se.err = errors.Wrap(err, "os.MkdirAll failed")
		return
	}

	if err := exportHTMLPages(se.pages, se.sourceDir, se.outputDir); se.err != nil {
		se.err = errors.Wrap(err, "exportHTMLPages failed")
		return
	}

	if err := exportCSSFile(se.outputDir); se.err != nil {
		se.err = errors.Wrap(err, "exportCSSFile failed")
		return
	}

	if err := exportStaticFiles(se.staticFiles, se.sourceDir, se.outputDir); se.err != nil {
		se.err = errors.Wrap(err, "exportStaticFiles failed")
		return
	}
}

func exportHTMLPages(pages core.Pages, sourceDir, outputDir string) error {
	for _, page := range pages {
		zap.L().With(zap.String("page", page.Name)).Info("exporting HTML file")

		filepath := strings.ReplaceAll(page.Filepath, ".md", ".html")
		filepath = strings.ReplaceAll(filepath, sourceDir, outputDir)

		if err := ioutil.WriteFile(filepath, []byte(page.HTML), 0644); err != nil {
			return errors.Wrap(err, "ioutil.WriteFile failed")
		}
	}

	return nil
}

func exportCSSFile(outputDir string) error {
	css := html_gen.GetMarkdownCSS()
	filepath := outputDir + "/markdown.css"

	zap.L().With(zap.String("file", "markdown.css")).Info("exporting CSS file")

	if err := ioutil.WriteFile(filepath, css, 0644); err != nil {
		return errors.Wrap(err, "ioutil.WriteFile failed")
	}

	return nil
}

func exportStaticFiles(staticFiles core.Files, sourceDir, outputDir string) error {
	for _, file := range staticFiles {
		filepath := strings.ReplaceAll(file.Path, sourceDir, outputDir)

		zap.L().With(zap.String("file", file.Name)).Info("exporting static file")

		if err := ioutil.WriteFile(filepath, []byte(file.Content), 0644); err != nil {
			return errors.Wrap(err, "ioutil.WriteFile failed")
		}
	}

	return nil
}
