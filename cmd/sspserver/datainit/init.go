package datainit

import (
	"github.com/geniusrabbit/adcorelib/adtype"
)

type initializer func(debug bool, urlGen adtype.URLGenerator)

var initializers = []initializer{}

// Initialize some other dependencies
func Initialize(debug bool, urlGen adtype.URLGenerator) {
	for _, ifnk := range initializers {
		ifnk(debug, urlGen)
	}
}
