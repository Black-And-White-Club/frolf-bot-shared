// observability/opentelemetry/guild/attributes.go
package guildmetrics

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
)

// Common attributes for guild metrics (return individual KeyValue instead of Set)

func successAttr(success bool) attribute.KeyValue {
	return attribute.Bool("success", success)
}

func sourceAttr(source string) attribute.KeyValue {
	return attribute.String("source", source)
}

func guildIDAttr(guildID sharedtypes.GuildID) attribute.KeyValue {
	return attribute.String("guild_id", string(guildID))
}

func operationAttr(operationName string) attribute.KeyValue {
	return attribute.String("operation", operationName)
}

func serviceAttr(serviceName string) attribute.KeyValue {
	return attribute.String("service", serviceName)
}
