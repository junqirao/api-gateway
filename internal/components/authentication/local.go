package authentication

import (
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	L = NewLocal()
)

type (
	Local struct {
		secret string // twice hashed
	}
)

func NewLocal(secret ...string) *Local {
	s := ""
	if len(secret) > 0 && secret[0] != "" {
		s = secret[0]
	} else {
		s = g.Cfg().MustGet(gctx.GetInitCtx(), "auth.secret", "").String()
	}
	return &Local{
		secret: s,
	}
}

func (l *Local) Compare(encode string, secretText string) bool {
	if l.secret == "" {
		return true
	}
	return l.Encode(encode) == secretText
}

func (l *Local) Encode(encode string) string {
	b := strings.Builder{}
	b.Write([]byte(encode))
	b.WriteString(".")
	b.Write([]byte(l.secret))
	return gmd5.MustEncrypt(b.String())
}
