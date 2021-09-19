package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"geniusrabbit.dev/sspserver/internal/ssp"
)

func Test_Options(t *testing.T) {
	server, err := NewServer(
		WithServiceName("test"),
		WithDebugMode(true),
		WithZoneAccessor(dummyTargetAccessor{}),
		WithAdvertisementSource(&ssp.Server{}),
	)
	assert.NoError(t, err)
	assert.Equal(t, "test", server.serviceName, "invalid service name initialisation")
	assert.True(t, server.debug, "invalid debug mode setup")
	assert.True(t, server.zoneAccessor != nil, "invalid zone accessor initialisation")
	assert.True(t, server.source != nil, "invalid SSP server initialisation")
	assert.True(t, server.logger != nil, "invalid logger initialisation")
}
