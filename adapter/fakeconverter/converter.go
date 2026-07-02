package fakeconverter

// Converter is a test double for domain/converter.Converter.
type Converter struct {
	Input  []byte
	Output []byte
	Err    error
}

func (c *Converter) Convert(src []byte) ([]byte, error) {
	c.Input = src
	return c.Output, c.Err
}
