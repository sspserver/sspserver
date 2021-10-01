//
// @project geniusrabbit::archivarious 2017 - 2018, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018, 2021
//

package pixelgenerator

import (
	"fmt"
	"net/url"

	"geniusrabbit.dev/sspserver/internal/eventtraking/events"
)

// PixelGenerator object
type PixelGenerator struct {
	hostname string
}

// NewPixelGenerator object
func NewPixelGenerator(hostname string) PixelGenerator {
	return PixelGenerator{
		hostname: hostname,
	}
}

// Event pixel
func (g PixelGenerator) Event(ev *events.Event, js bool) (a string, err error) {
	var (
		code = ev.Pack().Compress().URLEncode()
		u    = url.Values{"i": []string{code.String()}}
	)

	if err = code.ErrorObj(); err != nil {
		return
	}

	if js {
		a = fmt.Sprintf("//%s/t/px.js?%s", g.hostname, u.Encode())
	} else {
		a = fmt.Sprintf("//%s/t/px.gif?%s", g.hostname, u.Encode())
	}
	return
}

// EventDirect pixel
func (g PixelGenerator) EventDirect(ev *events.Event, direct string) (a string, err error) {
	var (
		code = ev.Pack().Compress().URLEncode()
		u    = url.Values{
			"i": []string{code.String()},
			"u": []string{direct},
		}
	)

	if err = code.ErrorObj(); err != nil {
		return
	}

	return fmt.Sprintf("//%s/go/m?%s", g.hostname, u.Encode()), nil
}

// Lead URL
func (g PixelGenerator) Lead(lead *events.LeadCode) (string, error) {
	return fmt.Sprintf("//%s/lead?l=%s", g.hostname, url.QueryEscape(lead.String())), nil
}
