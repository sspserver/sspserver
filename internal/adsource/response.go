//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adsource

import (
	"context"
)

// Response from different sources
type Response struct {
	request *BidRequest
	source  Sourcer
	items   []ResponserItemCommon
	err     error
	context context.Context
}

// NewResponse common object
func NewResponse(request *BidRequest, source Sourcer, items []ResponserItemCommon, err error) *Response {
	return &Response{
		request: request,
		source:  source,
		items:   items,
		err:     err,
	}
}

// NewErrorResponse object
func NewErrorResponse(request *BidRequest, err error) *Response {
	return NewResponse(request, nil, nil, err)
}

// AuctionID response
func (r *Response) AuctionID() string {
	return r.request.ID
}

// AuctionType of request
func (r *Response) AuctionType() AuctionType {
	return r.request.AuctionType
}

// Source of response
func (r *Response) Source() Sourcer {
	return r.source
}

// Request information
func (r *Response) Request() *BidRequest {
	return r.request
}

// AddItem to response
func (r *Response) AddItem(it ResponserItemCommon) {
	r.items = append(r.items, it)
}

// Item by impression code
func (r *Response) Item(impid string) ResponserItemCommon {
	for _, it := range r.items {
		if it.ImpressionID() == impid {
			return it
		}
	}
	return nil
}

// Ads list
func (r *Response) Ads() []ResponserItemCommon {
	return r.items
}

// Count of response items
func (r *Response) Count() int {
	return len(r.items)
}

// Validate response
func (r *Response) Validate() (err error) {
	if r.err != nil {
		return r.err
	}
	if r.Count() < 1 {
		return ErrResponseEmpty
	}
	for _, it := range r.items {
		if err = it.Validate(); err != nil {
			break
		}
	}
	return
}

// Error of the response
func (r *Response) Error() error {
	return r.err
}

// Context value
func (r *Response) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		oldContext := r.context
		r.context = ctx[0]
		return oldContext
	}
	return r.context
}

// Get context item by key
func (r *Response) Get(key string) interface{} {
	if r.context == nil {
		return nil
	}
	return r.context.Value(key)
}

func (r *Response) reset() {
	r.items = r.items[:0]
	r.err = nil
	r.context = nil
	r.source = nil
	r.request = nil
}

var (
	_ Responser = &Response{}
)
