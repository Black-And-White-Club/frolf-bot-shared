package scoremetrics

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
)

const (
	NoTagValue = "no_tag"
)

// Common attributes
func roundIDAttr(roundID sharedtypes.RoundID) attribute.KeyValue {
	_ = roundID
	return attribute.String("scope", "all_rounds")
}

func userIDAttr(userID sharedtypes.DiscordID) attribute.KeyValue {
	_ = userID
	return attribute.String("subject", "user")
}

func operationAttr(operationName string) attribute.KeyValue {
	return attribute.String("operation", operationName)
}

func handlerAttr(handlerName string) attribute.KeyValue {
	return attribute.String("handler", handlerName)
}

func tagNumberAttr(tagNumber int) attribute.KeyValue {
	return attribute.Int("tag_number", tagNumber)
}

func roundAttrs(roundID sharedtypes.RoundID) attribute.KeyValue {
	return roundIDAttr(roundID)
}

func roundUserAttrs(roundID sharedtypes.RoundID, userID sharedtypes.DiscordID) []attribute.KeyValue {
	return []attribute.KeyValue{
		roundIDAttr(roundID),
		userIDAttr(userID),
	}
}

func operationRoundAttrs(operationName string, roundID sharedtypes.RoundID) []attribute.KeyValue {
	return []attribute.KeyValue{
		operationAttr(operationName),
		roundIDAttr(roundID),
	}
}

func operationAttrs(operationName string) attribute.KeyValue {
	return operationAttr(operationName)
}

func handlerAttrs(handlerName string) attribute.KeyValue {
	return handlerAttr(handlerName)
}

func tagAttrs(roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber) []attribute.KeyValue {
	attrs := []attribute.KeyValue{roundIDAttr(roundID)}
	if tagNumber != nil && *tagNumber != 0 {
		attrs = append(attrs, tagNumberAttr(int(*tagNumber)))
	} else {
		attrs = append(attrs, attribute.String("tag_number", NoTagValue))
	}
	return attrs
}

func tagMovementAttrs(roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, fromUserID, toUserID sharedtypes.DiscordID) []attribute.KeyValue {
	_ = fromUserID
	_ = toUserID
	attrs := []attribute.KeyValue{
		roundIDAttr(roundID),
		attribute.String("movement", "user_to_user"),
	}

	if tagNumber != nil && *tagNumber != 0 {
		attrs = append(attrs, tagNumberAttr(int(*tagNumber)))
	} else {
		attrs = append(attrs, attribute.String("tag_number", NoTagValue))
	}
	return attrs
}
