//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package dynamic

import (
	"errors"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
)

func init() {
	endpoint.Register("dynamic", func(source endpoint.Sourcer, options ...interface{}) (endpoint.Endpoint, error) {
		var (
			endpoint = &_endpoint{source: source}
			err      error
		)

		for _, opt := range options {
			switch v := opt.(type) {
			case adsource.URLGenerator:
				endpoint.urlGen = v
			}
		}

		switch {
		case endpoint.source == nil:
			err = errors.New("dynamic: banner source object is not inited")
		case endpoint.urlGen == nil:
			err = errors.New("dynamic: URL generator is not inited")
		}

		return endpoint, err
	})
}
