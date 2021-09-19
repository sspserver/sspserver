//
// @project geniusrabbit::corelib 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package fasthttp

import (
	"mime/multipart"
	"strings"

	"github.com/valyala/fasthttp"
)

// ParamReader for echo request
type ParamReader struct {
	Form *multipart.Form
	*fasthttp.Args
}

// NewParamReader for echo request
func NewParamReader(ctx *fasthttp.RequestCtx) ParamReader {
	var (
		method = strings.ToUpper(string(ctx.Method()))
		form   *multipart.Form
	)
	switch method {
	case "POST", "PUT":
		form, _ = ctx.Request.MultipartForm()
	}
	return ParamReader{
		Form: form,
		Args: ctx.Request.URI().QueryArgs(),
	}
}

// Param from reuqst
func (r ParamReader) Param(key string) string {
	if r.Form != nil {
		if v, ok := r.Form.Value[key]; ok {
			if len(v) > 0 {
				return v[0]
			}
			return ""
		}
	}
	return string(r.Args.Peek(key))
}
