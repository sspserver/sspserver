package platform

import (
	"context"
	"time"

	"github.com/geniusrabbit/notificationcenter"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

// WinNotifier redeclared type
type WinNotifier struct {
	p notificationcenter.Publisher
}

// WinNotifications returns win notifier wrapper
func WinNotifications(p notificationcenter.Publisher) *WinNotifier {
	return &WinNotifier{p: p}
}

// Send URL win notify
func (w *WinNotifier) Send(ctx context.Context, url string) error {
	return w.p.Publish(ctx, &adsource.WinEvent{URL: url, Time: time.Now()})
}

// SendEvent win notify
func (w *WinNotifier) SendEvent(ctx context.Context, event *adsource.WinEvent) error {
	return w.p.Publish(ctx, event)
}
