package cloudreg

import (
	"context"

	"github.com/demdxx/cloudregistry"
)

var (
	ctxKeyCloudRegistryServerNumberConfig = struct{ s string }{s: "cloudreg_server_number_config"}
)

func GetCloudRegistryServerNumberConfig(ctx context.Context) *cloudregistry.SyncUInt64Value {
	return ctx.Value(ctxKeyCloudRegistryServerNumberConfig).(*cloudregistry.SyncUInt64Value)
}

func WithCloudRegistryServerNumberConfig(ctx context.Context, v *cloudregistry.SyncUInt64Value) context.Context {
	return context.WithValue(ctx, ctxKeyCloudRegistryServerNumberConfig, v)
}
