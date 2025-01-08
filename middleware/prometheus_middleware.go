package middleware

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total jumlah HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Durasi request HTTP",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	// Registrasi metrik ke Prometheus registry
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(requestDuration)
}

func normalizeEndpoint(endpoint string) string {
	re := regexp.MustCompile(`:[a-zA-Z0-9_]+`)
	return re.ReplaceAllString(endpoint, "{id}")
}

// PrometheusMiddleware adalah middleware untuk mencatat metrik Prometheus
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Lanjutkan ke handler berikutnya
		c.Next()

		// Setelah handler selesai, catat metrik
		status := fmt.Sprintf("%d", c.Writer.Status())
		duration := time.Since(start).Seconds()
		normalizedEndpoint := normalizeEndpoint(c.FullPath())
		httpRequestsTotal.WithLabelValues(c.Request.Method, normalizedEndpoint, status).Inc()
		requestDuration.WithLabelValues(c.Request.Method, normalizedEndpoint, status).Observe(duration)
	}
}
