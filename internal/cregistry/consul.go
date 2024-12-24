//go:build !consul
// +build !consul

package cregistry

import (
	"context"
	"net/url"

	"github.com/demdxx/cloudregistry"
	"github.com/demdxx/cloudregistry/consul"
)

func init() {
	registryConnector["consul"] = connectConsul
	registryConnector["consuls"] = connectConsul
}

func connectConsul(ctx context.Context, uri *url.URL) (cloudregistry.Registry, error) {
	return consul.Connect(ctx, consul.WithURI(uri.String()))
}
