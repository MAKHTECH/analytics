package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusClient реализация провайдера метрик на основе Prometheus
type PrometheusClient struct {
	eventCounter       *prometheus.CounterVec
	processingTimeHist *prometheus.HistogramVec
}

// NewPrometheusClient создает нового клиента Prometheus
func NewPrometheusClient() *PrometheusClient {
	eventCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_events_total",
			Help: "The total number of processed analytics events",
		},
		[]string{"event_type"},
	)

	processingTimeHist := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "analytics_processing_duration_seconds",
			Help:    "Histogram of analytics event processing time in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_type"},
	)

	return &PrometheusClient{
		eventCounter:       eventCounter,
		processingTimeHist: processingTimeHist,
	}
}

// RecordEvent записывает событие в метрики
func (p *PrometheusClient) RecordEvent(eventType string) {
	p.eventCounter.WithLabelValues(eventType).Inc()
}

// RecordProcessingTime записывает время обработки события
func (p *PrometheusClient) RecordProcessingTime(duration time.Duration) {
	p.processingTimeHist.WithLabelValues("all").Observe(duration.Seconds())
}
