// observability/otel/eventbusmetrics/attributes.go
package eventbusmetrics

import (
	"go.opentelemetry.io/otel/attribute"
)

// Common attributes
func topicAttr(topic string) attribute.KeyValue {
	return attribute.String("topic", topic)
}

func successAttr(success bool) attribute.KeyValue {
	return attribute.Bool("success", success)
}

// Combined attribute sets
func topicAttrs(topic string) attribute.KeyValue {
	return topicAttr(topic)
}

func topicSuccessAttrs(topic string, success bool) []attribute.KeyValue {
	return []attribute.KeyValue{
		topicAttr(topic),
		successAttr(success),
	}
}
