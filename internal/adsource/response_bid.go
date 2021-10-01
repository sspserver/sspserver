//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adsource

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	openrtb "github.com/bsm/openrtb"
	natresp "github.com/bsm/openrtb/native/response"
	"github.com/demdxx/gocast"
	"golang.org/x/net/html/charset"

	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models"
	"geniusrabbit.dev/sspserver/internal/models/types"
)

// BidResponse RTB record
type BidResponse struct {
	Src         Sourcer
	Req         *BidRequest
	Application *models.Application
	Target      models.Target
	BidResponse openrtb.BidResponse
	context     context.Context
	optimalBids []*openrtb.Bid
	ads         []ResponserItemCommon
}

// AuctionID response
func (r *BidResponse) AuctionID() string {
	return r.BidResponse.ID
}

// AuctionType of request
func (r *BidResponse) AuctionType() AuctionType {
	return r.Req.AuctionType
}

// Source of response
func (r *BidResponse) Source() Sourcer {
	return r.Src
}

// Prepare bid response
func (r *BidResponse) Prepare() {
	for i, seat := range r.BidResponse.SeatBid {
		for i, bid := range seat.Bid {
			replacer := strings.NewReplacer(
				"${AUCTION_AD_ID}", bid.AdID,
				"${AUCTION_ID}", r.BidResponse.ID,
				"${AUCTION_BID_ID}", r.BidResponse.BidID,
				"${AUCTION_IMP_ID}", bid.ImpID,
				"${AUCTION_PRICE}", fmt.Sprintf("%.6f", bid.Price),
				"${AUCTION_CURRENCY}", "USD",
			)

			// Custom direct detect
			if len(bid.AdMarkup) < 1 {
				var (
					ext map[string]interface{}
					url interface{}
				)
				if json.Unmarshal(bid.Ext, &ext); ext != nil {
					if url, _ = ext["url"]; url == nil {
						url, _ = ext["landingpage"]
					}

					if url != nil {
						switch v := url.(type) {
						case string:
							bid.AdMarkup = v
						case []string:
							if len(v) > 0 {
								bid.AdMarkup = v[0]
							}
						case []interface{}:
							if len(v) > 0 {
								bid.AdMarkup = gocast.ToString(v[0])
							}
						}
					} // end if
				}
			}

			if imp := r.Req.ImpressionByIDvariation(bid.ImpID); imp != nil {
				// Prepare date for bid W/H
				if bid.W == 0 && bid.H == 0 {
					bid.W, bid.H = imp.W, imp.H
				}

				if imp.IsDirect() {
					if strings.HasPrefix(bid.AdMarkup, `<?xml`) {
						// TODO postprocess direct error
						bid.AdMarkup, _ = decodePopMarkup([]byte(bid.AdMarkup))
					}
				}
			}

			bid.AdMarkup = replacer.Replace(bid.AdMarkup)

			if len(bid.NURL) > 0 {
				if u, err := url.QueryUnescape(bid.NURL); err == nil {
					bid.NURL = u
				}
				bid.NURL = replacer.Replace(bid.NURL)
			}

			seat.Bid[i] = bid
		}

		r.BidResponse.SeatBid[i] = seat
	} // end for

	bids := r.OptimalBids()
	for _, bid := range bids {
		imp := r.Req.ImpressionByIDvariation(bid.ImpID)
		if imp == nil {
			continue
		}
		if imp.IsDirect() {
			format := imp.FormatByType(types.FormatDirectType)
			if format == nil {
				continue
			}
			r.ads = append(r.ads, &ResponseBidItem{
				ItemID:     imp.ID,
				Src:        r.Src,
				Req:        r.Req,
				Imp:        imp,
				FormatType: types.FormatDirectType,
				RespFormat: format,
				Bid:        bid,
				Native:     nil,
				Data:       nil,
				BidPrice:   billing.MoneyFloat(bid.Price),
			})
			continue
		}
		for _, format := range imp.Formats() {
			if bid.ImpID != imp.IDByFormat(format) {
				continue
			}
			switch {
			case format.IsNative():
				native, err := decodeNativeMarkup([]byte(bid.AdMarkup))
				// TODO parse native request
				if err == nil {
					r.ads = append(r.ads, &ResponseBidItem{
						ItemID:     imp.ID,
						Src:        r.Src,
						Req:        r.Req,
						Imp:        imp,
						FormatType: types.FormatNativeType,
						RespFormat: format,
						Bid:        bid,
						Native:     native,
						BidPrice:   billing.MoneyFloat(bid.Price),
					})
				}
			case format.IsBanner() || format.IsProxy():
				r.ads = append(r.ads, &ResponseBidItem{
					ItemID:     imp.ID,
					Src:        r.Src,
					Req:        r.Req,
					Imp:        imp,
					FormatType: bannerFormatType(bid.AdMarkup),
					RespFormat: format,
					Bid:        bid,
					BidPrice:   billing.MoneyFloat(bid.Price),
				})
			}
			break
		}
	}
}

// Request information
func (r *BidResponse) Request() *BidRequest {
	return r.Req
}

// Ads list
func (r *BidResponse) Ads() []ResponserItemCommon {
	return r.ads
}

// Item by impression code
func (r *BidResponse) Item(impid string) ResponserItemCommon {
	for _, it := range r.Ads() {
		if it.ImpressionID() == impid {
			return it
		}
	}
	return nil
}

// ActionURL for rtb
func (r *BidResponse) ActionURL() string {
	return ""
}

// Price for response
func (r *BidResponse) Price() billing.Money {
	var price billing.Money
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			price += billing.MoneyFloat(bid.Price)
		}
	}
	return price
}

// Count bids
func (r *BidResponse) Count() int {
	return len(r.Bids())
}

// Validate response
func (r *BidResponse) Validate() error {
	err := r.BidResponse.Validate()
	if err == nil {
		for _, seat := range r.BidResponse.SeatBid {
			if seat.Group == 1 {
				return ErrResponseInvalidGroup
			}
		}
	}
	return err
}

// Error of the response
func (r *BidResponse) Error() error {
	return r.Validate()
}

// Bids list
func (r *BidResponse) Bids() []*openrtb.Bid {
	result := make([]*openrtb.Bid, 0, len(r.BidResponse.SeatBid))
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			result = append(result, &bid)
		}
	}
	return result
}

// OptimalBids list (the most expensive)
func (r *BidResponse) OptimalBids() []*openrtb.Bid {
	if len(r.optimalBids) > 0 {
		return r.optimalBids
	}

	var bids = map[string]*openrtb.Bid{}
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			if obid, ok := bids[bid.ImpID]; !ok || obid.Price < bid.Price {
				bids[bid.ImpID] = &bid
			}
		}
	}

	for _, b := range bids {
		r.optimalBids = append(r.optimalBids, b)
	}
	return r.optimalBids
}

// BidPosition returns index from OpenRTB bid
func (r *BidResponse) BidPosition(b *openrtb.Bid) int {
	var idx int
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			if bid.ImpID == b.ImpID {
				return idx
			}
			idx++
		}
	}
	return idx
}

// UpdateBid object
func (r *BidResponse) UpdateBid(b *openrtb.Bid) {
	for _, seat := range r.BidResponse.SeatBid {
		for j, bid := range seat.Bid {
			if bid.ImpID == b.ImpID {
				seat.Bid[j] = *b
			}
		}
	}
}

// Context of response
func (r *BidResponse) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		r.context = ctx[0]
	}
	return r.context
}

// Get context value
func (r *BidResponse) Get(key string) interface{} {
	if r.context != nil {
		return r.context.Value(key)
	}
	return nil
}

func decodePopMarkup(data []byte) (val string, err error) {
	var item struct {
		URL string `xml:"popunderAd>url"`
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&item); err == nil {
		val = item.URL
	}
	return val, err
}

func decodeNativeMarkup(data []byte) (*natresp.Response, error) {
	var (
		native struct {
			Native natresp.Response `json:"native"`
		}
		err = json.Unmarshal(data, &native)
	)
	if err != nil {
		err = json.Unmarshal(data, &native.Native)
	}
	if err != nil {
		return nil, err
	}
	return &native.Native, nil
}

func bannerFormatType(markup string) types.FormatType {
	if strings.HasPrefix(markup, "http://") ||
		strings.HasPrefix(markup, "https://") ||
		(strings.HasPrefix(markup, "//") && !strings.ContainsAny(markup, "\n\t")) ||
		strings.Contains(markup, "<iframe") {
		return types.FormatProxyType
	}
	return types.FormatBannerType
}

var (
	_ Responser = &BidResponse{}
)
