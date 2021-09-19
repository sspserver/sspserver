package ssp

import (
	"fmt"
	"time"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/adsourceexperiments"
)

// Option sets some property of the server
type Option func(srv *Server)

// WithBaseSource as default
func WithBaseSource(source interface{}) Option {
	return func(srv *Server) {
		switch src := source.(type) {
		case nil:
		case adsource.Source:
			srv.baseSource = adsourceexperiments.NewSimpleWrapper(src)
		case adsourceexperiments.SourceWrapper:
			srv.baseSource = src
		default:
			panic(fmt.Sprintf("Invalid base source type %T", source))
		}
	}
}

// WithSourceAccessor for the server
func WithSourceAccessor(sources adsource.SourceAccessor) Option {
	return func(srv *Server) {
		srv.sources = sources
	}
}

// WithTimeout of one request
func WithTimeout(timeout time.Duration) Option {
	return func(srv *Server) {
		srv.requestTimeout = timeout
	}
}

// WithMaxParallelRequests returns count of requests to external sources by one request
func WithMaxParallelRequests(maxParallelRequest int) Option {
	return func(srv *Server) {
		srv.maxParallelRequest = maxParallelRequest
	}
}
