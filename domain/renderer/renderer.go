package renderer

import "context"

// Renderer transforms an HTML document into a PDF.
type Renderer interface {
	Render(ctx context.Context, html []byte) ([]byte, error)
}
