package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	goldmarkadapter "github.com/damianoneill/md2pdf/adapter/goldmark"
	"github.com/damianoneill/md2pdf/infrastructure/playwright"
	"github.com/damianoneill/md2pdf/usecase/md2pdf"
)

func main() {
	output := flag.String("o", "", "output PDF path (default: replaces .md extension with .pdf)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: md2pdf [flags] <input.md>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)
	src, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("read %s: %v", inputPath, err)
	}

	outPath := *output
	if outPath == "" {
		if strings.HasSuffix(inputPath, ".md") {
			outPath = strings.TrimSuffix(inputPath, ".md") + ".pdf"
		} else {
			outPath = inputPath + ".pdf"
		}
	}

	converter := goldmarkadapter.New()

	renderer, err := playwright.New()
	if err != nil {
		log.Fatalf("init renderer: %v", err)
	}
	defer func() {
		if cerr := renderer.Close(); cerr != nil {
			log.Printf("close renderer: %v", cerr)
		}
	}()

	svc := md2pdf.NewConvertUseCase(converter, renderer)

	pdf, err := svc.Convert(context.Background(), src)
	if err != nil {
		if cerr := renderer.Close(); cerr != nil {
			log.Printf("close renderer: %v", cerr)
		}
		log.Fatalf("convert: %v", err)
	}

	if err := os.WriteFile(outPath, pdf, 0o644); err != nil {
		if cerr := renderer.Close(); cerr != nil {
			log.Printf("close renderer: %v", cerr)
		}
		log.Fatalf("write %s: %v", outPath, err)
	}

	fmt.Printf("wrote %s\n", outPath)
}
