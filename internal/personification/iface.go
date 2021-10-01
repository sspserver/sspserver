package personification

import "github.com/sspserver/udetect"

type (
	// UserInfo value
	UserInfo struct {
		User   *udetect.User
		Device *udetect.Device
		Geo    *udetect.Geo
	}

	// PredictRequest ...
	PredictRequest struct{}

	// PredictResponse ...
	PredictResponse struct{}

	// PredictPriceRequest ...
	PredictPriceRequest struct{}

	// PredictPriceResponse ...
	PredictPriceResponse struct{}
)

// UUID of the user
func (i *UserInfo) UUID() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.UUID.String()
}

// SessionID of the user
func (i *UserInfo) SessionID() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.SessionID
}

// Fingerprint of the iser
func (i *UserInfo) Fingerprint() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.FingerPrintID
}

// Country info
func (i *UserInfo) Country() *udetect.Geo {
	if i == nil || i.Geo == nil {
		return &udetect.GeoDefault
	}
	return i.Geo
}

// Ages of the user
func (i *UserInfo) Ages() (from, to int) {
	if i == nil || i.User == nil {
		return 0, 0
	}
	return i.User.AgeStart, i.User.AgeEnd
}

// ETag of the user
func (i *UserInfo) ETag() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.ETag
}

// Keywords of the user
func (i *UserInfo) Keywords() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.Keywords
}

// MostPossibleSex of the user
func (i *UserInfo) MostPossibleSex() int {
	if i == nil || i.User == nil {
		return 0
	}
	return i.User.MostPossibleSex()
}

// DeviceInfo get method
func (i *UserInfo) DeviceInfo() *udetect.Device {
	if i == nil || i.Device == nil {
		return &udetect.DeviceDefault
	}
	return i.Device
}

// GeoInfo get method
func (i *UserInfo) GeoInfo() *udetect.Geo {
	if i == nil || i.Geo == nil {
		return &udetect.GeoDefault
	}
	return i.Geo
}

// Person information block
type Person interface {
	// User info data
	UserInfo() *UserInfo

	// IsInited person in database
	IsInited() bool

	// Properties for domain
	Properties(name string) Properties

	// Predict what does he likes?
	Predict(req *PredictRequest) (*PredictResponse, error)

	// PredictPrice what minimal
	PredictPrice(req *PredictPriceRequest) (*PredictPriceResponse, error)
}

// Properties accessor
type Properties interface {
	// Get property by key
	Get(key string) interface{}

	// GetString property by key
	GetString(key string) string

	// GetIntSlice property by key
	GetIntSlice(key string) []int

	// Set property
	Set(key string, prop interface{})

	// Delete property by key
	Delete(key string)

	// Synchronise properties
	Synchronise() error
}
