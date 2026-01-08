# Router Implementation Guide
## Event-Driven Architecture Best Practices

This guide explains how to implement router-owned topic resolution across all routers in the frolf-bot ecosystem.

## Table of Contents
1. [Overview](#overview)
2. [Current Status](#current-status)
3. [Implementation Steps](#implementation-steps)
4. [Testing](#testing)
5. [Deployment](#deployment)

## Overview

### The Problem
Handlers were setting `metadata["topic"]` to control where messages get published. This violated separation of concerns and made routing logic implicit and hard to test.

### The Solution
Routers own topic resolution through explicit `getPublishTopic()` methods. Handlers remain infrastructure-agnostic.

### Current Architecture
```
Handler → CreateResultMessage(payload, topic) → Sets metadata["topic"] → Router reads metadata → Publishes
```

### Target Architecture
```
Handler → CreateResultMessage(payload, topic) → Router.getPublishTopic() → Publishes
```

## Current Status

### ✅ Phase 0: Complete (All 10 Routers)
All routers now have invariant checks that log "MESSAGE DROPPED" if topic resolution fails.

**Files Modified:**
- `discord-frolf-bot/app/user/watermill/router.go`
- `discord-frolf-bot/app/round/watermill/router.go`
- `discord-frolf-bot/app/score/watermill/router.go`
- `discord-frolf-bot/app/leaderboard/watermill/router.go`
- `discord-frolf-bot/app/guild/watermill/router.go`
- `frolf-bot/app/modules/user/infrastructure/router/router.go`
- `frolf-bot/app/modules/round/infrastructure/router/router.go`
- `frolf-bot/app/modules/score/infrastructure/router/router.go`
- `frolf-bot/app/modules/leaderboard/infrastructure/router/router.go`
- `frolf-bot/app/modules/guild/infrastructure/router/router.go`

### 🚧 Phase 1: In Progress (1/10 Routers)
Only `discord-frolf-bot/app/user/watermill/router.go` has router-owned topic resolution implemented.

**Remaining:**
- 4 discord-frolf-bot routers (round, score, leaderboard, guild)
- 5 frolf-bot routers (user, round, score, leaderboard, guild)

## Implementation Steps

### Step 1: Analyze Handler Mappings

For each router, create a mapping of input topics to output topics.

**Example for discord-frolf-bot/app/user:**

```bash
cd discord-frolf-bot/app/user/watermill/handlers
grep -r "CreateResultMessage" . | grep -o '"[^"]*\.v[0-9]"' | sort -u
```

Create a table:

| Handler Function | Input Topic | Output Topic(s) | Notes |
|-----------------|-------------|-----------------|-------|
| HandleUserCreated | user.created.v1 | discord.user.signup.role.add.v1 | Always returns this |
| HandleAddRole | discord.user.signup.role.add.v1 | discord.user.signup.role.added.v1 OR discord.user.signup.role.addition.failed.v1 | Success/failure paths |
| HandleRoleAdded | discord.user.signup.role.added.v1 | (none) | Returns nil |

### Step 2: Implement getPublishTopic()

Add the method BEFORE `RegisterHandlers()` in the router file:

```go
// getPublishTopic resolves the topic to publish for a given handler's returned message.
// This centralizes routing logic in the router (not in handlers or helpers).
func (r *UserRouter) getPublishTopic(handlerName string, msg *message.Message) string {
	// Extract base topic from handlerName format: "discord-user.{topic}"
	// Map handler input topic → output topic(s)

	switch {
	case handlerName == "discord-user."+userevents.UserCreatedV1:
		// HandleUserCreated always returns SignupAddRoleV1
		return shareduserevents.SignupAddRoleV1

	case handlerName == "discord-user."+shareduserevents.SignupAddRoleV1:
		// HandleAddRole returns either success or failure
		// Use metadata temporarily for complex cases
		return msg.Metadata.Get("topic")

	case handlerName == "discord-user."+shareduserevents.SignupRoleAddedV1:
		// HandleRoleAdded doesn't return messages (nil)
		return ""

	// ... add all handlers

	default:
		r.logger.Warn("unknown handler in topic resolution",
			attr.String("handler", handlerName),
		)
		// Fallback to metadata (graceful degradation during migration)
		return msg.Metadata.Get("topic")
	}
}
```

### Step 3: Update RegisterHandlers to Use getPublishTopic()

Find the message publishing loop in `RegisterHandlers()` and update it:

**Before:**
```go
for _, m := range messages {
    publishTopic := m.Metadata.Get("topic")

    if publishTopic == "" {
        r.logger.Error("router failed to resolve publish topic - MESSAGE DROPPED", ...)
        continue
    }

    if err := r.publisher.Publish(publishTopic, m); err != nil {
        return nil, fmt.Errorf("failed to publish to %s: %w", publishTopic, err)
    }
}
```

**After:**
```go
for _, m := range messages {
    // Router resolves topic (not metadata)
    publishTopic := r.getPublishTopic(handlerName, m)

    if publishTopic == "" {
        r.logger.Error("router failed to resolve publish topic - MESSAGE DROPPED", ...)
        continue
    }

    if err := r.publisher.Publish(publishTopic, m); err != nil {
        return nil, fmt.Errorf("failed to publish to %s: %w", publishTopic, err)
    }
}
```

### Step 4: Handle Complex Cases

For handlers that return multiple different topics based on logic:

**Option A: Use payload type inspection**
```go
case handlerName == "round."+roundevents.RoundCreationRequestedV1:
    // Unmarshal payload to determine result type
    var payload interface{}
    json.Unmarshal(msg.Payload, &payload)

    switch payload.(type) {
    case *roundevents.RoundValidationPassedPayloadV1:
        return roundevents.RoundValidationPassedV1
    case *roundevents.RoundValidationFailedPayloadV1:
        return roundevents.RoundValidationFailedV1
    default:
        return ""
    }
```

**Option B: Use transient metadata (acceptable during migration)**
```go
case handlerName == "round."+roundevents.RoundCreationRequestedV1:
    // Handler sets result_type metadata
    resultType := msg.Metadata.Get("result_type")
    switch resultType {
    case "validation_passed":
        return roundevents.RoundValidationPassedV1
    case "validation_failed":
        return roundevents.RoundValidationFailedV1
    default:
        // Fallback to topic metadata during migration
        return msg.Metadata.Get("topic")
    }
```

**Option C: Keep metadata fallback temporarily**
```go
case handlerName == "round."+roundevents.RoundCreationRequestedV1:
    // Complex handler - use metadata temporarily
    // TODO: Refactor to use result_type or payload inspection
    return msg.Metadata.Get("topic")
```

## Testing

### Unit Tests

Create `router_test.go` in each router package:

```go
package userrouter_test

import (
    "testing"

    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/stretchr/testify/assert"
)

func TestUserRouter_getPublishTopic(t *testing.T) {
    router := &UserRouter{logger: slog.Default()}

    tests := []struct {
        name          string
        handlerName   string
        msgMetadata   message.Metadata
        expectedTopic string
    }{
        {
            name:          "UserCreated returns SignupAddRole",
            handlerName:   "discord-user.user.created.v1",
            msgMetadata:   message.Metadata{},
            expectedTopic: "discord.user.signup.role.add.v1",
        },
        {
            name:          "RoleAdded returns empty (no publishing)",
            handlerName:   "discord-user.discord.user.signup.role.added.v1",
            msgMetadata:   message.Metadata{},
            expectedTopic: "",
        },
        // Add test for every handler
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            msg := &message.Message{Metadata: tt.msgMetadata}
            topic := router.getPublishTopic(tt.handlerName, msg)
            assert.Equal(t, tt.expectedTopic, topic)
        })
    }
}
```

### Integration Tests

Test end-to-end flows to ensure messages reach the correct destinations:

```go
func TestUserSignupFlow(t *testing.T) {
    // Setup test event bus
    bus := setupTestEventBus(t)

    // Publish SignupRequested
    bus.Publish("user.signup.requested.v1", signupPayload)

    // Assert UserCreated published
    msg := waitForMessage(t, bus, "user.created.v1", 5*time.Second)
    assert.NotNil(t, msg)

    // Assert SignupAddRole published (router should have routed it)
    msg = waitForMessage(t, bus, "discord.user.signup.role.add.v1", 5*time.Second)
    assert.NotNil(t, msg)
}
```

## Deployment

### Pre-Deployment Checklist

1. ✅ All routers have invariant checks (Phase 0)
2. ⬜ Router has `getPublishTopic()` method
3. ⬜ `RegisterHandlers()` uses `getPublishTopic()`
4. ⬜ Unit tests added for router
5. ⬜ Integration tests pass
6. ⬜ No "MESSAGE DROPPED" errors in staging

### Deployment Strategy

**Gradual Rollout:**

1. Deploy to staging
2. Monitor for "MESSAGE DROPPED" errors (should be zero)
3. Run integration tests in staging
4. Deploy to production with 10% traffic
5. Monitor for 24 hours
6. Increase to 50% traffic
7. Monitor for 48 hours
8. Deploy to 100% if no errors

### Monitoring

Watch these metrics:

```bash
# Check for dropped messages
kubectl logs -l app=discord-frolf-bot --tail=1000 | grep "MESSAGE DROPPED"

# Check message publish rate (should be unchanged)
kubectl logs -l app=discord-frolf-bot --tail=1000 | grep "publishing message" | wc -l

# Check for correlation ID preservation
kubectl logs -l app=discord-frolf-bot --tail=1000 | grep "correlation_id"
```

## Router-Specific Implementation Guides

### discord-frolf-bot/app/round (24 Handlers - Complex)

This router is the most complex with many conditional paths.

**Strategy:**
1. Start with simple handlers that always return the same topic
2. Use metadata fallback for complex handlers initially
3. Gradually refactor complex handlers to use `result_type` metadata or payload inspection

**Key Handlers:**
- `HandleRoundCreateRequested`: Multiple validation paths
- `HandleParticipantScoreUpdated`: Success/failure paths
- `HandleScorecardUploaded`: Parse success/failure

### discord-frolf-bot/app/score (4 Handlers - Simple)

**All handlers are straightforward:**
- `HandleScoreUpdateRequest` → publishes to score service
- `HandleScoreUpdateSuccess` → no publishing (Discord notification only)
- `HandleScoreUpdateFailure` → no publishing (Discord notification only)
- `HandleProcessRoundScoresFailed` → no publishing (Discord notification only)

### discord-frolf-bot/app/leaderboard (11 Handlers - Medium)

**Mix of simple and conditional:**
- Tag assign: success/failure paths
- Leaderboard retrieve: success/failure paths
- Tag swap: success/failure paths

### discord-frolf-bot/app/guild (8 Handlers - Simple)

**Mostly notification handlers (no publishing):**
- Config created/updated/deleted: Discord notifications only
- Few handlers actually publish follow-up events

### frolf-bot Routers

**All 5 frolf-bot routers follow similar patterns:**

1. Request handlers → publish to domain stream
2. Response handlers → publish back to Discord stream
3. Failure handlers → usually don't publish (just log)

**Implementation order (easiest to hardest):**
1. guild (simplest - mostly CRUD)
2. score (simple - few handlers)
3. user (medium - signup flow)
4. leaderboard (medium - tag management)
5. round (complex - many states)

## Next Steps After All Routers Complete

### Phase 2: Remove Metadata Dependency

Once all 10 routers implement `getPublishTopic()`:

1. Remove `metadata["topic"]` from helper functions:
```go
// frolf-bot-shared/utils/messages.go

// BEFORE:
newEvent.Metadata.Set("topic", topic)

// AFTER:
// Removed - router now owns topic resolution
// Keep topic_hint for debugging only
newEvent.Metadata.Set("topic_hint", topic)
```

2. Remove metadata fallback from routers:
```go
// BEFORE:
default:
    return msg.Metadata.Get("topic") // Fallback

// AFTER:
default:
    r.logger.Error("unmapped handler - this is a bug")
    return "" // Will trigger MESSAGE DROPPED
```

### Phase 3: Add Comprehensive Tests

1. Unit test every handler→topic mapping
2. Integration test full event flows
3. Add metrics for topic resolution

### Phase 4: Documentation

1. Update README with routing architecture
2. Create runbook for debugging message flow
3. Document how to add new handlers

## FAQs

**Q: Why not just keep using metadata["topic"]?**
A: It violates separation of concerns, is stringly-typed (fragile), and makes testing hard. Router-owned routing is a Watermill best practice.

**Q: What if a handler needs to publish to multiple topics?**
A: Return multiple messages. Each message will be resolved independently by `getPublishTopic()`.

**Q: Can I use metadata for complex routing logic?**
A: Temporarily yes (during migration). Long-term, use payload inspection or refactor handlers to be more explicit.

**Q: What happens if getPublishTopic() returns ""?**
A: The message won't be published, and "MESSAGE DROPPED" will be logged. This is intentional for handlers that don't publish.

**Q: How do I debug routing issues?**
A: Check logs for "MESSAGE DROPPED", "publishing message", and handler names. Use `topic_hint` metadata for debugging.

## References

- ADR: `/docs/adr/0001-router-owned-topic-resolution.md`
- Watermill docs: https://watermill.io/docs/messages-router/
- Migration plan: See original roadmap document

## Support

Questions? Check:
1. This guide first
2. Existing implementation in `discord-frolf-bot/app/user/watermill/router.go`
3. Team Slack channel
4. Create an issue in GitHub

---

Last Updated: 2026-01-07
Status: Phase 1 in progress (1/10 routers complete)
