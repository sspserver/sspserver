package personification

import (
	"context"

	"github.com/sspserver/udetect"
)

// Client interface
type Client interface {
	Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error)
}

// Connect to the udetect server
func Connect(tr udetect.Transport) Client {
	return udetect.NewClient(tr)
}

type DummyClient struct {
}

func (DummyClient) Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error) {
	return nil, nil
}
