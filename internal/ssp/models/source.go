//
// @project geniusrabbit::sspserver 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package models

import (
	"crypto/tls"
	"math"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Source source
type Source struct {
	sync.Mutex
	URL            string
	Driver         string
	HTTPMethod     string
	Headers        map[string]string
	RequestFormat  string
	AuctionType    int
	Weight         int
	Tags           []string
	GeoMinimalBids GeoBidSlice
	RPS            int
	MinBid         float64
	Targeting      map[string][]string
	Client         *http.Client
	rpsCurrent     int
	errorCounter   int
	lastPeriod     time.Time
}

// Init source
func (s *Source) Init() {
	s.Prepare()
}

// Prepare params
func (s *Source) Prepare() {
	if len(s.Tags) > 0 {
		sort.Strings(s.Tags)
	}

	if nil != s.Targeting {
		for key, val := range s.Targeting {
			sort.Strings(val)
			s.Targeting[key] = val
		}
	}
}

// Test targeting & RPS
func (s *Source) Test(tags []string, targeting map[string]string) bool {
	if s.RPS > 0 {
		s.Lock()
		defer s.Unlock()

		if s.errorCounter > 8 {
			if v := notmalise((float64(s.errorCounter) / 2.0) - 10.0); v > 0.01 && rand.Float64()*1.1 <= v {
				return false
			}
		}

		if now := time.Now(); now.Sub(s.lastPeriod).Seconds() >= 1 {
			s.lastPeriod = now
			s.rpsCurrent = 0
		}

		if s.rpsCurrent >= s.RPS {
			return false
		}
	}
	return s.testFilter(tags, targeting)
}

// Test targeting
func (s *Source) testFilter(tags []string, targeting map[string]string) (result bool) {
	// Compare tags
	if len(s.Tags) > 0 {
		result = false
		if len(tags) < 1 {
			return
		}
		for _, tag := range tags {
			if i := sort.SearchStrings(s.Tags, tag); i >= 0 && i < len(s.Tags) && s.Tags[i] == tag {
				result = true
				break
			}
		}
	} else {
		result = true
	}

	if !result {
		return
	}

	// Compare targeting
	if s.Targeting != nil && len(s.Targeting) > 0 {
		if targeting == nil || len(targeting) < 1 {
			return false
		}

		for key, val := range targeting {
			if trg := s.Targeting[key]; len(trg) > 0 {
				if i := sort.SearchStrings(trg, val); i < 0 || i >= len(trg) || trg[i] != val {
					result = false
					break
				}
			}
		} // end for
	}
	return
}

// GetClient of http
func (s *Source) GetClient() *http.Client {
	if s.Client == nil {
		s.updateClient(100 * time.Millisecond)
	}
	return s.Client
}

func (s *Source) updateClient(timeout time.Duration) {
	s.Client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           (&net.Dialer{Timeout: timeout, KeepAlive: 30 * time.Second}).DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			DisableKeepAlives:     false,
			IdleConnTimeout:       300 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 100 * time.Millisecond,
		},
		Timeout: timeout,
	}
}

// IncRequestCounter count
func (s *Source) IncRequestCounter() {
	s.rpsCurrent++
}

// IncErrorCounter of errors
func (s *Source) IncErrorCounter(v int) {
	s.errorCounter += v
	if s.errorCounter > 300 {
		s.errorCounter = 300
	} else if s.errorCounter < 0 {
		s.errorCounter = 0
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

var exp1 = math.Exp(1)

func notmalise(x float64) float64 {
	var pow = math.Pow(exp1, x/1.8)
	return pow / (1.0 + pow)
}
