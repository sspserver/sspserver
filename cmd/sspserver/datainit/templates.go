//go:build htmltemplates
// +build htmltemplates

package datainit

import (
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/sspserver/private/templates"
)

func init() {
	initializers = append(initializers, func(debug bool, urlGen adtype.URLGenerator) {
		templates.Debug = debug
		templates.URLGen = urlGen
	})
}
