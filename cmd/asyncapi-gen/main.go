// Package main generates an AsyncAPI specification from Go event definitions.
//
// Usage:
//
//	go run ./cmd/asyncapi-gen > asyncapi/asyncapi.yaml
//
// This generator uses the swaggest/go-asyncapi library to produce an AsyncAPI 2.4.0
// specification directly from the event struct definitions in the events package.
package main

import (
	"fmt"
	"os"

	guildevents "github.com/Black-And-White-Club/frolf-bot-shared/events/guild"
	leaderboardevents "github.com/Black-And-White-Club/frolf-bot-shared/events/leaderboard"
	roundevents "github.com/Black-And-White-Club/frolf-bot-shared/events/round"
	scoreevents "github.com/Black-And-White-Club/frolf-bot-shared/events/score"
	userevents "github.com/Black-And-White-Club/frolf-bot-shared/events/user"
	asyncapi "github.com/swaggest/go-asyncapi/reflector/asyncapi-2.4.0"
	"github.com/swaggest/go-asyncapi/spec-2.4.0"
)

func main() {
	schema := spec.AsyncAPI{}
	schema.Info.Title = "Frolf Bot Event API"
	schema.Info.Version = "1.0.0"
	schema.Info.Description = `Event-driven architecture for the Frolf Bot ecosystem.

This specification documents all events published and consumed across the Frolf Bot services:
- **frolf-bot** (Backend): Core business logic, database operations, scheduling
- **discord-frolf-bot** (Discord Bot): User interactions, Discord API integration

## Domains
- **Round**: Golf round lifecycle (creation, updates, deletion, participants, scoring)
- **Score**: Score processing and leaderboard updates
- **User**: User registration and tag number management
- **Leaderboard**: Weekly/monthly leaderboard calculations
- **Guild**: Discord server configuration

## Architecture
All events flow through NATS JetStream with the Event Notification pattern.
Each event contains a complete snapshot of the data needed for processing.`

	schema.DefaultContentType = "application/json"

	// Server configuration
	schema.AddServer("production", spec.Server{
		URL:         "nats://nats:4222",
		Description: "Production NATS JetStream server",
		Protocol:    "nats",
	})
	schema.AddServer("development", spec.Server{
		URL:         "nats://localhost:4222",
		Description: "Local development NATS server",
		Protocol:    "nats",
	})

	reflector := asyncapi.Reflector{}
	reflector.Schema = &schema

	mustNotFail := func(err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}

	// ==========================================================================
	// ROUND DOMAIN EVENTS
	// ==========================================================================

	// Round Creation Flow
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundCreationRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round creation initiated by user via Discord",
				Summary:     "Round Creation Request",
			},
			MessageSample: new(roundevents.CreateRoundRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundValidationPassedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round creation input passed validation",
				Summary:     "Round Validation Passed",
			},
			MessageSample: new(roundevents.RoundValidationPassedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundValidationFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round creation input failed validation",
				Summary:     "Round Validation Failed",
			},
			MessageSample: new(roundevents.RoundValidationFailedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundDateTimeParsedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Natural language datetime successfully parsed",
				Summary:     "DateTime Parsed",
			},
			MessageSample: new(roundevents.RoundDateTimeParsedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundEntityCreatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round entity created in memory",
				Summary:     "Round Entity Created",
			},
			MessageSample: new(roundevents.RoundEntityCreatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundStoredV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round persisted to database",
				Summary:     "Round Stored",
			},
			MessageSample: new(roundevents.RoundStoredPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundScheduledV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round reminders scheduled",
				Summary:     "Round Scheduled",
			},
			MessageSample: new(roundevents.RoundScheduledPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundCreatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round creation completed successfully",
				Summary:     "Round Created",
			},
			MessageSample: new(roundevents.RoundCreatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundCreationFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round creation failed",
				Summary:     "Round Creation Failed",
			},
			MessageSample: new(roundevents.RoundCreationFailedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundErrorV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Generic round operation error",
				Summary:     "Round Error",
			},
			MessageSample: new(roundevents.RoundErrorPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundEventMessageIDUpdateV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to update Discord message ID for round",
				Summary:     "Message ID Update Request",
			},
			MessageSample: new(roundevents.RoundMessageIDUpdatePayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundTraceEventV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Distributed tracing event for observability",
				Summary:     "Trace Event",
			},
			MessageSample: new(map[string]interface{}),
		},
	}))

	// Round Update Flow
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundUpdateRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round update requested by user",
				Summary:     "Round Update Request",
			},
			MessageSample: new(roundevents.UpdateRoundRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundUpdateValidatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round update validated",
				Summary:     "Round Update Validated",
			},
			MessageSample: new(roundevents.RoundUpdateValidatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundFetchedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round fetched from database for update",
				Summary:     "Round Fetched",
			},
			MessageSample: new(roundevents.RoundFetchedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundEntityUpdatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round entity updated in database",
				Summary:     "Round Entity Updated",
			},
			MessageSample: new(roundevents.RoundEntityUpdatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundUpdateSuccessV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round update completed successfully",
				Summary:     "Round Update Success",
			},
			MessageSample: new(roundevents.RoundUpdateSuccessPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundUpdateErrorV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round update failed",
				Summary:     "Round Update Error",
			},
			MessageSample: new(roundevents.RoundUpdateErrorPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundScheduleUpdatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round schedule updated",
				Summary:     "Round Schedule Updated",
			},
			MessageSample: new(roundevents.RoundScheduleUpdatePayloadV1),
		},
	}))

	// Round Delete Flow
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundDeleteRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round deletion requested by user",
				Summary:     "Round Delete Request",
			},
			MessageSample: new(roundevents.RoundDeleteRequestPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundDeleteValidatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round deletion request validated",
				Summary:     "Round Delete Validated",
			},
			MessageSample: new(roundevents.RoundDeleteValidatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundToDeleteFetchedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round fetched for deletion verification",
				Summary:     "Round To Delete Fetched",
			},
			MessageSample: new(roundevents.RoundToDeleteFetchedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundDeleteAuthorizedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User authorized to delete round",
				Summary:     "Round Delete Authorized",
			},
			MessageSample: new(roundevents.RoundDeleteAuthorizedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundDeletedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round deleted successfully",
				Summary:     "Round Deleted",
			},
			MessageSample: new(roundevents.RoundDeletedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundDeleteErrorV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round deletion failed",
				Summary:     "Round Delete Error",
			},
			MessageSample: new(roundevents.RoundDeleteErrorPayloadV1),
		},
	}))

	// Participant Events
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantJoinRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User requested to join a round",
				Summary:     "Participant Join Request",
			},
			MessageSample: new(roundevents.ParticipantJoinRequestPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantJoinValidatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Participant join request validated",
				Summary:     "Participant Join Validated",
			},
			MessageSample: new(roundevents.ParticipantJoinValidatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantJoinedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Participant successfully joined round",
				Summary:     "Participant Joined",
			},
			MessageSample: new(roundevents.ParticipantJoinedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantDeclinedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Participant declined to join round",
				Summary:     "Participant Declined",
			},
			MessageSample: new(roundevents.ParticipantDeclinedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantJoinErrorV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Error occurred while joining round",
				Summary:     "Participant Join Error",
			},
			MessageSample: new(roundevents.RoundParticipantJoinErrorPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantRemovalRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to remove participant from round",
				Summary:     "Participant Removal Request",
			},
			MessageSample: new(roundevents.ParticipantRemovalRequestPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantRemovedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Participant removed from round",
				Summary:     "Participant Removed",
			},
			MessageSample: new(roundevents.ParticipantRemovedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantRemovalErrorV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Error occurred while removing participant",
				Summary:     "Participant Removal Error",
			},
			MessageSample: new(roundevents.ParticipantRemovalErrorPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantStatusUpdateRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to update participant status",
				Summary:     "Participant Status Update Request",
			},
			MessageSample: new(roundevents.ParticipantStatusRequestPayloadV1),
		},
	}))

	// Score Events
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundScoreUpdateRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Score update requested for participant",
				Summary:     "Score Update Request",
			},
			MessageSample: new(roundevents.ScoreUpdateRequestPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundScoreUpdateValidatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Score update validated",
				Summary:     "Score Update Validated",
			},
			MessageSample: new(roundevents.ScoreUpdateValidatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundParticipantScoreUpdatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Participant score updated successfully",
				Summary:     "Participant Score Updated",
			},
			MessageSample: new(roundevents.ParticipantScoreUpdatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundScoreUpdateErrorV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Score update failed",
				Summary:     "Score Update Error",
			},
			MessageSample: new(roundevents.RoundScoreUpdateErrorPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundAllScoresSubmittedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "All participants have submitted scores",
				Summary:     "All Scores Submitted",
			},
			MessageSample: new(roundevents.AllScoresSubmittedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundScoresPartiallySubmittedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Some participants still missing scores",
				Summary:     "Scores Partially Submitted",
			},
			MessageSample: new(roundevents.ScoresPartiallySubmittedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.ProcessRoundScoresRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to process round scores for leaderboard",
				Summary:     "Process Round Scores Request",
			},
			MessageSample: new(roundevents.ProcessRoundScoresRequestPayloadV1),
		},
	}))

	// Lifecycle Events
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundStartedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round has started (scheduled time reached)",
				Summary:     "Round Started",
			},
			MessageSample: new(roundevents.RoundStartedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundStartedDiscordV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round started notification for Discord",
				Summary:     "Round Started Discord",
			},
			MessageSample: new(roundevents.DiscordRoundStartPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundFinalizedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round finalized in database",
				Summary:     "Round Finalized",
			},
			MessageSample: new(roundevents.RoundFinalizedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundFinalizedDiscordV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round finalized notification for Discord embed update",
				Summary:     "Round Finalized Discord",
			},
			MessageSample: new(roundevents.RoundFinalizedDiscordPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundReminderSentV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Round reminder sent",
				Summary:     "Round Reminder Sent",
			},
			MessageSample: new(roundevents.DiscordReminderPayloadV1),
		},
	}))

	// Retrieval Events
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.GetRoundRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to retrieve round details",
				Summary:     "Round Retrieval Request",
			},
			MessageSample: new(roundevents.GetRoundRequestPayloadV1),
		},
	}))

	// Tag Events
	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: roundevents.RoundTagNumberRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to look up user tag number",
				Summary:     "Tag Lookup Request",
			},
			MessageSample: new(roundevents.TagNumberRequestPayloadV1),
		},
	}))

	// ==========================================================================
	// SCORE DOMAIN EVENTS
	// ==========================================================================

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: scoreevents.ProcessRoundScoresRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to process scores for a round",
				Summary:     "Score Processing Request",
			},
			MessageSample: new(scoreevents.ProcessRoundScoresRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: scoreevents.ProcessRoundScoresSucceededV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Score processing completed",
				Summary:     "Score Processing Completed",
			},
			MessageSample: new(scoreevents.ProcessRoundScoresSucceededPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: scoreevents.ProcessRoundScoresFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Score processing failed",
				Summary:     "Score Processing Error",
			},
			MessageSample: new(scoreevents.ProcessRoundScoresFailedPayloadV1),
		},
	}))

	// ==========================================================================
	// USER DOMAIN EVENTS
	// ==========================================================================

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.UserCreationRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to create new user",
				Summary:     "User Creation Request",
			},
			MessageSample: new(userevents.UserCreationRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.UserCreatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User created successfully",
				Summary:     "User Created",
			},
			MessageSample: new(userevents.UserCreatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.UserCreationFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User creation failed",
				Summary:     "User Creation Error",
			},
			MessageSample: new(userevents.UserCreationFailedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.GetUserRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to retrieve user details",
				Summary:     "User Retrieval Request",
			},
			MessageSample: new(userevents.GetUserRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.GetUserResponseV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User details retrieved",
				Summary:     "User Retrieved",
			},
			MessageSample: new(userevents.GetUserResponsePayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.GetUserFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User not found",
				Summary:     "User Not Found",
			},
			MessageSample: new(userevents.GetUserFailedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.UserRoleUpdateRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to update user role",
				Summary:     "User Role Update Request",
			},
			MessageSample: new(userevents.UserRoleUpdateRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: userevents.UserRoleUpdatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "User role updated successfully",
				Summary:     "User Role Updated",
			},
			MessageSample: new(userevents.UserRoleUpdatedPayloadV1),
		},
	}))

	// ==========================================================================
	// LEADERBOARD DOMAIN EVENTS
	// ==========================================================================

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.LeaderboardUpdateRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to update leaderboard with new scores",
				Summary:     "Leaderboard Update Request",
			},
			MessageSample: new(leaderboardevents.LeaderboardUpdateRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.LeaderboardUpdatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Leaderboard updated successfully",
				Summary:     "Leaderboard Updated",
			},
			MessageSample: new(leaderboardevents.LeaderboardUpdatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.LeaderboardUpdateFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Leaderboard update failed",
				Summary:     "Leaderboard Update Error",
			},
			MessageSample: new(leaderboardevents.LeaderboardUpdateFailedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.GetLeaderboardRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to retrieve leaderboard",
				Summary:     "Leaderboard Retrieval Request",
			},
			MessageSample: new(leaderboardevents.GetLeaderboardRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.GetLeaderboardResponseV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Leaderboard retrieved successfully",
				Summary:     "Leaderboard Retrieved",
			},
			MessageSample: new(leaderboardevents.GetLeaderboardResponsePayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.TagSwapRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to swap tag numbers between users",
				Summary:     "Tag Swap Request",
			},
			MessageSample: new(leaderboardevents.TagSwapRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.TagSwapProcessedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Tag swap completed successfully",
				Summary:     "Tag Swap Completed",
			},
			MessageSample: new(leaderboardevents.TagSwapProcessedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: leaderboardevents.TagSwapFailedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Tag swap failed",
				Summary:     "Tag Swap Error",
			},
			MessageSample: new(leaderboardevents.TagSwapFailedPayloadV1),
		},
	}))

	// ==========================================================================
	// GUILD DOMAIN EVENTS
	// ==========================================================================

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: guildevents.GuildConfigRetrievalRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to retrieve guild configuration",
				Summary:     "Guild Config Retrieval Request",
			},
			MessageSample: new(guildevents.GuildConfigRetrievalRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: guildevents.GuildConfigRetrievedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Guild configuration retrieved successfully",
				Summary:     "Guild Config Retrieved",
			},
			MessageSample: new(guildevents.GuildConfigRetrievedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: guildevents.GuildConfigCreationRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to create guild configuration",
				Summary:     "Guild Config Creation Request",
			},
			MessageSample: new(guildevents.GuildConfigCreationRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: guildevents.GuildConfigCreatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Guild configuration created successfully",
				Summary:     "Guild Config Created",
			},
			MessageSample: new(guildevents.GuildConfigCreatedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: guildevents.GuildConfigUpdateRequestedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Request to update guild configuration",
				Summary:     "Guild Config Update Request",
			},
			MessageSample: new(guildevents.GuildConfigUpdateRequestedPayloadV1),
		},
	}))

	mustNotFail(reflector.AddChannel(asyncapi.ChannelInfo{
		Name: guildevents.GuildConfigUpdatedV1,
		Publish: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Guild configuration updated successfully",
				Summary:     "Guild Config Updated",
			},
			MessageSample: new(guildevents.GuildConfigUpdatedPayloadV1),
		},
	}))

	// Generate YAML output
	yaml, err := reflector.Schema.MarshalYAML()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating AsyncAPI spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(yaml))
}
