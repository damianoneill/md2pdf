package converter

// Converter transforms Markdown source bytes into an HTML document.
type Converter interface {
	Convert(src []byte) ([]byte, error)
}
