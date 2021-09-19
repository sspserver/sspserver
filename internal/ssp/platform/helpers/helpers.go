//
// @project GeniusRabbit rotator 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package helpers

import (
	"strings"

	"bitbucket.org/geniusrabbit/corelib/models"
	"geniusrabbit.dev/sspserver/internal/adsource"
)

// InitBaseFilterByClient data from model
func InitBaseFilterByClient(f *adsource.BaseFilter, client *models.RTBSource) {
	f.Secure = client.Secure
	f.Adblock = client.AdBlock
	f.PrivateBrowsing = client.PrivateBrowsing
	f.Devices = []int(client.Devices)
	f.OS = intArr2UintPositive([]int(client.OS))
	f.OSExclude = intArr2UintNegative([]int(client.OS))
	f.Browsers = intArr2UintPositive([]int(client.Browsers))
	f.BrowsersExclude = intArr2UintNegative([]int(client.Browsers))
	f.Categories = intArr2Uint(client.Categories)
	f.Countries = []string(client.Countries)
	f.Applications = intArr2UintPositive([]int(client.Applications))
	f.ApplicationsExclude = intArr2UintNegative([]int(client.Applications))
	f.Domains = strArrFromPoritive([]string(client.Domains))
	f.DomainsExclude = strArrFromNegative([]string(client.Domains))
	f.Zones = intArr2UintPositive([]int(client.Zones))
	f.ZonesExclude = intArr2UintNegative([]int(client.Zones))
	// f.ExternalZones = strArrFromPoritive([]string(client.ExternalZones))
	// f.ExternalZonesExclude = strArrFromNegative([]string(client.ExternalZones))
	f.Normalise()
}

///////////////////////////////////////////////////////////////////////////////
// Conversion methods
///////////////////////////////////////////////////////////////////////////////

func int32Arr2Int(a []int32) (arr []int) {
	if a == nil {
		return nil
	}
	for _, v := range a {
		arr = append(arr, int(v))
	}
	return
}

func int64Arr2Int(a []int64) (arr []int) {
	if a == nil {
		return nil
	}
	for _, v := range a {
		arr = append(arr, int(v))
	}
	return
}

func intArr2Uint(a []int) (arr []uint) {
	if a == nil {
		return nil
	}
	for _, v := range a {
		arr = append(arr, uint(v))
	}
	return
}

func intArr2UintPositive(a []int) (arr []uint) {
	for _, v := range a {
		if v > 0 {
			arr = append(arr, uint(v))
		}
	}
	return
}

func intArr2UintNegative(a []int) (arr []uint) {
	for _, v := range a {
		if v < 0 {
			arr = append(arr, uint(-v))
		}
	}
	return
}

func strArrFromPoritive(a []string) (arr []string) {
	for _, v := range a {
		if !strings.HasPrefix(v, "-") {
			arr = append(arr, v)
		}
	}
	return
}

func strArrFromNegative(a []string) (arr []string) {
	for _, v := range a {
		if strings.HasPrefix(v, "-") {
			arr = append(arr, strings.TrimPrefix(v, "-"))
		}
	}
	return
}
