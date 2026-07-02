package md2pdf_test

import (
	"bytes"
	"context"
	"testing"

	goldmarkadapter "github.com/damianoneill/md2pdf/adapter/goldmark"
	"github.com/damianoneill/md2pdf/infrastructure/playwright"
	"github.com/damianoneill/md2pdf/usecase/md2pdf"
)

// TestConvertUseCase_Convert_integration exercises the full Markdown → PDF
// pipeline with real implementations. It is skipped under -short.
func TestConvertUseCase_Convert_integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	src := []byte(`# Integration Test

A paragraph with **bold** and _italic_ text.

` + "```go\nfunc main() {}\n```" + `

` + "```mermaid\ngraph TD\n    A --> B\n```")

	renderer, err := playwright.New()
	if err != nil {
		t.Skipf("playwright unavailable: %v", err)
	}
	defer func() {
		if cerr := renderer.Close(); cerr != nil {
			t.Logf("close renderer: %v", cerr)
		}
	}()

	svc := md2pdf.NewConvertUseCase(goldmarkadapter.New(), renderer)

	pdf, err := svc.Convert(context.Background(), src)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF")) {
		t.Fatalf("output is not a PDF (got first bytes: %q)", pdf[:min(8, len(pdf))])
	}
}
