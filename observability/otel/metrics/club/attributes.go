package clubmetrics

import (
	"go.opentelemetry.io/otel/attribute"
)

func operationAttr(name string) attribute.KeyValue {
	return attribute.String("operation", name)
}

func serviceAttr(name string) attribute.KeyValue {
	return attribute.String("service", name)
}
