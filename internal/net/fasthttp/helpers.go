//
// @project geniusrabbit::corelib 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package fasthttp

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

// IsSecure request
func IsSecure(ctx *fasthttp.RequestCtx) bool {
	return bytes.EqualFold(ctx.Request.Header.Peek("X-Forwarded-Proto"), []byte("https"))
}

// IsSecureCF request
func IsSecureCF(ctx *fasthttp.RequestCtx) bool {
	return bytes.Contains(ctx.Request.Header.Peek("Cf-Visitor"), []byte(`"https"`)) || IsSecure(ctx)
}
