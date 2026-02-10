package sharedevents

// PointsAwardedV1 is published after round points have been calculated and persisted.
// This is a shared event consumed by display services (like Discord) to update embeds.
//
// Pattern: Event Notification
// Subject: round.points.awarded.v1
// Producer: leaderboard-service (after ProcessRound)
// Consumers: discord-service (update finalized embed with point values)
// Version: v1 (February 2026)
const PointsAwardedV1 = "round.points.awarded.v1"
