//go:build !etcd
// +build !etcd

package cregistry

import (
	"context"
	"net/url"

	"github.com/demdxx/cloudregistry"
	"github.com/demdxx/cloudregistry/etcd"
)

func init() {
	registryConnector["etcd"] = connectEtcd
	registryConnector["etcds"] = connectEtcd
}

func connectEtcd(ctx context.Context, uri *url.URL) (cloudregistry.Registry, error) {
	return etcd.Connect(ctx, etcd.WithURI(uri.String()))
}
