package response

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type JSON struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(r *ghttp.Request, ec *Code) {
	r.Response.WriteHeader(ec.status)
	r.Response.WriteJson(JSON{
		Code:    ec.Code(),
		Message: ec.Error(),
		Data:    r.GetHandlerResponse(),
	})
}

func WriteData(r *ghttp.Request, ec *Code, data ...interface{}) {
	var dat interface{}
	if len(data) > 0 && data[0] != nil {
		dat = data[0]
	}

	r.Response.WriteHeader(ec.status)
	r.Response.WriteJson(JSON{
		Code:    ec.Code(),
		Message: ec.Message(),
		Data:    dat,
	})
}

func Success(r *ghttp.Request, data ...interface{}) {
	WriteData(r, CodeDefaultSuccess, data...)
}

func Failed(r *ghttp.Request, err error, data ...interface{}) {
	code := CodeDefaultFailure
	if err != nil {
		code = CodeDefaultFailure.WithMessage(err.Error())
	}
	WriteData(r, code, data...)
}
