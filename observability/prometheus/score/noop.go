package scoremetrics

import sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func (n *NoOpMetrics) RecordScoreProcessingAttempt(roundID sharedtypes.RoundID)                    {}
func (n *NoOpMetrics) RecordScoreProcessingSuccess(roundID sharedtypes.RoundID)                    {}
func (n *NoOpMetrics) RecordScoreProcessingFailure(roundID sharedtypes.RoundID)                    {}
func (n *NoOpMetrics) RecordScoreProcessingDuration(roundID sharedtypes.RoundID, duration float64) {}
func (n *NoOpMetrics) RecordScoreCorrectionAttempt(roundID sharedtypes.RoundID)                    {}
func (n *NoOpMetrics) RecordScoreCorrectionSuccess(roundID sharedtypes.RoundID)                    {}
func (n *NoOpMetrics) RecordScoreCorrectionFailure(roundID sharedtypes.RoundID)                    {}
func (n *NoOpMetrics) RecordScoreCorrectionDuration(roundID sharedtypes.RoundID, duration float64) {}
func (n *NoOpMetrics) RecordLeaderboardUpdateAttempt(roundID sharedtypes.RoundID)                  {}
func (n *NoOpMetrics) RecordLeaderboardUpdateSuccess(roundID sharedtypes.RoundID)                  {}
func (n *NoOpMetrics) RecordLeaderboardUpdateFailure(roundID sharedtypes.RoundID)                  {}
func (n *NoOpMetrics) RecordLeaderboardUpdateDuration(roundID sharedtypes.RoundID, duration float64) {
}
func (n *NoOpMetrics) RecordDBQueryDuration(duration float64) {}
