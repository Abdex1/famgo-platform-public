# FamGo Platform - Migration Mapping Guide

## Overview

This document maps existing code from `C:\dev\FamGo\` (RidePool-STRPS) to the new enterprise architecture at `C:\dev\FamGo-platform\`.

## Existing Project Structure (Source)

```
C:\dev\FamGo\
├── backend/
│   ├── app/
│   │   ├── main.py               (FastAPI entry)
│   │   ├── config.py             (Configuration)
│   │   ├── models/               (Pydantic models)
│   │   ├── routes/               (API endpoints)
│   │   ├── services/             (Business logic)
│   │   ├── utils/                (Helpers)
│   │   └── websocket/            (Socket.IO handlers)
│   ├── seed_data/
│   ├── requirements.txt
│   └── .env
│
└── src/                          (React frontend)
    ├── components/
    │   ├── common/               (Shared components)
    │   ├── user/                 (Rider components)
    │   ├── driver/               (Driver components)
    │   └── admin/                (Admin components)
    ├── pages/
    ├── context/
    ├── services/
    ├── types/
    └── utils/
```

## Migration Mapping by Component

### 1. AUTHENTICATION LAYER

| Source | Destination | Notes |
|--------|-------------|-------|
| `backend/app/routes/auth.py` | `services/auth-service/internal/interfaces/rest/handlers/auth_handler.go` | Convert FastAPI to Go; Keep business logic |
| `backend/app/models/user.py` | `services/auth-service/internal/domain/entities/user.go` | Map Pydantic models to Go structs |
| `backend/app/utils/jwt_handler.py` | `services/auth-service/internal/infrastructure/security/jwt_manager.go` | JWT token management |
| `backend/app/services/auth_service.py` | `services/auth-service/internal/application/usecases/` | Use cases for auth operations |

**Status**: Phase 1 Target

---

### 2. RIDE SERVICE

| Source | Destination | Notes |
|--------|-------------|-------|
| `backend/app/routes/ride.py` | `services/ride-service/internal/interfaces/rest/handlers/` | REST handlers |
| `backend/app/models/ride.py` | `services/ride-service/internal/domain/entities/ride.go` | Ride entity mapping |
| `backend/app/services/ride_matching.py` | `services/dispatch-service/internal/application/` | Move to separate dispatch service |
| `backend/app/websocket/socket_handler.py` | `services/websocket-gateway/internal/handlers/` | WebSocket handlers for realtime updates |

**Status**: Phase 6 Target

---

### 3. DRIVER SERVICE

| Source | Destination | Notes |
|--------|-------------|-------|
| `backend/app/routes/driver.py` | `services/driver-service/internal/interfaces/rest/handlers/` | Driver profile routes |
| `backend/app/models/driver.py` | `services/driver-service/internal/domain/entities/driver.go` | Driver entity |
| GPS location updates | `services/gps-service/internal/` | Realtime GPS service (separate) |

**Status**: Phase 6 Target

---

### 4. PAYMENT & WALLET

| Source | Destination | Notes |
|--------|-------------|-------|
| `backend/app/services/payment_service.py` | `services/payment-service/internal/application/usecases/` | Payment processing |
| Fare calculation logic | `services/pricing-service/internal/` | Separate pricing service |
| Balance management | `services/wallet-service/internal/domain/` | Immutable ledger (not in current codebase) |

**Status**: Phase 10 Target

---

### 5. ADMIN DASHBOARD

| Source | Destination | Notes |
|--------|-------------|-------|
| `backend/app/routes/admin.py` | Multiple destinations | Admin APIs distributed across services |
| `src/components/admin/` | `apps/admin-dashboard/app/` | Next.js admin dashboard |
| Admin analytics logic | `services/analytics-service/internal/` | Analytics service |

**Status**: Phase 15 Target

---

### 6. RIDER WEB APP

| Source | Destination | Notes |
|--------|-------------|-------|
| `src/components/user/` | `apps/rider-web/app/` | Rider dashboard (Next.js) |
| `src/pages/` | `apps/rider-web/app/` | Next.js App Router |
| `src/context/AuthContext.tsx` | `packages/auth-client/src/` | Shared auth context |
| `src/services/api.ts` | `packages/api-client/src/` | API client SDK |

**Status**: Phase 15 Target

---

### 7. DRIVER WEB APP

| Source | Destination | Notes |
|--------|-------------|-------|
| `src/components/driver/` | `apps/driver-web/app/` | Driver dashboard (Next.js) |
| `src/components/common/Map.tsx` | `packages/ui-kit/src/` or `packages/maps-sdk/` | Shared map component |

**Status**: Phase 15 Target

---

### 8. MOBILE APP

| Source | Destination | Notes |
|--------|-------------|-------|
| React rider + driver components | `apps/flutter-mobile/lib/features/` | **CONVERT** React to Flutter |
| `src/types/` | `apps/flutter-mobile/lib/models/` | Dart models |
| `src/services/socket.ts` | `apps/flutter-mobile/packages/websocket_sdk/` | WebSocket SDK for Flutter |

**Status**: Phase 15 Target

**CONVERSION EFFORT**: High - requires rewriting React components in Flutter/Dart

---

### 9. SHARED UTILITIES & TYPES

| Source | Destination | Notes |
|--------|-------------|-------|
| `src/types/` | `shared/schemas/` | TypeScript types → JSON schemas + protobuf |
| `src/utils/` | `shared/utilities/` or `packages/geo-utils/` | General utilities |
| `backend/app/utils/` | `shared/utilities/` | Backend utilities |

**Status**: Phase 0-1 Target

---

### 10. CONFIGURATION & ENVIRONMENT

| Source | Destination | Notes |
|--------|-------------|-------|
| `backend/.env` | `env/development/.env`, `env/production/.env` | Environment configs |
| `backend/app/config.py` | `platform/config/` | Configuration management |

**Status**: Phase 0 (in progress)

---

## Phase-by-Phase Implementation

### Phase 0: Foundation ✓
- [x] Create directory structure
- [x] Set up root configs (package.json, tsconfig.json, turbo.json)
- [x] Create migration mapping
- [ ] Initialize git repository
- [ ] Set up CI/CD pipelines

### Phase 1: Core Infrastructure (Next)
- **Database**: PostgreSQL + PostGIS setup
- **Auth Service**: Rewrite `backend/app/routes/auth.py` in Go
- **Kafka**: Event bus setup
- **Kong Gateway**: API Gateway
- **Redis**: Cache layer

### Phase 2-5: Additional Infrastructure
- Service mesh (Istio/Linkerd)
- Feature flags
- Event governance
- Vault integration

### Phase 6: Core Services
- Ride Service
- Driver Service
- Dispatch Service (from `ride_matching.py`)
- GPS Service (realtime location)
- WebSocket Gateway

### Phase 7-10: Domain Services
- Pricing Service
- Payment Service
- Wallet Service (new - immutable ledger)
- Pooling Service
- Notification Service

### Phase 11-15: Advanced & Frontend
- Safety, Fraud, Analytics services
- Admin Dashboard (from `src/components/admin/`)
- Rider Web (from `src/components/user/`)
- Driver Web (from `src/components/driver/`)
- **Flutter Mobile** (convert from React)

### Phase 16-20: Production
- Kubernetes deployment
- Helm charts
- Observability (Prometheus, Grafana, Loki, Jaeger)
- CI/CD integration
- Production hardening

---

## Code Reuse Strategy

### What to Port As-Is
- ✓ Business logic algorithms (ride matching, fare calculation)
- ✓ Database schema concepts (adapt to PostgreSQL)
- ✓ API endpoint patterns (map to new service boundaries)
- ✓ UI components (rider/driver dashboards → Next.js)

### What to Refactor
- ✗ Architecture: Monolith → Microservices
- ✗ Language: Python FastAPI → Go (type safety, performance)
- ✗ Database: MongoDB → PostgreSQL + PostGIS (relational, geospatial)
- ✗ Communication: Sync only → Event-driven + Sync
- ✗ Infrastructure: Local Docker → Kubernetes

### What to Build New
- Event contracts & Kafka topics
- Immutable wallet ledger
- Safety & fraud services
- ML pipelines (demand, ETA, surge, pooling optimization)
- Observability infrastructure
- Multi-region deployment

---

## Dependency Management

### Backend Services → Go Modules
- Copy relevant Go code to `services/*/internal/`
- Use shared packages in `packages/` for common logic
- Vendor dependencies via `go.mod`

### Frontend Apps → pnpm Workspaces
- React dashboards as separate Next.js apps
- Shared UI kit in `packages/ui-kit/`
- Shared APIs client in `packages/api-client/`

### Data/Contracts → Protobuf & JSON Schema
- Define gRPC services in `shared/protobufs/`
- Define Kafka events in `shared/contracts/kafka/`
- Define REST in `shared/contracts/rest/`

---

## Migration Execution Checklist

### Phase 0 ✓
- [x] Create directory structure (119 dirs)
- [x] Root configs (package.json, tsconfig.json, turbo.json, .gitignore)
- [x] README with architecture overview
- [x] This migration mapping document

### Phase 1 (Next)
- [ ] Create auth-service Go skeleton
- [ ] Port JWT logic from `backend/app/utils/jwt_handler.py`
- [ ] Port user/auth models to Go
- [ ] Set up PostgreSQL migrations
- [ ] Configure Kong Gateway
- [ ] Set up Kafka topics

### Phase 2
- [ ] Create ride-service skeleton
- [ ] Port ride models and logic
- [ ] Create dispatch-service
- [ ] Implement ride matching algorithm

### Phase 3
- [ ] Create gps-service
- [ ] Port WebSocket handlers from `socket_handler.py`
- [ ] Set up Redis GEO
- [ ] Implement realtime location streaming

### Phase 4-5
- [ ] Remaining platform infrastructure
- [ ] Event bus implementation
- [ ] Saga orchestration setup

### Phase 6-10
- [ ] Remaining domain services
- [ ] Event contracts
- [ ] Integration tests

### Phase 11-15
- [ ] Next.js apps (admin, rider, driver)
- [ ] Flutter mobile app (convert React components)
- [ ] API client SDKs

### Phase 16-20
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] Observability setup
- [ ] CI/CD pipelines
- [ ] Production testing

---

## File Naming Conventions

### Go Services
```
services/auth-service/
├── internal/
│   ├── application/
│   │   ├── commands/
│   │   ├── queries/
│   │   ├── dto/
│   │   └── usecases/
│   ├── domain/
│   │   ├── entities/
│   │   ├── events/
│   │   ├── services/
│   │   └── valueobjects/
│   ├── infrastructure/
│   │   ├── grpc/
│   │   ├── kafka/
│   │   ├── postgres/
│   │   ├── redis/
│   │   └── security/
│   └── interfaces/
│       ├── rest/
│       │   ├── handlers/
│       │   ├── middleware/
│       │   └── routes/
│       ├── grpc/
│       └── websocket/
├── cmd/
│   └── api/
│       └── main.go
└── migrations/
    └── 001_initial_schema.sql
```

### Shared Contracts
```
shared/contracts/
├── events/
│   ├── ride/
│   │   ├── ride_created.go
│   │   ├── ride_matched.go
│   │   └── ride_completed.go
│   ├── payment/
│   │   ├── payment_requested.go
│   │   └── payment_completed.go
│   └── driver/
│       ├── driver_online.go
│       └── driver_location_updated.go
├── grpc/
│   ├── ride_service.proto
│   ├── auth_service.proto
│   └── ...
├── rest/
│   ├── auth_contract.ts
│   ├── ride_contract.ts
│   └── ...
└── kafka/
    └── event_envelope.go
```

---

## Success Criteria

✓ Phase 0 complete when:
- Directory structure exists
- All config files present
- Migration mapping documented
- Git repo initialized
- CI/CD basic pipeline working

✓ Phase 1 complete when:
- Auth service can issue JWT tokens
- PostgreSQL migrations run successfully
- Kong Gateway routes requests
- Kafka topics created
- Tests pass for all new services

✓ Full migration complete when:
- All 18+ services deployed to Kubernetes
- Admin, rider, driver dashboards working in Next.js
- Flutter mobile app functional on iOS & Android
- Observability stack collecting metrics
- All existing functionality ported with no regressions

---

## Questions & Decisions

### Q: Should we keep MongoDB?
**A**: No. Migrate to PostgreSQL + PostGIS for:
- Geospatial queries (PostGIS)
- ACID transactions
- Better cost at scale
- Industry standard for enterprise

### Q: Keep FastAPI?
**A**: No. Migrate to Go for:
- Type safety
- Concurrency model (goroutines)
- Deployment (single binary)
- Performance for realtime GPS

### Q: React or Flutter for mobile?
**A**: Flutter for:
- Single codebase (iOS + Android)
- Better offline support (African markets)
- Native performance
- Better geolocation APIs

### Q: Kafka vs RabbitMQ?
**A**: Kafka for:
- Event sourcing patterns
- Replay capability
- Stream processing
- Enterprise grade

---

**Next Step**: Begin Phase 1 - create Auth Service Go skeleton
