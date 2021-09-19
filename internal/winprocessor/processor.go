package winprocessor

import (
	nc "github.com/geniusrabbit/notificationcenter"
)

type processor struct{}

// New processor object
func New() nc.Receiver {
	return processor{}
}

func (p processor) Receive(msg nc.Message) error {
	return nil
}
