package getr

import (
	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

func init() {
	endpoint.Register("getr", func(source endpoint.Sourcer, options ...interface{}) (endpoint.Endpoint, error) {
		var (
			endpoint = &_endpoint{source: source}
			err      error
		)

		for _, opt := range options {
			switch v := opt.(type) {
			case types.FormatsAccessor:
				endpoint.formats = v
			case adsource.URLGenerator:
				endpoint.urlGen = v
			}
		}

		if endpoint.formats == nil {
			err = ErrInvalidFormatsAccessor
		}

		return endpoint, err
	})
}
