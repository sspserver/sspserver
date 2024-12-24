package cregistry

import (
	"context"
	"errors"
	"net/url"

	"github.com/demdxx/cloudregistry"
)

var ErrUnknownRegistry = errors.New("unknown registry")

type connector func(ctx context.Context, uri *url.URL) (cloudregistry.Registry, error)

var registryConnector = map[string]connector{}

// Connect to the registry
func Connect(ctx context.Context, uri string) (cloudregistry.Registry, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if c, ok := registryConnector[u.Scheme]; ok {
		return c(ctx, u)
	}
	return nil, ErrUnknownRegistry
}
