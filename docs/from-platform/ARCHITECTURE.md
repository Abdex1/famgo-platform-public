# FamGo Platform - Enterprise Architecture Document

**Version**: 1.0
**Status**: Phase 0 Complete
**Last Updated**: 2024
**Next Review**: After Phase 1

---

## 1. SYSTEM OVERVIEW

### Vision
Build an enterprise-grade, distributed mobility platform optimized for African markets with:
- Event-driven microservices architecture
- Real-time GPS tracking
- Intelligent ride pooling
- Fraud detection & safety
- Multi-provider payment integration
- AI/ML optimization engines

### Principles
- **Event-Driven**: All state changes emit Kafka events
- **Fault-Tolerant**: Circuit breakers, retries, DLQs
- **Observable**: Complete distributed tracing & metrics
- **Scalable**: Horizontal scaling via Kubernetes
- **Geospatial**: PostGIS + Redis GEO for location queries
- **Secure**: Vault, mTLS, RBAC, device fingerprinting
- **Pooling-First**: Economics optimized for shared rides

---

## 2. HIGH-LEVEL ARCHITECTURE

```
                              INTERNET
                                 │
                    Global CDN + WAF (Cloudflare)
                                 │
                         API EDGE / INGRESS
                                 │
        ┌────────────────────────────────────────────────┐
        │          API GATEWAY LAYER                    │
        │  Kong: Auth, Routing, Rate Limits, WAF       │
        └────────────┬─────────────────────────────────┘
                     │
        ┌────────────┴────────────────────────────────────────┐
        │                                                     │
        ▼                                                     ▼

┌──────────────────────┐                    ┌─────────────────────┐
│ SYNCHRONOUS LAYER    │                    │ ASYNCHRONOUS LAYER  │
├──────────────────────┤                    ├─────────────────────┤
│ REST APIs (/v1/*)    │                    │ Kafka Event Bus     │
│ gRPC (internal)      │                    │ Topics:             │
│ WebSocket (realtime) │                    │ - ride.created      │
└──────────┬───────────┘                    │ - ride.matched      │
           │                                │ - payment.*         │
           │                                │ - driver.*          │
           │                                └──────────┬──────────┘
           │                                           │
    ┌──────▼───────────────────────────────────────────▼─────────┐
    │              MICROSERVICES LAYER (18+ services)           │
    ├─────────────────────────────────────────────────────────────┤
    │                                                             │
    │  Core Domain Services:                                    │
    │  • Auth Service (JWT, OTP, RBAC)                         │
    │  • User Service (Rider profiles)                         │
    │  • Driver Service (Driver profiles, onboarding)          │
    │  • Ride Service (Trip lifecycle)                         │
    │                                                             │
    │  Matching & Optimization:                                 │
    │  • Dispatch Service (Driver matching algorithm)          │
    │  • Pooling Service (Route compatibility, occupancy)      │
    │  • Pricing Service (Fare calculation, surge)             │
    │  • GPS Service (Realtime location streaming)             │
    │                                                             │
    │  Finance & Payments:                                      │
    │  • Payment Service (Multi-provider: Telebirr, CBE)       │
    │  • Wallet Service (Immutable ledger)                     │
    │  • Subscription Service (Commute passes)                 │
    │                                                             │
    │  Operations:                                               │
    │  • Notification Service (SMS, Push, Email)               │
    │  • Analytics Service (BI, metrics aggregation)           │
    │  • Safety Service (SOS, monitoring, detection)           │
    │  • Fraud Service (Fraud scoring, prevention)             │
    │                                                             │
    │  Advanced Features:                                        │
    │  • Smart Pickup Service (AI location recommendations)    │
    │  • Voice Booking Service (IVR)                           │
    │  • WebSocket Gateway (Realtime communication)            │
    │                                                             │
    └──────┬──────────────────────────────────────────────────┬──┘
           │                                                  │
    ┌──────▼───────────────────────────────────────────────────▼─────┐
    │           PLATFORM ENGINEERING LAYER                         │
    ├───────────────────────────────────────────────────────────────┤
    │ • Event Bus & Kafka                                          │
    │ • Service Discovery & Registry                              │
    │ • API Policies & Rate Limiting                              │
    │ • Feature Flags & Configuration                             │
    │ • Saga Orchestration (Distributed transactions)             │
    │ • Circuit Breakers & Resilience                             │
    │ • Telemetry & Tracing                                       │
    └──────┬────────────────────────────────────────────────────┬───┘
           │                                                    │
    ┌──────▼────────────────────────────────────────────────────▼────┐
    │              DATA & INFRASTRUCTURE LAYER                      │
    ├──────────────────────────────────────────────────────────────┤
    │                                                               │
    │ Transactional:   PostgreSQL 14+ with PostGIS              │
    │ Cache:          Redis 7+ (sessions, GEO, locks)           │
    │ Events:         Kafka 3.5+                                │
    │ Analytics:      ClickHouse                                │
    │ Search:         Elasticsearch                             │
    │ Storage:        S3-compatible (object store)              │
    │ Embeddings:     PostgreSQL pgvector                       │
    │                                                               │
    └──────┬────────────────────────────────────────────────────┬───┘
           │                                                    │
    ┌──────▼────────────────────────────────────────────────────▼────┐
    │            OBSERVABILITY & SECURITY LAYER                    │
    ├──────────────────────────────────────────────────────────────┤
    │                                                               │
    │ Metrics:    Prometheus + Grafana                            │
    │ Logs:       Loki (log aggregation)                          │
    │ Traces:     Jaeger (distributed tracing)                   │
    │ Telemetry:  OpenTelemetry SDK                              │
    │ Errors:     Sentry (error tracking)                        │
    │                                                               │
    │ Security:   Vault (secrets), mTLS, RBAC, Audit logging    │
    │ WAF:        Cloudflare + Kong                              │
    │                                                               │
    └──────────────────────────────────────────────────────────────┘
           │
           ▼
    ┌──────────────────────────────────────────────────────────────┐
    │             AI/ML OPTIMIZATION LAYER                        │
    ├──────────────────────────────────────────────────────────────┤
    │ • Demand Prediction (LSTM, Prophet)                        │
    │ • ETA Prediction (XGBoost)                                 │
    │ • Surge Price Prediction (Ensemble)                        │
    │ • Pool Optimization (ILP, genetic algorithms)              │
    │ • Fraud Detection (Isolation Forest, LSTM)                 │
    │                                                               │
    └──────────────────────────────────────────────────────────────┘
```

---

## 3. SERVICE BOUNDARIES

### Core Domain Services

#### 1. Auth Service
**Responsibility**: Authentication, authorization, session management
**Technology**: Go, PostgreSQL, Redis
**Key Features**:
- JWT token generation & validation
- OTP (one-time password) for security
- Device fingerprinting
- Session management
- RBAC (3 roles: rider, driver, admin)
- Token refresh & expiration
- Audit logging

**API**:
```
POST   /v1/auth/register
POST   /v1/auth/login
POST   /v1/auth/refresh
POST   /v1/auth/logout
POST   /v1/auth/otp/request
POST   /v1/auth/otp/verify
GET    /v1/auth/me
```

**Events**:
- `auth.user.registered`
- `auth.user.logged_in`
- `auth.user.logged_out`
- `auth.token.refreshed`

---

#### 2. User Service
**Responsibility**: Rider profile & management
**Technology**: Go, PostgreSQL
**Key Features**:
- Profile management
- KYC (Know Your Customer) data
- Preferences
- Rating & feedback
- Ride history
- Wallet balance

---

#### 3. Driver Service
**Responsibility**: Driver profile, onboarding, verification
**Technology**: Go, PostgreSQL, S3
**Key Features**:
- Driver registration & onboarding
- Document management (license, insurance)
- Vehicle information
- Background checks
- Ratings & performance
- Availability status
- Earnings tracking

---

#### 4. Ride Service
**Responsibility**: Ride lifecycle & state management
**Technology**: Go, PostgreSQL, Kafka
**State Machine**:
```
REQUESTED → MATCHING → MATCHED → ACCEPTED → IN_PROGRESS → COMPLETED/CANCELLED
```

**Not responsible for**:
- Driver matching (Dispatch Service)
- GPS tracking (GPS Service)
- Pooling (Pooling Service)
- Pricing (Pricing Service)

---

#### 5. Dispatch Service
**Responsibility**: Intelligent driver matching
**Technology**: Go, PostgreSQL, Redis GEO
**Algorithm**:
1. Filter available drivers within 5km
2. Score each driver:
   - ETA to pickup
   - Rating
   - Acceptance rate
   - Online duration
3. Rank and assign top driver
4. Handle timeouts & rejections

**NOT responsible for**:
- Ride state (Ride Service)
- Route optimization (Pooling Service)
- Pricing (Pricing Service)

---

#### 6. Pooling Service
**Responsibility**: Ride pooling optimization
**Technology**: Go, PostgreSQL, PostGIS
**Features**:
- Route compatibility scoring
- Pool formation (max 3 passengers)
- Detour calculation (<10 min)
- Female-only pools (safety)
- Occupancy optimization
- Commute pass integration

---

#### 7. Pricing Service
**Responsibility**: Fare calculation
**Technology**: Go, PostgreSQL, Redis
**Formula**:
```
Fare = BaseFare +
       (Distance × DistanceRate) +
       (Duration × TimeRate) +
       (SurgeFactor × BaseFare) +
       Taxes -
       Discounts
```

**Features**:
- Distance-based pricing
- Time-based pricing
- Surge pricing (demand-based)
- Discounts & promotions
- Subscription pricing
- Pooling discounts

---

#### 8. Payment Service
**Responsibility**: Payment processing
**Technology**: Go, PostgreSQL, Vault
**Providers**:
- **Critical**: Telebirr, CBE Birr, Cash
- **Medium**: Chapa, PayPal, Stripe

**State Machine**:
```
REQUESTED → PROCESSING → COMPLETED/FAILED
```

**Features**:
- Multi-provider integration
- Retry logic
- Reconciliation
- Refunds
- Audit trail

---

#### 9. Wallet Service
**Responsibility**: Immutable balance ledger
**Technology**: Go, PostgreSQL
**Architecture**:
- Append-only transaction log
- Never mutate balances
- Supports reversals
- Real-time balance calculation

**Example**:
```
wallet_transactions (immutable):
├── ID: 1, user_id: 123, amount: +100, type: TOP_UP
├── ID: 2, user_id: 123, amount: -45.50, type: RIDE_FARE
├── ID: 3, user_id: 123, amount: +50, type: PROMOTION
└── ID: 4, user_id: 123, amount: -5 (reversal of ID:2)
```

---

#### 10. GPS Service
**Responsibility**: Real-time location tracking
**Technology**: Go, Redis GEO, WebSocket
**Features**:
- Driver location updates (every 2 seconds)
- Live tracking for passengers
- Geospatial queries (nearby drivers)
- WebSocket streaming
- Location history

**Data Structure**:
```
Redis GEO:
GEOADD drivers:geo 13.361 38.115 "driver:123"

WebSocket events:
driver.location.updated {
  driver_id: "123",
  lat: 13.361,
  lng: 38.115,
  timestamp: 1704067200000
}
```

---

#### 11. Notification Service
**Responsibility**: SMS, Push, Email notifications
**Technology**: Go, external APIs
**Providers**:
- SMS: Telebirr, Africa's Talking, Twilio
- Push: Firebase Cloud Messaging
- Email: SendGrid

---

#### 12. Safety Service
**Responsibility**: Safety & emergency features
**Technology**: Go, PostgreSQL, Kafka
**Features**:
- SOS panic button
- Emergency contact escalation
- Trip sharing (real-time link)
- Route deviation detection
- Speed monitoring
- Harsh braking detection
- Inactivity detection

---

#### 13. Fraud Service
**Responsibility**: Fraud detection & prevention
**Technology**: Go, Python (ML), PostgreSQL
**Features**:
- Emulator detection (fake devices)
- GPS spoofing detection (impossible speed)
- Suspicious payment detection
- Fake trip detection
- Abuse pattern detection (cancellations, ratings)
- Real-time fraud scoring

**ML Models**:
- Isolation Forest (anomalies)
- Random Forest (classification)
- LSTM (sequence detection)

---

#### 14. Analytics Service
**Responsibility**: Platform metrics & BI
**Technology**: Go, ClickHouse, Kafka
**Features**:
- Event aggregation
- Real-time dashboards
- Custom reports
- Trend analysis

---

#### 15. Subscription Service
**Responsibility**: Monthly passes & commute subscriptions
**Technology**: Go, PostgreSQL
**Features**:
- Subscription plans
- Commute pass (daily, weekly, monthly)
- Billing & renewals
- Usage tracking

---

#### 16. Smart Pickup Service
**Responsibility**: AI-powered pickup recommendations
**Technology**: Go, Python (ML), PostGIS
**Features**:
- Suggested pickup locations
- Geo-fence management
- Accessibility features

---

#### 17. Voice Booking Service
**Responsibility**: Voice-activated ride booking
**Technology**: Go, Google Cloud Speech API
**Features**:
- Speech-to-text
- NLU (natural language understanding)
- Confirmation flow
- Error handling

---

#### 18. WebSocket Gateway
**Responsibility**: Real-time bidirectional communication
**Technology**: Go, Gorilla WebSocket
**Features**:
- Connection pooling
- Message routing
- Pub/Sub (Redis backed)
- Broadcast capability

---

## 4. TECHNOLOGY DECISIONS

### Backend Services: Go
**Why Go?**
- Type safety (prevents runtime errors)
- Concurrency model (goroutines, channels)
- Single binary deployment
- High performance (native compilation)
- Memory efficient
- Excellent for microservices

**Frameworks**:
- stdlib (http, net/http)
- Gorilla (WebSocket, router)
- GORM (database ORM)

---

### Frontend: Next.js + React
**Why Next.js?**
- SSR/SSG for performance
- File-based routing
- API routes for backends
- Excellent TypeScript support
- Vercel deployment

**Dashboards**:
- Admin Dashboard (ops, metrics)
- Rider Dashboard (booking, history)
- Driver Dashboard (rides, earnings)

---

### Mobile: Flutter
**Why Flutter?**
- Single codebase (iOS + Android)
- Hot reload for development
- Native performance
- Excellent for offline (important in Africa)
- Growing adoption in Africa

**Modules**:
- Rider app (booking, tracking)
- Driver app (acceptance, navigation)
- Shared services (auth, payments)

---

### Database: PostgreSQL + PostGIS
**Why PostgreSQL?**
- ACID transactions
- JSON support
- PostGIS for geospatial queries
- pgvector for ML embeddings
- Mature, battle-tested
- Excellent for relational data

**Extensions**:
- PostGIS: Geographic queries
- pgvector: Vector embeddings
- uuid-ossp: UUID generation

---

### Cache: Redis
**Why Redis?**
- In-memory performance
- Redis GEO for spatial queries
- Session management
- Rate limiting
- Pub/Sub for messaging

---

### Event Streaming: Kafka
**Why Kafka?**
- Event sourcing patterns
- Event replay capability
- At-least-once delivery
- Partitioning for scalability
- Stream processing

---

## 5. DATA FLOW

### Ride Request Flow

```
1. Rider requests ride
   POST /v1/rides
   → Ride Service creates ride

2. Ride Service emits event
   ride.created event → Kafka
   
3. Dispatch Service consumes event
   Finds matching drivers
   → ride.matching.started event

4. Best driver selected
   → ride.driver.assigned event
   
5. Driver notified (push)
   Notification Service sends notification

6. Driver accepts
   → ride.accepted event
   
7. Ride starts
   GPS Service streams location
   passenger sees live tracking
   → ride.started event

8. Ride completes
   Pricing Service calculates fare
   Payment Service processes payment
   Wallet Service records transaction
   → ride.completed event

9. Feedback collected
   Both parties rate each other
   → feedback.submitted event

10. Analytics aggregated
    Events flow to ClickHouse
    Dashboard updated
```

---

### Event-Driven Payment Flow

```
ride.completed event
    ↓
Payment Service: Create payment request
    ↓
Emit: payment.requested
    ↓
Multiple consumers:
├─ Payment Provider (Telebirr, CBE)
├─ Fraud Service (check for fraud)
├─ Analytics Service (track metrics)
└─ Notification Service (notify user)
    ↓
Payment Provider responds
    ↓
Emit: payment.completed OR payment.failed
    ↓
If completed:
├─ Wallet Service: Record transaction
├─ Ride Service: Mark ride paid
└─ Notification Service: Send receipt
```

---

## 6. DEPLOYMENT TOPOLOGY

### Local Development
```
Docker Compose:
├─ PostgreSQL
├─ Redis
├─ Kafka
├─ Kong
├─ Service 1
├─ Service 2
└─ Monitoring (Prometheus, Grafana, Loki, Jaeger)
```

### Staging
```
Kubernetes (AWS EKS):
├─ RDS PostgreSQL (managed)
├─ ElastiCache Redis (managed)
├─ Managed Kafka (MSK)
├─ 3 zones for HA
├─ Auto-scaling (2-10 replicas per service)
├─ Load balancing (ALB)
└─ Monitoring (Prometheus, Grafana)
```

### Production
```
Kubernetes Multi-Region (AWS EKS):
├─ Primary Region:
│  ├─ RDS PostgreSQL (Multi-AZ)
│  ├─ ElastiCache Redis (Multi-AZ)
│  ├─ Managed Kafka (Multi-AZ)
│  ├─ 3 zones HA
│  ├─ Auto-scaling (5-50 replicas)
│  └─ Monitoring
│
├─ Secondary Region:
│  ├─ Read-only replica
│  ├─ Disaster recovery
│  └─ Failover automation
│
└─ Global:
   ├─ CloudFlare CDN
   ├─ Route 53 DNS
   └─ Global load balancing
```

---

## 7. SECURITY ARCHITECTURE

### Authentication
- JWT tokens (issued by Auth Service)
- Token expiry (15 min access, 7 day refresh)
- Device fingerprinting (prevent token theft)
- OTP for sensitive operations (payments)

### Authorization
- RBAC: 3 roles (rider, driver, admin)
- Service-to-service: mTLS
- API Gateway: Kong with JWT validation

### Data Protection
- TLS/SSL everywhere (in-transit)
- Sensitive data encrypted at-rest (Vault)
- Database encryption (AWS RDS encryption)
- PII masking in logs

### Secrets Management
- HashiCorp Vault (all secrets)
- Auto-rotation
- Audit trail
- Different secrets per environment

### Network Security
- VPC with private subnets
- Security groups (least privilege)
- WAF (Cloudflare + Kong)
- DDoS protection

---

## 8. OBSERVABILITY

### Metrics (Prometheus + Grafana)
```
Application Metrics:
├─ Request latency (p50, p95, p99)
├─ Error rates
├─ Request count
├─ Active connections

Business Metrics:
├─ Rides per hour
├─ Revenue per hour
├─ Driver acceptance rate
├─ Average fare
└─ Pool conversion rate

Infrastructure Metrics:
├─ CPU utilization
├─ Memory usage
├─ Disk I/O
└─ Network traffic
```

### Logs (Loki)
```
Service logs → Loki
Query by:
├─ Service name
├─ Request ID
├─ User ID
├─ Error level
└─ Timestamp range
```

### Traces (Jaeger)
```
Request flow:
API Gateway → Auth Service → Ride Service
            → Dispatch Service → GPS Service
            → Payment Service

Full trace shows:
├─ Service call sequence
├─ Timing (where is latency?)
├─ Errors (where did it fail?)
└─ Dependencies (what else was involved?)
```

---

## 9. DISASTER RECOVERY

### RTO & RPO
- **RTO** (Recovery Time Objective): < 1 hour
- **RPO** (Recovery Point Objective): < 5 minutes

### Backup Strategy
- Database: Continuous replication + snapshots
- Kafka: Multi-region replication
- Configuration: GitOps (stored in Git)

### Failover Process
1. Detect outage (health checks)
2. Failover DNS to secondary region
3. Promote read replica to primary
4. Replay Kafka events from last checkpoint
5. Verify data consistency
6. Restore traffic

---

## 10. SCALING STRATEGY

### Horizontal Scaling
- Services scale independently (Kubernetes HPA)
- Database: Read replicas (analytics), primary (writes)
- Redis: Cluster mode (multiple masters)
- Kafka: Partition by ride_id (parallelism)

### Bottleneck Management
- Database: Connection pooling, caching
- Cache: TTL policies, eviction
- Kafka: Parallel consumers
- GPS: WebSocket batching

### Capacity Planning
- Monitor: p95 latency, CPU utilization
- Alert at: 70% CPU
- Scale at: 80% CPU
- Plan ahead for growth

---

## 11. DESIGN PATTERNS

### Event Sourcing
- All changes are events
- Append-only event store (Kafka)
- Replay for state reconstruction

### CQRS (Command Query Responsibility Segregation)
- Commands: Modify state (Write)
- Queries: Read state (Read)
- Separate models for each

### Saga Pattern
- Distributed transactions across services
- Orchestration (central coordinator)
- Choreography (event-driven)

### Circuit Breaker
- Detect failures
- Fail fast
- Half-open state for recovery

---

## 12. CONSISTENCY MODEL

### Strong Consistency
- Auth operations (user registration)
- Payment processing
- Wallet transactions

### Eventual Consistency
- User profile updates
- Driver ratings
- Analytics aggregation

### Conflict Resolution
- Last-write-wins for most data
- Custom reconciliation for financial data

---

## 13. FUTURE EXTENSIBILITY

### Super-App Features (Phase 21+)
- Food delivery
- Parcel delivery
- Financial services (lending, insurance)
- Entertainment (concert tickets)
- E-commerce integration

### Architecture supports:
- Plugin-based services
- Shared payment infrastructure
- Shared user profiles
- Common notification system
- Unified analytics

---

## 14. COMPLIANCE & REGULATIONS

### Data Protection (GDPR-like)
- Consent management
- Right to be forgotten
- Data portability
- Audit trails

### Financial Compliance
- PCI-DSS (payment processing)
- AML/KYC (know your customer)
- Transaction logging
- Regulatory reporting

---

**Document Version**: 1.0
**Last Updated**: 2024
**Next Review**: After Phase 1 Architecture Validation
