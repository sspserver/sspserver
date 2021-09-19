package regdrive

import (
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/ssp/platform"
	"geniusrabbit.dev/sspserver/internal/ssp/platform/info"
)

const (
	protocol        = "regdrive"
	protocolVersion = "1.0"
)

type factory struct{}

func (factory) New(source *models.RTBSource, opts ...interface{}) (platform.Platformer, error) {
	return new(source, opts...)
}

func (factory) Info() info.Platform {
	return info.Platform{
		Name:        "Autoregistrator Driver",
		Protocol:    protocol,
		Versions:    []string{protocolVersion},
		Description: "",
		Options: []*info.FieldOption{
			{
				Name:     "request_mapping",
				Type:     info.FieldOptionStringStringMap,
				Required: false,
			},
			{
				Name:     "response_mapping",
				Type:     info.FieldOptionStringStringMap,
				Required: false,
			},
		},
	}
}

func init() {
	platform.Register(factory{})
}
