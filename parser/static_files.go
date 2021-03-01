package parser

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/lonnblad/go-service-doc/core"
	"github.com/lonnblad/go-service-doc/utils"
)

func (p *Parser) findStaticFiles() {
	zap.L().Info("search for static files")

	files, err := ioutil.ReadDir(p.sourceDir + "/static")
	if os.IsNotExist(err) {
		return
	}

	if err != nil {
		p.err = errors.Wrap(err, "ioutil.ReadDir failed")
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		var (
			file          core.File
			fileExtension string
		)

		switch {
		case strings.HasSuffix(f.Name(), ".svg"):
			file.ContentType = "image/svg+xml"
			fileExtension = ".svg"
		case strings.HasSuffix(f.Name(), ".png"):
			file.ContentType = "image/png"
			fileExtension = ".png"
		case strings.HasSuffix(f.Name(), ".ico"):
			file.ContentType = "image/ico"
			fileExtension = ".ico"
		default:
			continue
		}

		file.Path = p.outputDir + "/static/" + f.Name()
		file.Href = p.basepath + "/static/" + utils.ConvertToKebabCase(f.Name())

		if f.Name() == "favicon.ico" {
			p.faviconHref = file.Href
		}

		file.Name = strings.ReplaceAll(f.Name(), fileExtension, "")
		file.Name = utils.ConvertToCamelCase(file.Name)

		filepath := p.sourceDir + "/static/" + f.Name()

		file.Content, err = ioutil.ReadFile(filepath)
		if err != nil {
			p.err = errors.Wrap(err, "ioutil.ReadFile failed")
			return
		}

		p.staticFiles = append(p.staticFiles, file)
	}
}
