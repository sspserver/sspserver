package adsource

import (
	"sync"

	"geniusrabbit.dev/sspserver/internal/rand"
)

// Common pools
var (
	responseAdItemPool = sync.Pool{
		New: func() interface{} { return new(ResponseAdItem) },
	}
	responsePool = sync.Pool{
		New: func() interface{} { return new(Response) },
	}
	requestRTBPool = sync.Pool{
		New: func() interface{} { return new(RTBRequest) },
	}
)

///////////////////////////////////////////////////////////////////////////////
/// Response ad item sync pool
///////////////////////////////////////////////////////////////////////////////

// BorrowResponseAdItem object
func BorrowResponseAdItem() *ResponseAdItem {
	item := responseAdItemPool.Get().(*ResponseAdItem)
	item.ItemID = rand.UUID()
	return item
}

// ReturnResponseAdItem back to pool
func ReturnResponseAdItem(o *ResponseAdItem) {
	o.reset()
	responseAdItemPool.Put(o)
}

///////////////////////////////////////////////////////////////////////////////
/// Response sync pool
///////////////////////////////////////////////////////////////////////////////

// BorrowResponse object
func BorrowResponse(request *BidRequest, source Sourcer, items []ResponserItemCommon, err error) *Response {
	resp := responsePool.Get().(*Response)
	resp.request = request
	resp.source = source
	resp.items = items
	resp.err = err
	return resp
}

// ReturnResponse back to pool
func ReturnResponse(o *Response) {
	o.reset()
	responsePool.Put(o)
}

///////////////////////////////////////////////////////////////////////////////
/// Request RTB item sync pool
///////////////////////////////////////////////////////////////////////////////

// BorrowRTBRequest object
func BorrowRTBRequest() *RTBRequest {
	return requestRTBPool.Get().(*RTBRequest)
}

// ReturnRequestRTB back to pool
func ReturnRequestRTB(o *RTBRequest) {
	o.reset()
	requestRTBPool.Put(o)
}
