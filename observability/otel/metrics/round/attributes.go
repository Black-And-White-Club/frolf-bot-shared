package roundmetrics

import (
	"go.opentelemetry.io/otel/attribute"
)

// Common attributes for round metrics
func operationAttr(operation, service string) attribute.KeyValue {
	return attribute.String("operation", operation)
}

// func serviceAttr(service string) attribute.KeyValue {
// 	return attribute.String("service", service)
// }

func locationAttr(location string) attribute.KeyValue {
	return attribute.String("location", location)
}

func dbOperationAttr(operation string) attribute.KeyValue {
	return attribute.String("operation", operation)
}

func handlerNameAttr(handlerName string) attribute.KeyValue {
	return attribute.String("handler_name", handlerName)
}

func operationAttrs(operation, service string) []attribute.KeyValue {
	return []attribute.KeyValue{
		operationAttr(operation, service),
	}
}

func locationAttrs(location string) []attribute.KeyValue {
	return []attribute.KeyValue{
		locationAttr(location),
	}
}

func dbOperationAttrs(operation string) []attribute.KeyValue {
	return []attribute.KeyValue{
		dbOperationAttr(operation),
	}
}

func handlerNameAttrs(handlerName string) []attribute.KeyValue {
	return []attribute.KeyValue{
		handlerNameAttr(handlerName),
	}
}
