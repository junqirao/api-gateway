package authentication

import (
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"api-gateway/internal/components/response"
)

var (
	L = NewLocal()
)

type (
	Local struct {
		secret string
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

func (l *Local) Compare(secretText string, label ...string) bool {
	if l.secret == "" {
		return true
	}
	return l.Encode(label...) == secretText
}

func (l *Local) Encode(label ...string) string {
	b := strings.Builder{}
	labelStr := ""
	if len(label) > 0 && label[0] != "" {
		labelStr = label[0]
	}
	if labelStr != "" {
		b.Write([]byte(labelStr))
		b.WriteString(".")
	}
	b.Write([]byte(l.secret))
	return gmd5.MustEncrypt(b.String())
}

func (l *Local) Middleware(r *ghttp.Request) {
	authStr := r.GetHeader("Authorization")
	if authStr == "" {
		response.WriteJSON(r, response.CodeUnauthorized.WithDetail("Authorization header is required"))
		return
	}
	if !l.Compare(authStr) {
		response.WriteJSON(r, response.CodePermissionDeny)
		return
	}
	r.Middleware.Next()
}
