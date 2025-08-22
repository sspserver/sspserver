package commands

import (
	"context"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/eventtraking/pixelgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/urlgenerator"
	nc "github.com/geniusrabbit/notificationcenter/v2"

	"github.com/sspserver/sspserver/cmd/sspserver/appcontext"
)

type EventType = eventgenerator.EventType

func initEventsGenerator[EventT EventType](
	ctx context.Context,
	serviceName, serviceHostname string,
	adServerConf appcontext.AdServerConfig,
	eventAllocator eventgenerator.Allocator[EventT],
) (context.Context, adtype.URLGenerator, eventstream.Stream, error) {
	type (
		LeadT     = *eventgenerator.TestLead
		UserInfoT = *eventgenerator.TestUserInfo
	)

	// Event flow processor
	eventGenerator := eventgenerator.New(
		serviceName,
		eventAllocator,
		func() *eventgenerator.TestUserInfo { return &eventgenerator.TestUserInfo{} },
	)
	eventStream := eventstream.New(
		nc.PublisherByName(eventsStreamName),
		nil, eventGenerator,
	)
	ctx = eventstream.WithStream(ctx, eventStream)

	// Win processor store into the context of requests
	winStream := eventstream.WinNotifications(nc.PublisherByName(winStreamName))
	ctx = eventstream.WithWins(ctx, winStream)

	// URL generator object
	urlGenerator := (&urlgenerator.Generator[EventT, LeadT, UserInfoT]{
		EventGenerator: eventGenerator,
		PixelGenerator: pixelgenerator.NewPixelGenerator[EventT, LeadT](adServerConf.TrackerHost),

		Schema:         "",
		ServiceDomain:  serviceHostname,
		CDNDomain:      adServerConf.CDNDomain,
		LibDomain:      adServerConf.LibDomain,
		ClickPattern:   "/click?c={code}",
		DirectPattern:  "/direct?c={code}",
		WinPattern:     "/win?c={code}",
		PriceParamName: appcontext.PriceParamName,
	}).Init()

	return ctx, urlGenerator, eventStream, nil
}
