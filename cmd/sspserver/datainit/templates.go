//go:build htmltemplates
// +build htmltemplates

package datainit

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adstdendpoints/templates"
)

func init() {
	initializers = append(initializers, func(debug bool, urlGen adtype.URLGenerator) {
		templates.Debug = debug
		templates.URLGen = urlGen
	})
}
