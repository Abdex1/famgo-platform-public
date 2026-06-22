# Notification Service

## Responsibilities

- Multi-channel delivery: push (FCM), SMS (Twilio), and email (SMTP).
- Ride lifecycle notification templates (requested, accepted, started, completed, cancelled).
- Admin bulk broadcast and user preference management.
- Consumes domain events and translates them into channel-specific payloads.

## Database Tables (proposed)

| Table | Purpose |
|---|---|
| `notifications` | Outbound message queue with delivery status |
| `notification_preferences` | Per-user channel opt-in/out |
| `device_tokens` | Push registration tokens |
| `notification_templates` | Localized template bodies by event type |

## Publishes

| Event | Topic | Consumers |
|---|---|---|
| Notification sent | `notification.sent.v1` | analytics, audit |
| Notification failed | `notification.failed.v1` | retry worker, DLQ |

## Consumes

| Event | Source | Action |
|---|---|---|
| `ride.created.v1` | ride-service | Notify rider that matching started |
| `dispatch.driver.assigned.v1` | dispatch-service | Push driver assignment to rider |
| `ride.started.v1` | ride-service | Trip-in-progress alerts |
| `ride.completed.v1` | ride-service | Receipt and rating prompt |
| `ride.cancelled.v1` | ride-service | Cancellation notice |
| `payment.failed.v1` | payment-service | Payment retry / support alert |

## External Dependencies

- **Firebase Cloud Messaging** — mobile push
- **Twilio** — SMS
- **SMTP** — email
- **PostgreSQL** — delivery log and preferences

## Notes

Richxcame used synchronous HTTP hooks (`POST /notifications/ride/{event}`). FamGo target standardizes on Kafka consumers with idempotent template rendering and retry/DLQ policies from `packages/kafka-sdk`.
