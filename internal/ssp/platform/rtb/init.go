package rtb

import (
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/ssp/platform"
	"geniusrabbit.dev/sspserver/internal/ssp/platform/info"
)

const protocol = "openrtb"

type factory struct{}

func (factory) New(source *models.RTBSource, opts ...interface{}) (platform.Platformer, error) {
	return New(source, opts...)
}

func (factory) Info() info.Platform {
	return info.Platform{
		Name:        "OpenRTB",
		Protocol:    protocol,
		Versions:    []string{"2.3", "2.4", "2.5"},
		Description: "",
		Docs: []info.Documentation{
			{
				Title: "OpenRTB (Real-Time Bidding)",
				Link:  "https://www.iab.com/guidelines/real-time-bidding-rtb-project/",
			},
		},
		Subprotocols: []info.Subprotocol{
			{
				Name:     "VAST",
				Protocol: "vast",
			},
			{
				Name:     "OpenNative",
				Protocol: "opennative",
			},
		},
	}
}

func init() {
	platform.Register(factory{})
}
