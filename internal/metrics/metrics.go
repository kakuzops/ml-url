package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	urlShorteningTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_shortening_total",
			Help: "Total number of URLs shortened",
		},
	)

	urlRedirectsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "url_redirects_total",
			Help: "Total number of URL redirects",
		},
	)

	activeURLs = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_urls",
			Help: "Current number of active shortened URLs",
		},
	)
)

// MetricsMiddleware retorna um middleware Gin que coleta métricas
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Prossegue com o request
		c.Next()

		// Não coletar métricas para o endpoint do prometheus
		if c.Request.URL.Path == "/metrics" {
			return
		}

		// Registra a duração
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)

		// Registra o total de requests
		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(c.Writer.Status())),
		).Inc()

		// Registra métricas específicas
		switch {
		case c.FullPath() == "/shorten" && c.Request.Method == "POST":
			if c.Writer.Status() == 201 {
				urlShorteningTotal.Inc()
				activeURLs.Inc()
			}
		case c.FullPath() == "/:shortURL" && c.Request.Method == "GET":
			if c.Writer.Status() == 301 {
				urlRedirectsTotal.Inc()
			}
		}
	}
}

// IncrementActiveURLs incrementa o contador de URLs ativas
func IncrementActiveURLs() {
	activeURLs.Inc()
}

// DecrementActiveURLs decrementa o contador de URLs ativas
func DecrementActiveURLs() {
	activeURLs.Dec()
} 