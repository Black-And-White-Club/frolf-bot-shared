package guildmetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.

type NoOpMetrics struct{}

func NewNoop() GuildMetrics {
	return &NoOpMetrics{}
}

func (n *NoOpMetrics) RecordGuildCreated(ctx context.Context, success bool, guildID sharedtypes.GuildID, source string) {
}

func (n *NoOpMetrics) RecordGuildDeleted(ctx context.Context, success bool, guildID sharedtypes.GuildID, source string) {
}

func (n *NoOpMetrics) RecordOperationAttempt(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationSuccess(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationFailure(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationDuration(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
}
