//
// @project GeniusRabbit rotator 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2019
//

package banner

import (
	"errors"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/endpoint"
)

func init() {
	endpoint.Register("banner", func(source endpoint.Sourcer, options ...interface{}) (endpoint.Endpoint, error) {
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
			err = errors.New("banner: banner source object is not inited")
		case endpoint.urlGen == nil:
			err = errors.New("banner: URL generator is not inited")
		}

		return endpoint, err
	})
}
