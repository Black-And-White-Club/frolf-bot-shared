// observability/prometheus/eventbusmetrics/metrics.go
package eventbusmetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type EventBusMetrics interface {
	RecordMessagePublish(topic string)
	RecordMessagePublishError(topic string)
	RecordMessageProcess(topic string, success bool)
	RecordMessageProcessingTime(topic string, duration float64)
	RecordMessageSubscribe(topic string)
	RecordMessageSubscribeError(topic string)
}

type eventBusMetrics struct {
	messagePublishCounter        *prometheus.CounterVec
	messageProcessCounter        *prometheus.CounterVec
	messageProcessingTimeHist    *prometheus.HistogramVec
	messageSubscribeCounter      *prometheus.CounterVec
	messageSubscribeErrorCounter *prometheus.CounterVec
}

func NewEventBusMetrics(registry *prometheus.Registry, prefix string) EventBusMetrics {
	messagePublishCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "eventbus",
			Name:      "messages_published_total",
			Help:      "Number of messages published, partitioned by topic",
		},
		[]string{"topic"},
	)

	messageProcessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "eventbus",
			Name:      "messages_processed_total",
			Help:      "Number of messages processed, partitioned by topic and success",
		},
		[]string{"topic", "success"},
	)

	messageProcessingTimeHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "eventbus",
			Name:      "message_processing_time_seconds",
			Help:      "Time taken to process a message, partitioned by topic",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	messageSubscribeCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "eventbus",
			Name:      "messages_subscribed_total",
			Help:      "Number of messages subscribed, partitioned by topic",
		},
		[]string{"topic"},
	)

	messageSubscribeErrorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "eventbus",
			Name:      "messages_subscribe_errors_total",
			Help:      "Number of subscription errors, partitioned by topic",
		},
		[]string{"topic"},
	)

	registry.MustRegister(
		messagePublishCounter,
		messageProcessCounter,
		messageProcessingTimeHist,
		messageSubscribeCounter,
		messageSubscribeErrorCounter,
	)

	return &eventBusMetrics{
		messagePublishCounter:        messagePublishCounter,
		messageProcessCounter:        messageProcessCounter,
		messageProcessingTimeHist:    messageProcessingTimeHist,
		messageSubscribeCounter:      messageSubscribeCounter,
		messageSubscribeErrorCounter: messageSubscribeErrorCounter,
	}
}

func (m *eventBusMetrics) RecordMessagePublish(topic string) {
	m.messagePublishCounter.WithLabelValues(topic).Inc()
}

func (m *eventBusMetrics) RecordMessageProcess(topic string, success bool) {
	m.messageProcessCounter.WithLabelValues(topic, fmt.Sprintf("%t", success)).Inc()
}

func (m *eventBusMetrics) RecordMessageProcessingTime(topic string, duration float64) {
	m.messageProcessingTimeHist.WithLabelValues(topic).Observe(duration)
}

func (m *eventBusMetrics) RecordMessagePublishError(topic string) {
	m.messagePublishCounter.WithLabelValues(topic).Inc()
}

func (m *eventBusMetrics) RecordMessageSubscribe(topic string) {
	m.messageSubscribeCounter.WithLabelValues(topic).Inc()
}

func (m *eventBusMetrics) RecordMessageSubscribeError(topic string) {
	m.messageSubscribeErrorCounter.WithLabelValues(topic).Inc()
}
