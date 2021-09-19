//
// @project GeniusRabbit 2015
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015
//

package wrappers

import (
	"net/http"
)

func NoCache(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hd := w.Header()
		hd.Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		hd.Set("Pragma", "no-cache")                                   // HTTP 1.0.
		hd.Set("Expires", "0")                                         // Proxies
		hd.Set("Vary", "*")
		h(w, r)
	}
}
