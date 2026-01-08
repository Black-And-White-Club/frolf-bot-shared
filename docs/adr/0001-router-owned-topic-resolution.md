# ADR 001: Router-Owned Topic Resolution

## Status
**In Progress** - Partial implementation complete

## Context

Previously, handlers set `metadata["topic"]` to indicate where messages should be published. This approach had several issues:

1. **Separation of Concerns Violation**: Handlers (business logic layer) were making infrastructure decisions about message routing
2. **Stringly-Typed**: Topic resolution relied on string metadata, which is fragile and error-prone
3. **Implicit Routing**: The routing logic was hidden in handlers and helpers, making it hard to understand message flow
4. **Testing Difficulty**: Hard to test routing logic in isolation from handler logic

## Decision

Routers own topic resolution via `getPublishTopic()` method. Handlers remain infrastructure-agnostic, returning messages without routing knowledge.

### Architecture Pattern

```go
// Router owns routing logic
func (r *UserRouter) getPublishTopic(handlerName string, msg *message.Message) string {
    switch {
    case handlerName == "discord-user."+userevents.UserCreatedV1:
        return shareduserevents.SignupAddRoleV1
    case handlerName == "discord-user."+shareduserevents.RoleUpdateButtonPressV1:
        return userevents.UserRoleUpdateRequestedV1
    default:
        return msg.Metadata.Get("topic") // Fallback during migration
    }
}

// Handler returns messages without routing knowledge
func (h *UserHandlers) HandleUserCreated(msg *message.Message) ([]*message.Message, error) {
    // Business logic only
    resultMsg, _ := h.Helper.CreateResultMessage(msg, payload, targetTopic)
    return []*message.Message{resultMsg}, nil
}
```

### Migration Strategy

**Phase 0: Safety Nets** ✅ COMPLETED
- Added invariant checks to all 10 routers
- Log "MESSAGE DROPPED" errors if topic resolution fails
- Continue publishing but skip failed resolutions

**Phase 1: Router Ownership** 🚧 IN PROGRESS
- Implement `getPublishTopic()` in each router (1/10 complete: discord-frolf-bot user router)
- Update helper to keep `metadata["topic"]` temporarily for backward compatibility
- Deploy with gradual rollout
- Monitor for zero dropped messages

**Phase 2: Remove Metadata** ⏳ PENDING
- After all routers implement `getPublishTopic()`
- Remove `metadata["topic"]` from helper functions
- Remove metadata fallbacks from `getPublishTopic()` methods

## Consequences

### Positive
- **Type-Safe**: Routing logic uses constants, not strings
- **Explicit**: Easy to grep and understand message flow
- **Testable**: Router logic can be unit tested independently
- **Maintainable**: Adding new handlers requires updating router (compile-time check)

### Negative
- **Requires Router Updates**: Adding new handlers means updating router code
- **Migration Effort**: 10 routers need updates

### Mitigation
- Unit tests enforce correctness of handler→topic mappings
- Documentation guides adding new handlers
- Gradual migration with backward compatibility

## Implementation Status

### Completed ✅
1. **Phase 0: All 10 routers** - Invariant checks added
   - discord-frolf-bot: user, round, score, leaderboard, guild
   - frolf-bot: user, round, score, leaderboard, guild

2. **Phase 1: 1/10 routers** - `getPublishTopic()` implemented
   - discord-frolf-bot/app/user/watermill/router.go ✅

### In Progress 🚧
1. **Phase 1: Remaining 9 routers** - Need `getPublishTopic()` implementation

### Pending ⏳
1. **Phase 2**: Remove metadata["topic"] from helpers
2. **Phase 3**: Add unit tests for all routers
3. **Phase 4**: Remove metadata fallbacks

## Handler→Topic Mapping Reference

### discord-frolf-bot/app/user/watermill/router.go ✅ COMPLETED
| Input Topic (Handler) | Output Topic | Notes |
|----------------------|--------------|-------|
| user.created.v1 | discord.user.signup.role.add.v1 | Always |
| discord.user.signup.role.add.v1 | discord.user.signup.role.added.v1 OR discord.user.signup.role.addition.failed.v1 | Success/Failure |
| discord.user.role.update.button.press.v1 | user.role.update.requested.v1 | Always |
| All others | (none - return nil) | No publishing |

### discord-frolf-bot/app/round/watermill/router.go ⏳ TODO
Requires analysis of 24 handlers - see handlers directory for mapping

### discord-frolf-bot/app/score/watermill/router.go ⏳ TODO
Requires analysis of 4 handlers

### discord-frolf-bot/app/leaderboard/watermill/router.go ⏳ TODO
Requires analysis of 11 handlers

### discord-frolf-bot/app/guild/watermill/router.go ⏳ TODO
Requires analysis of 8 handlers

### frolf-bot routers ⏳ TODO
Similar pattern for all 5 frolf-bot routers

## References

- Watermill documentation: [Handlers return messages](https://watermill.io/docs/messages-router/#handlers)
- ThreeDotsLabs EDA patterns
- Migration roadmap: `/docs/migration/event-driven-architecture.md`

## Authors

- Implementation: Claude Sonnet 4.5
- Reviewed by: Jace (Project Owner)

## Date

Created: 2026-01-07
Last Updated: 2026-01-07
