package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lonnblad/go-service-doc/core"
	go_gen "github.com/lonnblad/go-service-doc/go-pkg-gen"
	html_gen "github.com/lonnblad/go-service-doc/html-gen"
)

func init() {
	encoderConf := zap.NewProductionEncoderConfig()

	encoderConf.TimeKey = "timestamp"
	encoderConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConf),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevel(),
		),
	)
	zap.ReplaceGlobals(logger)
}

func main() {
	serviceFilename := flag.String("s", "service.md", "Main Markdown file for the service.")
	dir := flag.String("d", "docs", "Directory where to get markdown files.")
	output := flag.String("o", "docs", "Directory where to write output.")
	basePath := flag.String("p", "/docs", "Base path for the generated documentation.")

	flag.Parse()

	if err := os.MkdirAll(*output, os.ModePerm); err != nil {
		zap.L().
			With(zap.Error(err), zap.String("output", *output)).
			Error("Failed to create output directory")
	}

	renderer := serviceDocRenderer{
		dir:             *dir,
		output:          *output,
		basePath:        *basePath,
		serviceFilename: *serviceFilename,
		uniqueLinks:     make(map[string]bool),
	}

	renderer.findMDFiles(*dir, *basePath)
	renderer.parseMarkdown()
	renderer.buildAndExportHTMLPages()

	renderer.findSVGFiles(*dir, *basePath)

	renderer.buildAndExportGoPkg()
	zap.L().Info("done")
}

func convertKebabCaseToCamel(str string) string {
	parts := strings.Split(str, "-")
	for idx, part := range parts {
		if idx == 0 {
			continue
		}
		parts[idx] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

type serviceDocRenderer struct {
	dir             string
	output          string
	basePath        string
	serviceFilename string
	serviceName     string
	serviceTitle    string
	uniqueLinks     map[string]bool
	pages           core.Pages
	staticFiles     core.Files
	css             string
}

func (sdr *serviceDocRenderer) findMDFiles(dir, basePath string) {
	zap.L().Info("search for MD files")
	files, err := ioutil.ReadDir(dir + "/")
	if err != nil {
		log.Print(err)
		return
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		page := core.Page{}
		page.WebPath = strings.ReplaceAll(f.Name(), ".md", "")
		page.Name = convertKebabCaseToCamel(page.WebPath)
		if f.Name() != sdr.serviceFilename {
			page.WebPath = basePath + "/" + page.WebPath
		} else {
			page.WebPath = basePath
			sdr.serviceName = page.Name
		}

		page.Filepath = dir + "/" + f.Name()
		sdr.pages = append(sdr.pages, page)
	}
}

func (sdr *serviceDocRenderer) findSVGFiles(dir, basePath string) {
	zap.L().Info("search for SVG files")
	files, err := ioutil.ReadDir(dir + "/static")
	if err != nil {
		log.Print(err)
		return
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".svg") {
			continue
		}

		file := core.File{ContentType: "image/svg+xml"}
		file.Path = basePath + "/static/" + f.Name()

		file.Name = strings.ReplaceAll(f.Name(), ".svg", "")
		file.Name = convertKebabCaseToCamel(file.Name)

		filepath := dir + "/static/" + f.Name()
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			zap.L().
				With(zap.Error(err), zap.String("file", filepath)).
				Error("Failed to read file")
		}

		file.Content = string(content)
		sdr.staticFiles = append(sdr.staticFiles, file)
	}
}

func (sdr *serviceDocRenderer) parseMarkdown() {
	for idx, page := range sdr.pages {
		zap.L().With(zap.String("page", page.Name)).Info("parsing markdown")

		content, err := ioutil.ReadFile(page.Filepath)
		if err != nil {
			zap.L().
				With(zap.Error(err), zap.String("file", page.Filepath)).
				Error("Failed to read file")
		}

		exts := blackfriday.NoIntraEmphasis | blackfriday.HardLineBreak | blackfriday.HeadingIDs | blackfriday.FencedCode
		markdown := blackfriday.Run(content)
		page.Markdown = string(markdown)

		node := blackfriday.New(blackfriday.WithExtensions(exts)).Parse(content)
		node.Walk(sdr.walker(&page))
		sdr.pages[idx] = page
	}
}

func (sdr *serviceDocRenderer) walker(page *core.Page) blackfriday.NodeVisitor {
	return func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Level > 2 ||
			node.Type != blackfriday.Heading ||
			!entering ||
			node.HeadingID == "" ||
			string(node.FirstChild.Literal) == "" {
			return blackfriday.GoToNext
		}

		link := fmt.Sprintf("%s#%s", page.WebPath, node.HeadingID)
		if exists := sdr.uniqueLinks[link]; exists {
			zap.L().
				With(zap.String("link", link)).
				Fatal("Link already exists")
		}

		sdr.uniqueLinks[link] = true
		h := core.Header{
			Title: string(node.FirstChild.Literal),
			Link:  link,
		}

		if node.Level == 1 {
			if sdr.serviceName == page.Name && sdr.serviceTitle == "" {
				sdr.serviceTitle = h.Title
			}
			page.Headers = append(page.Headers, h)
		} else {
			idx := len(page.Headers) - 1
			page.Headers[idx].Headers = append(page.Headers[idx].Headers, h)
		}

		return blackfriday.GoToNext
	}
}

func (sdr *serviceDocRenderer) buildAndExportHTMLPages() {
	sdr.pages.SortByName(sdr.serviceName)

	for idx, page := range sdr.pages {
		zap.L().With(zap.String("page", page.Name)).Info("building HTML page")

		bs, err := html_gen.New().
			WithAPITitle(sdr.serviceTitle).
			WithPages(sdr.pages).
			WithDocument(page.Markdown).
			Build()
		if err != nil {
			zap.L().
				With(zap.Error(err)).
				Fatal("Failed to generate HTML page")
		}

		zap.L().With(zap.String("page", page.Name)).Info("exporting HTML file")

		filepath := strings.ReplaceAll(page.Filepath, ".md", ".html")
		filepath = strings.ReplaceAll(filepath, sdr.dir, sdr.output)

		if err := ioutil.WriteFile(filepath, bs, 0644); err != nil {
			zap.L().
				With(zap.Error(err), zap.String("filepath", filepath)).
				Fatal("Failed to write file")
		}

		// page.HTML = strings.ReplaceAll(string(bs), "\n", "")
		page.HTML = string(bs)
		sdr.pages[idx] = page
	}

	css := html_gen.GetMarkdownCSS()
	filepath := sdr.output + "/markdown.css"

	if err := ioutil.WriteFile(filepath, css, 0644); err != nil {
		zap.L().
			With(zap.Error(err), zap.String("filepath", filepath)).
			Fatal("Failed to write file")
	}

	sdr.css = strings.ReplaceAll(string(css), "\n", "")
}

func (sdr *serviceDocRenderer) buildAndExportGoPkg() {
	zap.L().Info("building go pkg")

	bs, err := go_gen.New().
		WithPages(sdr.pages).
		WithStaticFiles(sdr.staticFiles).
		WithCSS(sdr.css).
		WithBasePath(sdr.basePath).
		Build()
	if err != nil {
		zap.L().
			With(zap.Error(err)).
			Fatal("Failed to generate go pkg")
	}

	filepath := sdr.output + "/" + "docs.go"

	zap.L().Info("exporting go pkg")
	if err := ioutil.WriteFile(filepath, bs, 0644); err != nil {
		zap.L().
			With(zap.Error(err), zap.String("file", filepath)).
			Fatal("Failed to write file")
	}
}
