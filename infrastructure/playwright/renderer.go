package playwright

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	pw "github.com/mxschmitt/playwright-go"

	"github.com/damianoneill/md2pdf/domain/renderer"
)

// Compile-time assertion: Renderer satisfies domain/renderer.Renderer.
var _ renderer.Renderer = &Renderer{}

// Renderer renders an HTML document to PDF using a headless Chromium browser
// via Playwright.
type Renderer struct {
	run     *pw.Playwright
	browser pw.Browser
}

// chromePath returns the path to a locally installed Chrome/Chromium if the
// playwright-bundled browser is not available. Returns nil if no local browser
// is found; playwright then falls back to its bundled browser.
func chromePath() *string {
	candidates := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return pw.String(p)
		}
	}
	return nil
}

// New starts Playwright and launches a headless Chromium browser.
// Call Close when done to release resources.
func New() (*Renderer, error) {
	run, err := pw.Run()
	if err != nil {
		return nil, fmt.Errorf("start playwright: %w", err)
	}

	browser, err := run.Chromium.Launch(pw.BrowserTypeLaunchOptions{
		Headless:       pw.Bool(true),
		ExecutablePath: chromePath(),
	})
	if err != nil {
		run.Stop() //nolint:errcheck
		return nil, fmt.Errorf("launch chromium: %w", err)
	}

	return &Renderer{run: run, browser: browser}, nil
}

// Render writes html to a temporary file, loads it in the headless browser,
// and returns the resulting PDF bytes.
func (r *Renderer) Render(ctx context.Context, html []byte) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f, err := os.CreateTemp("", "md2pdf-*.html")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(f.Name()) //nolint:errcheck // best-effort temp file cleanup

	if _, err := f.Write(html); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("close temp file: %w", err)
	}

	page, err := r.browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("new page: %w", err)
	}
	defer func() { _ = page.Close() }()

	fileURL := "file://" + filepath.ToSlash(f.Name())
	if _, err := page.Goto(fileURL, pw.PageGotoOptions{
		WaitUntil: pw.WaitUntilStateNetworkidle,
	}); err != nil {
		return nil, fmt.Errorf("load page: %w", err)
	}

	pdf, err := page.PDF(pw.PagePdfOptions{
		PrintBackground: pw.Bool(true),
		Format:          pw.String("A4"),
		Margin: &pw.Margin{
			Top:    pw.String("20mm"),
			Bottom: pw.String("20mm"),
			Left:   pw.String("15mm"),
			Right:  pw.String("15mm"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}

	return pdf, nil
}

// Close shuts down the browser and Playwright runtime.
func (r *Renderer) Close() error {
	if err := r.browser.Close(); err != nil {
		return fmt.Errorf("close browser: %w", err)
	}
	if err := r.run.Stop(); err != nil {
		return fmt.Errorf("stop playwright: %w", err)
	}
	return nil
}
