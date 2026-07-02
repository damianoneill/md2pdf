package playwright

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	pw "github.com/mxschmitt/playwright-go"

	"github.com/damianoneill/md2pdf/domain/renderer"
)

// a4ContentViewportPx is the viewport width that simulates the PDF printable area.
// A4 = 210mm at 96dpi = 794px; minus 15mm left/right margins (57px each) = 680px.
// Setting the viewport to this value before measuring tables makes scrollWidth
// reflect the same layout constraints that Chromium applies during PDF generation.
const a4ContentViewportPx = 680

// a4HeightViewportPx is the A4 page height in CSS pixels at 96dpi.
// 297mm × (96px/25.4mm) ≈ 1123px. Used only to satisfy SetViewportSize;
// table width measurements are unaffected by viewport height.
const a4HeightViewportPx = 1123

// a4ContentHeightPx is the usable content height in CSS pixels for A4 paper
// with 20mm top/bottom PDF margins (76px each) and 20px body padding each side.
// 1123 − 76 − 76 − 40 = 931px.
// Matches: PagePdfOptions Margin Top/Bottom=20mm and @media print body padding=20px.
// Update this constant if either value changes.
const a4ContentHeightPx = 931

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

	// Emulate print media so @media print CSS applies during measurement,
	// then constrain the viewport to the PDF printable width so scrollWidth
	// values reflect the same layout constraints used by page.PDF().
	if err := page.EmulateMedia(pw.PageEmulateMediaOptions{
		Media: pw.MediaPrint,
	}); err != nil {
		return nil, fmt.Errorf("emulate print media: %w", err)
	}
	if err := page.SetViewportSize(a4ContentViewportPx, a4HeightViewportPx); err != nil {
		return nil, fmt.Errorf("set viewport: %w", err)
	}

	if err := scaleWideTables(page); err != nil {
		return nil, err
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

// scaleWideTables shrinks any table whose natural width exceeds the body
// content area using CSS zoom (which affects layout, unlike transform: scale).
// It temporarily forces white-space:nowrap on cells to measure the true
// unconstrained table width before applying zoom.
// Tables that fit on one page after zooming get break-inside:avoid; taller
// tables fall back to break-inside:auto so they span pages naturally (row-level
// breaking is still prevented by the tr { break-inside: avoid } CSS rule).
func scaleWideTables(page pw.Page) error {
	script := fmt.Sprintf(`(function(maxH) {
		var body = document.body;
		var cs = window.getComputedStyle(body);
		var cw = body.clientWidth
			- parseFloat(cs.paddingLeft)
			- parseFloat(cs.paddingRight);
		document.querySelectorAll('table').forEach(function(t) {
			// Force no-wrap to measure true unconstrained table width.
			var cells = t.querySelectorAll('td, th');
			cells.forEach(function(c) { c.style.whiteSpace = 'nowrap'; });
			var tw = t.scrollWidth;
			cells.forEach(function(c) { c.style.whiteSpace = ''; });
			if (tw > cw) {
				t.style.zoom = (cw / tw).toFixed(4);
			}
			// Keep short tables on one page; let tall tables span pages naturally.
			t.style.breakInside = t.offsetHeight <= maxH ? 'avoid' : 'auto';
		});
	})(%d)`, a4ContentHeightPx)
	if _, err := page.Evaluate(script); err != nil {
		return fmt.Errorf("scale wide tables: %w", err)
	}
	return nil
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
