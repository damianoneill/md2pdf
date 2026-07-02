package md2pdf

import (
	"context"
	"fmt"

	"github.com/damianoneill/md2pdf/domain/converter"
	"github.com/damianoneill/md2pdf/domain/renderer"
)

// ConvertUseCase orchestrates the Markdown → HTML → PDF pipeline.
type ConvertUseCase struct {
	converter converter.Converter
	renderer  renderer.Renderer
}

// NewConvertUseCase returns a ConvertUseCase wired with the given converter and renderer.
func NewConvertUseCase(c converter.Converter, r renderer.Renderer) *ConvertUseCase {
	return &ConvertUseCase{converter: c, renderer: r}
}

// Convert transforms Markdown src into a PDF document.
func (s *ConvertUseCase) Convert(ctx context.Context, src []byte) ([]byte, error) {
	html, err := s.converter.Convert(src)
	if err != nil {
		return nil, fmt.Errorf("convert markdown: %w", err)
	}
	return s.renderer.Render(ctx, html)
}
