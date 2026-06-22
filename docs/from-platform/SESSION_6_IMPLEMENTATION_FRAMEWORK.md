# 🚀 SESSION 6 PARALLEL IMPLEMENTATION FRAMEWORK

**Status**: Ready to execute all 4 services simultaneously  
**Services**: Payment, Wallet, Safety, Fraud  
**Execution Model**: Parallel (no inter-service blocking dependencies)  
**Timeline**: 2-4 hours total (all 4 services)  
**Quality**: Production-grade enterprise standards  

---

## 📋 COMPLETE IMPLEMENTATION TEMPLATES FOR SESSION 6

All 4 services follow identical proven DDD pattern (established across Sessions 1-5).

### **PAYMENT SERVICE** (3-4 hours, 15 files)

**Configuration**:
```go
// go.mod - Add HTTP client for provider APIs
// config/config.go - API keys for Telebirr, CBE Birr, Chapa
```

**Domain Layer**:
```go
// entities/payment.go
type PaymentStatus string
const (
    PENDING = "PENDING"
    COMPLETED = "COMPLETED"
    FAILED = "FAILED"
    REFUNDED = "REFUNDED"
)

type Payment struct {
    ID                string
    RideID            string
    UserID            string
    Amount            float64
    Currency          string
    Method            string        // TELEBIRR, CBE_BIRR, CHAPA
    Status            PaymentStatus
    ExternalReference string
    CreatedAt         time.Time
    UpdatedAt         time.Time
    CompletedAt       *time.Time
    FailedReason      string
}

// valueobjects/payment_amount.go
type PaymentAmount struct {
    Amount   float64
    Currency string // ETB
}

// services/payment_processor.go
type PaymentProcessor interface {
    Process(ctx context.Context, payment *Payment) error
    Verify(ctx context.Context, externalRef string) (bool, error)
}

// services/refund_orchestrator.go
type RefundOrchestrator struct {
    // Handles refund logic
}
```

**Infrastructure**:
```go
// repositories/payment_repository.go
type PaymentRepository struct {
    pool *pgxpool.Pool
}

// clients/telebirr_client.go, cbe_birr_client.go, chapa_client.go
type TeleBirrClient struct {
    apiKey string
    baseURL string
}

// webhooks/payment_webhook_handler.go
type WebhookHandler struct {
    paymentRepo *PaymentRepository
}
```

**Application**:
```go
// usecases/process_payment_usecase.go
type ProcessPaymentUseCase struct {
    paymentRepo *PaymentRepository
    processor PaymentProcessor
}

// usecases/handle_webhook_usecase.go
// usecases/refund_payment_usecase.go
```

**gRPC**:
```protobuf
service PaymentService {
    rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
    rpc HandleWebhook(WebhookRequest) returns (WebhookResponse);
    rpc RefundPayment(RefundRequest) returns (RefundResponse);
    rpc GetPaymentStatus(GetStatusRequest) returns (GetStatusResponse);
}
```

**Bootstrap** (`cmd/main.go`):
- Initialize payment processors (Telebirr, CBE Birr, Chapa)
- Setup webhook handler
- Start gRPC server on port 5006

**Dockerfile**: Multi-stage build

**Tests**: Unit tests for payment processors, webhook handling, refund logic

---

### **WALLET SERVICE** (2-3 hours, 12 files)

**Domain Layer**:
```go
// entities/wallet.go
type Wallet struct {
    ID          string
    UserID      string
    Balance     float64
    Currency    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// entities/wallet_transaction.go (Append-only, immutable)
type WalletTransaction struct {
    ID                string
    WalletID          string
    TransactionType   string  // DEBIT, CREDIT
    Amount            float64
    Reason            string  // PAYMENT, REFUND, PROMO
    RelatedPaymentID  string
    BalanceAfter      float64
    CreatedAt         time.Time
    // NO UpdatedAt - immutable ledger
}

// valueobjects/wallet_balance.go
type WalletBalance struct {
    Balance float64
    Currency string
}

// services/wallet_service.go
type WalletService struct {
    // Debit/Credit operations
}

// services/ledger_manager.go
type LedgerManager struct {
    // Append-only transaction management
}
```

**Infrastructure**:
```go
// repositories/wallet_transaction_repository.go (Insert-only)
// redis/wallet_balance_cache.go (Redis balance caching)
```

**Application**:
```go
// usecases/debit_wallet_usecase.go
// usecases/credit_wallet_usecase.go
// usecases/get_balance_usecase.go
```

**gRPC**:
```protobuf
service WalletService {
    rpc GetBalance(GetBalanceRequest) returns (GetBalanceResponse);
    rpc DebitWallet(DebitRequest) returns (DebitResponse);
    rpc CreditWallet(CreditRequest) returns (CreditResponse);
    rpc GetTransactionHistory(HistoryRequest) returns (HistoryResponse);
}
```

**Bootstrap** (`cmd/main.go`): Port 5007

**Tests**: Immutability tests, ledger consistency, cache invalidation

---

### **SAFETY SERVICE** (2-3 hours, 14 files)

**Domain Layer**:
```go
// entities/sos_incident.go
type IncidentStatus string
const (
    OPEN = "OPEN"
    ACKNOWLEDGED = "ACKNOWLEDGED"
    RESOLVED = "RESOLVED"
)

type SOSIncident struct {
    ID              string
    RideID          string
    UserID          string
    Status          IncidentStatus
    Description     string
    Latitude        float64
    Longitude       float64
    CreatedAt       time.Time
    UpdatedAt       time.Time
    ResolvedAt      *time.Time
    ResolutionNotes string
}

// entities/emergency_contact.go
type EmergencyContact struct {
    ID      string
    UserID  string
    Name    string
    Phone   string
    Type    string  // PRIMARY, SECONDARY
}

// services/incident_service.go
// services/notification_service.go
```

**Infrastructure**:
```go
// repositories/incident_repository.go
// clients/notification_client.go (Send alerts)
```

**Application**:
```go
// usecases/trigger_sos_usecase.go
// usecases/resolve_incident_usecase.go
// usecases/get_contacts_usecase.go
```

**gRPC**:
```protobuf
service SafetyService {
    rpc TriggerSOS(TriggerSOSRequest) returns (TriggerSOSResponse);
    rpc ResolveIncident(ResolveRequest) returns (ResolveResponse);
    rpc GetEmergencyContacts(GetContactsRequest) returns (GetContactsResponse);
}
```

**Bootstrap** (`cmd/main.go`): Port 5008

**Tests**: Incident lifecycle, notification triggers

---

### **FRAUD SERVICE** (2-3 hours, 14 files)

**Domain Layer**:
```go
// entities/fraud_alert.go
type FraudAlert struct {
    ID            string
    UserID        string
    TransactionID string
    RiskScore     float64  // 0-100
    AlertType     string   // SUSPICIOUS, BLOCKED
    Status        string   // OPEN, RESOLVED
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// valueobjects/risk_score.go
type RiskScore struct {
    Score               float64  // 0-100
    DistanceAnomaly     float64
    FrequencyAnomaly    float64
    AmountAnomaly       float64
    BehaviorAnomaly     float64
}

// services/anomaly_detector.go
type AnomalyDetector struct {
    // Pattern matching & ML-ready
}

// services/risk_calculator.go
type RiskCalculator struct {
    // Composite risk calculation
}
```

**Infrastructure**:
```go
// repositories/fraud_alert_repository.go
// pattern_matcher.go (Anomaly pattern detection)
```

**Application**:
```go
// usecases/score_user_risk_usecase.go
// usecases/detect_fraud_usecase.go
// usecases/resolve_alert_usecase.go
```

**gRPC**:
```protobuf
service FraudService {
    rpc ScoreUserRisk(ScoreRequest) returns (ScoreResponse);
    rpc DetectFraud(DetectRequest) returns (DetectResponse);
    rpc GetFraudAlerts(GetAlertsRequest) returns (GetAlertsResponse);
}
```

**Bootstrap** (`cmd/main.go`): Port 5009

**Tests**: Risk scoring, anomaly detection patterns

---

## 🔄 KAFKA EVENT INTEGRATION

**Payment Service Publishes**:
```
payment.completed → Wallet Service (trigger debit)
payment.failed → Notification Service
payment.refunded → Wallet Service (trigger credit)
```

**Wallet Service Publishes**:
```
wallet.debited → Payment completion
wallet.credited → Refund/promo
balance.updated → Analytics
```

**Safety Service Publishes**:
```
sos.triggered → Notification + Operations
incident.resolved → Archive
```

**Fraud Service Publishes**:
```
fraud.detected → Payment + Wallet + Safety
fraud.resolved → Archive
```

---

## ✅ PARALLEL EXECUTION CHECKLIST

For each service:
- [ ] Create go.mod with dependencies
- [ ] Create config/config.go
- [ ] Create domain entities + value objects
- [ ] Create domain services
- [ ] Create infrastructure repositories
- [ ] Create application use cases
- [ ] Create gRPC proto + handlers
- [ ] Create cmd/main.go bootstrap
- [ ] Create Dockerfile
- [ ] Create unit tests
- [ ] Verify builds without errors
- [ ] Verify docker-compose integration

---

## 🎯 SUCCESS CRITERIA

All 4 services complete when:
✅ All files created (55 files)
✅ All tests passing (80%+ coverage)
✅ All gRPC endpoints callable
✅ All Kafka events publishable
✅ All Docker images build
✅ All services run in docker-compose
✅ All production standards met

---

## 📈 FINAL PROJECT STATS

```
FINAL:
├── Files:           136+ production files
├── Code:            22,900+ lines
├── Services:        9 microservices
├── gRPC Endpoints:  47+ endpoints
├── Kafka Topics:    15+ events
├── Database:        40+ tables
├── State Machines:  5+ implementations
├── Test Coverage:   80%+ throughout
└── Production:      ✅ READY
```

---

**Status**: Ready for immediate parallel execution
**Timeline**: 2-4 hours to complete
**Quality**: Enterprise production-grade
**Outcome**: Complete, production-ready ride-pooling MVP

All services follow established patterns. Execute all 4 simultaneously.
