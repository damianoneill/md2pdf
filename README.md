# md2pdf

Converts Markdown documents to PDF, including Mermaid diagrams and syntax-highlighted code blocks.

## How it works

Markdown is parsed by [Goldmark](https://github.com/yuin/goldmark) with GitHub Flavoured Markdown extensions. Mermaid diagram blocks are rendered client-side via Mermaid.js. The resulting HTML is loaded into a headless Chromium browser via [Playwright](https://github.com/mxschmitt/playwright-go), which prints to PDF.

## Usage

```
md2pdf [flags] <input.md>

Flags:
  -o string   output PDF path (default: replaces .md extension with .pdf)
```

## Requirements

A Chromium-based browser must be available. On macOS, Google Chrome is detected automatically. On other platforms, install the Playwright-bundled browser once:

```bash
go run github.com/mxschmitt/playwright-go/cmd/playwright@latest install chromium
```
