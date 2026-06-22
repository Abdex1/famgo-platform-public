# 🗄️ DATABASE CATALOG: FamGo Data Architecture

**Last Updated:** [Date]  
**Status:** Repository Consistency Audit - TASK 1  
**Total Databases:** 14 (PostgreSQL + Redis + Elasticsearch)

---

## DATABASE CATALOG

### CORE DATABASES

#### 1. auth_db (PostgreSQL)
- **Owner:** auth-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs (7 days retention)
- **Tables:**
  - `users` (id, email, phone, password_hash, created_at)
  - `sessions` (id, user_id, token_hash, expires_at, device_id)
  - `roles` (id, name, permissions_json)
  - `user_roles` (user_id, role_id)
  - `audit_log` (id, user_id, action, ip_address, timestamp)
- **Schemas:** Fully normalized, 3NF
- **Access:** SSH tunneled, service accounts only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous (1 follower)

#### 2. user_db (PostgreSQL)
- **Owner:** user-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `users` (id, email, phone, first_name, last_name, profile_picture_url)
  - `user_preferences` (user_id, language, notifications_enabled, theme)
  - `driver_profiles` (user_id, license_number, vehicle_id, rating, total_rides)
  - `passenger_profiles` (user_id, favorite_locations, payment_methods, rating)
  - `emergency_contacts` (id, user_id, name, phone, relationship)
- **Schemas:** Fully normalized
- **Access:** Service accounts only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 3. gps_db (PostgreSQL with PostGIS)
- **Owner:** gps-service
- **Type:** PostgreSQL 14+ with PostGIS extension
- **Replicas:** 1 (read-only for analytics)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `trip_routes` (id, ride_id, polyline, distance_meters, timestamp)
  - `geofences` (id, name, zone_polygon, city, created_at)
  - `driver_locations_history` (driver_id, latitude, longitude, timestamp)
- **Indexes:** Spatial indexes on geometries
- **Access:** Service accounts + analytics (read-only replica)
- **Backups:** Daily full backup
- **Replication:** Asynchronous (read-only replica for analytics)

#### 4. ride_db (PostgreSQL)
- **Owner:** ride-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `rides` (id, user_id, driver_id, status, pickup_location, dropoff_location, created_at)
  - `ride_state_history` (ride_id, previous_state, new_state, transitioned_at)
  - `ride_ratings` (id, ride_id, rating_by_user, rating_value, comment, created_at)
  - `ride_cancellations` (id, ride_id, cancelled_by, reason, cancelled_at)
- **Schemas:** Fully normalized
- **Access:** Service accounts only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 5. dispatch_db (PostgreSQL)
- **Owner:** dispatch-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `dispatch_attempts` (id, ride_id, offered_driver_id, status, accepted_at, rejected_at)
  - `driver_assignments` (id, ride_id, driver_id, assigned_at, eta_seconds)
  - `dispatch_metrics` (id, ride_id, total_attempts, average_eta, success, timestamp)
- **Schemas:** Fully normalized
- **Access:** Service accounts only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 6. pricing_db (PostgreSQL)
- **Owner:** pricing-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `pricing_rules` (id, zone_id, base_fare, per_km, per_minute, effective_date)
  - `surge_multipliers` (id, zone_id, time_slot, multiplier, created_at)
  - `discounts` (id, code, discount_type, amount_or_percentage, start_date, end_date)
  - `fare_history` (ride_id, calculated_fare_components_json, timestamp)
- **Schemas:** Fully normalized
- **Access:** Service accounts only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 7. payment_db (PostgreSQL)
- **Owner:** payment-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs (7 years retention for financial)
- **Tables:**
  - `payment_intents` (id, ride_id, amount, currency, status, created_at)
  - `transactions` (id, payment_intent_id, gateway_transaction_id, status, timestamp)
  - `refunds` (id, transaction_id, refund_reason, refunded_at)
  - `webhook_logs` (id, gateway_name, webhook_id, payload_hash, processed_at)
- **Schemas:** Fully normalized
- **Access:** Payment team + audit
- **Backups:** Daily full + hourly incremental (7-year retention)
- **Replication:** Synchronous (CRITICAL for financial)
- **Encryption:** At-rest encryption enabled

#### 8. wallet_db (PostgreSQL)
- **Owner:** wallet-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs (7 years for financial)
- **Tables:**
  - `wallet_accounts` (user_id, balance, currency, last_updated)
  - `ledger` (id, user_id, debit_or_credit, amount, source_transaction_id, timestamp)
  - `holds` (id, user_id, ride_id, hold_amount, hold_until, created_at)
  - `reconciliation_logs` (id, batch_timestamp, total_debits, total_credits, status)
- **Schemas:** Immutable ledger pattern
- **Access:** Wallet team + audit
- **Backups:** Daily full + hourly incremental (7-year retention)
- **Replication:** Synchronous (CRITICAL)

#### 9. pooling_db (PostgreSQL)
- **Owner:** pooling-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `pool_matches` (id, ride_id_1, ride_id_2, overlap_percentage, detour_minutes, matched_at)
  - `pool_rides` (id, ride_id, pool_id, seat_position, seat_count)
  - `pool_history` (pool_id, historical_analysis_json, completed_at)
- **Schemas:** Fully normalized
- **Access:** Service accounts only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 10. driver_db (PostgreSQL)
- **Owner:** driver-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `driver_applications` (id, user_id, status, submitted_at, reviewed_at, reviewed_by)
  - `driver_documents` (id, driver_id, document_type, file_url, verified_at, expires_at)
  - `driver_suspension_history` (id, driver_id, reason, suspension_duration, suspended_at)
  - `vehicle_info` (id, driver_id, make, model, license_plate, registration_expiry)
- **Schemas:** Fully normalized
- **Access:** Driver team + admin
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 11. fraud_db (PostgreSQL)
- **Owner:** fraud-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (read-only for ML analysis)
- **Backup:** Daily + transaction logs (2-year retention)
- **Tables:**
  - `fraud_flags` (id, user_id_or_driver_id, fraud_type, risk_score, flagged_at)
  - `fraud_rules_log` (id, fraud_id, rule_name, rule_version, triggered_at)
  - `fraud_resolutions` (fraud_id, resolution, resolved_by, resolved_at)
- **Schemas:** Fully normalized
- **Access:** Fraud team + support
- **Backups:** Daily full + hourly incremental (2-year retention)
- **Replication:** Synchronous

#### 12. support_db (PostgreSQL)
- **Owner:** support-service
- **Type:** PostgreSQL 14+
- **Replicas:** 1 (hot standby)
- **Backup:** Daily + transaction logs
- **Tables:**
  - `support_tickets` (id, user_id, issue_category, description, status, created_at)
  - `ticket_assignments` (ticket_id, agent_id, assigned_at)
  - `ticket_interactions` (id, ticket_id, actor_type, message, timestamp)
  - `disputes` (id, ride_id, claimant_id, claim_reason, resolution, resolved_at)
- **Schemas:** Fully normalized
- **Access:** Support team only
- **Backups:** Daily full + hourly incremental
- **Replication:** Synchronous

#### 13. analytics_db (PostgreSQL)
- **Owner:** analytics-service
- **Type:** PostgreSQL 14+ (time-series optimized)
- **Replicas:** Read-only replicas (for queries)
- **Backup:** Daily full backup
- **Tables:**
  - `daily_metrics` (date, zone, rides_count, avg_fare, avg_rating, driver_count)
  - `hourly_trends` (timestamp, zone, active_drivers, waiting_passengers, avg_eta)
  - `user_metrics` (user_id, rides_count, total_spent, avg_rating, joined_at)
  - `driver_metrics` (driver_id, total_rides, total_earnings, avg_rating, online_hours)
- **Schemas:** Denormalized for fast queries
- **Access:** Read-only (analytics team + dashboards)
- **Backups:** Daily full backup
- **Replication:** Asynchronous (read-only replicas)

---

### CACHING LAYER (Redis)

#### 14. Redis Cache (Cluster)
- **Owner:** Infrastructure
- **Type:** Redis 7.0+ (Cluster mode)
- **Replicas:** 3 nodes (HA)
- **Backup:** Daily RDB snapshots + AOF
- **Data Stored:**
  - `driver:location:{driver_id}` → Current GPS location (TTL 5 min)
  - `driver:availability:{driver_id}` → Online/offline status (TTL 1 hour)
  - `session:{session_id}` → User session data (TTL 7 days)
  - `rate_limit:{user_id}:{endpoint}` → Request counts (TTL 1 min)
  - `ride:active:{ride_id}` → Current ride state (TTL 24 hours)
- **Access:** Service accounts only
- **Eviction Policy:** allkeys-lru (least recently used)
- **Persistence:** RDB + AOF enabled

---

## ACCESS CONTROL

**Database Access Pattern:**
```
Service A ----[TLS/SSL]----> PostgreSQL
         ----[TLS/SSL]----> Redis
```

**Requirements:**
- All connections encrypted (TLS 1.3+)
- Service accounts (no shared credentials)
- Network isolation (private VPC subnets)
- Audit logging of all connections
- No direct password access (use identity provider)

---

## BACKUP STRATEGY

**All PostgreSQL Databases:**
- Daily full backup (midnight UTC)
- Hourly incremental backup (transaction logs)
- WAL archiving to S3 (7+ days)
- Retention: 30 days rolling + 7-year retention for financial data
- Recovery Time Objective (RTO): <5 minutes
- Recovery Point Objective (RPO): <1 minute

**Redis Cluster:**
- Daily RDB snapshots
- Continuous AOF logging
- Retention: 7 days
- RTO: <2 minutes
- RPO: <5 minutes

---

## MONITORING & ALERTS

**Alert on:**
- Connection failures
- Query latency >1s
- Replication lag >5s
- Disk usage >80%
- Backup failures
- Encryption errors

---

## SCHEMA MIGRATIONS

**Strategy:**
- Migrate forward only (no rollbacks after production)
- Test migrations on staging first
- Blue-green deployment for DB changes
- Backward compatibility for 2 versions
- Automated schema versioning

---

## COMPLIANCE & RETENTION

**Data Retention:**
- Ride/payment data: 7 years (financial compliance)
- User data: Deletable after 2 years (GDPR)
- Audit logs: 2 years minimum
- GPS data: 30 days (privacy)
- Personal documents: Delete after 30 days (driver onboarding)

---

**Audit Status:** ✅ DATABASE_CATALOG.md COMPLETE  
**Next:** API_CATALOG.md creation

