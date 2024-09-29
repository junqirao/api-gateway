package utils

import (
	"io"
)

type (
	NopCloseBuf struct {
		content   []byte
		index     int
		origin    io.ReadCloser
		originEOF bool
	}
)

// NewNopCloseBuf constructor
func NewNopCloseBuf() *NopCloseBuf {
	return &NopCloseBuf{}
}

// Read from origin or buffer, when origin is EOF
// it will close origin and ignore error
func (b *NopCloseBuf) Read(p []byte) (n int, err error) {
	// read buffer
	if b.originEOF {
		return b.readBuffered(p)
	}

	// read origin
	n, err = b.origin.Read(p)
	if err == io.EOF {
		b.originEOF = true
		_ = b.origin.Close()
	}
	if n > 0 {
		b.content = append(b.content, p[:n]...)
		b.index += n
	}
	return
}

func (b *NopCloseBuf) readBuffered(p []byte) (n int, err error) {
	if b.index >= len(b.content) {
		return 0, io.EOF
	}
	n = copy(p, b.content[b.index:])
	b.index += n
	if b.index >= len(b.content) {
		return n, io.EOF
	}
	return n, nil
}

// Close do not thing
func (b *NopCloseBuf) Close() error {
	return nil
}

func (b *NopCloseBuf) SetOrigin(r io.ReadCloser) {
	b.origin = r
	return
}

func (b *NopCloseBuf) Reset() {
	b.content = b.content[:0]
	b.index = 0
	b.origin = nil
	b.originEOF = false
}

func (b *NopCloseBuf) ResetIndex() {
	b.index = 0
}
