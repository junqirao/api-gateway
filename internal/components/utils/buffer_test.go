package utils

import (
	"bytes"
	"io"
	"testing"
)

type testBuf struct {
	*bytes.Buffer
}

func (t testBuf) Close() error {
	return nil
}

func TestBuffer(t *testing.T) {
	buf := NewNopCloseBuf()
	b := []byte("hello")
	from := &testBuf{Buffer: &bytes.Buffer{}}
	length := len(b)
	from.Write(b)
	buf.SetOrigin(from)

	for i := 0; i < 2; i++ {
		bb := &bytes.Buffer{}
		written, err := io.Copy(bb, buf)
		if err != nil {
			t.Fatal(err)
			return
		}

		t.Logf("read %d bytes", written)

		if bb.String() != "hello" {
			t.Fatal("fail")
			return
		}
		buf.ResetIndex()
	}

	target := &bytes.Buffer{}
	n, err := io.Copy(target, io.LimitReader(buf, int64(length)))
	if err != nil {
		t.Fatal(err)
		return
	}
	if n != int64(length) {
		t.Fatalf("expect read %d bytes, got %d", length, n)
	}
}
