package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/demdxx/gocast"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/ssp/platform"
	"geniusrabbit.dev/sspserver/internal/ssp/platform/info"
	// "geniusrabbit.dev/sspserver/cmd/rotator/api/v1/dsp/platform"
	// "geniusrabbit.dev/sspserver/cmd/rotator/api/v1/info"
)

// PlatformListHandler returns list of platforms with short info
func PlatformListHandler(ctx *fasthttp.RequestCtx) {
	var list []info.Platform

	platform.Each(func(name string, fact platform.Factory) {
		list = append(list, fact.Info().Short())
	})

	ctx.SetStatusCode(http.StatusOK)

	errorResponse(
		ctx,
		json.NewEncoder(ctx).Encode(list),
		http.StatusInternalServerError,
	)
}

// PlatformInfoHandler returns detail infotmation about paticular platform
func PlatformInfoHandler(ctx *fasthttp.RequestCtx) {
	fact := platform.ByProtocol(gocast.ToString(ctx.UserValue("code")))
	if fact != nil {
		ctx.SetStatusCode(http.StatusOK)

		errorResponse(
			ctx,
			json.NewEncoder(ctx).Encode(fact.Info()),
			http.StatusInternalServerError,
		)
	}

	ctx.SetStatusCode(http.StatusNotFound)
	ctx.Write([]byte(`{"response": "undefined platform"}`))
}

func errorResponse(ctx *fasthttp.RequestCtx, err error, code int) {
	if err != nil {
		ctx.SetStatusCode(code)
		json.NewEncoder(ctx).Encode(map[string]interface{}{"error": err.Error()})
	}
}
