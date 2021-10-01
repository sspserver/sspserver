//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package rtb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bsm/openrtb"
	"go.uber.org/zap"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/eventtraking/events"
	"geniusrabbit.dev/sspserver/internal/eventtraking/eventstream"
	"geniusrabbit.dev/sspserver/internal/models"
	counter "geniusrabbit.dev/sspserver/internal/ssp/errorcounter"
	"geniusrabbit.dev/sspserver/internal/ssp/openlatency"
	"geniusrabbit.dev/sspserver/internal/ssp/optimizer"
	"geniusrabbit.dev/sspserver/internal/ssp/platform"
	"geniusrabbit.dev/sspserver/internal/ssp/platform/dataencode"
)

var (
	errInvalidRTBSource      = errors.New("SSP base platform: invalid RTB source model")
	errInvalidWinEventStream = errors.New("SSP base platform: invalid win event stream")
	errInvalidEventStream    = errors.New("SSP base platform: invalid event stream")
	errInvalidLogger         = errors.New("SSP base platform: invalid logger object")
)

const (
	headerRequestOpenRTBVersion        = "X-Openrtb-Version"
	headerRequestOpenRTBVersionDefault = "2.5"
)

// Platform object
type Platform struct {
	// Original source model
	source *models.RTBSource

	// Request headers
	Headers map[string]string

	// Extend current playform interface
	Extend platformExt

	// Request prepare interfaces
	requestPreparer  platformRequestPreparer
	responsePreparer platformResponsePreparer
	marshalel        dataencode.Encoder

	// Requests RPS counter
	rpsCurrent     counter.Counter
	errorCounter   counter.ErrorCounter
	metricsCounter *openlatency.MetricsCounter

	// Revenue share this is the rate of publisher revenue
	revenueShareReduce float64

	// win events stream
	winEvents platform.WinNotifier

	// bid action stream
	bidActions eventstream.Stream

	// Client of HTTP requests
	Client *http.Client

	// Datetime of last first request
	lastPeriod time.Time

	// optimizator object
	optimizator *optimizer.Optimizer

	logger *zap.Logger
}

// New client object by client
func New(source *models.RTBSource, opts ...interface{}) (*Platform, error) {
	var (
		extend      platformExt
		eventer     *platform.WinNotifier
		eventStream eventstream.Stream
		optim       *optimizer.Optimizer
		logger      *zap.Logger
	)

	for _, opt := range opts {
		switch o := opt.(type) {
		case platformExt:
			extend = o
		case *platform.WinNotifier:
			eventer = o
		case eventstream.Stream:
			eventStream = o
		case *zap.Logger:
			logger = o
		case *optimizer.Optimizer:
			optim = o
		}
	}

	switch {
	case source == nil:
		return nil, errInvalidRTBSource
	case eventer == nil:
		return nil, errInvalidWinEventStream
	case eventStream == nil:
		return nil, errInvalidEventStream
	case logger == nil:
		return nil, errInvalidLogger
	}

	platformObj := &Platform{
		Extend:         extend,
		source:         source,
		optimizator:    optim,
		metricsCounter: openlatency.NewMetricsCounter(),
		marshalel:      dataencode.NewEncoder(source.RequestType),
		logger:         logger,
	}

	if err := platformObj.preparePlatform(); err != nil {
		return nil, err
	}
	return platformObj, nil
}

// ID of platform
func (p *Platform) ID() uint64 {
	return p.source.ID
}

// CompanyID number
func (p *Platform) CompanyID() uint64 {
	return p.source.Company.ID
}

// RevenueShareReduceFactor from 0. to 1.
func (p *Platform) RevenueShareReduceFactor() float64 {
	return p.revenueShareReduce
}

// RequestStrategy description
func (p *Platform) RequestStrategy() adsource.RequestStrategy {
	return adsource.AsynchronousRequestStrategy
}

// InitFromClient data from model
func (p *Platform) preparePlatform() error {
	p.logInfo("preparePlatform", p.source.ID, p.source.URL)

	p.Headers = map[string]string{}
	p.revenueShareReduce = p.source.RevenueShareReduceFactor()

	if p.optimizator == nil {
		p.optimizator = &optimizer.Optimizer{}
		p.optimizator.Init()
	}
	if p.source.Method == "" {
		p.source.Method = "POST"
	}
	if headers := p.source.Headers.Data(); headers != nil {
		for k, v := range headers {
			p.Headers[k] = v.String
		}
	}

	// Apply extension
	if p.Extend != nil {
		if rp, ok := p.Extend.(platformRequestPreparer); ok {
			p.requestPreparer = rp
		}
		if rp, ok := p.Extend.(platformResponsePreparer); ok {
			p.responsePreparer = rp
		}
		if rm, ok := p.Extend.(dataencode.Encoder); ok {
			p.marshalel = rm
		}
		return p.Extend.InitPlatform(p)
	}
	return nil
}

// Test current request for compatibility
func (p *Platform) Test(request *adsource.BidRequest) bool {
	if p.source.RPS > 0 {
		if !p.source.Options.ErrorsIgnore && !p.errorCounter.Next() {
			return false
		}

		if now := time.Now(); now.Sub(p.lastPeriod).Seconds() >= 1 {
			p.lastPeriod = now
			p.rpsCurrent.Set(0)
		}

		if p.rpsCurrent.Get() >= int64(p.source.RPS) {
			return false
		}
	}

	if !p.source.Test(request) {
		return false
	}

	// Check formats targeting
	for _, f := range request.Formats() {
		if !p.optimizator.Test(
			uint(f.ID),
			byte(request.GeoID()),
			request.LanguageID(),
			request.DeviceID(),
			request.OSID(),
			request.BrowserID(),
			p.source.MinimalWeight,
		) {
			return false
		}
	}

	return true
}

// Bid request for standart system filter
func (p *Platform) Bid(request *adsource.BidRequest) (response adsource.Responser) {
	p.rpsCurrent.Inc(1)

	httpRequest, err := p.request(request)
	if err != nil {
		return adsource.NewErrorResponse(request, err)
	}

	var resp *http.Response
	resp, err = p.getClient().Do(httpRequest)
	if err != nil {
		p.processHTTPReponse(resp, err)
		p.logDebug("bid", p.source.URL, err)
		return adsource.NewErrorResponse(request, err)
	}

	if http.StatusNoContent == resp.StatusCode {
		return adsource.NewErrorResponse(request, platform.ErrNoCampaignsStatus)
	}

	if http.StatusOK != resp.StatusCode {
		p.processHTTPReponse(resp, nil)
		p.logDebug("bid", p.source.URL, http.StatusText(resp.StatusCode), resp.StatusCode)
		return adsource.NewErrorResponse(request, ErrInvalidResponseStatus)
	}

	defer resp.Body.Close()
	if response, err = p.unmarshal(request, resp.Body); p.source.Options.Trace && err != nil {
		response = adsource.NewErrorResponse(request, err)
		p.logError("bid", err)
	}

	p.processHTTPReponse(resp, err)
	return
}

// ProcessResponseItem result
func (p *Platform) ProcessResponseItem(response adsource.Responser, it adsource.ResponserItem) {
	if response == nil || response.Error() != nil {
		return
	}
	for _, ad := range response.Ads() {
		switch bid := ad.(type) {
		case *adsource.ResponseBidItem:
			if len(bid.Bid.NURL) > 0 {
				p.winEvents.Send(response.Context(), bid.Bid.NURL)
				p.logInfo("ping", bid.Bid.NURL)
			}
			p.bidActions.Send(events.Impression, events.StatusUndefined, response, bid)
		default:
			// Dummy...
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Implementation of platform.State interface
///////////////////////////////////////////////////////////////////////////////

// GetState of the driver object
func (p *Platform) GetState() interface{} {
	return p.optimizator
}

// SetState from driver object
func (p *Platform) SetState(state interface{}) {
	if state != nil {
		if opt, _ := state.(*optimizer.Optimizer); opt != nil {
			p.optimizator = opt
		}
	}
}

// GetStateData of the platform
func (p *Platform) GetStateData() (io.Reader, error) {
	var (
		buff = bytes.NewBuffer(nil)
		err  = json.NewEncoder(buff).Encode(p.optimizator)
	)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// LoadStateData of the platform
func (p *Platform) LoadStateData(data io.Reader) error {
	return json.NewDecoder(data).Decode(p.optimizator)
}

///////////////////////////////////////////////////////////////////////////////
/// Implementation of platform.Metrics interface
///////////////////////////////////////////////////////////////////////////////

// Metrics information of the platform
func (p *Platform) Metrics() *openlatency.MetricsInfo {
	info := &openlatency.MetricsInfo{
		ID:       p.source.ID,
		Protocol: p.source.Protocol,
		QPSLimit: p.source.RPS,
	}
	p.metricsCounter.FillMetrics(info)
	return info
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (p *Platform) getClient() *http.Client {
	if p.Client == nil {
		timeout := time.Millisecond * time.Duration(p.source.Timeout)
		if timeout < 1 {
			timeout = time.Millisecond * 150
		}
		p.Client = &http.Client{Timeout: timeout}
	}
	return p.Client
}

// prepare request for RTB
func (p *Platform) request(request *adsource.BidRequest) (req *http.Request, err error) {
	var (
		bidRequest *adsource.RTBRequest
		data       io.Reader
	)

	if bidRequest = request.RTBBidRequest(p.getRequestOptions()...); bidRequest == nil {
		return nil, nil
	}

	bidRequest.RTBRequest.TMax = int(p.getClient().Timeout / time.Millisecond)

	// Prepare data for request
	if data, err = p.marshal(bidRequest); err != nil {
		return nil, err
	}

	// Create new request
	if req, err = http.NewRequest(p.source.Method, p.source.URL, data); err != nil {
		return nil, err
	}

	p.fillRequest(request, req)
	return req, err
}

func (p *Platform) marshal(request *adsource.RTBRequest) (_ io.Reader, err error) {
	var data []byte
	if data, err = p.marshalel.Marshal(request.RTBRequest); err != nil {
		return
	}
	return bytes.NewBuffer(data), err
}

func (p *Platform) unmarshal(request *adsource.BidRequest, r io.Reader) (response *adsource.BidResponse, err error) {
	var bidResp openrtb.BidResponse

	switch p.source.RequestType {
	case RequestTypeJSON:
		if p.source.Options.Trace {
			var data []byte
			if data, err = ioutil.ReadAll(r); err == nil {
				p.logError("unmarshal", p.source.URL, string(data))
				err = json.Unmarshal(data, &bidResp)
			}
		} else {
			err = json.NewDecoder(r).Decode(&bidResp)
		}
	case RequestTypeXML, RequestTypeProtobuff:
		err = fmt.Errorf("Request body type not supported")
	default:
		err = fmt.Errorf("Undefined request type")
	}

	// Check response for support HTTPS
	if err == nil && request.IsSecure() {
		for _, seat := range bidResp.SeatBid {
			for _, bid := range seat.Bid {
				if strings.Contains(bid.AdMarkup, "http://") {
					err = ErrResponseAreNotSecure
					break
				}
			}
			if err != nil {
				break
			}
		} // end for
	}

	// Extends prepare
	if err == nil && p.responsePreparer != nil {
		err = p.responsePreparer.PrepareResponse(&bidResp, request)
	}

	// Build response
	if err == nil {
		response = &adsource.BidResponse{
			Src:         p,
			Req:         request,
			BidResponse: bidResp,
		}
		response.Prepare()
	}
	return response, err
}

// fillRequest of HTTP
func (p *Platform) fillRequest(request *adsource.BidRequest, httpReq *http.Request) {
	httpReq.Header.Set("Content-Type", p.marshalel.ContentType())
	httpReq.Header.Set(headerRequestOpenRTBVersion, headerRequestOpenRTBVersionDefault)
	httpReq.Header.Set(openlatency.HTTPHeaderRequestTimemark, strconv.FormatInt(openlatency.RequestInitTime(request.Time()), 10))

	// Fill default headers
	if p.Headers != nil {
		for key, value := range p.Headers {
			httpReq.Header.Set(key, value)
		}
	}
}

// @link https://golang.org/src/net/http/status.go
func (p *Platform) processHTTPReponse(resp *http.Response, err error) {
	switch {
	case err != nil || resp == nil ||
		(resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent):
		if err == http.ErrHandlerTimeout {

		}
		p.errorCounter.Inc()
	default:
		p.errorCounter.Dec()
	}
}

func (p *Platform) getRequestOptions() []adsource.BidRequestRTBOption {
	return []adsource.BidRequestRTBOption{
		adsource.WithRTBOpenNativeVersion("1.1"),
		adsource.WithFormatFilter(p.source.TestFormat),
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Logs
///////////////////////////////////////////////////////////////////////////////

func (p *Platform) logInfo(block string, params ...interface{}) {
	p.logger.Info(
		strings.TrimRight(fmt.Sprintln(params...), " \n\t"),
		zap.String("block", block),
	)
}

func (p *Platform) logDebug(block string, params ...interface{}) {
	p.logger.Debug(
		strings.TrimRight(fmt.Sprintln(params...), " \n\t"),
		zap.String("block", block),
	)
}

func (p *Platform) logError(block string, params ...interface{}) {
	if len(params) > 0 && params[0] != nil {
		p.logger.Error(
			strings.TrimRight(fmt.Sprintln(params...), " \n\t"),
			zap.String("block", block),
		)
	}
}

var (
	_ platform.Platformer = (*Platform)(nil)
	_ platform.State      = (*Platform)(nil)
	_ platform.Metrics    = (*Platform)(nil)
)
