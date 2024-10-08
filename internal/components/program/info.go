package program

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	Info struct {
		Expr        string `json:"expr"` // base64 encoded
		Name        string `json:"name"`
		ServiceName string `json:"service_name"`
	}
)

func (i *Info) TryDecode(ctx context.Context) (res string) {
	res = i.Expr
	if i.Expr != "" {
		v, err := gbase64.DecodeToString(i.Expr)
		if err == nil {
			res = v
		} else {
			g.Log().Warningf(ctx, "decode program expr failed: %v", err)
		}
	}
	return
}

func (i *Info) Decode() (res string, err error) {
	return gbase64.DecodeToString(i.Expr)
}

func (i *Info) TryCompile(ctx context.Context) error {
	_, err := NewProgram(i.Name, i.TryDecode(ctx))
	return err
}
