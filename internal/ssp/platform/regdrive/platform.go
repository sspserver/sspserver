package regdrive

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/notificationcenter"
	log "github.com/sirupsen/logrus"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/events"
	"geniusrabbit.dev/sspserver/internal/eventstream"
	"geniusrabbit.dev/sspserver/internal/models"
	counter "geniusrabbit.dev/sspserver/internal/ssp/errorcounter"
	"geniusrabbit.dev/sspserver/internal/ssp/openlatency"
	"geniusrabbit.dev/sspserver/internal/ssp/platform"
	"geniusrabbit.dev/sspserver/internal/ssp/platform/dataencode"
)

var (
	requestKeys = []string{
		"imp.id",
		"track_code", "imp.extid",
		"request.id",
		"request.extid",
		"sub_id", "imp.subid",
		"partner.id",
		"browser.referrer_url",
		"user.firstname",
		"user.lastname",
		"user.username",
		"user.email",
		"user.password",
		"user.age",
		"user.birthday",
		"user.birthday.year",
		"user.birthday.month",
		"user.birthday.day",
		"user.gender",
		"user.search_gender",
		"ua",
		"ip",
		"geo.country",
		"geo.city",
		"geo.zip",
	}
)

const (
	headerRequestVersion = "X-" + protocol + "-Version"
)

// Response items list...
const (
	ResponseItemLink = "link"
)

type platformDriver struct {
	// Original source model
	source *models.RTBSource

	// Contains mapping or external fields into internal format
	// * This parameter is optional, beside in some case we can reaceive just redirect or PLAINTEXT
	requestMapping  map[string]MappingItem
	requestDefaults map[string]interface{}
	responseMapping map[string]string
	responseCheck   ResponseCheck

	// contains information per one impression
	defaultPricePerImpression billing.Money

	// requestEncoder represents data processing encoder
	requestEncoder dataencode.Encoder

	// Requests RPS counter
	rpsCurrent     counter.Counter
	errorCounter   counter.ErrorCounter
	metricsCounter *openlatency.MetricsCounter

	// Revenue share this is the rate of publisher revenue
	revenueShareReduce float64

	// win events stream
	winEvents notificationcenter.Publisher

	// bid action stream
	bidActions eventstream.Stream

	// Client of HTTP requests
	httpClient *http.Client

	// Datetime of last first request
	lastPeriod time.Time

	logger log.FieldLogger
}

// New client object by client
func new(source *models.RTBSource, opts ...interface{}) (_ *platformDriver, err error) {
	var (
		eventer     notificationcenter.Publisher
		eventStream eventstream.Stream
		logger      *log.Entry
	)

	for _, opt := range opts {
		switch o := opt.(type) {
		case notificationcenter.Publisher:
			eventer = o
		case eventstream.Stream:
			eventStream = o
		case *log.Entry:
			logger = o
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

	var conf struct {
		RequestMapping  map[string]MappingItem `json:"request_mapping"`
		RequestDefaults map[string]interface{} `json:"request_defaults"`
		ResponseMapping map[string]string      `json:"response_mapping"`
		ResponseCheck   ResponseCheck          `json:"response_check"`
		DefaultPrice    billing.Money          `json:"default_price"`
	}

	if err = source.Config.UnmarshalTo(&conf); err != nil {
		return nil, err
	}

	return &platformDriver{
		source:                    source,
		requestMapping:            conf.RequestMapping,
		requestDefaults:           conf.RequestDefaults,
		responseMapping:           conf.ResponseMapping,
		responseCheck:             conf.ResponseCheck,
		defaultPricePerImpression: conf.DefaultPrice,
		requestEncoder:            dataencode.NewEncoder(source.RequestType),
		metricsCounter:            openlatency.NewMetricsCounter(),
		winEvents:                 eventer,
		bidActions:                eventStream,
		logger:                    logger,
	}, nil
}

// ID of platform
func (p *platformDriver) ID() uint64 {
	return p.source.ID
}

// CompanyID number
func (p *platformDriver) CompanyID() uint64 {
	return p.source.Company.ID
}

// RevenueShareReduceFactor from 0. to 1.
func (p *platformDriver) RevenueShareReduceFactor() float64 {
	return p.revenueShareReduce
}

// RequestStrategy description
func (p *platformDriver) RequestStrategy() adsource.RequestStrategy {
	return adsource.SingleRequestStrategy
}

// Test current request for compatibility
func (p *platformDriver) Test(request *adsource.BidRequest) bool {
	if !request.FormatTypeMask().HasOneType().IsAutoregisterType() {
		return false
	}

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
	return true
}

// SetTimeout for sourcer
func (p *platformDriver) SetTimeout(timeout time.Duration) {
	p.source.Timeout = int(timeout / time.Millisecond)
}

// Bid request for standart system filter
func (p *platformDriver) Bid(request *adsource.BidRequest) (response adsource.Responser) {
	p.rpsCurrent.Inc(1)

	var (
		err         error
		resp        *http.Response
		httpRequest *http.Request
		reqImp      = &request.Imps[0]
	)
	if httpRequest, err = p.request(request, reqImp); err != nil {
		return adsource.NewErrorResponse(request, err)
	}

	if resp, err = p.getHTTPClient().Do(httpRequest); err != nil {
		switch r := err.(type) {
		case *redirectError:
			err = nil
			resp = &http.Response{
				StatusCode: http.StatusFound,
				Request:    httpRequest,
				Body:       newReadCloserBuffer([]byte(r.TargetURL)),
			}
		default:
			p.processHTTPReponse(resp, err)
			p.logDebug("bid", httpRequest.URL.String(), err)
			return adsource.NewErrorResponse(request, err)
		}
	}

	if resp.StatusCode == http.StatusNoContent {
		return adsource.NewErrorResponse(request, platform.ErrNoCampaignsStatus)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		p.processHTTPReponse(resp, nil)
		p.logDebug("bid", httpRequest.URL.String(), http.StatusText(resp.StatusCode), resp.StatusCode)
		return adsource.NewErrorResponse(request, errInvalidResponseStatus)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if response, err = p.unmarshal(request, reqImp, resp); err != nil {
		response = adsource.NewErrorResponse(request, err)
		if p.source.Options.Trace {
			p.logError("bid", err)
		}
	}

	p.processHTTPReponse(resp, err)
	return
}

// ProcessResponseItem result of the item or Error processing
// * General various errors and invalid responses processed on the top level, so here need to process only driver specific errors and responses
func (p *platformDriver) ProcessResponseItem(response adsource.Responser, it adsource.ResponserItem) {
	if response == nil || response.Error() != nil {
		return
	}

	switch ad := it.(type) {
	case adsource.ResponserItem:
		if nurl := ad.ContentItemString(adsource.ContentItemNotifyWinURL); len(nurl) > 0 {
			p.logWinURL(response.Context(), nurl)
			p.logInfo("ping", nurl)
		}
		p.bidActions.Send(events.Impression, events.StatusUndefined, response, ad)
	case adsource.ResponserMultipleItem:
		for _, ad := range ad.Ads() {
			if nurl := ad.ContentItemString(adsource.ContentItemNotifyWinURL); len(nurl) > 0 {
				p.logWinURL(response.Context(), nurl)
				p.logInfo("ping", nurl)
			}
			p.bidActions.Send(events.Impression, events.StatusUndefined, response, ad)
		}
	default:
		// ! something wrong happend! It's unappropriate response
		// TODO add the metrics and log
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Internal driver methods
///////////////////////////////////////////////////////////////////////////////

func (p *platformDriver) getHTTPClient() *http.Client {
	if p.httpClient == nil {
		timeout := time.Millisecond * time.Duration(p.source.Timeout)
		if timeout < 1 {
			timeout = time.Millisecond * 100
		}
		p.httpClient = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			// CheckRedirect: p.redirectPrevent,
		}
	}
	return p.httpClient
}

func (p *platformDriver) request(request *adsource.BidRequest, imp *adsource.Impression) (httpReq *http.Request, err error) {
	var requestURL string
	if requestURL, err = p.urlPrepare(request, imp, p.source.Method != http.MethodGet); err != nil {
		return
	}

	switch p.source.Method {
	case http.MethodGet, "":
		httpReq, err = http.NewRequest(http.MethodGet, requestURL, nil)
	case http.MethodPost, http.MethodPut:
		var (
			data   []byte
			params map[string]string
		)
		if params, err = p.requestDataMap(request, imp); err != nil {
			return
		}
		if params == nil || len(params) < 1 {
			httpReq, err = http.NewRequest(p.source.Method, requestURL, nil)
		} else if data, err = p.requestEncoder.Marshal(params); err == nil {
			httpReq, err = http.NewRequest(p.source.Method, requestURL, bytes.NewBuffer(data))
		}
	default:
		// TODO add log and metrics
		err = fmt.Errorf(protocol+": invalid request method [%s], supports only [GET, POST, PUT]", p.source.Method)
	}
	if err == nil {
		if p.source.Method != http.MethodGet {
			httpReq.Header.Set("Content-Type", p.requestEncoder.ContentType())
		}
		httpReq.Header.Set(headerRequestVersion, protocolVersion)
		httpReq.Header.Set(openlatency.HTTPHeaderRequestTimemark, strconv.FormatInt(openlatency.RequestInitTime(request.Time()), 10))

		// Fill default headers
		if p.source.Headers.Data() != nil {
			for key, value := range p.source.Headers.Data() {
				httpReq.Header.Set(key, value.String)
			}
		}
	}

	return
}

func (p *platformDriver) unmarshal(request *adsource.BidRequest, imp *adsource.Impression, resp *http.Response) (response adsource.Responser, err error) {
	if resp.StatusCode == http.StatusFound {
		if resp.Body == nil {
			return nil, errInvalidRedirectBody
		}
	}

	var (
		contentType = strings.ToLower(resp.Header.Get("Content-Type"))
		data        map[string]interface{}
		data2       map[string]interface{}
	)

	switch {
	case contentType == "" || strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/javascript") || strings.HasSuffix(contentType, "/json"):
		err = json.NewDecoder(resp.Body).Decode(&data)
	case strings.HasPrefix(contentType, "application/xml") || strings.HasPrefix(contentType, "text/xml") || strings.HasSuffix(contentType, "/xml"):
		err = xml.NewDecoder(resp.Body).Decode(&data)
	case strings.HasPrefix(contentType, "plain/text") || strings.HasPrefix(contentType, "text/html"):
		var urlData []byte
		if urlData, err = ioutil.ReadAll(resp.Body); err == nil {
			// NOTE this is the hack as many tools generates invalid content type response.
			// They are assume that integration platform must aware that the response was JSON
			// It's the legacy of incorrectly writed code
			if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
				err = nil
				data = map[string]interface{}{
					ResponseItemLink: strings.TrimSpace(string(urlData)),
				}
			}
		}
	default:
		err = errors.New(errInvalidResponseContentType.Error() + " not " + contentType)
	}

	if err != nil {
		return nil, err
	}

	if err = p.responseCheck.Validate(data); err != nil {
		return nil, err
	}

	if data2 = p.responseDataMap(data); data2 == nil || len(data2) < 1 {
		return nil, newDataError(data, errInvalidResponseMapping)
	}

	if v, ok := data2[ResponseItemLink]; !ok || gocast.IsEmpty(v) {
		return nil, newDataError(data2, errInvalidResponseData)
	}

	response = adsource.NewResponse(request, p, []adsource.ResponserItemCommon{
		&ResponseItem{
			imp:      imp,
			link:     gocast.ToString(data2[ResponseItemLink]),
			data:     data2,
			bidPrice: p.defaultPricePerImpression,
			source:   p,
		},
	}, nil)

	return
}

// @link https://golang.org/src/net/http/status.go
func (p *platformDriver) processHTTPReponse(resp *http.Response, err error) {
	switch {
	case err != nil || resp == nil ||
		(resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent):
		if err == http.ErrHandlerTimeout {
			p.metricsCounter.IncTimeout()
		}
		p.errorCounter.Inc()
	default:
		p.errorCounter.Dec()
	}
}

func (p *platformDriver) urlPrepare(request *adsource.BidRequest, imp *adsource.Impression, post bool) (string, error) {
	data, params, err := p.requestOriginalDataMap(request, imp, post)
	if err != nil {
		return "", err
	}
	sourceURL := p.source.URL
	if params != nil && len(params) > 0 {
		if strings.ContainsRune(sourceURL, '?') {
			sourceURL += "&" + params.Encode()
		} else {
			sourceURL += "?" + params.Encode()
		}
	}
	return strings.NewReplacer(data...).Replace(sourceURL), nil
}

func (p *platformDriver) requestDataMap(request *adsource.BidRequest, imp *adsource.Impression) (req map[string]string, err error) {
	req = map[string]string{}
	for key, targetValue := range p.requestMapping {
		v := p.itemByRequest(key, request, imp)
		if v, err = targetValue.PrepareValue(v); err != nil {
			return nil, err
		}
		if v != "" {
			req[targetValue.Target] = v
		}
	}
	for key, val := range p.requestDefaults {
		if _, ok := req[key]; !ok {
			req[key] = gocast.ToString(val)
		}
	}
	return req, err
}

func (p *platformDriver) requestOriginalDataMap(request *adsource.BidRequest, imp *adsource.Impression, post bool) (req []string, params url.Values, err error) {
	req = []string{}
	for _, key := range requestKeys {
		if !strings.Contains(p.source.URL, "{"+key+"}") {
			if post {
				continue
			}
			if item, ok := p.requestMapping[key]; ok {
				v := p.itemByRequest(key, request, imp)
				if v, err = item.PrepareValue(v); err != nil {
					return
				}
				if v != "" {
					if params == nil {
						params = url.Values{}
					}
					params.Set(item.Target, v)
				}
			}
			continue
		}

		v := p.itemByRequest(key, request, imp)
		if item, ok := p.requestMapping[key]; ok {
			if v, err = item.PrepareValue(v); err != nil {
				return nil, nil, err
			}
		}
		req = append(req, "{"+key+"}", url.QueryEscape(v))
	}
	return
}

func (p *platformDriver) itemByRequest(key string, request *adsource.BidRequest, imp *adsource.Impression) (v string) {
	var (
		user    = request.UserInfo()
		geo     = request.GeoInfo()
		browser = request.BrowserInfo()
	)
	switch key {
	case "imp.id":
		v = imp.ID
	case "track_code", "imp.extid":
		v = imp.ExtID
	case "request.id":
		v = request.ID
	case "request.extid":
		v = request.ExtID
	case "sub_id", "imp.subid":
		v = imp.SubID1
	case "partner.id":
		v = strconv.FormatUint(p.source.ID, 10)
	case "browser.referrer_url":
		v = browser.Ref
	case "browser.ua", "ua":
		v = browser.UA
	case "user.username":
		v = user.Username
	case "user.firstname":
		v = gocast.ToString(imp.Get("firstname"))
	case "user.lastname":
		v = gocast.ToString(imp.Get("lastname"))
	case "user.email":
		v = user.Email
	case "user.password":
		v = gocast.ToString(imp.Get("new_password"))
	case "user.age":
		if user.AgeStart > 0 {
			v = strconv.Itoa(user.AgeStart)
		}
	case "user.birthday":
		v = user.Birthday
	case "user.birthday.year":
		if user.Birthday != "" {
			v = user.BirthdayTime().Format("2006")
		}
	case "user.birthday.month":
		if user.Birthday != "" {
			v = user.BirthdayTime().Format("01")
		}
	case "user.birthday.day":
		if user.Birthday != "" {
			v = user.BirthdayTime().Format("02")
		}
	case "user.gender":
		v = user.Gender
	case "user.search_gender":
		v = gocast.ToString(imp.Get("search_gender"))
	case "geo.ip", "ip":
		v = geo.IP.String()
	case "geo.country":
		if geo.Country != "**" && geo.Country != "A1" {
			v = geo.Country
		}
	case "geo.city":
		v = geo.City
	case "geo.zip":
		v = geo.Zip
	}
	return v
}

func (p *platformDriver) responseDataMap(data map[string]interface{}) (resp map[string]interface{}) {
	if data == nil {
		return nil
	}
	if p.responseMapping == nil || len(p.responseMapping) < 1 {
		resp, _ = gocast.ToSiMap(data, "", false)
		return resp
	}

	resp = map[string]interface{}{}
	for to, key := range p.responseMapping {
		if v, ok := elementByKey(key, data); ok {
			resp[to] = v
		}
	}

	return
}

func (p *platformDriver) redirectPrevent(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		return nil
	}
	return newRedirectError(req.URL.String())
}

///////////////////////////////////////////////////////////////////////////////
/// Logs
///////////////////////////////////////////////////////////////////////////////

func (p *platformDriver) logWinURL(ctx context.Context, url string) {
	if p.winEvents != nil {
		p.winEvents.Publish(ctx, adsource.WinEvent{URL: url})
	}
}

func (p *platformDriver) logInfo(block string, params ...interface{}) {
	p.logger.WithField("block", block).Infoln(
		strings.TrimRight(fmt.Sprintln(params...), " \n\t"),
	)
}

func (p *platformDriver) logDebug(block string, params ...interface{}) {
	p.logger.WithField("block", block).Debugln(
		strings.TrimRight(fmt.Sprintln(params...), " \n\t"),
	)
}

func (p *platformDriver) logError(block string, params ...interface{}) {
	if len(params) > 0 && params[0] != nil {
		p.logger.WithField("block", block).Error(
			strings.TrimRight(fmt.Sprintln(params...), " \n\t"),
		)
	}
}

var (
	_ platform.Platformer = (*platformDriver)(nil)
)
