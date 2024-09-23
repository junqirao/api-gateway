package utils

import (
	"bytes"
)

type (
	NopCloseBuf struct {
		*bytes.Buffer
	}
)

// NewNopCloseBuf constructor
func NewNopCloseBuf() *NopCloseBuf {
	return &NopCloseBuf{Buffer: &bytes.Buffer{}}
}

// Close do not thing
func (b *NopCloseBuf) Close() error {
	return nil
}
