package rtbevents

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

// ErrInvalidMultipleItemAsSingle invalid multiple item as single action
var ErrInvalidMultipleItemAsSingle = errors.New("can`t convert multipleitem to single action")

// Event representation of internal adserver event like impression, click, lead, etc.
type Event struct {
	Time     int64       `json:"tm,omitempty"` // Timestamp
	Delay    uint64      `json:"dl,omitempty"` // Delay of preparation of Ads in Nanosecinds
	Duration uint64      `json:"d,omitempty"`  // Duration in Nanoseconds
	Event    events.Type `json:"e,omitempty"`  // Event code string
	Status   uint8       `json:"st,omitempty"` // Status: 0 - undefined, 1 - success, 2 - failed, 3 - compromised

	// Source
	AuctionID    string `json:"auc,omitempty"`     // Internal Auction ID
	AuctionType  uint8  `json:"auctype,omitempty"` // Aution type 1 - First price, 2 - Second price
	ImpID        string `json:"imp,omitempty"`     // Sub ID of request for paticular impression spot
	ImpAdID      string `json:"impad,omitempty"`   // Specific ID for paticular ad impression
	ExtAuctionID string `json:"eauc,omitempty"`    // RTB Request/Response ID
	ExtImpID     string `json:"eimp,omitempty"`    // RTB Imp ID
	ExtTargetID  string `json:"extz,omitempty"`    // RTB Zone ID (tagid)
	SourceID     uint64 `json:"sid,omitempty"`     // Advertisement Source ID

	// State Location
	Platform      int    `json:"pl,omitempty"`  // Where displaid? 0 – undefined, 1 – web site, 2 – native app, 3 – game
	Domain        string `json:"dm,omitempty"`  // If not web site then "bundle"
	ApplicationID uint64 `json:"app,omitempty"` // application ID (registered in the system)
	AdUnitID      uint64 `json:"z,omitempty"`   // -- // --
	FormatID      uint64 `json:"fmt,omitempty"` // Format ID
	AdWidth       int    `json:"aw,omitempty"`  // -- // --
	AdHeight      int    `json:"ah,omitempty"`  // -- // --
	SourceURL     string `json:"su,omitempty"`  // Advertisement source URL (iframe, image, video, direct)
	WinURL        string `json:"wu,omitempty"`  // Win URL used for RTB confirmation
	URL           string `json:"u,omitempty"`   // Non modified target URL

	// Money section
	// We have a different situations how we are pay for the traffic
	// For example
	// 	1) If we have external traffic resource:
	//		ViewPrice = 0.01 CPV - Have to pay advertiser
	//		LeadPrice = 40$ - How much will we pay for the lead
	//		PurchaseViewPrice = 0.001 CPV - how much we will pay for this traffic
	//	2) If external second price auction type
	//		ViewPrice = 0.02 CPV - Have to pay advertiser
	//		PurchaseViewPrice = 0.01 CPV - how much we will pay for this traffic
	//  3) If we are buying the traffic for CPA
	//		LeadPrice = 20$ - Have to pay advertiser
	//		PurchaseViewPrice = 10$ - Have to pay to the source
	PricingModel uint    `json:"pm,omitempty"`   // Display As CPM/CPC/CPA/CPI
	ECPM         float64 `json:"ecpm,omitempty"` // Effective Cost per Mille
	// Prurchase price from SSP or other TRAFFIC sources (menetisation of income requests)
	// We are buying the place of advertisement display
	PurchaseViewPrice   int64   `json:"pvpr,omitempty"`   // Price of the view of source traffic cost
	PurchaseClickPrice  int64   `json:"pcpr,omitempty"`   // Price of the click of source traffic cost
	PotentialViewPrice  int64   `json:"ptvpr,omitempty"`  // Price of the view of source of ads oitential
	PotentialClickPrice int64   `json:"ptcpr,omitempty"`  // Price of the click of source of ads oitential
	ViewPrice           int64   `json:"vpr,omitempty"`    // Price per view
	ClickPrice          int64   `json:"cpr,omitempty"`    // Price per click
	CompetitorSourceID  uint64  `json:"cmsrc,omitempty"`  // Competitor source ID
	CompetitorECPM      float64 `json:"cmecpm,omitempty"` // Competitor ECPM or auction

	// User IDENTITY
	UDID        string `json:"udi,omitempty"` // Unique Device ID (IDFA)
	UUID        string `json:"uui,omitempty"` // User
	SessionID   string `json:"ses,omitempty"` // -- // --
	Fingerprint string `json:"fpr,omitempty"` // -- // --
	ETag        string `json:"etg,omitempty"` // -- // --

	// Targeting
	CarrierID       uint   `json:"car,omitempty"`  // -- // --
	Country         string `json:"cc,omitempty"`   // Country Code ISO-2
	Latitude        string `json:"lt,omitempty"`   // -- // --
	Longitude       string `json:"lg,omitempty"`   // -- // --
	Language        string `json:"lng,omitempty"`  // en-US
	IPString        string `json:"ip,omitempty"`   // -- // --
	Referer         string `json:"ref,omitempty"`  // Referal link
	Page            string `json:"page,omitempty"` // Page link
	UserAgent       string `json:"ua,omitempty"`   // -- // --
	DeviceType      uint   `json:"dvt,omitempty"`  // Device type 0 - Undefined, 1 - Desktop, etc.
	DeviceID        uint   `json:"dv,omitempty"`   // Device ID
	OSID            uint   `json:"os,omitempty"`   // OS ID
	BrowserID       uint   `json:"br,omitempty"`   // Browser ID
	Categories      string `json:"c,omitempty"`    // -- // --
	Adblock         uint   `json:"ab,omitempty"`   // Ext params
	PrivateBrowsing uint   `json:"prv,omitempty"`  // Private Mode
	Robot           uint   `json:"rt,omitempty"`   // Robot traffic
	Proxy           uint   `json:"pt,omitempty"`   // Proxy traffic
	Backup          uint   `json:"bt,omitempty"`   // Backup Display Type
	X               int    `json:"x,omitempty"`    // X - coord of addisplay or click position
	Y               int    `json:"y,omitempty"`    // Y - coord of addisplay or click position
	Width           int    `json:"w,omitempty"`    // Width - available space
	Height          int    `json:"h,omitempty"`    // Height - available space

	SubID1 string `json:"sd1,omitempty"`
	SubID2 string `json:"sd2,omitempty"`
	SubID3 string `json:"sd3,omitempty"`
	SubID4 string `json:"sd4,omitempty"`
	SubID5 string `json:"sd5,omitempty"`
}

// AllocateEvent object
func AllocateEvent() *Event {
	return &Event{}
}

// String imprlementation of Stringer interface
func (event *Event) String() string {
	if json, err := json.MarshalIndent(event, "", "  "); err == nil {
		return string(json)
	}
	return ""
}

// EventType of event
func (event *Event) EventType() events.Type {
	return event.Event
}

// EventURL of event return URL of event target
func (event *Event) EventURL() string {
	return event.URL
}

// SetEventPurchaseViewPrice set price of the event
func (event *Event) SetEventPurchaseViewPrice(price int64) error {
	if event.ViewPrice < price {
		return fmt.Errorf("price is higher than view price %f < %f",
			billing.Money(event.ViewPrice).Float64(), billing.Money(price).Float64())
	}
	event.PurchaseViewPrice = price
	return nil
}

// Fill event object from response and Ad item objects
func (event *Event) Fill(service string, eventType events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) error {
	var (
		r            = response.Request()
		imp          = it.Impression()
		sourceID     uint64
		targetSpotID uint64
	)

	if src := it.Source(); src != nil {
		sourceID = src.ID()
	}

	if sourceID == 0 && response.Source() != nil {
		sourceID = response.Source().ID()
	}

	if imp != nil && imp.Target != nil {
		targetSpotID = imp.Target.ID()
	}

	if _, ok := it.(adtype.ResponserMultipleItem); ok {
		return ErrInvalidMultipleItemAsSingle
	}

	*event = Event{
		Time:     time.Now().UnixNano(),
		Delay:    0,
		Duration: 0,
		Event:    eventType, // Action code (tech param, Do not store)
		Status:   status,

		// Source
		AuctionID:    r.ID,                          // ID of last auction
		AuctionType:  uint8(response.AuctionType()), // Aution type 1 - First price, 2 - Second price
		ImpID:        it.ImpressionID(),             // Sub ID of request for paticular impression spot
		ImpAdID:      it.ID(),                       // Specific ID for paticular ad impression
		ExtAuctionID: r.ExtID,                       // External auction ID
		ExtImpID:     it.ExtImpressionID(),          // External auction Imp ID
		SourceID:     sourceID,                      // Advertisement Source ID

		// State Location
		Platform:      0,                           // Where displaid? 0 – undefined, 1 – web site, 2 – native app, 3 – game
		Domain:        r.DomainName(),              //
		ApplicationID: r.AppID(),                   // Place target
		AdUnitID:      targetSpotID,                // -- // --
		FormatID:      it.Format().ID,              // Format object ID
		AdWidth:       positiveNumber(it.Width()),  // -- // --
		AdHeight:      positiveNumber(it.Height()), // -- // --
		SourceURL:     "",                          // Advertisement source URL (iframe, image, video, direct)
		WinURL:        "",                          // Win URL used for RTB confirmation
		URL:           it.ActionURL(),              // Non modified target URL

		// Money
		PricingModel:        it.PricingModel().UInt(),                    // Display As CPM/CPC/CPA/CPI
		ECPM:                it.ECPM().Float64(),                         // Effective Cost Per Mille (1000 views)
		PurchaseViewPrice:   it.PurchasePrice(adtype.ActionView).I64(),   // Price of of the view of source traffic cost
		PurchaseClickPrice:  it.PurchasePrice(adtype.ActionClick).I64(),  // Price of of the click of source traffic cost
		PotentialViewPrice:  it.PotentialPrice(adtype.ActionView).I64(),  // Price of of the view of source traffic cost including descrepancy correction
		PotentialClickPrice: it.PotentialPrice(adtype.ActionClick).I64(), // Price of of the click of source traffic cost including descrepancy correction
		ViewPrice:           it.FinalPrice(adtype.ActionView).I64(),      // Price per view with total comissions and with descrepancy correction
		ClickPrice:          it.FinalPrice(adtype.ActionClick).I64(),     // Price per click with total comissions and with descrepancy correction
		CompetitorSourceID:  it.Second().GetSourceID(),                   // Competitor source ID
		CompetitorECPM:      it.Second().GetECPM().Float64(),             // Competitor ECPM or auction

		// User IDENTITY
		UDID:        r.DeviceInfo().IFA,         // Unique Device ID (IDFA)
		UUID:        r.UserInfo().ID,            // User
		SessionID:   r.UserInfo().SessionID,     // -- // --
		Fingerprint: r.UserInfo().FingerPrintID, //
		ETag:        r.UserInfo().ETag,          //

		// Targeting
		CarrierID:       r.CarrierInfo().ID,
		Country:         r.GeoInfo().Country,
		Language:        r.BrowserInfo().PrimaryLanguage,
		Referer:         r.BrowserInfo().Ref,
		IPString:        r.GeoInfo().IP.String(),
		UserAgent:       r.BrowserInfo().UA,
		DeviceID:        r.DeviceInfo().ID,
		OSID:            r.DeviceInfo().OS.ID,
		BrowserID:       uint(r.BrowserInfo().ID),
		Categories:      "",
		Adblock:         b2u(r.IsAdblock()),
		PrivateBrowsing: b2u(r.IsPrivateBrowsing()),
		Robot:           b2u(r.IsRobot()),
		Proxy:           b2u(r.IsProxy()),
		Backup:          b2u(it.IsBackup()),
		X:               positiveNumber(imp.X),
		Y:               positiveNumber(imp.Y),
		Width:           positiveNumber(r.Width()),
		Height:          positiveNumber(r.Height()),

		SubID1: imp.SubID1,
		SubID2: imp.SubID2,
		SubID3: imp.SubID3,
		SubID4: imp.SubID4,
		SubID5: imp.SubID5,
	}

	return nil
}

// Validate event object
func (event *Event) Validate() error {
	return nil
}

// DateTime object
func (event *Event) DateTime() time.Time {
	return time.Unix(0, event.Time)
}

// SetDateTime timestamp
func (event *Event) SetDateTime(tmNano int64) {
	event.Time = tmNano
}

// Pack object to event Code
func (event *Event) Pack() events.Code {
	return events.ObjectCode(event, streamCodeEncoder)
}

// Unpack event object
func (event *Event) Unpack(data []byte, unpuckFnc ...events.EventUnpacFunc) error {
	code := events.CodeObj(data, nil)
	if len(unpuckFnc) > 0 && unpuckFnc[0] != nil {
		code = unpuckFnc[0](code)
	}
	return code.DecodeObject(event, streamCodeDecoder)
}

// PrepareURL by event
func (event *Event) PrepareURL(url string) string {
	if !strings.Contains(url, "{") || !strings.Contains(url, "}") {
		return url
	}
	replacer := strings.NewReplacer(
		"{country}", event.Country,
		"{lang}", event.Language,
		"{domain}", event.Domain,
		"{impid}", event.ImpID,
		"{aucid}", event.AuctionID,
		"{auctype}", types.AuctionType(event.AuctionType).Name(),
		"{platform}", types.PlatformType(event.Platform).Name(),
		"{zone_id}", gocast.Str(event.AdUnitID),
		"{adunit_id}", gocast.Str(event.AdUnitID),
		"{pm}", types.PricingModel(event.PricingModel).Name(),
		"{udid}", event.UDID,
		"{uuid}", event.UUID,
		"{sessid}", event.SessionID,
		"{fingerprint}", event.Fingerprint,
		"{etag}", event.ETag,
		"{ip}", event.IPString,
		"{carrier_id}", gocast.Str(event.CarrierID),
		"{latitude}", event.Latitude,
		"{longitude}", event.Longitude,
		"{device_type}", types.PlatformType(event.DeviceType).Name(),
		"{device_id}", gocast.Str(event.DeviceID),
		"{os_id}", gocast.Str(event.OSID),
		"{browser_id}", gocast.Str(event.BrowserID),
	)
	return replacer.Replace(url)
}

// Release event object
func (event *Event) Release() {}
