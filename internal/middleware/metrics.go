//
// @project geniusrabbit::rotator 2016 - 2017, 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2017, 2019, 2021
//

package middleware

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/sspserver/internal/personification"
)

// CollectMetrics and send to metrics
func CollectMetrics(metric string, spy Spy, next func(p personification.Person, ctx *fasthttp.RequestCtx)) fasthttp.RequestHandler {
	var (
		buckets      = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
		metricsCount = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: metric + "_count",
			Help: "Count of requests by country",
		}, []string{"country"})
		metricTiming = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "spy",
			Name:      metric + "_duration_seconds",
			Help:      "Histogram of response time for handler in seconds by country",
			Buckets:   buckets,
		}, []string{"country", "method", "status_code"})
	)

	return spy(func(p personification.Person, ctx *fasthttp.RequestCtx) {
		var (
			info    = p.UserInfo()
			country = info.Country().Country
			start   = time.Now()
		)

		if len(country) != 2 || country == "**" { // Undefined country
			country = "A1"
		}

		metricsCount.WithLabelValues(country).Inc()

		next(p, ctx)
		duration := time.Since(start)
		statusCode := strconv.Itoa(ctx.Response.StatusCode())
		metricTiming.WithLabelValues(country, string(ctx.Method()), statusCode).Observe(duration.Seconds())
	})
}

// CollectSimpleMetrics and send to metrics
func CollectSimpleMetrics(metric string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	var (
		buckets      = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
		metricsCount = promauto.NewCounter(prometheus.CounterOpts{
			Name: metric + "_count",
			Help: "Count of requests by country",
		})
		metricTiming = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "simple",
			Name:      metric + "_duration_seconds",
			Help:      "Histogram of response time for handler in seconds",
			Buckets:   buckets,
		}, []string{"method", "status_code"})
	)

	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		metricsCount.Inc()
		next(ctx)
		duration := time.Since(start)
		statusCode := strconv.Itoa(ctx.Response.StatusCode())
		metricTiming.WithLabelValues(string(ctx.Method()), statusCode).Observe(duration.Seconds())
	}
}
