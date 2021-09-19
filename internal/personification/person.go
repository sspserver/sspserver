package personification

import (
	"github.com/sspserver/udetect"
)

type person struct {
	request  *udetect.Request
	userInfo UserInfo
}

// User info data
func (p *person) UserInfo() *UserInfo {
	return &p.userInfo
}

// IsInited person in database
func (p *person) IsInited() bool { return false }

// Properties for domain
func (p *person) Properties(name string) Properties { return nil }

// Predict what does he likes?
func (p *person) Predict(req *PredictRequest) (*PredictResponse, error) {
	return nil, nil
}

// PredictPrice what minimal
func (p *person) PredictPrice(req *PredictPriceRequest) (*PredictPriceResponse, error) {
	return nil, nil
}

var EmptyPerson = &person{}
