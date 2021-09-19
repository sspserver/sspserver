//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package eventgenerator

import (
	"errors"
	"strconv"
	"time"

	"github.com/demdxx/gocast"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/events"
	"geniusrabbit.dev/sspserver/internal/models"
)

// Errors set
var (
	ErrInvalidMultipleItemAsSingle = errors.New("Cant convert multipleitem to single action")
)

// Generator object
type Generator interface {
	// Event object by response
	Event(event events.Type, status uint8, response adsource.Responser, it adsource.ResponserItem) (*events.Event, error)

	// Events object list
	Events(event events.Type, status uint8, response adsource.Responser, it adsource.ResponserItemCommon) []*events.Event

	// UserInfo event object by response
	UserInfo(response adsource.Responser, it adsource.ResponserItem) (*events.UserInfo, error)
}

type generator struct {
	service string
}

// New generator object
func New(service string) Generator {
	return generator{service: service}
}

// Event object by response
func (g generator) Event(event events.Type, status uint8, response adsource.Responser, it adsource.ResponserItem) (*events.Event, error) {
	var (
		r             = response.Request()
		imp           = it.Impression()
		sourceID      uint64
		accessPointID uint64
		zoneID        uint64
	)

	if src := it.Source(); src != nil {
		sourceID = src.ID()
	}

	if response.AccessPoint() != nil {
		accessPointID = response.AccessPoint().ID()
	}

	if imp != nil && imp.Target != nil {
		zoneID = imp.Target.ID()
	}

	if _, ok := it.(adsource.ResponserMultipleItem); ok {
		return nil, ErrInvalidMultipleItemAsSingle
	}

	return &events.Event{
		Time:     time.Now().UnixNano(),
		Delay:    0,
		Duration: 0,         //
		Service:  g.service, // Service
		Event:    event,     // Action code (tech param, Do not store)
		Status:   status,    //
		// Accounts link information
		Project:           0,               // Project network ID
		PublisherCompany:  imp.CompanyID(), // -- // --
		AdvertiserCompany: it.CompanyID(),  // -- // --
		// Source
		AuctionID:    r.ID,                          // ID of last auction
		AuctionType:  uint8(response.AuctionType()), // Aution type 1 - First price, 2 - Second price
		ImpID:        it.ImpressionID(),             // Sub ID of request for paticular impression spot
		ImpAdID:      it.ID(),                       // Specific ID for paticular ad impression
		ExtAuctionID: r.ExtID,                       // External auction ID
		ExtImpID:     it.ExtImpressionID(),          // External auction Imp ID
		Source:       sourceID,                      // Advertisement Source ID
		Network:      it.NetworkName(),              // Source Network Name or Domain (Cross sails)
		AccessPoint:  accessPointID,                 // Access Point ID to own Advertisement
		// State Location
		Platform:    0,                 // Where displaid? 0 – undefined, 1 – web site, 2 – native app, 3 – game
		Domain:      r.DomainName(),    //
		Application: uint64(r.AppID()), // Place target
		Zone:        zoneID,            // -- // --
		Campaign:    it.CampaignID(),   // Campaign info
		FormatID:    it.Format().ID,    // Format object ID
		AdID:        it.AdID(),         // -- // --
		AdW:         it.Width(),        // -- // --
		AdH:         it.Height(),       // -- // --
		SourceURL:   "",                // Advertisement source URL (iframe, image, video, direct)
		WinURL:      "",                // Win URL used for RTB confirmation
		URL:         it.ActionURL(),    // Non modified target URL
		Jumper:      0,                 // Jumper Page ID
		// Money
		PricingModel:       it.PricingModel().UInt(),                   // Display As CPM/CPC/CPA/CPI
		PurchaseViewPrice:  imp.PurchasePrice(models.ActionImpression), // Price of of the view of source traffic cost
		PurchaseClickPrice: imp.PurchasePrice(models.ActionClick),      // Price of of the click of source traffic cost
		PurchaseLeadPrice:  imp.PurchasePrice(models.ActionLead),       // Price of of the lead of source traffic cost
		ViewPrice:          it.Price(models.ActionImpression),          // Price per view
		ClickPrice:         it.Price(models.ActionClick),               // Price per click
		LeadPrice:          it.Price(models.ActionLead),                // Price per lead
		Competitor:         it.Second().GetCampaignID(),                // Competitor compaign ID
		CompetitorSource:   it.Second().GetSourceID(),                  // Competitor source ID
		CompetitorECPM:     it.Second().GetECPM(),                      // Competitor ECPM or auction
		Revenue:            int(it.Revenue() * 100),                    // Should be percent in money type
		Potential:          int(it.Potential() * 100),                  // Potential revenue in %s
		// User IDENTITY
		UDID:        r.DeviceInfo().IFA,         // Unique Device ID (IDFA)
		UUID:        r.UserInfo().ID,            // User
		SessionID:   r.UserInfo().SessionID,     // -- // --
		Fingerprint: r.UserInfo().FingerPrintID, //
		ETag:        r.UserInfo().ETag,          //
		// Targeting
		Carrier:         r.CarrierInfo().ID,
		Country:         r.GeoInfo().Country,
		Language:        r.BrowserInfo().PrimaryLanguage,
		Referer:         r.BrowserInfo().Ref,
		IPString:        r.GeoInfo().IP.String(),
		UserAgent:       r.BrowserInfo().UA,
		Device:          r.DeviceInfo().ID,
		OS:              r.DeviceInfo().OS.ID,
		Browser:         r.BrowserInfo().ID,
		Categories:      "",
		Adblock:         b2u(r.IsAdblock()),
		PrivateBrowsing: b2u(r.IsPrivateBrowsing()),
		Robot:           0,
		Proxy:           0,
		Backup:          0,
		X:               imp.X,
		Y:               imp.Y,
		W:               r.Width(),
		H:               r.Height(),

		SubID1: imp.SubID1,
		SubID2: imp.SubID2,
		SubID3: imp.SubID3,
		SubID4: imp.SubID4,
		SubID5: imp.SubID5,
	}, nil
}

// Events object list
func (g generator) Events(event events.Type, status uint8, response adsource.Responser, it adsource.ResponserItemCommon) (events []*events.Event) {
	var items []adsource.ResponserItem

	if mit, _ := it.(adsource.ResponserMultipleItem); items != nil {
		for _, it := range mit.Ads() {
			if event, err := g.Event(event, status, response, it); err == nil {
				events = append(events, event)
			}
		}
	} else if event, err := g.Event(event, status, response, it.(adsource.ResponserItem)); err == nil {
		events = append(events, event)
	}
	return
}

// UserInfo event object by response
func (g generator) UserInfo(response adsource.Responser, it adsource.ResponserItem) (*events.UserInfo, error) {
	var (
		r       = response.Request()
		imp     = it.Impression()
		user    = r.UserInfo()
		geo     = r.GeoInfo()
		browser = r.BrowserInfo()
	)
	if user.Email == "" {
		return nil, nil
	}
	return &events.UserInfo{
		Time:      time.Now().UnixNano(),
		AuctionID: r.ID, // ID of last auction
		// User IDENTITY
		UDID:      r.DeviceInfo().IFA, // Unique Device ID (IDFA)
		UUID:      user.ID,            // User
		SessionID: user.SessionID,     // -- // --
		// Personal information
		Age:           user.AvgAge(),
		Gender:        byte(user.Sex()),
		SearchGender:  sex(gocast.ToString(imp.Get("search_gender"))),
		Email:         user.Email,
		Phone:         user.GetDataItemOrDefault("phone", ""),
		MessangerType: user.GetDataItemOrDefault("messanger_type", ""),
		Messanger:     user.GetDataItemOrDefault("messanger", ""),
		Postcode:      geo.Zip,
		Facebook:      user.GetDataItemOrDefault("sn.facebook", ""),
		Twitter:       user.GetDataItemOrDefault("sn.twitter", ""),
		Linkedin:      user.GetDataItemOrDefault("sn.linkedin", ""),
		// Location info
		Country:   geo.Country,                               // Country Code ISO-2
		City:      geo.City,                                  // City Code
		Latitude:  strconv.FormatFloat(geo.Lat, 'G', -1, 64), // -- // --
		Longitude: strconv.FormatFloat(geo.Lon, 'G', -1, 64), // -- // --
		Language:  browser.PrimaryLanguage,                   // en-US
	}, nil
}

func b2u(v bool) uint {
	if v {
		return 1
	}
	return 0
}

func sex(s string) (sx byte) {
	switch s {
	case "male", "m", "M":
		sx = byte(adsource.UserSexMale)
	case "female", "f", "F":
		sx = byte(adsource.UserSexFemale)
	}
	return
}
