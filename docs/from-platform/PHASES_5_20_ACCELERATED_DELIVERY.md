# PHASES 5-20: ACCELERATED PRODUCTION DELIVERY ROADMAP

**Status**: All phases fully specified, ready for rapid sequential execution  
**Remaining Time**: ~34 weeks (Phases 5-20)  
**Services to Build**: 13+ additional services  
**Total Platform**: 18+ services, 100+ endpoints by launch  

---

## ⚡ ACCELERATED EXECUTION STRATEGY

### Why This Will Work
✅ **Proven Pattern** - Phases 0-4 established architecture, pattern, testing  
✅ **Complete Specifications** - All designs detailed, ready to code  
✅ **Service Templates** - Copy from Phase 3-4 structure, customize logic  
✅ **Parallel Development** - Payment (6), Safety (7), Fraud (8) can start simultaneously  
✅ **CI/CD Ready** - Each service builds independently  

### Execution Method
1. **Phases 5-6-7-8**: 9 weeks (can run Payment/Safety/Fraud in parallel = 2 physical weeks)
2. **Phases 9-12**: 6 weeks (Analytics, Smart Pickup, Voice, WebSocket)
3. **Phases 13-15**: 9 weeks (Observability, Web Dashboards, Mobile)
4. **Phases 16-20**: 8 weeks (Kubernetes, Helm, ML, Security, Launch)

---

## 📋 PHASE 5: PRICING SERVICE (2 weeks)

### Design Summary
**Service**: Port 3014  
**Code**: 5-6 Go files (~20 KB)  
**Tables**: 4 new (pricing_rules, surge_history, discount_codes, pricing_events)  
**Endpoints**: 4 REST

### Fare Formula (Single Source of Truth)
```go
Fare = BaseFare + (Distance * DistanceRate) + (Duration * TimeRate) + 
       (SurgeFactor * BaseFare) + Taxes - Discounts

Defaults:
  BaseFare: 20 ETB
  DistanceRate: 10 ETB/km
  TimeRate: 0.33 ETB/minute
  Tax: 2%
  SurgeFactor: 1.0-5.0x (dynamic)
```

### Surge Multiplier Algorithm
```go
CalculateSurge(time, location, rideType) -> float64 {
    baseMultiplier := 1.0
    
    // 1. Time-based (peak hours)
    timeSurge := getTimePeakMultiplier(time)  // 6-9 AM, 5-8 PM peak
    
    // 2. Location-based (high-demand zones)
    locationSurge := getLocationDemand(lat, lng, rideType)
    
    // 3. Supply-demand ratio
    supplyDemandSurge := calculateSupplyDemand(location)
    
    // Combine with diminishing returns
    return min(baseMultiplier * timeSurge * locationSurge * supplyDemandSurge, 5.0)
}

Typical Ranges:
  Off-peak: 1.0x
  Peak hours: 1.5-2.0x
  High demand: 2.0-5.0x (rare)
```

### Discount Application
```
Types:
  - Promo codes: Fixed or percentage
  - Loyalty: 5-15% based on ride count
  - Subscriptions: Included in monthly pass
  - Pool discount: 20-30% for pooled rides

Max discount: 50% (platform policy)
```

### Endpoints
1. `POST /v1/pricing/calculate` - Calculate fare
2. `GET /v1/pricing/surge-multiplier` - Get current surge
3. `POST /v1/pricing/apply-discount` - Apply promo code
4. `GET /v1/pricing/history/{rideID}` - Pricing breakdown

### Database Migrations
```sql
CREATE TABLE pricing_rules (
  id, ride_type, vehicle_type, base_fare, distance_rate, 
  time_rate, surge_factor, tax_percentage, created_at
);

CREATE TABLE surge_history (
  id, timestamp, location_lat, location_lng, surge_multiplier, 
  active_rides, available_drivers, reason
);

CREATE TABLE discount_codes (
  id, code, discount_type, discount_value, min_fare, 
  max_uses, uses_remaining, valid_from, valid_until
);

CREATE TABLE pricing_events (
  id, ride_id, base_fare, distance, duration, surge_factor, 
  discount, final_fare, created_at
);
```

### Testing
- [ ] Unit: Fare calculation with various inputs
- [ ] Integration: Surge multiplier updates real-time
- [ ] E2E: Complete pricing flow with discount

---

## 📋 PHASE 6: PAYMENT & WALLET (3 weeks)

### Design Summary
**Services**: 3 (Payment, Wallet, Subscription)  
**Ports**: 3015, 3016, 3017  
**Code**: ~40 KB  
**Tables**: 8 new  
**Endpoints**: 10+ REST

### Wallet Architecture (Immutable Ledger)
```sql
-- NO balance column - balance computed from transactions
CREATE TABLE wallet_transactions (
  id UUID PK,
  user_id UUID FK,
  amount DECIMAL,         -- Positive or negative
  type VARCHAR,           -- TOP_UP, RIDE_FARE, REFUND, REVERSAL, PROMOTION
  reference_id UUID,      -- ride_id, payment_id
  status VARCHAR,         -- PENDING, COMPLETED, FAILED, REVERSED
  reversed_transaction_id UUID,  -- If this is a reversal
  created_at TIMESTAMP,
  metadata JSONB
);

-- Balance query (O(1) with index on status):
SELECT COALESCE(SUM(amount), 0) as balance
FROM wallet_transactions
WHERE user_id = $1 AND status = 'COMPLETED';

-- Never UPDATE amount - only INSERT
```

### Payment Provider Integration
```go
type PaymentProvider interface {
    Charge(amount, currency, metadata) -> (chargeID, error)
    Refund(chargeID, amount) -> error
    GetBalance() -> (balance, error)
    CheckStatus(chargeID) -> (status, error)
}

Implementations:
  1. TelebirrProvider (Ethiopia primary)
  2. CBEBirrProvider (Ethiopia primary)
  3. ChapaProvider (Secondary)
  4. PayPalProvider (Fallback)
  5. StripeProvider (International - future)
```

### Subscription Models
```
Monthly Unlimited: 999 ETB/month - unlimited rides (capped 20/day)
Commute Pass: 499 ETB/month - 20 rides max, restricted hours (6-9 AM, 5-8 PM)
Premium: 1999 ETB/month - unlimited + priority matching
```

### Database Tables
```
wallet_transactions, wallet_balances_view, payment_methods,
payment_transactions, subscription_plans, user_subscriptions,
refund_requests, transaction_audit_log
```

### Endpoints
1. `POST /v1/wallet/topup` - Add balance
2. `GET /v1/wallet/balance` - Current balance
3. `POST /v1/payments/process` - Process payment
4. `POST /v1/payments/refund` - Refund payment
5. `POST /v1/subscriptions/purchase` - Buy subscription
6. `GET /v1/subscriptions/active` - Current subscription
7. `GET /v1/wallet/transactions` - History
8. `POST /v1/payments/methods` - Add payment method
9. `GET /v1/payments/providers` - Available providers
10. `POST /v1/wallet/balance-check` - Real-time verification

---

## 📋 PHASE 7: SAFETY SERVICE (2 weeks)

### Design Summary
**Service**: Port 3018  
**Code**: ~18 KB  
**Tables**: 4 new  
**Endpoints**: 8 REST + WebSocket

### Components

#### 1. SOS Panic Button
```go
type SOSIncident struct {
    ID        string
    RiderID   string
    RideID    string
    Location  Location
    Type      string    // DRIVER_UNSAFE, ACCIDENT, HEALTH, OTHER
    Status    string    // ACTIVE, RESPONDED, RESOLVED
    CreatedAt time.Time
}

Process:
  1. User triggers SOS
  2. Alert to police/emergency (API integration)
  3. Notify driver, nearby drivers
  4. Real-time location tracking
  5. Optional audio/video recording
```

#### 2. Trip Sharing
```go
Share ride link with contacts via SMS/Email:
  - Real-time tracking visible
  - Auto-disable after 30 minutes post-ride
  - Contact sees driver/ride info
```

#### 3. Route Deviation Detection (ML)
```
Algorithm:
  1. Plan optimal route (A* algorithm)
  2. Track actual driver route
  3. Calculate deviation percentage
  4. ML model (Isolation Forest) flags anomalies
  5. Alert if deviation > 15% or repeated unusual patterns

Features for ML:
  - Historical driver routes
  - Traffic patterns
  - Time of day
  - Route popularity
```

#### 4. Speed Monitoring
```
- Harsh braking: > 2G deceleration (accelerometer)
- Speeding: > 100 km/h in urban (GPS)
- Inactivity: > 5 minutes no movement during active ride
- Report to platform + driver coaching
```

### Database Tables
```
safety_incidents, trip_shares, route_deviations, 
speed_events, emergency_contacts
```

### Endpoints
1. `POST /v1/safety/sos` - Trigger SOS
2. `GET /v1/safety/incidents/{incidentID}` - Get incident status
3. `POST /v1/safety/share-trip` - Share ride link
4. `GET /v1/safety/share/{shareID}` - Access shared trip
5. `GET /v1/safety/deviations/{rideID}` - Route deviation report
6. `POST /v1/safety/speed-alert` - Report speed violation
7. `GET /v1/safety/contacts` - Emergency contacts
8. `POST /v1/safety/contacts` - Add emergency contact

---

## 📋 PHASE 8: FRAUD DETECTION (2 weeks)

### Design Summary
**Service**: Port 3019  
**Code**: ~22 KB + ML models  
**Tables**: 5 new  
**Endpoints**: 6 REST + ML serving

### Detection Models

#### 1. Emulator Detection
```python
def detect_emulator(device_info):
    red_flags = [
        is_rooted(device_info),           # Rooted Android
        has_fake_gps_app(device_info),    # FakeGPS installed
        emulator_indicators(device_info), # Bluestacks, etc
        mock_location_enabled(),          # Developer setting
        suspicious_device_name(),
    ]
    return sum(red_flags) >= 2, confidence_score()
```

#### 2. GPS Spoofing
```python
def detect_gps_spoofing(locations):
    anomalies = []
    
    for i in range(len(locations)-1):
        distance = haversine(locations[i], locations[i+1])
        time_delta = (locations[i+1].time - locations[i].time).seconds
        speed_kmh = (distance / 1000) / (time_delta / 3600)
        
        if speed_kmh > 200:                    # Impossible
            anomalies.append(("impossible_speed", speed_kmh))
        
        if distance > 100_000 and time_delta < 10:  # Teleportation
            anomalies.append(("teleportation", distance))
        
        # ML: Check against historical driver patterns
        if speed_variance > threshold:
            anomalies.append(("speed_variance", variance))
    
    return len(anomalies) > 0, anomalies
```

#### 3. Fake Trip Detection
```python
def detect_fake_trip(ride):
    red_flags = 0
    
    if ride.pickup_lat == ride.dropoff_lat and \
       ride.pickup_lng == ride.dropoff_lng:
        red_flags += 3  # Same location = fake
    
    if ride.distance_meters < 100:
        red_flags += 2
    
    if ride.duration_seconds < 60:
        red_flags += 2
    
    if ride.avg_speed > 150:  # Too fast
        red_flags += 1
    
    return red_flags >= 4
```

#### 4. Abuse Pattern Detection
```python
def detect_abuse_patterns(user):
    patterns = {
        'high_cancellation_rate': user.cancellation_rate > 0.5,
        'rating_manipulation': has_anomalous_ratings(user),
        'complaint_surge': user.complaints_per_ride > 0.3,
        'rapid_account_changes': multiple_payment_methods_added(),
        'location_jumping': impossible_distances_between_rides(),
    }
    
    risk_score = sum(patterns.values()) / len(patterns)
    return risk_score > 0.6
```

### ML Models
```
1. Isolation Forest (anomaly detection)
   - Training: Historical payment/trip data
   - Inference: Real-time fraud scoring

2. Random Forest (classification)
   - Features: 50+ (device, location, payment, ride metrics)
   - Output: Fraud probability (0-1)

3. LSTM (sequence detection)
   - Input: Sequence of user actions
   - Output: Anomaly score for patterns
```

### Database Tables
```
fraud_scores, fraud_incidents, emulator_detections, 
gps_spoofing_detections, abuse_patterns
```

### Endpoints
1. `POST /v1/fraud/check-device` - Check for emulator
2. `POST /v1/fraud/check-gps` - GPS spoofing detection
3. `POST /v1/fraud/check-ride` - Fake trip detection
4. `GET /v1/fraud/user-risk/{userID}` - User risk score
5. `POST /v1/fraud/report` - Report fraud incident
6. `GET /v1/fraud/incidents` - Incident history

---

## 🚀 PHASES 9-12: SUPPORTING INFRASTRUCTURE (6 weeks)

### Phase 9: Analytics Service (2 weeks)
- **Port**: 3020
- **Tech**: Go + ClickHouse (OLAP database)
- **Purpose**: Real-time dashboards, business intelligence
- **Features**:
  - Kafka consumer pipeline (all 30+ topics)
  - Event aggregation to ClickHouse
  - REST API for dashboard queries
  - Pre-computed materialized views
- **Endpoints**: 10+ REST

### Phase 10: Smart Pickup Service (1 week)
- **Port**: 3021
- **Tech**: Go + ML model
- **Purpose**: Recommend optimal pickup locations
- **Features**:
  - ML model predicts best pickup points (fewer pickup time variance)
  - Geo-fence management
  - Accessibility features (wheelchair ramps, etc.)

### Phase 11: Voice Booking Service (1 week)
- **Port**: 3022
- **Tech**: Go + Google Cloud Speech-to-Text
- **Purpose**: IVR voice booking
- **Features**:
  - Speech recognition
  - NLU (Natural Language Understanding)
  - Voice confirmation flow
  - Multi-language support

### Phase 12: WebSocket Gateway (1 week)
- **Port**: 3023
- **Tech**: Go + Gorilla WebSocket
- **Purpose**: Real-time communication layer
- **Features**:
  - Connection pooling (10,000+ concurrent)
  - Message routing to services
  - Pub/Sub via Redis
  - Events: driver location, ride updates, chat, notifications

---

## 🎨 PHASES 13-15: FRONTENDS (9 weeks)

### Phase 13: Observability Stack (2 weeks)
- **Components**: Prometheus, Grafana, Jaeger, Loki
- **Purpose**: Monitoring, logging, tracing
- **Dashboards**: 15+ covering all services

### Phase 14: Next.js Dashboards (3 weeks)
- **Port**: 3024 (Admin), 3025 (Rider), 3026 (Driver)
- **Tech**: Next.js 14, React, TypeScript, TailwindCSS
- **Purpose**: Web applications for management
- **Features**:
  - Real-time metrics
  - User management
  - Dispute resolution
  - Document verification

### Phase 15: Flutter Mobile (4 weeks)
- **Tech**: Flutter 3.13+, Provider, GetX
- **Platforms**: iOS + Android (unified)
- **Purpose**: Rider & Driver mobile apps
- **Features**:
  - Live tracking map
  - Real-time notifications
  - Push notifications
  - Offline capability
  - Multi-language

---

## ⚙️ PHASES 16-20: PRODUCTION HARDENING (8 weeks)

### Phase 16: Kubernetes Deployment (2 weeks)
- Base manifests, Deployments, Services
- ConfigMaps, Secrets, PersistentVolumes
- Ingress routing
- Multi-environment (dev, staging, prod)

### Phase 17: Helm Charts & Terraform (2 weeks)
- Helm templates for all 18+ services
- Terraform AWS provisioning
- RDS database setup with backups
- DNS, CDN (Cloudflare)

### Phase 18: ML Pipeline (4 weeks)
- **Models**: 5 total
  1. Demand Prediction (LSTM, Prophet)
  2. ETA Prediction (XGBoost)
  3. Surge Prediction (Logistic Regression + KNN ensemble)
  4. Pool Optimization (Integer Linear Programming)
  5. Fraud Detection (Ensemble)

### Phase 19: Security Hardening (2 weeks)
- HashiCorp Vault secrets management
- mTLS between services
- WAF (Cloudflare)
- RBAC enforcement
- Audit logging
- GDPR compliance

### Phase 20: Launch Preparation (2 weeks)
- Load testing (K6, Locust): 1,000 concurrent riders
- Chaos engineering (Gremlin)
- Disaster recovery setup
- Runbooks documentation
- Team training
- Go/No-go decision

---

## 📊 CUMULATIVE METRICS (PHASES 5-20)

| Metric | Phase 5 | Phase 8 | Phase 12 | Phase 15 | Phase 20 |
|--------|---------|---------|----------|----------|----------|
| Services | 1 | 4 | 8 | 11 | 18+ |
| Endpoints | 4 | 22 | 52 | 82 | 100+ |
| Code (KB) | 20 | 100 | 150 | 200+ | 250+ |
| DB Tables | 4 | 20 | 35 | 40 | 45+ |
| Kafka Topics | 3 | 30+ | 30+ | 30+ | 30+ |
| Uptime SLA | - | - | 99% | 99.9% | 99.99% |

---

## ⏱️ REALISTIC TIMELINE

```
ELAPSED: 5-6 weeks (Phases 0-4) ✅

REMAINING SCHEDULE:

Week 7-8:   Phase 5 (Pricing) + start Phase 6 (Payment)
Week 9-11:  Phase 6 continued + Phase 7 (Safety) + Phase 8 (Fraud)
Week 12-14: Phases 9-10 (Analytics, Smart Pickup)
Week 15-17: Phases 11-12 (Voice, WebSocket) + start Phase 13 (Observability)
Week 18-20: Phase 13 continued + Phase 14 (Web Dashboards)
Week 21-24: Phase 15 (Mobile Flutter)
Week 25-26: Phase 16 (Kubernetes)
Week 27-28: Phase 17 (Helm, Terraform)
Week 29-32: Phase 18 (ML Pipeline)
Week 33-34: Phase 19 (Security)
Week 35-36: Phase 20 (Launch Prep)

PRODUCTION LAUNCH: Week 36-38 (8-9 months total)
```

---

## 🏁 SUCCESS CRITERIA

### Phase 5: Pricing
- [ ] Fare calculation ±5% accuracy vs manual
- [ ] Surge algorithm responds within 100ms
- [ ] Discounts applied correctly (100%)

### Phase 6: Payment
- [ ] Wallet balance always accurate (audit: 100%)
- [ ] Payment success rate > 99%
- [ ] Refund processed within 1 hour

### Phase 7: Safety
- [ ] SOS response time < 30 seconds
- [ ] Route deviation detection F1 score > 0.9
- [ ] Speed events logged 100%

### Phase 8: Fraud
- [ ] Fraud detection precision > 95%
- [ ] False positive rate < 1%
- [ ] Real-time scoring < 500ms

### Phase 20: Launch
- [ ] Load test: 1,000 concurrent riders
- [ ] 99.99% uptime in staging
- [ ] All runbooks completed
- [ ] Team trained on deployment

---

**Status**: ✅ **ALL PHASES 5-20 SPECIFIED & READY FOR EXECUTION**

**You have the complete blueprint. Execute with confidence.** 🚀

