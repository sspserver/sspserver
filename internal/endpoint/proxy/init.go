//go:build !htmltemplates
// +build !htmltemplates

package proxy

import "github.com/geniusrabbit/adcorelib/httpserver/extensions/endpoint"

func New() endpoint.Endpoint { return nil }
