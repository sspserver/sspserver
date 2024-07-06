// Code generated by qtc from "ad_banner.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line private/templates/ad_banner.qtpl:2
package templates

//line private/templates/ad_banner.qtpl:2
import (
	"github.com/geniusrabbit/adcorelib/adtype"
)

//line private/templates/ad_banner.qtpl:7
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line private/templates/ad_banner.qtpl:7
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line private/templates/ad_banner.qtpl:7
func streamadRenderBanner(qw422016 *qt422016.Writer, resp adtype.Responser, it adtype.ResponserItem) {
//line private/templates/ad_banner.qtpl:9
	urlStr := URLGen.MustClickURL(it, resp)
	asset := it.MainAsset()
	format := it.Format()

//line private/templates/ad_banner.qtpl:13
	if asset.IsImage() {
//line private/templates/ad_banner.qtpl:13
		qw422016.N().S(`<a href="`)
//line private/templates/ad_banner.qtpl:14
		qw422016.N().S(urlStr)
//line private/templates/ad_banner.qtpl:14
		qw422016.N().S(`" target="_blank"><img src="`)
//line private/templates/ad_banner.qtpl:14
		qw422016.N().S(asset.Path)
//line private/templates/ad_banner.qtpl:14
		qw422016.N().S(`"onload="u`)
//line private/templates/ad_banner.qtpl:15
		qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:15
		qw422016.N().S(`(1);v`)
//line private/templates/ad_banner.qtpl:15
		qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:15
		qw422016.N().S(`(1)"onerror="u`)
//line private/templates/ad_banner.qtpl:16
		qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:16
		qw422016.N().S(`(0);v`)
//line private/templates/ad_banner.qtpl:16
		qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:16
		qw422016.N().S(`(0)"`)
//line private/templates/ad_banner.qtpl:17
		if format.IsFixed() {
//line private/templates/ad_banner.qtpl:17
			qw422016.N().S(`width="`)
//line private/templates/ad_banner.qtpl:18
			qw422016.N().D(format.Width)
//line private/templates/ad_banner.qtpl:18
			qw422016.N().S(`" height="`)
//line private/templates/ad_banner.qtpl:18
			qw422016.N().D(format.Height)
//line private/templates/ad_banner.qtpl:18
			qw422016.N().S(`"`)
//line private/templates/ad_banner.qtpl:19
		}
//line private/templates/ad_banner.qtpl:19
		qw422016.N().S(`/></a>`)
//line private/templates/ad_banner.qtpl:21
	} else {
//line private/templates/ad_banner.qtpl:22
		if asset.IsVideo() {
//line private/templates/ad_banner.qtpl:22
			qw422016.N().S(`<a href="`)
//line private/templates/ad_banner.qtpl:23
			qw422016.N().S(urlStr)
//line private/templates/ad_banner.qtpl:23
			qw422016.N().S(`" target="_blank"><videoonload="u`)
//line private/templates/ad_banner.qtpl:24
			qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:24
			qw422016.N().S(`(1);v`)
//line private/templates/ad_banner.qtpl:24
			qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:24
			qw422016.N().S(`(1)"onerror="u`)
//line private/templates/ad_banner.qtpl:25
			qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:25
			qw422016.N().S(`(0);v`)
//line private/templates/ad_banner.qtpl:25
			qw422016.N().D(int(it.AdID()))
//line private/templates/ad_banner.qtpl:25
			qw422016.N().S(`(0)"`)
//line private/templates/ad_banner.qtpl:26
			if format.IsFixed() {
//line private/templates/ad_banner.qtpl:26
				qw422016.N().S(`width="`)
//line private/templates/ad_banner.qtpl:27
				qw422016.N().D(format.Width)
//line private/templates/ad_banner.qtpl:27
				qw422016.N().S(`" height="`)
//line private/templates/ad_banner.qtpl:27
				qw422016.N().D(format.Height)
//line private/templates/ad_banner.qtpl:27
				qw422016.N().S(`"`)
//line private/templates/ad_banner.qtpl:28
			}
//line private/templates/ad_banner.qtpl:28
			qw422016.N().S(`autoplay loop><source src="`)
//line private/templates/ad_banner.qtpl:29
			qw422016.N().S(asset.Path)
//line private/templates/ad_banner.qtpl:29
			qw422016.N().S(`" type="video/mp4"></video></a>`)
//line private/templates/ad_banner.qtpl:31
		} else {
//line private/templates/ad_banner.qtpl:31
			qw422016.N().S(`Undefined asset type`)
//line private/templates/ad_banner.qtpl:32
			qw422016.E().V(asset.Type)
//line private/templates/ad_banner.qtpl:33
		}
//line private/templates/ad_banner.qtpl:34
	}
//line private/templates/ad_banner.qtpl:35
}

//line private/templates/ad_banner.qtpl:35
func writeadRenderBanner(qq422016 qtio422016.Writer, resp adtype.Responser, it adtype.ResponserItem) {
//line private/templates/ad_banner.qtpl:35
	qw422016 := qt422016.AcquireWriter(qq422016)
//line private/templates/ad_banner.qtpl:35
	streamadRenderBanner(qw422016, resp, it)
//line private/templates/ad_banner.qtpl:35
	qt422016.ReleaseWriter(qw422016)
//line private/templates/ad_banner.qtpl:35
}

//line private/templates/ad_banner.qtpl:35
func adRenderBanner(resp adtype.Responser, it adtype.ResponserItem) string {
//line private/templates/ad_banner.qtpl:35
	qb422016 := qt422016.AcquireByteBuffer()
//line private/templates/ad_banner.qtpl:35
	writeadRenderBanner(qb422016, resp, it)
//line private/templates/ad_banner.qtpl:35
	qs422016 := string(qb422016.B)
//line private/templates/ad_banner.qtpl:35
	qt422016.ReleaseByteBuffer(qb422016)
//line private/templates/ad_banner.qtpl:35
	return qs422016
//line private/templates/ad_banner.qtpl:35
}
