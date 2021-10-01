//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adsource

import "geniusrabbit.dev/sspserver/internal/eventtraking/events"

// URLGenerator of advertisement
type URLGenerator interface {
	// CDNURL returns full URL to path
	CDNURL(path string) string

	// PixelURL generator from response of ite
	PixelURL(event events.Type, status uint8, item ResponserItem, response Responser, js bool) (string, error)

	// PixelDirectURL generator from response of item
	PixelDirectURL(event events.Type, status uint8, item ResponserItem, response Responser, direct string) (string, error)

	// PixelLead URL
	PixelLead(item ResponserItem, response Responser, js bool) (string, error)

	// MustClickURL generator from response of item
	MustClickURL(item ResponserItem, response Responser) string

	// ClickURL generator from response of item
	ClickURL(item ResponserItem, response Responser) (string, error)

	// ClickRouterURL returns router pattern
	ClickRouterURL() string

	// DirectURL generator from response of item
	DirectURL(event events.Type, item ResponserItem, response Responser) (string, error)

	// DirectRouterURL returns router pattern
	DirectRouterURL() string
}
