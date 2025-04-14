// observability/otel/eventbusmetrics/constructor.go
package eventbusmetrics

import "go.opentelemetry.io/otel/metric"

// NewEventBusMetrics creates a new EventBusMetrics implementation using OpenTelemetry.
// It requires an OTEL Meter instance.
func NewEventBusMetrics(meter metric.Meter, prefix string) (EventBusMetrics, error) {
	// Helper function to create metric names with prefix
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_eventbus_" + name
		}
		return "eventbus_" + name
	}

	var err error
	m := &eventBusMetrics{meter: meter}

	// Message Publish Metrics
	m.messagePublishCounter, err = meter.Int64Counter(
		metricName("messages_published_total"),
		metric.WithDescription("Number of messages published, partitioned by topic"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.messagePublishErrorCounter, err = meter.Int64Counter(
		metricName("messages_publish_errors_total"),
		metric.WithDescription("Number of message publish errors, partitioned by topic"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Message Process Metrics
	m.messageProcessCounter, err = meter.Int64Counter(
		metricName("messages_processed_total"),
		metric.WithDescription("Number of messages processed, partitioned by topic and success"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.messageProcessingTimeHistogram, err = meter.Float64Histogram(
		metricName("message_processing_time_seconds"),
		metric.WithDescription("Time taken to process a message, partitioned by topic"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// Message Subscribe Metrics
	m.messageSubscribeCounter, err = meter.Int64Counter(
		metricName("messages_subscribed_total"),
		metric.WithDescription("Number of messages subscribed, partitioned by topic"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.messageSubscribeErrorCounter, err = meter.Int64Counter(
		metricName("messages_subscribe_errors_total"),
		metric.WithDescription("Number of subscription errors, partitioned by topic"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
