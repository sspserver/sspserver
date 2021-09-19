package personification

type (
	// UserInfo value
	UserInfo struct{}

	// PredictRequest ...
	PredictRequest struct{}

	// PredictResponse ...
	PredictResponse struct{}

	// PredictPriceRequest ...
	PredictPriceRequest struct{}

	// PredictPriceResponse ...
	PredictPriceResponse struct{}
)

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
