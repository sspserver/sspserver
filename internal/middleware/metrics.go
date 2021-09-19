//
// @project geniusrabbit::rotator 2016 - 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2017, 2019
//

package middleware

import (
	"github.com/valyala/fasthttp"

	"bitbucket.org/geniusrabbit/bigbrother/client"
	"github.com/geniusrabbit/notificationcenter"
)

// CollectMetrics and send to metrics
func CollectMetrics(metric string, spy Spy, next func(p client.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler {
	var (
		metrics          = notificationcenter.PublisherByName("metrics")
		metricsCountKey  = metric + ".count"
		metricsTimingKey = metric + ".timing"
	)

	return spy(func(p client.Person, ctx *fasthttp.RequestCtx) {
		var (
			info              = p.UserInfo()
			country           = info.Country().Code2
			_metricsCountKey  = metricsCountKey
			_metricsTimingKey = metricsTimingKey
		)

		if len(country) != 2 || country == "**" { // Undefined country
			country = "A1"
		}

		_metricsCountKey += ",country=" + country
		_metricsTimingKey += ",country=" + country

		metrics.Send(_metricsCountKey, map[string]func() error{
			_metricsTimingKey: func() error { next(p, ctx); return nil },
		})
	}) // end func
}

// CollectSimpleMetrics and send to metrics
func CollectSimpleMetrics(metric string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	var (
		metrics          = notificationcenter.PublisherByName("metrics")
		metricsCountKey  = metric + ".count"
		metricsTimingKey = metric + ".timing"
	)

	return func(ctx *fasthttp.RequestCtx) {
		metrics.Send(metricsCountKey, map[string]func() error{
			metricsTimingKey: func() error { next(ctx); return nil },
		})
	} // end func
}
