package testutils

import "io"

// MockReader is a custom io.Reader that modifies the data being read.
type MockReader struct {
	Reader io.Reader
}

func NewMockReader(r io.Reader) *MockReader {
	return &MockReader{Reader: r}
}

func (r *MockReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		p[0] ^= 0xFF // Corrupt the first byte of the data
	}
	return n, err
}
