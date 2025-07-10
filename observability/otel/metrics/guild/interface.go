package guildmetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// GuildMetrics defines metrics specific to guild operations
type GuildMetrics interface {
	RecordGuildCreated(ctx context.Context, success bool, guildID sharedtypes.GuildID, source string)
	RecordGuildDeleted(ctx context.Context, success bool, guildID sharedtypes.GuildID, source string)
	RecordOperationAttempt(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string)
	RecordOperationSuccess(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string)
	RecordOperationFailure(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string)
	RecordOperationDuration(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string, duration time.Duration)
	RecordHandlerAttempt(ctx context.Context, handlerName string)
	RecordHandlerSuccess(ctx context.Context, handlerName string)
	RecordHandlerFailure(ctx context.Context, handlerName string)
	RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration)
}
