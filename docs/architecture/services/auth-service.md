# Auth Service

## Responsibilities

- User registration, login, OTP verification, and password reset.
- JWT access/refresh token issuance and rotation.
- Role-based access control for rider, driver, and admin personas.
- Trusted service-account tokens for inter-service calls (e.g., notification hooks).

## Database Tables (proposed)

| Table | Purpose |
|---|---|
| `users` | Credentials, role, status, email verification state |
| `otp_verification` | One-time codes for registration and password reset |
| `refresh_tokens` | Rotating refresh token store (hashed) |
| `jwt_keys` | Optional signing key rotation metadata |

## Publishes

| Event | Topic | Consumers |
|---|---|---|
| Login succeeded | `auth.login.succeeded.v1` | audit, analytics |
| Login failed | `auth.login.failed.v1` | fraud-service, audit |
| User registered | `user.registered.v1` | user-service, driver-service (onboarding) |

## Consumes

| Event | Source | Action |
|---|---|---|
| *(none — upstream identity provider)* | — | Auth is the root of most identity flows |

## External Dependencies

- **PostgreSQL** — primary identity store
- **Redis** (optional) — session cache, OTP rate limiting

## Notes

Downstream services verify JWTs via `famgo/shared/pkg/auth` or direct auth-service introspection. No ride or payment domain logic belongs here.
