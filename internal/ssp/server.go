//
// @project geniusrabbit::sspserver 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package ssp

import (
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/adsourceexperiments"
	"geniusrabbit.dev/sspserver/internal/auction"
	"geniusrabbit.dev/sspserver/internal/context/ctxlogger"
	"geniusrabbit.dev/sspserver/internal/gtracing"
)

// Error set...
var (
	ErrSourcesCantBeNil = errors.New("[SSP] seurces cant be nil")
)

const (
	minimalTimeout          = time.Millisecond * 10
	minimalParallelRequests = 1
)

// Server describes the abstraction which can control what request
// where should be sent in which driver
type Server struct {
	// Main source which called everytime
	baseSource adsourceexperiments.SourceWrapper

	// Source list of external platforms
	sources adsource.SourceAccessor

	// RequestTimeout duration
	requestTimeout time.Duration

	// MaxParallelRequest number
	maxParallelRequest int

	// Metrics accessor
	metrics Metrics
}

// NewServer SSP inited with options
func NewServer(options ...Option) (*Server, error) {
	var server Server

	for _, opt := range options {
		opt(&server)
	}

	if server.sources == nil {
		return nil, ErrSourcesCantBeNil
	}

	if server.requestTimeout < minimalTimeout {
		server.requestTimeout = minimalTimeout
	}

	if server.maxParallelRequest < minimalParallelRequests {
		server.maxParallelRequest = minimalParallelRequests
	}

	return &server, nil
}

// Bid request for standart system filter
func (srv *Server) Bid(request *adsource.BidRequest) (response adsource.Responser) {
	var (
		err     error
		count   = srv.maxParallelRequest
		tube    = make(chan adsource.Responser, srv.maxParallelRequest)
		span, _ = gtracing.StartSpanFromContext(request.Context, "ssp.bid")
		referee auction.Referee
		timeout bool
	)

	if span != nil {
		ext.Component.Set(span, "ssp")
		oldContext := request.Context
		request.Context = opentracing.ContextWithSpan(oldContext, span)
		defer func() {
			request.Context = oldContext
			span.Finish()
		}()
	}

	// Base request to internal DB
	if src := srv.getMainSource(); srv != nil && srv.testSource(src, request) {
		startTime := time.Now()
		response := src.Bid(request)
		srv.metrics.IncrementBidRequestCount(src, time.Now().Sub(startTime))

		if response.Error() == nil {
			referee.Push(response.Ads()...)
			// TODO update minimal bids by response
			// TODO release response
		} else {
			srv.metrics.IncrementBidErrorCount(src, response.Error())
		}
	}

	// Source request loop
	for iterator := srv.sources.Iterator(request); ; {
		src := iterator.Next()
		if src == nil {
			break
		}

		if srv.testSource(src, request) {
			count--
			go func(src adsource.Source) {
				startTime := time.Now()
				response := src.Bid(request)
				srv.metrics.IncrementBidRequestCount(src, time.Now().Sub(startTime))
				tube <- response

				if response.Error() != nil {
					srv.metrics.IncrementBidErrorCount(src, response.Error())
				}
			}(src)

			if src.RequestStrategy() == adsource.SingleRequestStrategy {
				break
			}
		}

		if count < 1 {
			break
		}
	}

	// Auction loop
	if count < srv.maxParallelRequest {
		timer := time.NewTimer(srv.requestTimeout)
		for ; count < srv.maxParallelRequest; count++ {
			select {
			case resp := <-tube:
				if e := resp.Error(); e != nil {
					err = e
				} else {
					referee.Push(resp.Ads()...)
				}
			case <-timer.C:
				count = srv.maxParallelRequest
				timeout = true
			}
		}

		if !timeout {
			timer.Stop()
		}
	}

	if items := referee.MatchRequest(request); len(items) > 0 {
		response = adsource.BorrowResponse(request, nil, items, nil)
		err = nil
	} else {
		response = adsource.NewEmptyResponse(request, err)
	}
	return response
}

// ProcessResponse when need to fix the result and process all counters
func (srv *Server) ProcessResponse(response adsource.Responser) {
	// Pricess prices of campaigns
	for _, it := range response.Ads() {
		if it.Validate() != nil {
			continue
		}

		switch ad := it.(type) {
		case adsource.ResponserItem:
			srv.processAdResponse(response, ad)
		case adsource.ResponserMultipleItem:
			for _, it := range ad.Ads() {
				srv.processAdResponse(response, it)
			}
		default:
			ctxlogger.Get(response.Context()).
				Warn("Unsupportable respont item type", zap.String("type", fmt.Sprintf("%T", it)))
		}
	}
}

// SetRequestTimeout of the simple request
func (srv *Server) SetRequestTimeout(timeout time.Duration) {
	if timeout < minimalTimeout {
		timeout = minimalTimeout
	}
	if srv.requestTimeout != timeout {
		srv.requestTimeout = timeout
		srv.sources.SetTimeout(timeout)
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (srv *Server) processAdResponse(response adsource.Responser, ad adsource.ResponserItem) {
	if src := ad.Source(); src != nil {
		src.ProcessResponseItem(response, ad)
	}
}

func (srv *Server) testSource(src adsource.Source, request *adsource.BidRequest) bool {
	return src != nil && request.SourceFilterCheck(src.ID()) && src.Test(request)
}

func (srv *Server) getMainSource() adsource.Source {
	if srv.baseSource == nil {
		return nil
	}
	return srv.baseSource.Next()
}
