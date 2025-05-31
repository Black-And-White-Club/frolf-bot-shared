// observability/opentelemetry/leaderboard/attributes.go
package leaderboardmetrics

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
)

// Common attributes for leaderboard metrics (return individual KeyValue instead of Set)

func successAttr(success bool) attribute.KeyValue {
	return attribute.Bool("success", success)
}

func sourceAttr(source string) attribute.KeyValue {
	return attribute.String("source", source)
}

func roundIDAttr(roundID sharedtypes.RoundID) attribute.KeyValue {
	return attribute.String("round_id", roundID.String())
}

func tagNumberAttr(tagNumber sharedtypes.TagNumber) attribute.KeyValue {
	return attribute.Int64("tag_number", int64(tagNumber))
}

func operationAttr(operationName string) attribute.KeyValue {
	return attribute.String("operation", operationName)
}

func serviceAttr(serviceName string) attribute.KeyValue {
	return attribute.String("service", serviceName)
}

func availableAttr(available bool) attribute.KeyValue {
	return attribute.Bool("available", available)
}

func requestorIDAttr(requestorID sharedtypes.DiscordID) attribute.KeyValue {
	return attribute.String("requestor_id", string(requestorID))
}

func targetIDAttr(targetID sharedtypes.DiscordID) attribute.KeyValue {
	return attribute.String("target_id", string(targetID))
}

func reasonAttr(reason string) attribute.KeyValue {
	return attribute.String("reason", reason)
}

func userIDAttr(userID sharedtypes.DiscordID) attribute.KeyValue {
	return attribute.String("user_id", string(userID))
}

func oldTagAttr(oldTag sharedtypes.TagNumber) attribute.KeyValue {
	return attribute.Int64("old_tag", int64(oldTag))
}

func newTagAttr(newTag sharedtypes.TagNumber) attribute.KeyValue {
	return attribute.Int64("new_tag", int64(newTag))
}

func handlerAttr(handlerName string) attribute.KeyValue {
	return attribute.String("handler", handlerName)
}
