package rtbevents

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/stretchr/testify/assert"
)

func TestNumberString(t *testing.T) {
	var tests = []struct {
		name   string
		value  string
		def    string
		target string
	}{
		{
			name:   "normal_value",
			value:  "13414",
			def:    "0",
			target: "13414",
		},
		{
			name:   "value_has_tail",
			value:  "13414 \t\n",
			def:    "0",
			target: "13414",
		},
		{
			name:   "value_has_prefix",
			value:  " \t\n13414",
			def:    "0",
			target: "13414",
		},
		{
			name:   "invalid_value",
			value:  "13414d",
			def:    "0",
			target: "0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if v := fixNumberString(test.value, test.def); v != test.target {
				t.Errorf("Invalid number string (%s) should be (%s)", test.value, test.target)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	t.Run("standart", func(t *testing.T) {
		var (
			e1 = newEvent()
			e2 Event
		)
		code := e1.Pack().Compress().URLEncode()
		assert.NoError(t, code.ErrorObj())
		assert.NoError(t, e2.Unpack(code.Data(), func(code events.Code) events.Code {
			return code.URLDecode().Decompress()
		}))
		t.Log("Code size:", len(code.Data()))

		if !reflect.DeepEqual(e1, e2) {
			t.Error("invalid code decode")
		}
	})
}

func newEvent() Event {
	return Event{
		Time:     time.Now().UnixNano(),
		Delay:    1,
		Duration: 1,
		Event:    events.SourceWin,
		Status:   1,
		// Source
		AuctionID:    "1234-123456-123456-1234",
		ImpID:        "1234-123456-123456-1234",
		ImpAdID:      "1234-123456-123456-1234",
		ExtAuctionID: "1234-123456-123456-1234",
		ExtImpID:     "1234-123456-123456-1234",
		ExtTargetID:  "codename",
		SourceID:     1,
		// State Location
		Platform:      1,
		Domain:        "domain.com",
		ApplicationID: 11,
		ZoneID:        12,
		AdWidth:       101,
		AdHeight:      22,
		SourceURL:     "http://as.com",
		WinURL:        "http://win.com",
		URL:           "http://win.com",
		FormatID:      1,
		// Money
		PricingModel:       1,
		PurchaseViewPrice:  2000,
		PurchaseClickPrice: 1000,
		PurchaseLeadPrice:  1000,
		ViewPrice:          1000,
		ClickPrice:         1000,
		LeadPrice:          19992,
		CompetitorSourceID: 1000,
		CompetitorECPM:     1000000,
		// User IDENTITY
		UDID:        "1234-123456-123456-1234",
		UUID:        "1234-123456-123456-1234",
		SessionID:   "1234-123456-123456-1234",
		Fingerprint: "1234-123456-123456-1234",
		ETag:        "etaGddk*dk0a",
		// Targeting
		CarrierID:       11,
		Country:         "US",
		Latitude:        "1000",
		Longitude:       "-111",
		Language:        "en-EN",
		IPString:        "127.0.0.1",
		Referer:         "http://code.com",
		Page:            "",
		UserAgent:       "UA1",
		DeviceType:      1,
		DeviceID:        1,
		OSID:            122,
		BrowserID:       12,
		Categories:      "sss,ddd",
		Adblock:         1,
		PrivateBrowsing: 1,
		Robot:           1,
		Proxy:           1,
		Backup:          1,
		X:               10,
		Y:               10,
		Width:           10,
		Height:          10,
		// SubIDs
		SubID1: "sub1",
		SubID2: "sub2-period",
		SubID3: "sub3-date",
		SubID4: "sub4-some-category",
		SubID5: "sub5-user_id",
	}
}

func fixNumberString(n, def string) string {
	n = strings.TrimSpace(n)
	if strings.Trim(n, "0123456789") != "" {
		n = def
	}
	return n
}
