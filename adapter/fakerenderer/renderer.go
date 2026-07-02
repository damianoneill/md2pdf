package fakerenderer

import "context"

// Renderer is a test double for domain/renderer.Renderer.
type Renderer struct {
	Input  []byte
	Output []byte
	Err    error
}

func (r *Renderer) Render(_ context.Context, html []byte) ([]byte, error) {
	r.Input = html
	return r.Output, r.Err
}
