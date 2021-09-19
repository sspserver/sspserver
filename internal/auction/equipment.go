//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package auction

import (
	"geniusrabbit.dev/sspserver/internal/adsource"
	"geniusrabbit.dev/sspserver/internal/billing"
)

type equipmentSlice []adsource.ResponserItemCommon

func (l equipmentSlice) Len() int      { return len(l) }
func (l equipmentSlice) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l equipmentSlice) Less(i, j int) bool {
	e1, _ := l[i].(adsource.ResponserMultipleItem)
	e2, _ := l[j].(adsource.ResponserMultipleItem)
	if e1 != nil && (e2 == nil || e1.Count() > e2.Count()) {
		return true
	}
	if e2 != nil && (e1 == nil || e1.Count() < e2.Count()) {
		return false
	}
	return avgBid(e1, l[i]) > avgBid(e2, l[j])
}

func avgBid(mit adsource.ResponserMultipleItem, it adsource.ResponserItemCommon) billing.Money {
	if mit != nil {
		return mit.AuctionCPMBid() / billing.Money(mit.Count())
	}
	return it.AuctionCPMBid()
}
