package prometheus

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/prometheus/client_golang/prometheus"

	"api-gateway/internal/components/utils"
)

// metrics define
var (
	statusCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_http_status_total",
		Help: "Request status count.",
	}, []string{"status"})
	reqCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gateway_http_request_total",
		Help: "Request count.",
	})
	serviceTimeCostCounter = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "gateway_http_service_time_cost",
			Help:       "Request time cost by service.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
)

var (
	// Middleware for goframe http server.
	// collecting metrics
	Middleware = func(r *ghttp.Request) {
		start := time.Now()
		r.Middleware.Next()

		// collect
		reqCounter.Inc()
		serviceTimeCostCounter.WithLabelValues(utils.ParseRoutingKey(r.Request.RequestURI)).
			Observe(float64(time.Now().Sub(start).Milliseconds()))
		statusCounter.WithLabelValues(fmt.Sprintf("%d", r.Response.Status)).Inc()
	}
)

func init() {
	prometheus.MustRegister(statusCounter, reqCounter, serviceTimeCostCounter)
}
