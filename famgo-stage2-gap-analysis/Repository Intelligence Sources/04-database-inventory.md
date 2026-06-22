# Database Inventory — richxcame/ride-hailing

Source: README, `docs/DATABASE_OPERATIONS.md`, `docs/API.md` (domain model references), `docker-compose.yml`, `.env.example`.

## Primary System of Record

| Item | Value |
|---|---|
| Engine | PostgreSQL 15 |
| Image (compose) | `postgis/postgis:15-3.4` — PostGIS extension included |
| Database name | `ridehailing` |
| Connection pooling | `DB_MAX_CONNS=25`, `DB_MIN_CONNS=5` (configurable per environment) |
| SSL | `DB_SSLMODE` (disabled in dev; production checklist calls for `require` or `verify-full`) |
| Read replicas | Mentioned in the README tech-stack line ("connection pooling, read replicas") but no replica topology, connection routing, or replica-specific configuration is documented elsewhere in the reviewed files |
| Resilience | Optional database-specific circuit breaker (`DB_BREAKER_ENABLED`, failure/success thresholds, timeout, interval) |
| Query timeout | `DB_QUERY_TIMEOUT` (default 10s) |

## Migration Tooling

| Item | Value |
|---|---|
| Tool | `golang-migrate` |
| Location | `db/migrations/` |
| Naming | `NNNNNN_description.up.sql` / `NNNNNN_description.down.sql` (sequential, zero-padded) |
| Count (per README) | 18 migrations |
| Management commands | `make migrate-create`, `make migrate-up`, `make migrate-down`, `make migrate-version`, `make migrate-force VERSION=N` |
| Testing | `scripts/test-migrations.sh` (supports `--verbose`, `--clean`, and rollback testing) |

## Domain Entities Referenced Across the API Surface

The migration files themselves were not individually retrievable in this session (see Caveat below), so the following table lists entities **inferred from request/response models named in `docs/API.md`** rather than confirmed directly against schema DDL:

| Entity (model) | Where referenced |
|---|---|
| `User` (`models.User`) | Auth, Admin, Mobile |
| `Ride` (`models.Ride`, `RideRequest`, `RideRatingRequest`) | Rides, Mobile, Admin, Analytics |
| `Wallet`, `WalletTransaction` | Payments |
| `Payment` | Payments |
| `PromoCode` (`internal/promos.PromoCode`) | Promos, Analytics |
| Referral records | Promos, Analytics |
| Ride type / fare configuration | Promos |
| `FavoriteLocation` | Mobile |
| `Notification` | Notifications |
| Fraud `Alert`, `UserRiskProfile` | Fraud |
| ETA prediction records | ML ETA |

No ERD or explicit column-level schema was found in the documentation set; the above should be treated as a domain map, not a verified schema.

## Seed Data Tiers (from `docs/DATABASE_OPERATIONS.md`)

| Tier | Scale | Script |
|---|---|---|
| Light (dev) | 11 users, 9 rides, 5 payments | `make db-seed` / `scripts/seed-database.sql` |
| Medium (test) | 50 users, 200 rides, promo codes, referrals, driver location history | `scripts/seed-medium.sql` |
| Heavy (load) | 1,000 users, 5,000 full-lifecycle rides, ML ETA predictions, wallet transactions | `scripts/seed-heavy.sql` |

`make db-reset` drops, recreates, migrates, and reseeds in one step.

## Backup, Restore & Disaster Recovery

| Capability | Mechanism |
|---|---|
| Full backup | `scripts/backup-database.sh` with `--compress`, `--encrypt`, `--storage s3\|gcs`, `--retention` options |
| Restore | `scripts/restore-database.sh` supporting latest/file/remote/timestamp/new-database/validate-only modes |
| Automated backups | Cron example (daily 2 AM) and a Kubernetes CronJob (`deploy/cronjobs/database-backup-cronjob.yaml`) |
| Backup health checks | `scripts/check-backup-health.sh` with optional Slack/email alerting |
| Point-in-Time Recovery | Separate runbook referenced as `docs/DATABASE_PITR.md`; `scripts/pitr-restore.sh` and `pg_create_restore_point` usage mentioned |
| Disaster recovery runbook | Referenced as `docs/DISASTER_RECOVERY.md` (not independently inspected in this session) |

## Secondary Data Stores (not the system of record, but persistence-adjacent)

| Store | Role |
|---|---|
| Redis 7 | Cache layer, geospatial driver-location index, Pub/Sub transport (see Realtime Inventory) |
| Kong's own PostgreSQL (`kong-database`) | Stores Kong gateway configuration only — operationally and logically separate from the application database |
| Sentry's self-hosted data stack (Postgres + Redis + ClickHouse + Kafka + Zookeeper) | Stores error-tracking data only, when running the self-hosted Sentry option in `docker-compose.yml`; entirely separate from application data |

## Monitoring Queries Documented for Operators

The operations guide includes example SQL for: active connection count, long-running query detection (>5 min), and database size — these are operational health checks rather than schema documentation.

## Caveat

GitHub's robots policy blocked direct retrieval of the `db/migrations/` directory listing during this session, so individual migration filenames, table definitions, indexes, and constraints could not be enumerated. This inventory is built from operational documentation and API-layer domain model references rather than from the SQL DDL itself.
