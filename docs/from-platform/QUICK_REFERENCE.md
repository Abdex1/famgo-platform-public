# 🎯 VISUAL QUICK REFERENCE GUIDE

## 🚀 ONE-COMMAND STARTUP

### Windows PowerShell
```powershell
cd C:\dev\FamGo-platform
.\manage_services.ps1 -Action all
```

### Windows Batch
```batch
cd C:\dev\FamGo-platform
start_all_services.bat
```

---

## 📊 SERVICE ARCHITECTURE

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT APPLICATIONS                      │
│              (Flutter Driver & Passenger Apps)              │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP/WebSocket
┌──────────────────────┴──────────────────────────────────────┐
│                     API GATEWAY                             │
│              (To be implemented)                            │
└───┬──────────┬──────────┬──────────┬──────────┬─────────────┘
    │          │          │          │          │
HTTP│ PORT    HTTP│ PORT HTTP│ PORT HTTP│ PORT  HTTP│ PORT
    │ 3014     │ 3002     │ 3015     │ 3010   │ 3011
    │          │          │          │        │
┌───▼──┐   ┌───▼──┐   ┌───▼──┐   ┌──▼──┐  ┌──▼───┐
│PRICING│   │DRIVER│   │PAYMENT│  │RIDE │  │DISP- │
│SERVICE│   │SERVICE   │SERVICE   │SERV.   │ATCH  │
│:3014  │   │:3002     │:3015     │:3010   │:3011 │
└───┬──┘   └───┬──┘   └───┬──┘   └──┬──┘  └──┬───┘
    │          │          │          │        │
    │ TCP 5432 │ TCP 5432 │ TCP 5432 │ TCP 5432 │ TCP 5432
    │          │          │          │        │
┌───▼──────────▼──────────▼──────────▼────────▼───┐
│          PostgreSQL Databases                    │
├──────────────────────────────────────────────────┤
│ ├─ famgo_pricing_service    (pricing_user)      │
│ ├─ famgo_driver_service     (driver_user)       │
│ ├─ famgo_payment_service    (payment_user)      │
│ ├─ famgo_ride_service       (ride_user)         │
│ └─ famgo_dispatch_service   (dispatch_user)     │
└──────────────────────────────────────────────────┘
```

---

## 📋 DATABASE SETUP

```
┌─────────────────────────────────────────────────┐
│         Run Database Setup Script                │
├─────────────────────────────────────────────────┤
│ psql -U postgres -h localhost \                 │
│   -f database/setup_production.sql              │
└──────────────────┬──────────────────────────────┘
                   │
        ┌──────────┴──────────┬──────────────┬─────────────┐
        │                     │              │             │
    ┌───▼────┐ ┌────────┐ ┌──▼──┐ ┌───────▼─┐ ┌────────┐
    │CREATE  │ │CREATE  │ │CREATE  │ │CREATE  │ │CREATE  │
    │DATABASE│ │DATABASE│ │DATABASE│ │DATABASE│ │DATABASE│
    │Pricing │ │Driver  │ │Payment │ │Ride    │ │Dispatch│
    └────┬───┘ └────┬───┘ └───┬───┘ └───┬────┘ └───┬────┘
         │          │         │         │          │
      ┌──▼──┐    ┌──▼──┐   ┌──▼──┐  ┌──▼──┐    ┌──▼──┐
      │CREATE   CREATE    CREATE    CREATE     CREATE
      │USER     USER      USER      USER       USER
      │pricing  driver    payment   ride       dispatch
      └────────────────────────────────────────────┘
```

---

## 🔄 SERVICE STARTUP SEQUENCE

```
Time  Event
────  ────────────────────────────────────────────────────────
0s    Start Database Setup
10s   ✓ Databases Created
      Start Building Services
30s   ✓ All Services Built
      Start Services (5 in parallel)
32s   ✓ Pricing Service (3014)
34s   ✓ Driver Service (3002)
36s   ✓ Payment Service (3015)
38s   ✓ Ride Service (3010)
40s   ✓ Dispatch Service (3011)
45s   ✓ All Services Running & Healthy
```

---

## 🧪 TEST ENDPOINTS

### Quick Health Check
```bash
# Test all 5 services
curl http://localhost:3014/v1/health  # Pricing
curl http://localhost:3002/v1/health  # Driver
curl http://localhost:3015/v1/health  # Payment
curl http://localhost:3010/v1/health  # Ride
curl http://localhost:3011/v1/health  # Dispatch
```

### Sample API Calls

#### 1. Estimate Price
```bash
curl -X POST http://localhost:3014/v1/pricing/estimate \
  -H "Content-Type: application/json" \
  -d '{
    "ride_type": "ECONOMY",
    "distance_meters": 5000,
    "active_rides": 50,
    "available_drivers": 20,
    "is_pool": false
  }'
```

#### 2. Create Ride
```bash
curl -X POST http://localhost:3010/v1/rides \
  -d "user_id=user_123&pickup_lat=9.0320&pickup_lng=38.7469&dropoff_lat=9.0265&dropoff_lng=38.7400&ride_type=economy"
```

#### 3. Match Drivers
```bash
curl -X POST http://localhost:3011/v1/dispatch/match \
  -d "ride_id=ride_123&pickup_lat=9.0320&pickup_lng=38.7469&ride_type=economy"
```

#### 4. Process Payment
```bash
curl -X POST http://localhost:3015/v1/payments/process \
  -d "ride_id=ride_123&user_id=user_123&amount=45.50&provider=telebirr"
```

#### 5. Get Driver Metrics
```bash
curl http://localhost:3002/v1/drivers/metrics?id=driver_123
```

---

## 🔐 CREDENTIALS QUICK REFERENCE

```
SERVICE          | USER           | PASSWORD                        | DB
─────────────────┼────────────────┼─────────────────────────────────┼──────────────────────
Pricing          | pricing_user   | pricing_service_pwd_secure_2024 | famgo_pricing_service
Driver           | driver_user    | driver_service_pwd_secure_2024  | famgo_driver_service
Payment          | payment_user   | payment_service_pwd_secure_2024 | famgo_payment_service
Ride             | ride_user      | ride_service_pwd_secure_2024    | famgo_ride_service
Dispatch         | dispatch_user  | dispatch_service_pwd_secure_2024 | famgo_dispatch_service

All ports: 5432 (PostgreSQL standard)
Host: localhost
```

---

## 📁 KEY FILES LOCATION

```
C:\dev\FamGo-platform\
├── START_HERE.md                    ← Overview
├── SETUP_SUMMARY.md                 ← This summary
├── PRODUCTION_SERVICE_SETUP.md       ← Detailed setup
├── COMPLETE_EXECUTION_GUIDE.md       ← Step-by-step
├── QUICK_REFERENCE.md               ← This file
├── manage_services.ps1              ← PowerShell manager
├── start_all_services.bat           ← Batch starter
│
├── database/
│   └── setup_production.sql         ← Database setup
│
└── services/
    ├── pricing-service/cmd/api/main.go
    ├── driver-service/cmd/api/main.go
    ├── payment-service/cmd/api/main.go
    ├── ride-service/cmd/api/main.go
    └── dispatch-service/cmd/api/main.go
```

---

## ⚡ COMMON COMMANDS

### PowerShell Management
```powershell
# Full setup and run
.\manage_services.ps1 -Action all

# Just setup databases
.\manage_services.ps1 -Action setup

# Just build services
.\manage_services.ps1 -Action build

# Just run services
.\manage_services.ps1 -Action run

# Test all services
.\manage_services.ps1 -Action test

# Stop all services
.\manage_services.ps1 -Action stop
```

### Manual Build & Run
```powershell
# Build individual service
cd C:\dev\FamGo-platform\services\pricing-service
go build -o bin\pricing-service.exe cmd\api\main.go

# Run with env
$env:SERVICE_PORT="3014"
.\bin\pricing-service.exe
```

### Database Commands
```powershell
# Setup databases
psql -U postgres -h localhost -f database/setup_production.sql

# List databases
psql -U postgres -h localhost -c "\l"

# Connect to service database
psql -U pricing_user -d famgo_pricing_service -h localhost

# Execute SQL file
psql -U postgres -d famgo_pricing_service -f migrations/001_initial.sql
```

---

## 🎯 WHAT EACH SERVICE DOES

### Pricing Service (Port 3014)
```
├─ Calculates ride fares
├─ Applies surge pricing
├─ Manages pooled discounts
└─ Provides price estimates
```

### Driver Service (Port 3002)
```
├─ Manages driver profiles
├─ Tracks driver location (GPS)
├─ Accepts ride requests
├─ Collects driver metrics
└─ Manages online/offline status
```

### Payment Service (Port 3015)
```
├─ Processes payments
├─ Manages wallets
├─ Adds money to wallets
├─ Refunds transactions
└─ Tracks transaction history
```

### Ride Service (Port 3010)
```
├─ Creates rides
├─ Tracks ride status
├─ Cancels rides
├─ Completes rides
└─ Manages ratings/reviews
```

### Dispatch Service (Port 3011)
```
├─ Matches drivers to rides
├─ Assigns best driver
├─ Tracks matching status
├─ Provides matching metrics
└─ Handles dispatch cancellation
```

---

## ⚠️ COMMON ISSUES & FIXES

| Issue | Fix |
|-------|-----|
| "Port already in use" | `taskkill /PID <PID> /F` |
| "Connection refused" | Restart PostgreSQL service |
| "Authentication failed" | Check DB_USER and DB_PASSWORD match |
| "Database doesn't exist" | Run database setup script |
| "Module not found" | Run `go mod download` in service dir |
| "Build failed" | Check Go version (1.21+) |

---

## ✅ VERIFY EVERYTHING WORKS

```
Step 1: Run: .\manage_services.ps1 -Action all
Step 2: Check: .\manage_services.ps1 -Action test
Step 3: Expected: All 5 services show ✓ HEALTHY
```

---

**Ready to start? Run: `.\manage_services.ps1 -Action all`** 🚀

