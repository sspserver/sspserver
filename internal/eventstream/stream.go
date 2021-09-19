//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package eventstream

import (
	"context"

	nc "github.com/geniusrabbit/notificationcenter"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/eventgenerator"
	"geniusrabbit.dev/sspserver/internal/events"
)

// Stream accessor
type Stream interface {
	// SendEvent native action
	SendEvent(ctx context.Context, event *events.Event) error

	// Send response
	Send(event events.Type, status uint8, response adsource.Responser, it adsource.ResponserItem) error
}

type stream struct {
	events    nc.Publisher
	userInfo  nc.Publisher
	generator eventgenerator.Generator
}

// New stream object
func New(events, userInfo nc.Publisher, generator eventgenerator.Generator) Stream {
	return &stream{
		events:    events,
		userInfo:  userInfo,
		generator: generator,
	}
}

// SendEvent native action
func (s *stream) SendEvent(ctx context.Context, event *events.Event) error {
	return s.events.Publish(ctx, event)
}

// Send response
func (s *stream) Send(event events.Type, status uint8, response adsource.Responser, it adsource.ResponserItem) (err error) {
	var (
		info *events.UserInfo
		ctx  = response.Context()
	)
	for _, event := range s.generator.Events(event, status, response, it) {
		if err = s.SendEvent(ctx, event); err != nil {
			break
		}
	}
	if err == nil {
		if info, err = s.generator.UserInfo(response, it); info != nil && err == nil {
			err = s.userInfo.Publish(ctx, info)
		}
	}
	return err
}
