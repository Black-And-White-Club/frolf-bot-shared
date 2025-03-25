package attr

import (
	"context"
	"log/slog"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/ThreeDotsLabs/watermill/message"
)

// LogAttr is our custom attribute type that wraps slog.Attr
// This provides abstraction from the underlying logging implementation
type LogAttr struct {
	attr slog.Attr
}

// ToSlogAttr converts our LogAttr to slog.Attr for internal use
func (la LogAttr) ToSlogAttr() slog.Attr {
	return la.attr
}

// --- Attribute constructors now return LogAttr ---

// CorrelationIDFromMsg creates a LogAttr for the correlation ID from message metadata.
func CorrelationIDFromMsg(msg *message.Message) LogAttr {
	return LogAttr{slog.String("correlation_id", msg.Metadata.Get("correlation_id"))}
}

// UserID creates a LogAttr for a user ID.
func UserID(userID sharedtypes.DiscordID) LogAttr {
	return LogAttr{slog.String("user_id", string(userID))}
}

// MessageID creates a LogAttr for a Watermill message ID.
func MessageID(msg *message.Message) LogAttr {
	return LogAttr{slog.String("message_id", msg.UUID)}
}

// Error creates a LogAttr for an error.
func Error(err error) LogAttr {
	return LogAttr{slog.Any("error", err)}
}

// EventName creates a LogAttr for an event name.
func EventName(event string) LogAttr {
	return LogAttr{slog.String("event", event)}
}

// Topic creates a LogAttr for a topic.
func Topic(topic string) LogAttr {
	return LogAttr{slog.String("topic", topic)}
}

// String creates a generic string attribute.
func String(key, value string) LogAttr {
	return LogAttr{slog.String(key, value)}
}

// Int creates a generic int attribute.
func Int(key string, value int) LogAttr {
	return LogAttr{slog.Int(key, value)}
}

// Int64 creates a generic int64 attribute.
func Int64(key string, value int64) LogAttr {
	return LogAttr{slog.Int64(key, value)}
}

// Uint64 creates a generic uint64 attribute.
func Uint64(key string, value uint64) LogAttr {
	return LogAttr{slog.Uint64(key, value)}
}

// Float64 creates a generic float64 attribute.
func Float64(key string, value float64) LogAttr {
	return LogAttr{slog.Float64(key, value)}
}

// Bool creates a generic bool attribute.
func Bool(key string, value bool) LogAttr {
	return LogAttr{slog.Bool(key, value)}
}

// Time creates a generic time attribute.
func Time(key string, value time.Time) LogAttr {
	return LogAttr{slog.Time(key, value)}
}

// Duration creates a generic duration attribute.
func Duration(key string, value time.Duration) LogAttr {
	return LogAttr{slog.Duration(key, value)}
}

// Any creates a generic interface attribute.
func Any(key string, value interface{}) LogAttr {
	return LogAttr{slog.Any(key, value)}
}

// Group creates a grouped attribute.
func Group(key string, attrs ...LogAttr) LogAttr {
	slogAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		slogAttrs[i] = attr.ToSlogAttr()
	}
	return LogAttr{slog.Group(key, slogAttrs...)}
}

// ExtractCorrelationID pulls the correlation ID from context, or sets a default.
func ExtractCorrelationID(ctx context.Context) LogAttr {
	if ctx == nil {
		return LogAttr{slog.String("correlation_id", "unknown")} // Safe fallback
	}

	// Assume correlation_id is stored in context under this key
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		return LogAttr{slog.String("correlation_id", correlationID)}
	}

	return LogAttr{slog.String("correlation_id", "unknown")} // Default if missing
}
