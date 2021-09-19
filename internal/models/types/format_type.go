//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package types

// FormatType value
type FormatType int

// Format types
const (
	FormatInvalidType      FormatType = -1
	FormatUndefinedType    FormatType = 0
	FormatDirectType       FormatType = 1
	FormatProxyType        FormatType = 2
	FormatVideoType        FormatType = 3 // It's kinde of integrated into video player
	FormatBannerType       FormatType = 4
	FormatBannerHTML5Type  FormatType = 5
	FormatNativeType       FormatType = 6
	FormatAutoregisterType FormatType = 7
	FormatCustomType       FormatType = 8
)

// FormatTypeList of types
var FormatTypeList = []FormatType{
	FormatInvalidType,
	FormatUndefinedType,
	FormatDirectType,
	FormatProxyType, // One of the banner types which works by iframe
	FormatVideoType,
	FormatBannerType,
	FormatBannerHTML5Type,
	FormatNativeType,
	FormatAutoregisterType,
	FormatCustomType,
}

// Name by format type
func (t FormatType) Name() string {
	switch t {
	case FormatInvalidType:
		return `invalid`
	case FormatDirectType:
		return `direct`
	case FormatProxyType:
		return `proxy`
	case FormatVideoType:
		return `video`
	case FormatBannerType:
		return `banner`
	case FormatBannerHTML5Type:
		return `html5`
	case FormatNativeType:
		return `native`
	case FormatAutoregisterType:
		return `autoreg`
	case FormatCustomType:
		return `custom`
	}
	return `undefined`
}

// IsInvalid type of format
func (t FormatType) IsInvalid() bool {
	return t == FormatInvalidType
}

// IsDirect format type
func (t FormatType) IsDirect() bool {
	return t == FormatDirectType
}

// IsBanner format type
func (t FormatType) IsBanner() bool {
	return t == FormatBannerType
}

// IsBannerHTML5 format type
func (t FormatType) IsBannerHTML5() bool {
	return t == FormatBannerHTML5Type
}

// IsAutoregisterType format type
func (t FormatType) IsAutoregisterType() bool {
	return t == FormatAutoregisterType
}

// IsProxy format type
func (t FormatType) IsProxy() bool {
	return t == FormatProxyType
}
