//
// @project GeniusRabbit corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package types

import "strings"

// PricingModel value
type PricingModel uint8

// PricingModel consts
const (
	PricingModelUndefined PricingModel = iota
	PricingModelCPM
	PricingModelCPC
	PricingModelCPA
)

func (pm PricingModel) String() string {
	return pm.Name()
}

// Name value
func (pm PricingModel) Name() string {
	switch pm {
	case PricingModelCPM:
		return `CPM`
	case PricingModelCPC:
		return `CPC`
	case PricingModelCPA:
		return `CPA`
	}
	return `undefined`
}

// IsCPM model
func (pm PricingModel) IsCPM() bool {
	return pm == PricingModelCPM
}

// IsCPC model
func (pm PricingModel) IsCPC() bool {
	return pm == PricingModelCPC
}

// IsCPA model
func (pm PricingModel) IsCPA() bool {
	return pm == PricingModelCPA
}

// UInt value
func (pm PricingModel) UInt() uint {
	return uint(pm)
}

// PricingModelByName string
func PricingModelByName(model string) PricingModel {
	switch strings.ToUpper(model) {
	case `CPM`:
		return PricingModelCPM
	case `CPC`:
		return PricingModelCPC
	case `CPA`:
		return PricingModelCPA
	}
	return PricingModelUndefined
}
