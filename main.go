package main

import (
	"flag"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lonnblad/go-service-doc/exporting/golang"
	"github.com/lonnblad/go-service-doc/exporting/simple"
	"github.com/lonnblad/go-service-doc/parser"
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
	sourceDir := flag.String("d", "docs", "Directory where to get markdown files.")
	outputDir := flag.String("o", "docs", "Directory where to write output.")
	basepath := flag.String("p", "/docs", "Base path for the generated documentation.")

	flag.Parse()

	mdParser := parser.NewParser().
		WithSourceDir(*sourceDir).
		WithOutputDir(*outputDir).
		WithBasepath(*basepath).
		ServiceFilename(*serviceFilename)

	mdParser.Run()

	if err := mdParser.Error(); err != nil {
		zap.L().With(zap.Error(err)).
			Error("parser returned an error")

		return
	}

	pages := mdParser.Pages()
	staticFiles := mdParser.StaticFiles()
	searchPage := mdParser.SearchPage()

	simpleExporter := simple.NewExporter().
		WithSourceDir(*sourceDir).
		WithOutputDir(*outputDir).
		WithPages(pages).
		WithStaticFiles(staticFiles)

	simpleExporter.Run()

	if err := simpleExporter.Error(); err != nil {
		zap.L().With(zap.Error(err)).
			Error("exporting simple returned an error")

		return
	}

	goExporter := golang.NewExporter().
		WithOutputDir(*outputDir).
		WithBasepath(*basepath).
		WithPages(pages).
		WithStaticFiles(staticFiles).
		WithSearchPage(searchPage)

	goExporter.Run()

	if err := goExporter.Error(); err != nil {
		zap.L().With(zap.Error(err)).
			Error("exporting golang returned an error")

		return
	}

	zap.L().Info("done")
}
