# Event-Driven Architecture Migration Status

## 🎯 Quick Summary

**Goal:** Migrate from metadata-driven routing to router-owned topic resolution across all 10 routers.

**Status:** Phase 0 complete ✅ | Phase 1: 50% complete (5/10 routers) 🚧

**Risk:** LOW - Backward compatible, incremental deployment

**Impact:** Zero user-facing changes, improves code quality and testability

---

## 📊 Current State

### ✅ COMPLETED

#### Phase 0: Safety Nets (ALL 10 ROUTERS)
All routers now have invariant checks that prevent silent message loss.

**What was done:**
- Added "MESSAGE DROPPED" logging if topic resolution fails
- Added correlation_id and handler to all publish logs
- Changed warnings to errors for missing topics

**Files modified:**
```
discord-frolf-bot/app/user/watermill/router.go
discord-frolf-bot/app/round/watermill/router.go
discord-frolf-bot/app/score/watermill/router.go
discord-frolf-bot/app/leaderboard/watermill/router.go
discord-frolf-bot/app/guild/watermill/router.go
frolf-bot/app/modules/user/infrastructure/router/router.go
frolf-bot/app/modules/round/infrastructure/router/router.go
frolf-bot/app/modules/score/infrastructure/router/router.go
frolf-bot/app/modules/leaderboard/infrastructure/router/router.go
frolf-bot/app/modules/guild/infrastructure/router/router.go
```

**Code pattern added to all routers:**
```go
for _, m := range messages {
    publishTopic := m.Metadata.Get("topic")

    // INVARIANT: Topic must be resolvable
    if publishTopic == "" {
        r.logger.Error("router failed to resolve publish topic - MESSAGE DROPPED",
            attr.String("handler", handlerName),
            attr.String("msg_uuid", m.UUID),
            attr.String("correlation_id", m.Metadata.Get("correlation_id")),
        )
        continue // Skip but don't fail entire batch
    }

    r.logger.InfoContext(ctx, "publishing message",
        attr.String("topic", publishTopic),
        attr.String("handler", handlerName),
        attr.String("correlation_id", m.Metadata.Get("correlation_id")),
    )

    if err := r.publisher.Publish(publishTopic, m); err != nil {
        return nil, fmt.Errorf("failed to publish to %s: %w", publishTopic, err)
    }
}
```

#### Phase 1: Router-Owned Routing (5/10 ROUTERS)

**Completed:**
- `discord-frolf-bot/app/user/watermill/router.go` ✅
  - Added `getPublishTopic()` method (lines 78-129)
  - Updated `RegisterHandlers()` to use router resolution (line 171)
  - Handles 9 event types with explicit mapping

- `discord-frolf-bot/app/guild/watermill/router.go` ✅
  - Added `getPublishTopic()` method (lines 73-119)
  - Updated `RegisterHandlers()` to use router resolution (line 168)
  - Handles 8 event types with explicit mapping
  - Simple implementation: 7 handlers return nil, only 1 publishes

- `discord-frolf-bot/app/score/watermill/router.go` ✅
  - Added `getPublishTopic()` method (lines 73-104)
  - Updated `RegisterHandlers()` to use router resolution (line 137)
  - Handles 4 event types with explicit mapping
  - 1 conditional case (uses metadata fallback temporarily)

- `discord-frolf-bot/app/leaderboard/watermill/router.go` ✅
  - Added `getPublishTopic()` method (lines 76-147)
  - Updated `RegisterHandlers()` to use router resolution (line 202)
  - Handles 14 event types with explicit mapping
  - 1 conditional case (uses metadata fallback temporarily)
  - 2 handlers return nil (error cases)

- `discord-frolf-bot/app/round/watermill/router.go` ✅
  - Added `getPublishTopic()` method (lines 77-193)
  - Updated `RegisterHandlers()` to use router resolution (line 273)
  - Handles 23 event types with explicit mapping
  - Complex implementation: Covers creation, update, participation, scoring, scorecard import, deletion, lifecycle, tag handling, and reminders
  - 10 handlers publish events, 13 handlers return nil (Discord API operations)
  - Includes score override bridging (CorrectScore service)

**Helper file migration annotations added:**
- `frolf-bot-shared/utils/messages.go` updated with TODO comments
- Kept `metadata["topic"]` for backward compatibility
- Added `topic_hint` metadata for debugging

---

## 🚧 REMAINING WORK

### Phase 1: Router-Owned Routing (6/10 REMAINING)

Each router needs:
1. `getPublishTopic(handlerName string, msg *message.Message) string` method
2. Update `RegisterHandlers()` to call `getPublishTopic()` instead of `m.Metadata.Get("topic")`
3. Handler→topic mapping for all event types

#### Remaining Routers (Priority Order):

**1. discord-frolf-bot/app/round/watermill/router.go** (HARD - 24 handlers)
- Complex validation flows
- Multiple state transitions
- Scorecard import flows
- Estimated effort: 2-3 hours

**2-6. frolf-bot routers** (MEDIUM - 5 routers)
- Similar patterns to discord-frolf-bot routers
- Mostly domain logic handlers
- Estimated effort: 3-4 hours total

---

## 📝 IMPLEMENTATION TEMPLATE

### For Each Router:

**Step 1:** Analyze handlers
```bash
cd <router-directory>/handlers
grep -r "CreateResultMessage" . | grep -o '"[^"]*\.v[0-9]"' | sort -u
```

**Step 2:** Create handler→topic mapping table

**Step 3:** Add `getPublishTopic()` method BEFORE `RegisterHandlers()`:
```go
// getPublishTopic resolves the topic to publish for a given handler's returned message.
func (r *RouterName) getPublishTopic(handlerName string, msg *message.Message) string {
    switch {
    case handlerName == "prefix."+events.TopicV1:
        return outputevents.OutputTopicV1
    // ... map all handlers
    default:
        r.logger.Warn("unknown handler in topic resolution",
            attr.String("handler", handlerName),
        )
        return msg.Metadata.Get("topic") // Fallback during migration
    }
}
```

**Step 4:** Update the publish loop in `RegisterHandlers()`:
```go
// Change this line:
publishTopic := m.Metadata.Get("topic")

// To this:
publishTopic := r.getPublishTopic(handlerName, m)
```

**Step 5:** Test in staging - ensure zero "MESSAGE DROPPED" errors

---

## 🧪 TESTING CHECKLIST

### Per Router:
- [ ] Unit test added for `getPublishTopic()` mapping
- [ ] All handlers have test cases
- [ ] Integration test for main flows
- [ ] Zero "MESSAGE DROPPED" in staging logs

### System-Wide:
- [ ] Message publish rate unchanged
- [ ] All correlation IDs preserved
- [ ] End-to-end flows work (signup, round creation, scoring)

---

## 🚀 DEPLOYMENT STRATEGY

### Phase 1 (Current):
1. Deploy routers with `getPublishTopic()` one at a time
2. Keep `metadata["topic"]` in helpers (backward compatible)
3. Monitor staging for 24 hours after each router
4. Deploy to production with canary (10% → 50% → 100%)

### Phase 2 (After All Routers Complete):
1. Remove `metadata["topic"]` from helpers:
   - `frolf-bot-shared/utils/messages.go:64` - Remove line
   - `frolf-bot-shared/utils/messages.go:120` - Remove line
2. Remove metadata fallbacks from all `getPublishTopic()` methods
3. Deploy with same canary strategy

### Phase 3 (Polish):
1. Add comprehensive router unit tests
2. Add metrics for topic resolution
3. Update documentation

---

## 📚 DOCUMENTATION

**Created:**
- ✅ `/docs/adr/0001-router-owned-topic-resolution.md` - Architecture decision record
- ✅ `/docs/development/router-implementation-guide.md` - Detailed implementation guide
- ✅ `/docs/MIGRATION_STATUS.md` - This file

**Reference Implementation:**
- ✅ `discord-frolf-bot/app/user/watermill/router.go` - Complete example

---

## ⚠️ CRITICAL NOTES

### DON'T:
- ❌ Remove `metadata["topic"]` from helpers yet (breaks other routers)
- ❌ Change handler signatures (no handler code changes needed)
- ❌ Deploy without testing in staging first
- ❌ Batch multiple router changes in one deploy

### DO:
- ✅ One router at a time
- ✅ Monitor logs after each deployment
- ✅ Keep metadata fallback in `getPublishTopic()` during migration
- ✅ Test correlation_id preservation
- ✅ Update unit tests for each router

---

## 📋 QUICK START FOR NEW SESSION

```bash
# 1. Pick next router (discord-frolf-bot/app/round)
cd discord-frolf-bot/app/round/watermill

# 2. Read existing implementation example
cat ../../user/watermill/router.go | grep -A 50 "getPublishTopic"

# 3. Analyze handlers
grep -r "CreateResultMessage" handlers/ | grep -o '"[^"]*\.v[0-9]"'

# 4. Implement getPublishTopic() following the pattern

# 5. Test in staging
kubectl logs -l app=discord-frolf-bot | grep "MESSAGE DROPPED"
```

---

## 🎓 UNDERSTANDING THE PATTERN

### Before (Metadata-Driven):
```go
// Handler (business logic layer) decides routing
func HandleUserCreated(msg) {
    resultMsg := CreateResultMessage(payload, "discord.user.signup.role.add.v1")
    // Helper sets: metadata["topic"] = "discord.user.signup.role.add.v1"
    return []*message.Message{resultMsg}
}

// Router blindly publishes to metadata["topic"]
publishTopic := m.Metadata.Get("topic")
publisher.Publish(publishTopic, m)
```

### After (Router-Owned):
```go
// Handler (business logic layer) is infrastructure-agnostic
func HandleUserCreated(msg) {
    resultMsg := CreateResultMessage(payload, shareduserevents.SignupAddRoleV1)
    // Helper sets: topic_hint (for debugging only)
    return []*message.Message{resultMsg}
}

// Router (infrastructure layer) owns routing decisions
func (r *Router) getPublishTopic(handlerName, msg) string {
    switch handlerName {
    case "discord-user.user.created.v1":
        return shareduserevents.SignupAddRoleV1  // Type-safe constant
    }
}
publishTopic := r.getPublishTopic(handlerName, m)
publisher.Publish(publishTopic, m)
```

### Why This Matters:
- **Separation of Concerns**: Business logic doesn't know about infrastructure
- **Type Safety**: Uses constants instead of strings
- **Testability**: Can unit test routing logic independently
- **Explicitness**: grep-able, clear message flow

---

## 📊 PROGRESS TRACKER

```
Phase 0: Safety Nets
[✅✅✅✅✅✅✅✅✅✅] 10/10 routers (100%)

Phase 1: Router-Owned Routing
[✅✅✅✅⬜⬜⬜⬜⬜⬜] 4/10 routers (40%)
 ✅ discord-frolf-bot/user
 ✅ discord-frolf-bot/guild
 ✅ discord-frolf-bot/score
 ✅ discord-frolf-bot/leaderboard
 ⬜ discord-frolf-bot/round (next - hardest)
 ⬜ frolf-bot/user
 ⬜ frolf-bot/guild
 ⬜ frolf-bot/score
 ⬜ frolf-bot/leaderboard
 ⬜ frolf-bot/round

Phase 2: Remove Metadata
[⬜] 0/1 (pending all routers)

Phase 3: Polish
[⬜] Tests, metrics, docs
```

---

## 🔧 TROUBLESHOOTING

### "MESSAGE DROPPED" Errors After Deployment
**Cause:** `getPublishTopic()` returned empty string
**Fix:** Check handler mapping - likely missing case in switch statement

### Messages Not Arriving at Destination
**Cause:** `getPublishTopic()` returning wrong topic
**Fix:** Check handler→topic mapping, compare with metadata["topic"] value in logs

### Tests Failing After Migration
**Cause:** Test expectations based on old metadata-driven pattern
**Fix:** Update tests to match new router-owned pattern

---

## ✅ ACCEPTANCE CRITERIA

### Per Router:
- ✅ Zero "MESSAGE DROPPED" errors in staging (24hr)
- ✅ Zero "MESSAGE DROPPED" errors in production (48hr)
- ✅ Message publish rate unchanged
- ✅ Integration tests pass
- ✅ Correlation IDs preserved through flows

### System-Wide:
- ✅ All 10 routers implement `getPublishTopic()`
- ✅ All routers use router resolution (not metadata)
- ✅ No regressions in message delivery
- ✅ Documentation complete

---

## 📅 TIMELINE ESTIMATE

- **Phase 1 Completion:** 8-10 hours development + testing
  - Guild router: 30 min ✅
  - Score router: 30 min ✅
  - Leaderboard router: 1 hour ✅
  - Round router (discord): 2-3 hours
  - All frolf-bot routers: 3-4 hours
  - Testing: 1-2 hours

- **Phase 2 Completion:** 1-2 hours
  - Remove metadata from helpers
  - Remove fallbacks from routers
  - Regression testing

- **Phase 3 Completion:** 2-3 hours
  - Comprehensive unit tests
  - Metrics implementation
  - Documentation updates

**Total:** 11-15 hours to complete entire migration

---

## 🎯 SUCCESS METRICS

- Zero message loss (no "MESSAGE DROPPED" errors)
- Zero routing errors in production
- 100% test coverage for router topic resolution
- Documentation complete and reviewed
- Team trained on new patterns

---

Last Updated: 2026-01-07
Next Action: Implement `discord-frolf-bot/app/round/watermill/router.go` (HARD - 24 handlers, complex validation flows)
Owner: Jace
Status: Phase 1 In Progress (40% complete)
