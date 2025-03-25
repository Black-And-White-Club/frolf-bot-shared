package usermetrics

import sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func (n *NoOpMetrics) RecordOperationAttempt(operationName string, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) RecordOperationDuration(operationName string, duration float64)            {}
func (n *NoOpMetrics) RecordOperationFailure(operationName string, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) RecordOperationSuccess(operationName string, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) UserCreationByTag(tag int)                                                 {}
func (n *NoOpMetrics) RecordUserCreation(userType, source, status string)                        {}
func (n *NoOpMetrics) RecordUserRoleUpdateAttempt(userID sharedtypes.DiscordID, role sharedtypes.UserRoleEnum) {
}
func (n *NoOpMetrics) RecordUserRoleUpdateFailure(userID sharedtypes.DiscordID, role sharedtypes.UserRoleEnum) {
}
func (n *NoOpMetrics) RecordUserRoleUpdateSuccess(userID sharedtypes.DiscordID, role sharedtypes.UserRoleEnum) {
}
func (n *NoOpMetrics) DBQueryDuration(duration float64)                                   {}
func (n *NoOpMetrics) UserCreationDuration(duration float64)                              {}
func (n *NoOpMetrics) RecordTagAvailabilityCheck(available bool, tag int)                 {}
func (n *NoOpMetrics) RecordUserRetrieval(success bool, userID sharedtypes.DiscordID)     {}
func (n *NoOpMetrics) RecordUserRoleRetrieval(success bool, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) RecordUserRetrievalDuration(duration float64)                       {}
func (n *NoOpMetrics) RecordPermissionCheck(role sharedtypes.UserRoleEnum, allowed bool, action string, resource string) {
}
func (n *NoOpMetrics) RecordRoleUpdate(oldRole, newRole sharedtypes.UserRoleEnum, context string, userID sharedtypes.DiscordID) {
}
func (n *NoOpMetrics) RecordHandlerAttempt(handlerName string)                    {}
func (n *NoOpMetrics) RecordHandlerDuration(handlerName string, duration float64) {}
func (n *NoOpMetrics) RecordHandlerFailure(handlerName string)                    {}
func (n *NoOpMetrics) RecordHandlerSuccess(handlerName string)                    {}
