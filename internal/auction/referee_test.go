//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

// go test -v -run=^$ -bench=. -cpuprofile=cpu.out
// go tool pprof --pdf auction.test cpu.out > cpu0.pdf

package auction

import (
	"fmt"
	"sort"
	"testing"

	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/billing"
	"geniusrabbit.dev/sspserver/internal/models"
)

type titem struct {
	ImpID string
	Bid   int64
}

type tsItem struct {
	Multi bool
	ImpID string
	Bid   billing.Money
}

func (i tsItem) String() string {
	if i.Multi {
		return fmt.Sprintf("@%s:%s", i.ImpID, i.Bid.String())
	}
	return fmt.Sprintf("%s:%s", i.ImpID, i.Bid.String())
}

type test struct {
	Name  string
	Rings []Ring
	Bids  []billing.Money
	Scope []adsource.ResponserItemCommon
}

func (t test) TestResponse(res []adsource.ResponserItemCommon) bool {
	return testBids(t.Bids, res)
}

func testBids(bids []billing.Money, res []adsource.ResponserItemCommon) bool {
	var newRes []adsource.ResponserItemCommon
	for _, it := range res {
		switch v := it.(type) {
		case adsource.ResponserMultipleItem:
			for _, it := range v.Ads() {
				newRes = append(newRes, it)
			}
		default:
			newRes = append(newRes, it)
		}
	}

	if len(bids) != len(newRes) {
		return false
	}

	sort.Slice(newRes, func(i, j int) bool { return newRes[i].AuctionCPMBid() < newRes[j].AuctionCPMBid() })
	sort.Slice(bids, func(i, j int) bool { return bids[i] < bids[j] })

	for i, b := range bids {
		if b != newRes[i].AuctionCPMBid() {
			return false
		}
	}
	return true
}

func respToBids(res []adsource.ResponserItemCommon) (bids []tsItem) {
	for _, it := range res {
		switch v := it.(type) {
		case adsource.ResponserMultipleItem:
			for _, it := range v.Ads() {
				bids = append(bids, tsItem{
					Multi: true,
					ImpID: it.ImpressionID(),
					Bid:   it.AuctionCPMBid(),
				})
			}
		default:
			bids = append(bids, tsItem{
				ImpID: it.ImpressionID(),
				Bid:   it.AuctionCPMBid(),
			})
		}
	}
	return
}

func TestRefereeMatch(t *testing.T) {
	tests := []test{
		{
			Name:  "Simple",
			Rings: []Ring{{ID: "1", Count: 3}},
			Bids:  []billing.Money{mi(7), mi(2), mi(2)},
			Scope: []adsource.ResponserItemCommon{
				newItem("1", 1),
				newItem("1", 2),
				newItem("1", 7),
				newItem("1", 2),
				newItem("1", 1),
			},
		},
		{
			Name:  "Multy Simple 2",
			Rings: []Ring{{ID: "1", Count: 3}, {ID: "2", Count: 1}},
			Bids:  []billing.Money{mi(7), mi(4), mi(3), mi(2)},
			Scope: []adsource.ResponserItemCommon{
				newItem("1", 1),
				newItem("1", 2),
				newItem("1", 7),
				newItem("1", 3),
				newItem("2", 2),
				newItem("2", 1),
				newItem("2", 4),
			},
		},
		{
			Name:  "Multy Simple 3",
			Rings: []Ring{{ID: "1", Count: 3}, {ID: "2", Count: 1}, {ID: "3", Count: 2}},
			Bids:  []billing.Money{mi(7), mi(4), mi(3), mi(2), mi(3), mi(9)},
			Scope: []adsource.ResponserItemCommon{
				newItem("1", 1),
				newItem("1", 2),
				newItem("1", 7),
				newItem("1", 3),
				newItem("2", 2),
				newItem("2", 1),
				newItem("2", 4),
				newItem("3", 3),
				newItem("3", 1),
				newItem("3", 9),
			},
		},
		{
			Name:  "Multy 3",
			Rings: []Ring{{ID: "1", Count: 3}, {ID: "2", Count: 1}, {ID: "3", Count: 2}},
			Bids:  []billing.Money{mi(1), mi(2), mi(3), mi(7), mi(3), mi(1)},
			Scope: []adsource.ResponserItemCommon{
				newMultipleItem(
					titem{ImpID: "1", Bid: 1},
					titem{ImpID: "2", Bid: 2},
					titem{ImpID: "3", Bid: 3},
				),
				newItem("1", 2),
				newItem("1", 7),
				newItem("1", 3),
				newItem("2", 3),
				newItem("2", 1),
				newItem("3", 1),
				newItem("3", 1),
			},
		},
		{
			Name:  "Multy 3*2",
			Rings: []Ring{{ID: "1", Count: 3}, {ID: "2", Count: 1}, {ID: "3", Count: 2}},
			Bids:  []billing.Money{mi(1), mi(2), mi(3), mi(7), mi(1), mi(1)},
			Scope: []adsource.ResponserItemCommon{
				newMultipleItem(
					titem{ImpID: "1", Bid: 1},
					titem{ImpID: "2", Bid: 2},
					titem{ImpID: "3", Bid: 3},
				),
				newMultipleItem(
					titem{ImpID: "1", Bid: 1},
					titem{ImpID: "3", Bid: 1},
				),
				newItem("1", 2),
				newItem("1", 7),
				newItem("1", 3),
				newItem("2", 3),
				newItem("2", 1),
				newItem("3", 1),
				newItem("3", 1),
			},
		},
	}

	for _, ts := range tests {
		t.Run(fmt.Sprintf("%s (%d ring's)", ts.Name, len(ts.Rings)), func(t *testing.T) {
			ref := Referee{}
			ref.Push(ts.Scope...)

			if res := ref.Match(ts.Rings...); !ts.TestResponse(res) {
				bids := respToBids(res)
				t.Errorf("Fail result test: %v => %v", ts.Bids, bids)
			}
		})
	}
}

func TestReplacement(t *testing.T) {
	bids := []billing.Money{mi(7), mi(3), mi(1)}
	resp := collectReplacement(newMultipleItem(
		titem{ImpID: "1", Bid: 1},
		titem{ImpID: "2", Bid: 2},
		titem{ImpID: "3", Bid: 3},
	), []adsource.ResponserItemCommon{
		newItem("1", 7),
		newItem("1", 3),
		newItem("2", 3),
		newItem("1", 2),
		newItem("2", 1),
		newItem("3", 1),
		newItem("3", 1),
	})

	if !testBids(bids, resp) {
		t.Errorf("Fail result test: %v => %v", bids, respToBids(resp))
	}
}

func BenchmarkRefereeMatch(b *testing.B) {
	tt := test{
		Name:  "Multy 3*2",
		Rings: []Ring{{ID: "1", Count: 3}, {ID: "2", Count: 1}, {ID: "3", Count: 2}},
		Bids:  []billing.Money{mi(1), mi(2), mi(3), mi(7), mi(1), mi(1)},
		Scope: []adsource.ResponserItemCommon{
			newMultipleItem(
				titem{ImpID: "1", Bid: 1},
				titem{ImpID: "2", Bid: 2},
				titem{ImpID: "3", Bid: 3},
			),
			newMultipleItem(
				titem{ImpID: "1", Bid: 1},
				titem{ImpID: "3", Bid: 1},
			),
			newItem("1", 2),
			newItem("1", 7),
			newItem("1", 3),
			newItem("2", 3),
			newItem("2", 1),
			newItem("3", 1),
			newItem("3", 1),
		},
	}

	for n := 0; n < b.N; n++ {
		ref := Referee{}
		ref.Push(tt.Scope...)
		ref.Match(tt.Rings...)
	}
}

func newItem(impid string, bid int64) adsource.ResponserItem {
	return &adsource.ResponseAdItem{
		Src:         nil,
		Req:         nil,
		Imp:         &adsource.Impression{ID: impid},
		Campaign:    &models.Campaign{}, //{Opt: [8]uint8{0, uint8(types.PricingModelCPC)}},
		Ad:          &models.Ad{},
		BidECPM:     billing.MoneyInt(bid),
		BidPrice:    0,
		CPMBidPrice: 0,
	}
}

func newMultipleItem(bids ...titem) adsource.ResponserMultipleItem {
	block := &adsource.ResponseItemBlock{}
	for _, bid := range bids {
		block.Items = append(block.Items, newItem(bid.ImpID, bid.Bid))
	}
	return block
}

func mi(v int) billing.Money {
	return billing.MoneyInt(int64(v))
}
