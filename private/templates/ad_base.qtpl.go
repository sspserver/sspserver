// Code generated by qtc from "ad_base.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// Advertisement base elements
//

//line templates/ad_base.qtpl:4
package templates

//line templates/ad_base.qtpl:4
import (
	"geniusrabbit.dev/sspserver/internal/events"
	"geniusrabbit.dev/sspserver/internal/adsource"
)

//line templates/ad_base.qtpl:10
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/ad_base.qtpl:10
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/ad_base.qtpl:11
var Debug bool
var ServiceDomain string
var URLGen adsource.URLGenerator

//line templates/ad_base.qtpl:16
func streamadActionScript(qw422016 *qt422016.Writer) {
//line templates/ad_base.qtpl:16
	qw422016.N().S(`<script type="text/javascript">var t = new Date();function e(u, st) {var delta = new Date() - t;var qPixel = new Image();qPixel.src = u+'&r='+st+'&d='+delta;};</script>`)
//line templates/ad_base.qtpl:25
}

//line templates/ad_base.qtpl:25
func writeadActionScript(qq422016 qtio422016.Writer) {
//line templates/ad_base.qtpl:25
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/ad_base.qtpl:25
	streamadActionScript(qw422016)
//line templates/ad_base.qtpl:25
	qt422016.ReleaseWriter(qw422016)
//line templates/ad_base.qtpl:25
}

//line templates/ad_base.qtpl:25
func adActionScript() string {
//line templates/ad_base.qtpl:25
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/ad_base.qtpl:25
	writeadActionScript(qb422016)
//line templates/ad_base.qtpl:25
	qs422016 := string(qb422016.B)
//line templates/ad_base.qtpl:25
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/ad_base.qtpl:25
	return qs422016
//line templates/ad_base.qtpl:25
}

//line templates/ad_base.qtpl:27
func streamadHeader(qw422016 *qt422016.Writer) {
//line templates/ad_base.qtpl:27
	qw422016.N().S(`<!DOCTYPE html><html><head><meta name="viewport" content="width=device-width, initial-scale=1"><meta charset="utf-8" /><style type="text/css">*, body, html { margin: 0; padding: 0; border:none; }body, html { width: 100%; height: 100%; background: transparent }iframe[seamless] {background-color: transparent;border: 0px none transparent;padding: 0px;overflow: hidden;margin: 0;}</style></head><body>`)
//line templates/ad_base.qtpl:43
	streamadActionScript(qw422016)
//line templates/ad_base.qtpl:44
}

//line templates/ad_base.qtpl:44
func writeadHeader(qq422016 qtio422016.Writer) {
//line templates/ad_base.qtpl:44
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/ad_base.qtpl:44
	streamadHeader(qw422016)
//line templates/ad_base.qtpl:44
	qt422016.ReleaseWriter(qw422016)
//line templates/ad_base.qtpl:44
}

//line templates/ad_base.qtpl:44
func adHeader() string {
//line templates/ad_base.qtpl:44
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/ad_base.qtpl:44
	writeadHeader(qb422016)
//line templates/ad_base.qtpl:44
	qs422016 := string(qb422016.B)
//line templates/ad_base.qtpl:44
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/ad_base.qtpl:44
	return qs422016
//line templates/ad_base.qtpl:44
}

//line templates/ad_base.qtpl:47
func streamadFooter(qw422016 *qt422016.Writer) {
//line templates/ad_base.qtpl:47
	qw422016.N().S(`</body></html>`)
//line templates/ad_base.qtpl:49
}

//line templates/ad_base.qtpl:49
func writeadFooter(qq422016 qtio422016.Writer) {
//line templates/ad_base.qtpl:49
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/ad_base.qtpl:49
	streamadFooter(qw422016)
//line templates/ad_base.qtpl:49
	qt422016.ReleaseWriter(qw422016)
//line templates/ad_base.qtpl:49
}

//line templates/ad_base.qtpl:49
func adFooter() string {
//line templates/ad_base.qtpl:49
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/ad_base.qtpl:49
	writeadFooter(qb422016)
//line templates/ad_base.qtpl:49
	qs422016 := string(qb422016.B)
//line templates/ad_base.qtpl:49
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/ad_base.qtpl:49
	return qs422016
//line templates/ad_base.qtpl:49
}

// Generate pixel base code

//line templates/ad_base.qtpl:53
func streamadPixel(qw422016 *qt422016.Writer, adID, spotID, campID int, tag string) {
//line templates/ad_base.qtpl:53
	qw422016.N().S(`<script type="text/javascript">function u`)
//line templates/ad_base.qtpl:55
	qw422016.N().D(adID)
//line templates/ad_base.qtpl:55
	qw422016.N().S(`(st){}</script>`)
//line templates/ad_base.qtpl:57
}

//line templates/ad_base.qtpl:57
func writeadPixel(qq422016 qtio422016.Writer, adID, spotID, campID int, tag string) {
//line templates/ad_base.qtpl:57
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/ad_base.qtpl:57
	streamadPixel(qw422016, adID, spotID, campID, tag)
//line templates/ad_base.qtpl:57
	qt422016.ReleaseWriter(qw422016)
//line templates/ad_base.qtpl:57
}

//line templates/ad_base.qtpl:57
func adPixel(adID, spotID, campID int, tag string) string {
//line templates/ad_base.qtpl:57
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/ad_base.qtpl:57
	writeadPixel(qb422016, adID, spotID, campID, tag)
//line templates/ad_base.qtpl:57
	qs422016 := string(qb422016.B)
//line templates/ad_base.qtpl:57
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/ad_base.qtpl:57
	return qs422016
//line templates/ad_base.qtpl:57
}

// Generate pixel base code for adresult item

//line templates/ad_base.qtpl:61
func streamadPixelItem(qw422016 *qt422016.Writer, ad adsource.ResponserItem, resp adsource.Responser) {
//line templates/ad_base.qtpl:62
	if ad != nil && resp != nil {
//line templates/ad_base.qtpl:62
		qw422016.N().S(`<script type="text/javascript">`)
//line templates/ad_base.qtpl:64
		var u, _ = URLGen.PixelURL(events.Impression, events.StatusSuccess, ad, resp, false)

//line templates/ad_base.qtpl:65
		var v, _ = URLGen.PixelURL(events.View, events.StatusSuccess, ad, resp, false)

//line templates/ad_base.qtpl:65
		qw422016.N().S(`function u`)
//line templates/ad_base.qtpl:66
		qw422016.N().D(int(ad.AdID()))
//line templates/ad_base.qtpl:66
		qw422016.N().S(`(st){e('`)
//line templates/ad_base.qtpl:66
		qw422016.N().S(u)
//line templates/ad_base.qtpl:66
		qw422016.N().S(`',st)}function v`)
//line templates/ad_base.qtpl:67
		qw422016.N().D(int(ad.AdID()))
//line templates/ad_base.qtpl:67
		qw422016.N().S(`(st){e('`)
//line templates/ad_base.qtpl:67
		qw422016.N().S(v)
//line templates/ad_base.qtpl:67
		qw422016.N().S(`',st)}</script>`)
//line templates/ad_base.qtpl:69
	}
//line templates/ad_base.qtpl:70
}

//line templates/ad_base.qtpl:70
func writeadPixelItem(qq422016 qtio422016.Writer, ad adsource.ResponserItem, resp adsource.Responser) {
//line templates/ad_base.qtpl:70
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/ad_base.qtpl:70
	streamadPixelItem(qw422016, ad, resp)
//line templates/ad_base.qtpl:70
	qt422016.ReleaseWriter(qw422016)
//line templates/ad_base.qtpl:70
}

//line templates/ad_base.qtpl:70
func adPixelItem(ad adsource.ResponserItem, resp adsource.Responser) string {
//line templates/ad_base.qtpl:70
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/ad_base.qtpl:70
	writeadPixelItem(qb422016, ad, resp)
//line templates/ad_base.qtpl:70
	qs422016 := string(qb422016.B)
//line templates/ad_base.qtpl:70
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/ad_base.qtpl:70
	return qs422016
//line templates/ad_base.qtpl:70
}

// { %code var uu, _ = URLGen.PixelLead(ad, resp, false)  % }
// <img src="{ %s= uu % }" />
//
//

//line templates/ad_base.qtpl:76
func streampreloader(qw422016 *qt422016.Writer) {
//line templates/ad_base.qtpl:76
	qw422016.N().S(`<style>.loading {position: absolute;height: 100%;width: 100%;top: 0;left: 0;background: #fefefe;display: block;z-index: 1000;}.loading .progress {position: fixed;display: block;width: 100%;height: 1.5pt;background: deepskyblue;}.loading .progress:before {content: "";position: absolute;left: 0;top: 0;width: 100%;height: 100%;transform: translateX(-100%);background: #ccc;animation: progress 3s ease infinite;}.loading .badge {position: absolute;left: 50%;top: 50%;display: block;padding: 3pt;margin: -15pt 0 0 -15pt;border: 1.5pt solid #ddd;border-radius: 5pt;font-family: Helvetica,sans-serif;font-size: 12pt;color: #ddd;}@-webkit-keyframes progress {50% {-webkit-transform: translateX(0%);transform: translateX(0%);}100% {-webkit-transform: translateX(100%);transform: translateX(100%);}}@keyframes progress {50% {-webkit-transform: translateX(0%);transform: translateX(0%);}100% {-webkit-transform: translateX(100%);transform: translateX(100%);}}}</style><div id="loadingBlock" class="loading"><div class="progress"></div><div class="badge">ADS</div></div>`)
//line templates/ad_base.qtpl:146
}

//line templates/ad_base.qtpl:146
func writepreloader(qq422016 qtio422016.Writer) {
//line templates/ad_base.qtpl:146
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/ad_base.qtpl:146
	streampreloader(qw422016)
//line templates/ad_base.qtpl:146
	qt422016.ReleaseWriter(qw422016)
//line templates/ad_base.qtpl:146
}

//line templates/ad_base.qtpl:146
func preloader() string {
//line templates/ad_base.qtpl:146
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/ad_base.qtpl:146
	writepreloader(qb422016)
//line templates/ad_base.qtpl:146
	qs422016 := string(qb422016.B)
//line templates/ad_base.qtpl:146
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/ad_base.qtpl:146
	return qs422016
//line templates/ad_base.qtpl:146
}