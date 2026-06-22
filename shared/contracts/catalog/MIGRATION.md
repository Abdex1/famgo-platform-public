# 📋 MIGRATION GUIDE: Event & Schema Evolution

**Status:** Task 2 Phase 2.4 Complete  
**Location:** shared/contracts/catalog/  
**Current Version:** v1 (no migrations needed yet)

---

## MIGRATION CHECKLIST

Use this checklist when evolving an event schema.

### PRE-MIGRATION (Planning)

- [ ] Identify event(s) that need changes
- [ ] Determine if change is breaking or backward-compatible
- [ ] Draft new schema
- [ ] Review with tech lead
- [ ] Estimate consumer impact

### SCHEMA CHANGE (Non-Breaking)

**Process:** Add optional field to v1 (no version bump)

**Checklist:**
- [ ] Add field with `omitempty` JSON tag
- [ ] Field is optional (consumers don't require it)
- [ ] Update SCHEMAS.md with new field
- [ ] Test: Existing consumers still work with new field
- [ ] Test: New consumers see new field
- [ ] Deploy updated publisher
- [ ] Update consumers at their own pace (no deadline)

**Example:**
```go
type PaymentCompleted struct {
    PaymentID  string  `json:"payment_id"`
    Amount     float64 `json:"amount"`
    Tip        float64 `json:"tip,omitempty"` // ✅ New optional field
}
```

### BREAKING CHANGE (Version Bump: v1 → v2)

**Process:** Create v2 event with breaking changes

**Checklist:**

1. **CREATE NEW VERSION**
   - [ ] Create shared/contracts/events/{domain}/v2/ directory
   - [ ] Copy v1 event file to v2/
   - [ ] Modify v2 schema (breaking changes)
   - [ ] Verify v2 schema compiles

2. **UPDATE REGISTRIES**
   - [ ] Add v2 event to shared/contracts/events/catalog/event-types.go
   - [ ] Add v2 topic to shared/contracts/events/topics/topics.go
   - [ ] Update SCHEMAS.md with v2 schema
   - [ ] Update VERSIONS.md with migration guide
   - [ ] Update this MIGRATION.md with specific migration steps

3. **UPDATE PUBLISHER**
   - [ ] Publish v2 events (keep v1 running)
   - [ ] Test: v2 events produced correctly
   - [ ] Monitor: All v2 events in logs

4. **NOTIFY CONSUMERS**
   - [ ] Create GitHub issue with timeline
   - [ ] Post in #engineering Slack
   - [ ] Tag all consuming services
   - [ ] Example message:
     ```
     🚀 New Event Version: payment.events.v2
     
     Breaking Changes:
     - PaymentCompleted.Amount moved to Money object
     
     Migration Timeline:
     - Now: v2 events available
     - Day 30: v1 events deprecated
     - Day 31: v1 removed
     
     Actions:
     - Update your service to consume v2
     - See: MIGRATION.md for details
     - Ask in #engineering if questions
     ```

5. **SUPPORT MIGRATION**
   - [ ] Answer consumer questions
   - [ ] Review consumer PRs that migrate
   - [ ] Help debug migration issues

6. **TRACK MIGRATION**
   - [ ] Create spreadsheet: Services → Migration Status
   - [ ] Update weekly in team standup
   - [ ] Follow up with slow migrants

7. **FINALIZE**
   - [ ] All consumers migrated to v2
   - [ ] Remove v1 event support (if possible)
   - [ ] Archive old code
   - [ ] Document in CHANGELOG

---

## EVENT-SPECIFIC MIGRATIONS

### Payment Events

#### No migrations yet (v1 current)

**When migration needed:**
When we need to add new payment methods or change payment structure.

**Migration steps:**
1. Create shared/contracts/events/payment/v2/
2. Update PaymentCompleted schema
3. Follow general breaking change process above

---

### Ride Events

#### No migrations yet (v1 current)

**When migration needed:**
When ride structure changes significantly.

**Migration steps:**
1. Create shared/contracts/events/ride/v2/
2. Update RideRequested, RideStarted, RideCompleted schemas
3. Coordinate migration with ride-service, dispatch-service, driver-service

---

### Auth Events

#### No migrations yet (v1 current)

**When migration needed:**
When authentication flow changes.

**Migration steps:**
1. Create shared/contracts/events/auth/v2/
2. Update auth event schemas
3. Coordinate migration with all services using auth

---

## CONSUMER MIGRATION TEMPLATE

**Use this template when migrating a consumer service:**

### Service: [SERVICE_NAME]
**Event:** [EVENT_NAME]
**Version:** v1 → v2
**Status:** [Planning / In Progress / Complete]

**Steps:**
1. [ ] Read MIGRATION.md
2. [ ] Read new v2 schema in SCHEMAS.md
3. [ ] Update code to handle v2 schema
4. [ ] Update tests to use v2 schema
5. [ ] Run locally with v2 events
6. [ ] Create PR with changes
7. [ ] Code review ✅
8. [ ] Deploy to staging
9. [ ] Test with v2 events in staging
10. [ ] Deploy to production
11. [ ] Monitor logs for success
12. [ ] Confirm migration complete

**PR Template:**
```markdown
## Migrate to [EVENT_NAME] v2

### Changes
- Updated [Consumer code] to handle v2 schema
- Updated [Tests] to use new schema
- Added [Error handling] for unknown fields

### Testing
- [x] Local testing with v2 events
- [x] Staging tests passed
- [x] No v1 event failures

### Migration Checklist
- [x] Read MIGRATION.md
- [x] Understand breaking changes
- [x] Updated all code paths
- [x] Tests passing
- [x] Ready for production

Fixes: [Reference to GitHub issue]
```

---

## ROLLBACK PROCEDURES

### If Migration Fails

**Immediate Actions:**
1. [ ] Stop deploying v2 consumers
2. [ ] Revert affected services to v1
3. [ ] Investigate failure
4. [ ] Document issue
5. [ ] Plan retry

**Investigation:**
- Check logs for error patterns
- Test v2 event locally
- Verify schema in production
- Check topic configuration

**Retry:**
- Fix identified issue
- Deploy to staging again
- Full testing
- Deploy to production

---

## COMMUNICATION TEMPLATE

### Announcing New Version

```
Subject: 🚀 [Event Name] v2 Available - Migration Required

Hi team,

We're introducing a new version of the [Event Name] event.

### What's Changing
[Describe breaking changes]

### Migration Timeline
- TODAY: v2 events available
- [DATE]: v1 events deprecated
- [DATE+1]: v1 events removed

### Required Actions
1. Review MIGRATION.md
2. Update your event consumer
3. Test with v2 events
4. Deploy your service
5. Confirm v2 working

### Need Help?
- Question about schema? → See SCHEMAS.md
- Question about migration? → See MIGRATION.md
- Stuck? → Ask in #engineering

### Timeline
Questions by [DATE]: We can delay migration if needed
```

### Reminding Consumers

```
Subject: ⏰ Reminder: [Event Name] v2 Migration Deadline [DATE]

Hi [Service] team,

Friendly reminder: [Event Name] v1 is being deprecated on [DATE].

Current Status:
- [Service 1]: ✅ Migrated
- [Service 2]: ⏳ In progress
- [Your Service]: ❌ Not started

Please migrate by [DATE-7 days] to have time for testing.

Questions? → #engineering
```

### Final Deprecation Notice

```
Subject: 🔴 [Event Name] v1 Deprecation - Remove by [DATE+1]

Hi team,

As of today, [Event Name] v1 is officially deprecated.

Remaining tasks:
- [ ] All services consuming v2
- [ ] No more v1 event production
- [ ] Remove v1 code from repository

Please remove v1 support from your service:
- Remove v1 event handlers
- Remove v1 topic subscriptions
- Update documentation

Deadline: [DATE+1]

Questions? → #engineering
```

---

## MIGRATION SUCCESS CRITERIA

### For Non-Breaking Changes
- ✅ New field added to v1
- ✅ Old consumers work (ignore new field)
- ✅ New consumers see new field
- ✅ All tests passing
- ✅ Production monitoring shows no errors

### For Breaking Changes
- ✅ v2 version created
- ✅ v2 schema documented
- ✅ All consumers notified
- ✅ Consumers have 30-day window
- ✅ All consumers migrated to v2
- ✅ v1 removed from production
- ✅ Code archived for history

---

## FREQUENTLY ASKED QUESTIONS

**Q: How do I know if my change is breaking?**
A: If existing consumers will fail to parse the new schema, it's breaking. Add fields → non-breaking. Remove/change fields → breaking.

**Q: Can I skip versions (v1 → v3)?**
A: No. Always increment sequentially: v1 → v2 → v3.

**Q: How long is the migration window?**
A: Standard: 30 days. Can be extended if consumers need more time (decide case-by-case).

**Q: What if a consumer doesn't migrate?**
A: After deprecation date, v1 events stop being produced. Non-migrating consumer will fail. Escalate to that team.

**Q: Can I have v1 and v2 running simultaneously?**
A: Yes! That's the point. Publisher produces v2, but still supports v1 consumers for 30 days. Gives time to migrate.

---

**Migration Guide:** ✅ COMPLETE & READY

