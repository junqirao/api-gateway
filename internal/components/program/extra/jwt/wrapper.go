package jwt

type (
	Wrapper struct {
		Err    error
		Claims map[string]interface{}
	}
)

func ParseToken(s string) *Wrapper {
	w := &Wrapper{}
	w.Claims, w.Err = parseTokenMap(s)
	return w
}

// MustSuccess for internal/components/program/program.go:55 resultWrapper
func (w *Wrapper) MustSuccess() error {
	return w.Err
}
