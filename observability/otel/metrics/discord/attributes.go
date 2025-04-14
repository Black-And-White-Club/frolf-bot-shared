// observability/otel/discordmetrics/attributes.go
package discordmetrics

import (
	"go.opentelemetry.io/otel/attribute"
)

// Common attributes
func endpointAttr(endpoint string) attribute.KeyValue {
	return attribute.String("endpoint", endpoint)
}

func errorTypeAttr(errorType string) attribute.KeyValue {
	return attribute.String("error_type", errorType)
}

func eventTypeAttr(eventType string) attribute.KeyValue {
	return attribute.String("event_type", eventType)
}

func disconnectReasonAttr(reason string) attribute.KeyValue {
	return attribute.String("reason", reason)
}

// // Combined attribute sets
// func endpointAttrs(endpoint string) attribute.KeyValue {
// 	return endpointAttr(endpoint)
// }

func endpointErrorAttrs(endpoint string, errorType string) []attribute.KeyValue {
	return []attribute.KeyValue{
		endpointAttr(endpoint),
		errorTypeAttr(errorType),
	}
}

// func eventTypeAttrs(eventType string) attribute.KeyValue {
// 	return eventTypeAttr(eventType)
// }

// func disconnectAttrs(reason string) attribute.KeyValue {
// 	return disconnectReasonAttr(reason)
// }
