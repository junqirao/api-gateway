package authentication

import (
	"testing"
)

func TestLocal(t *testing.T) {
	l1 := NewLocal("test")
	encoded := l1.Encode("aaa")
	t.Logf("l1 encoded: %s", encoded)
	if !l1.Compare(encoded, "aaa") {
		t.Fatal("compare failed")
	}
	l2 := NewLocal("test1")
	encoded = l2.Encode("aaa")
	t.Logf("l2 encoded: %s", encoded)
	if !l2.Compare(encoded, "aaa") {
		t.Fatal("compare failed")
	}
	if l1.Compare(encoded, "aaa") {
		t.Fatal("compare failed")
	}
}
