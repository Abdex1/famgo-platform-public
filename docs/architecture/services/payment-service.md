# Payment Service

## Responsibilities

- Wallet balances, ride fare capture, refunds, and settlement with external PSP (e.g., Stripe).
- Idempotent payment authorization and capture tied to ride completion.
- Webhook intake for provider events (`payment_intent.succeeded`, etc.).
- Publishes payment outcomes for wallet, fraud, and ride domains.

## Database Tables (proposed)

| Table | Purpose |
|---|---|
| `wallets` | Per-user balance and currency |
| `wallet_transactions` | Ledger entries (debit/credit) |
| `payments` | Ride-linked payment attempts and status |
| `payment_methods` | Tokenized cards / provider references |
| `refunds` | Refund records linked to payments |

## Publishes

| Event | Topic | Consumers |
|---|---|---|
| Payment authorized | `payment.authorized.v1` | ride-service, wallet-service |
| Payment completed | `payment.completed.v1` | ride-service, wallet-service, fraud-service |
| Payment failed | `payment.failed.v1` | ride-service, notification-service |
| Wallet debited | `wallet.debited.v1` | analytics |
| Wallet credited | `wallet.credited.v1` | analytics |

## Consumes

| Event | Source | Action |
|---|---|---|
| `ride.completed.v1` | ride-service | Capture fare against rider wallet or card |
| `ride.cancelled.v1` | ride-service | Release holds or issue partial refunds |
| `subscription.renewal_due.v1` | subscription-service | Recurring billing (future) |

## External Dependencies

- **Stripe** (or `packages/payment-sdk`) — card processing and webhooks
- **PostgreSQL** — ledger and payment state
- **Redis** — idempotency keys for webhook replay protection

## Notes

Richxcame reference used synchronous HTTP notification hooks; FamGo target moves settlement to Kafka with outbox-style durability.
