package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/eventbus"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// ErrorEventPayload defines the structure of an error event payload.
type ErrorEventPayload struct {
	CorrelationID string    `json:"correlation_id"`
	Message       string    `json:"message"`
	Error         string    `json:"error,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Context       string    `json:"context"`
	ErrorType     string    `json:"error_type,omitempty"`
}

// DefaultErrorTopic is the default topic for error events.
const DefaultErrorTopic = "error.frolf.bot"

// CreateErrorEventPayload creates an ErrorEventPayload with function context.
func CreateErrorEventPayload(correlationID, message string, err error, ctx ...string) ErrorEventPayload {
	payload := ErrorEventPayload{
		CorrelationID: correlationID,
		Message:       message,
		Timestamp:     time.Now().UTC(),
		Context:       strings.Join(ctx, " | "),
	}

	// Capture caller info
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		payload.Context += " | unknown"
		fmt.Println("WARNING: runtime.Caller() failed, function context missing")
	} else {
		functionName := runtime.FuncForPC(pc).Name()
		payload.Context += fmt.Sprintf(" | %s (%s:%d)", functionName, file, line)
	}

	if err != nil {
		payload.Error = err.Error()
		payload.ErrorType = reflect.TypeOf(err).String()
	} else {
		payload.ErrorType = "unknown"
	}

	return payload
}

// ErrorReporter handles error reporting.
type ErrorReporter struct {
	EventBus     eventbus.EventBus
	Logger       slog.Logger
	ErrorTopic   string
	DefaultLevel string
}

// NewErrorReporter creates a new ErrorReporter.
func NewErrorReporter(eventbus eventbus.EventBus, logger slog.Logger, topic, defaultLevel string) *ErrorReporter {
	if topic == "" {
		topic = DefaultErrorTopic
	}
	if defaultLevel == "" {
		defaultLevel = "ERROR"
	}
	return &ErrorReporter{
		EventBus:     eventbus,
		Logger:       logger,
		ErrorTopic:   topic,
		DefaultLevel: defaultLevel,
	}
}

// levelMap maps log level strings to slog.Level values.
var levelMap = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

// ReportError logs the error and publishes it to the event bus.
func (er *ErrorReporter) ReportError(correlationID string, msg string, err error, ctx ...string) {
	payload := CreateErrorEventPayload(correlationID, msg, err, ctx...)

	payloadBytes, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		er.Logger.Error("failed to marshal error payload", "error", marshalErr)
		return
	}

	// Use structured logging instead of With(map[string]string{})
	attrs := []slog.Attr{
		slog.String("application", "discord-frolf-bot"),
		slog.String("correlation_id", correlationID),
		slog.String("error_type", payload.ErrorType),
	}
	level := levelMap[er.DefaultLevel]

	// Convert attrs to []any
	anyAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		anyAttrs[i] = attr
	}

	er.Logger.With(anyAttrs...).Log(context.Background(), level, string(payloadBytes))

	// Publish the error event
	errorMsg := message.NewMessage(correlationID, payloadBytes)
	errorMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
	if publishErr := er.EventBus.Publish(er.ErrorTopic, errorMsg); publishErr != nil {
		er.Logger.Error("failed to publish error event", "error", publishErr)
	}
}
