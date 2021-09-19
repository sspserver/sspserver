// Code generated by qtc from "test_banners.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line templates/test_banners.qtpl:1
package templates

//line templates/test_banners.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/test_banners.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/test_banners.qtpl:1
func StreamTestBanners(qw422016 *qt422016.Writer, zoneID uint64) {
//line templates/test_banners.qtpl:3
	var sizes = []struct {
		title string
		w, h  int
	}{
		{
			title: "leaderboard",
			w:     728, h: 90,
		},
		{
			title: "Large Leaderboard",
			w:     728, h: 210,
		},
		{
			title: "Large Leaderboard",
			w:     720, h: 300,
		},
		{
			title: "banner",
			w:     468, h: 60,
		},
		{
			title: "half banner",
			w:     234, h: 60,
		},
		{
			title: "mobile banner",
			w:     320, h: 50,
		},
		{
			title: "small rectangle",
			w:     180, h: 150,
		},
		{
			title: "small square",
			w:     200, h: 200,
		},
		{
			title: "Small Banner",
			w:     230, h: 33,
		},
		{
			title: "square",
			w:     250, h: 250,
		},
		{
			title: "300x100",
			w:     300, h: 100,
		},
		{
			title: "medium rectangle",
			w:     300, h: 250,
		},
		{
			title: "large rectangle",
			w:     336, h: 280,
		},
		{
			title: "unstandart",
			w:     315, h: 300,
		},
		{
			title: "half page",
			w:     300, h: 600,
		},
		{
			title: "wide skyscraper",
			w:     160, h: 600,
		},
		{
			title: "Button",
			w:     120, h: 30,
		},
		{
			title: "skyscraper",
			w:     120, h: 600,
		},
		{
			title: "portrait",
			w:     300, h: 1050,
		},
		{
			title: "large leaderboard",
			w:     970, h: 90,
		},
		{
			title: "billboard",
			w:     970, h: 250,
		},
	}

//line templates/test_banners.qtpl:92
	qw422016.N().S(`<!DOCTYPE html><html><head><meta name="viewport" content="width=device-width, initial-scale=1"><meta charset="utf-8" /><style type="text/css">*, body, html {margin: 0;padding: 0;border: none;}body, html {width: 100%;height: 100%;}iframe[seamless] {background-color: transparent;border: 0px none transparent;padding: 0px;overflow: hidden;margin: 0;}</style></head><body>`)
//line templates/test_banners.qtpl:120
	for _, size := range sizes {
//line templates/test_banners.qtpl:120
		qw422016.N().S(`<div style="float:left;display:block;padding:10px"><h4>`)
//line templates/test_banners.qtpl:122
		qw422016.N().S(size.title)
//line templates/test_banners.qtpl:122
		qw422016.N().S(`(`)
//line templates/test_banners.qtpl:122
		qw422016.N().D(size.w)
//line templates/test_banners.qtpl:122
		qw422016.N().S(`x`)
//line templates/test_banners.qtpl:122
		qw422016.N().D(size.h)
//line templates/test_banners.qtpl:122
		qw422016.N().S(`)</h4><iframe src="/ad/`)
//line templates/test_banners.qtpl:123
		qw422016.N().D(int(zoneID))
//line templates/test_banners.qtpl:123
		qw422016.N().S(`.html?w=`)
//line templates/test_banners.qtpl:123
		qw422016.N().D(size.w)
//line templates/test_banners.qtpl:123
		qw422016.N().S(`&amp;h=`)
//line templates/test_banners.qtpl:123
		qw422016.N().D(size.h)
//line templates/test_banners.qtpl:123
		qw422016.N().S(`"style="width:`)
//line templates/test_banners.qtpl:124
		qw422016.N().D(size.w)
//line templates/test_banners.qtpl:124
		qw422016.N().S(`px;height:`)
//line templates/test_banners.qtpl:124
		qw422016.N().D(size.h)
//line templates/test_banners.qtpl:124
		qw422016.N().S(`px;border:1px solid #ccc"scrolling="no"allowtransparency="true"allowfullscreen="true"frameborder="0"marginheight="0"marginwidth="0"vspace="0"hspace="0"></iframe></div>`)
//line templates/test_banners.qtpl:134
	}
//line templates/test_banners.qtpl:134
	qw422016.N().S(`</body></html>`)
//line templates/test_banners.qtpl:137
}

//line templates/test_banners.qtpl:137
func WriteTestBanners(qq422016 qtio422016.Writer, zoneID uint64) {
//line templates/test_banners.qtpl:137
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/test_banners.qtpl:137
	StreamTestBanners(qw422016, zoneID)
//line templates/test_banners.qtpl:137
	qt422016.ReleaseWriter(qw422016)
//line templates/test_banners.qtpl:137
}

//line templates/test_banners.qtpl:137
func TestBanners(zoneID uint64) string {
//line templates/test_banners.qtpl:137
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/test_banners.qtpl:137
	WriteTestBanners(qb422016, zoneID)
//line templates/test_banners.qtpl:137
	qs422016 := string(qb422016.B)
//line templates/test_banners.qtpl:137
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/test_banners.qtpl:137
	return qs422016
//line templates/test_banners.qtpl:137
}