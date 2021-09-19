package httpserver

import "geniusrabbit.dev/sspserver/internal/adsource"

// Source of advertisement
type Source interface {
	// Bid request for standart system filter
	Bid(request *adsource.BidRequest) adsource.Responser

	// ProcessResponse when need to fix the result and process all counters
	ProcessResponse(response adsource.Responser)
}
