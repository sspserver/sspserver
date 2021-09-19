package personification

import (
	"github.com/sspserver/udetect"
)

// Client interface
type Client interface {
	Detect(req *udetect.Request) (*udetect.Response, error)
}
