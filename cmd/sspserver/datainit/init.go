package datainit

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/storage/loader"
)

// ErrUnsupportedDataType in case of data type not supported by paticular loader
var ErrUnsupportedDataType = errors.New("unsupported data type")

// DataLoaderAccessorFnk returns data accessor of the data type
type DataLoaderAccessorFnk func(ctx context.Context, dataType string) (loader.DataAccessor, error)

// DataLoaderFactoryFnk returns general type data accessor
type DataLoaderFactoryFnk func(ctx context.Context, u *url.URL) DataLoaderAccessorFnk

var dataLoaderAccessor = map[string]DataLoaderFactoryFnk{}

// Connect to the dataLoader accessor
func Connect(ctx context.Context, urlConnect string) (DataLoaderAccessorFnk, error) {
	u, err := url.Parse(urlConnect)
	if err != nil {
		return nil, err
	}
	accessor := dataLoaderAccessor[u.Scheme]
	if accessor == nil {
		return nil, fmt.Errorf("unsupported data accessor type [%s]", u.Scheme)
	}
	return accessor(ctx, u), nil
}

type initializer func(debug bool, urlGen adtype.URLGenerator)

var initializers = []initializer{}

// Initialize some other dependencies
func Initialize(debug bool, urlGen adtype.URLGenerator) {
	for _, ifnk := range initializers {
		ifnk(debug, urlGen)
	}
}
