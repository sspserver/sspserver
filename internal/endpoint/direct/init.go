package direct

import (
	"geniusrabbit.dev/sspserver/cmd/sspserver/appcontext"
	"geniusrabbit.dev/sspserver/internal/endpoint"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

func init() {
	endpoint.Register("direct", func(source endpoint.Sourcer, options ...interface{}) (endpoint.Endpoint, error) {
		var (
			endpoint = &_endpoint{source: source}
			err      error
		)

		for _, opt := range options {
			switch v := opt.(type) {
			case types.FormatsAccessor:
				endpoint.formats = v
			case *appcontext.Config:
				endpoint.superFailoverURL = v.AdServer.Logic.Direct.DefaultURL
			}
		}

		if endpoint.formats == nil {
			err = ErrInvalidFormatsAccessor
		}

		return endpoint, err
	})
}
