package program

import (
	"net/http"
)

type headerWrapper struct {
	h http.Header
}

func newHeaderWrapper(h http.Header) *headerWrapper {
	return &headerWrapper{h}
}

func (h *headerWrapper) Add(key, value string) bool {
	if h.h == nil {
		return false
	}
	h.h.Add(key, value)
	return true
}

func (h *headerWrapper) Set(key, value string) bool {
	if h.h == nil {
		return false
	}
	h.h.Set(key, value)
	return true
}

func (h *headerWrapper) Get(key string) string {
	if h.h == nil {
		return ""
	}
	return h.h.Get(key)
}
