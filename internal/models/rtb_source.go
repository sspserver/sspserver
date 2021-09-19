//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package models

import (
	"github.com/geniusrabbit/gosql"

	"geniusrabbit.dev/sspserver/internal/models/types"
)

// RTBRequestType contains type of representation of request information
type RTBRequestType int

// Request types
const (
	RTBRequestTypeUndefined       RTBRequestType = 0
	RTBRequestTypeJSON            RTBRequestType = 1
	RTBRequestTypeXML             RTBRequestType = 2
	RTBRequestTypeProtoBUFF       RTBRequestType = 3
	RTBRequestTypePOSTFormEncoded RTBRequestType = 4 // application/x-www-form-urlencoded
	RTBRequestTypePLAINTEXT       RTBRequestType = 5
)

type companyStore interface {
	CompanyByID(id uint64) *Company
}

// RTBSourceOptions flags
type RTBSourceOptions struct {
	ErrorsIgnore bool
	Trace        bool
}

// RTBSource describe the source of external DSP platform or similar exchange protocol.
// All that sources have similar options and very common prefilter configurations
type RTBSource struct {
	ID      uint64
	Company *Company

	Protocol    string         // rtb as default
	URL         string         // RTB client request URL
	Method      string         // HTTP method GET, POST, ect; Default POST
	RequestType RTBRequestType // 1 - json, 2 - xml, 3 - ProtoBUFF, 4 - MultipleFormaData, 5 - PLAINTEXT
	Headers     gosql.Hstore

	AuctionType int              // default: 0 – first price type, 1 – second price type
	RPS         int              // 0 – unlimit
	Timeout     int              // In milliseconds
	Options     RTBSourceOptions //
	Filter      types.BaseFilter //

	Accuracy           float64 // Price accuracy for auction in percentages
	RevenueShareReduce float64 // % 100, 80%, 65.5%
	MinimalWeight      float64

	Flags  gosql.Hstore
	Config gosql.NullableJSON
}

// // RTBSourceFromModel convert database model to specified model
// func RTBSourceFromModel(cl *models.RTBSource, comp *Company) (src *RTBSource) {
// 	if comp == nil {
// 		return nil
// 	}

// 	var (
// 		opt = RTBSourceOptions{
// 			ErrorsIgnore: cl.Flags.GetBool("errors_ignore"),
// 			Trace:        cl.Flags.GetBool("trace"),
// 		}
// 		filter = types.BaseFilter{
// 			Secure:          cl.Secure,
// 			Adblock:         cl.AdBlock,
// 			PrivateBrowsing: cl.PrivateBrowsing,
// 			IP:              cl.IP,
// 		}
// 	)

// 	filter.Set(types.FieldFormat, cl.Formats)
// 	filter.Set(types.FieldDeviceTypes, cl.DeviceTypes)
// 	filter.Set(types.FieldDevices, cl.Devices)
// 	filter.Set(types.FieldOS, cl.OS)
// 	filter.Set(types.FieldBrowsers, cl.Browsers)
// 	filter.Set(types.FieldCategories, cl.Categories)
// 	filter.Set(types.FieldCountries, cl.Countries)
// 	filter.Set(types.FieldLanguages, cl.Languages)
// 	filter.Set(types.FieldZones, cl.Zones)
// 	filter.Set(types.FieldDomains, cl.Domains)

// 	return &RTBSource{
// 		ID:                 cl.ID,
// 		Company:            comp,
// 		Protocol:           strings.ToLower(cl.Protocol),
// 		URL:                cl.URL,
// 		Method:             strings.ToUpper(cl.Method),
// 		RequestType:        cl.RequestType,
// 		Headers:            cl.Headers,
// 		AuctionType:        cl.AuctionType,
// 		RPS:                cl.RPS,
// 		Timeout:            cl.Timeout,
// 		Options:            opt,
// 		Filter:             filter,
// 		Accuracy:           cl.Accuracy,
// 		RevenueShareReduce: cl.RevenueShareReduce,
// 		MinimalWeight:      cl.MinimalWeight,
// 		Flags:              cl.Flags,
// 		Config:             cl.Config,
// 	}
// }

// Test RTB source
func (s *RTBSource) Test(t types.TargetPointer) bool {
	return s.Filter.Test(t)
}

// TestFormat available in filter
func (s *RTBSource) TestFormat(f *types.Format) bool {
	return s.Filter.TestFormat(f)
}

// RevenueShareReduceFactor from 0. to 1.
func (s *RTBSource) RevenueShareReduceFactor() float64 {
	return s.RevenueShareReduce / 100.0
}
