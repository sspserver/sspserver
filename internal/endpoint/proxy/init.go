//go:build !htmltemplates
// +build !htmltemplates

package proxy

import "geniusrabbit.dev/adcorelib/httpserver/extensions/endpoint"

func New() endpoint.Endpoint { return nil }
