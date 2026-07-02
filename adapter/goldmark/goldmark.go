package goldmark

import (
	"bytes"
	"fmt"
	"html/template"

	gm "github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/mermaid"

	"github.com/damianoneill/md2pdf/domain/converter"
)

// Compile-time assertion: Converter satisfies domain/converter.Converter.
var _ converter.Converter = &Converter{}

// pageTmpl wraps a rendered HTML body in a full document with GitHub Markdown CSS.
var pageTmpl = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.8.1/github-markdown.min.css">
<style>
body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 45px;
}
@media print {
  body { padding: 20px; }
  table { font-size: 0.75em; border-collapse: collapse; }
  tr { break-inside: avoid; }
  td, th { overflow-wrap: break-word; min-width: 0; }
  h1, h2, h3, h4, h5, h6 { break-after: avoid; }
}
</style>
</head>
<body class="markdown-body">
{{.}}
</body>
</html>`))

// Converter converts Markdown to a complete HTML document using Goldmark with
// GFM extensions, syntax highlighting (Chroma), and Mermaid diagram support.
type Converter struct {
	md gm.Markdown
}

// New returns a Converter ready for use.
func New() *Converter {
	md := gm.New(
		gm.WithExtensions(
			extension.GFM,
			highlighting.Highlighting,
			&mermaid.Extender{},
		),
		gm.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		gm.WithRendererOptions(
			gmhtml.WithUnsafe(),
			// WithUnsafe allows raw HTML passthrough from the markdown source.
			// Safe here: this tool processes local files under user control,
			// not untrusted web content.
		),
	)
	return &Converter{md: md}
}

// Convert transforms src Markdown into a complete HTML document.
func (c *Converter) Convert(src []byte) ([]byte, error) {
	var body bytes.Buffer
	if err := c.md.Convert(src, &body); err != nil {
		return nil, fmt.Errorf("goldmark: %w", err)
	}

	var out bytes.Buffer
	if err := pageTmpl.Execute(&out, template.HTML(body.String())); err != nil { //nolint:gosec // G203: intentional — goldmark output with WithUnsafe is trusted local content
		return nil, fmt.Errorf("html template: %w", err)
	}
	return out.Bytes(), nil
}
