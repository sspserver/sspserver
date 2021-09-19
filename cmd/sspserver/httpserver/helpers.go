package httpserver

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

func peekOneFromQuery(query *fasthttp.Args, keys ...string) string {
	for _, key := range keys {
		if v := query.Peek(key); len(v) > 0 {
			return string(v)
		}
	}
	return ""
}

func getRoutes(endp endpointer) (routes map[string][]string) {
	if rt, _ := endp.(endpointRouting); rt != nil {
		routes = rt.Routing()
	}
	if routes == nil || len(routes) < 1 {
		routes = map[string][]string{
			":zone": []string{http.MethodGet},
		}
	}
	return routes
}
