package regdrive

import (
	"bytes"
	"io"
)

type bufferCloser bytes.Buffer

func newReadCloserBuffer(data []byte) io.ReadCloser {
	return (*bufferCloser)(bytes.NewBuffer(data))
}

func (buff *bufferCloser) Read(p []byte) (int, error) {
	return (*bytes.Buffer)(buff).Read(p)
}

func (buff *bufferCloser) Close() error {
	return nil
}
