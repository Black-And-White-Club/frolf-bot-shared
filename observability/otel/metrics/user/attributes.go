package usermetrics

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
)

// Common attributes for user metrics (return individual KeyValue instead of Set)

func handlerAttr(handlerName string) attribute.KeyValue {
	return attribute.String("handler", handlerName)
}

func userTypeAttr(userType string) attribute.KeyValue {
	return attribute.String("user_type", userType)
}

func sourceAttr(source string) attribute.KeyValue {
	return attribute.String("source", source)
}

func userIDAttr(userID sharedtypes.DiscordID) attribute.KeyValue {
	return attribute.String("user_id", string(userID))
}

func tagNumberAttr(tagNumber sharedtypes.TagNumber) attribute.KeyValue {
	return attribute.Int64("tag_number", int64(tagNumber))
}

func roleAttr(role sharedtypes.UserRoleEnum) attribute.KeyValue {
	return attribute.String("role", role.String())
}

func actionAttr(action string) attribute.KeyValue {
	return attribute.String("action", action)
}

func resourceAttr(resource string) attribute.KeyValue {
	return attribute.String("resource", resource)
}

func operationAttr(operationName string) attribute.KeyValue {
	return attribute.String("operation", operationName)
}

func handlerAttrs(handlerName string) attribute.KeyValue {
	return handlerAttr(handlerName)
}
