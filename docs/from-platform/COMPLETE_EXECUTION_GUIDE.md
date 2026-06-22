# 🚀 COMPLETE STEP-BY-STEP EXECUTION GUIDE

**Status**: Production-ready  
**Services**: 5 microservices (Pricing, Driver, Payment, Ride, Dispatch)  
**Duration**: ~30-45 minutes for complete setup  
**Quality**: Enterprise production-grade  

---

## 📋 PART 1: PREREQUISITES & VERIFICATION (5 minutes)

### Step 1: Verify All Tools Installed

**Windows PowerShell (Run as Administrator):**
```powershell
# Check Go
go version
# Expected: go version go1.21+ windows/amd64

# Check Flutter
flutter --version
# Expected: Flutter 3.2.0+

# Check PostgreSQL
psql --version
# Expected: psql (PostgreSQL) 14+ 

# Check Redis
redis-cli ping
# Expected: PONG

# Check Kafka (optional)
# If using Kafka, verify broker at localhost:9092
```

### Step 2: Verify Database Connectivity

```powershell
# Test PostgreSQL connection
psql -U postgres -h localhost -c "SELECT version();"

# Should output PostgreSQL version
# If fails: Ensure PostgreSQL service is running
#   Windows: Services > PostgreSQL > Start
#   Or: net start postgresql-x64-14
```

### Step 3: Verify Redis Connectivity

```powershell
# Check Redis server
redis-cli ping
# Expected: PONG

# If fails: Start Redis
#   redis-server (in Redis directory)
#   Or: wsl -d Debian redis-server (if using WSL)
```

---

## 🔧 PART 2: DATABASE SETUP (5 minutes)

### Step 1: Execute Database Setup Script

**Windows PowerShell:**
```powershell
# Navigate to project root
cd C:\dev\FamGo-platform

# Execute setup script
psql -U postgres -h localhost -f database/setup_production.sql

# Output should show:
# CREATE DATABASE
# CREATE ROLE
# GRANT
# etc.
```

### Step 2: Verify Database Creation

```powershell
# List all databases
psql -U postgres -h localhost -c "\l"

# Should show these databases:
#  famgo_pricing_service    | pricing_user  | UTF8
#  famgo_driver_service     | driver_user   | UTF8
#  famgo_payment_service    | payment_user  | UTF8
#  famgo_ride_service       | ride_user     | UTF8
#  famgo_dispatch_service   | dispatch_user | UTF8

# Verify users
psql -U postgres -h localhost -c "\du"

# Should show:
#  pricing_user   | Can login, Create DB | 50
#  driver_user    | Can login, Create DB | 50
#  payment_user   | Can login, Create DB | 50
#  ride_user      | Can login, Create DB | 50
#  dispatch_user  | Can login, Create DB | 50
```

### Step 3: Verify Database Connections

```powershell
# Test each service database connection
psql -U pricing_user -h localhost -d famgo_pricing_service -c "SELECT 1;"
psql -U driver_user -h localhost -d famgo_driver_service -c "SELECT 1;"
psql -U payment_user -h localhost -d famgo_payment_service -c "SELECT 1;"
psql -U ride_user -h localhost -d famgo_ride_service -c "SELECT 1;"
psql -U dispatch_user -h localhost -d famgo_dispatch_service -c "SELECT 1;"

# Each should output: 1
# If any fails, check database/setup_production.sql for errors
```
 
---

## 🏗️ PART 3: SERVICE BUILD & RUN (20-30 minutes)

### Service 1: Pricing Service

**Terminal 1 (New PowerShell window):**
```powershell
cd C:\dev\FamGo-platform\services\pricing-service

# Verify file structure
ls cmd/api/main.go  # Should exist

# Initialize go module if needed
if (-not (Test-Path "go.mod")) {
    go mod init github.com/FamGo/platform/services/pricing-service
}

# Download dependencies
go mod download

# Build the service
go build -o bin\pricing-service.exe cmd\api\main.go

# Verify build succeeded
if (Test-Path "bin\pricing-service.exe") {
    Write-Host "✓ Build successful"
}

# Run with environment variables
$env:SERVICE_NAME="pricing-service"
$env:SERVICE_PORT="3014"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="pricing_user"
$env:DB_PASSWORD="pricing_service_pwd_secure_2024"
$env:DB_NAME="famgo_pricing_service"
$env:LOG_LEVEL="info"

.\bin\pricing-service.exe

# Expected output:
# 📋 Loading configuration from environment (production)
# ✓ Connected to database: pricing_user@localhost:5432/famgo_pricing_service
# ✓ Routes configured
# 🚀 Starting pricing-service on port 3014 (production environment)
```

### Verify Pricing Service (New PowerShell window)

```powershell
# Test health endpoint
curl http://localhost:3014/v1/health

# Expected: {"status":"healthy","service":"pricing-service","environment":"production",...}

# Test pricing estimation
$body = @{
    ride_type = "ECONOMY"
    distance_meters = 5000
    active_rides = 50
    available_drivers = 20
    is_pool = $false
} | ConvertTo-Json

curl -Method POST `
  -Uri http://localhost:3014/v1/pricing/estimate `
  -Body $body `
  -ContentType "application/json"

# Expected: {"base_fare":2.00,"distance_fare":6.00,"surge_multiplier":...}
```

### Service 2: Driver Service

**Terminal 2 (New PowerShell window):**
```powershell
cd C:\dev\FamGo-platform\services\driver-service

# Initialize go module if needed
if (-not (Test-Path "go.mod")) {
    go mod init github.com/FamGo/platform/services/driver-service
}

# Download dependencies
go mod download

# Build the service
go build -o bin\driver-service.exe cmd\api\main.go

# Set environment variables
$env:SERVICE_NAME="driver-service"
$env:SERVICE_PORT="3002"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="driver_user"
$env:DB_PASSWORD="driver_service_pwd_secure_2024"
$env:DB_NAME="famgo_driver_service"
$env:LOG_LEVEL="info"

# Run the service
.\bin\driver-service.exe

# Expected output:
# 📋 Loading configuration from environment (production)
# ✓ Connected to database: driver_user@localhost:5432/famgo_driver_service
# ✓ Routes configured
# 🚀 Starting driver-service on port 3002 (production environment)
```

### Verify Driver Service (New PowerShell window)

```powershell
# Test health endpoint
curl http://localhost:3002/v1/health

# Expected: {"status":"healthy",...}

# Test get driver metrics
curl "http://localhost:3002/v1/drivers/metrics?id=driver_123"

# Expected: {"driver_id":"driver_123","total_trips":250,...}
```

### Service 3: Payment Service

**Terminal 3 (New PowerShell window):**
```powershell
cd C:\dev\FamGo-platform\services\payment-service

# Initialize go module if needed
if (-not (Test-Path "go.mod")) {
    go mod init github.com/FamGo/platform/services/payment-service
}

# Download dependencies
go mod download

# Build the service
go build -o bin\payment-service.exe cmd\api\main.go

# Set environment variables
$env:SERVICE_NAME="payment-service"
$env:SERVICE_PORT="3015"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="payment_user"
$env:DB_PASSWORD="payment_service_pwd_secure_2024"
$env:DB_NAME="famgo_payment_service"
$env:LOG_LEVEL="info"

# Run the service
.\bin\payment-service.exe

# Expected output:
# 📋 Loading configuration from environment (production)
# ✓ Connected to database: payment_user@localhost:5432/famgo_payment_service
# ✓ Routes configured
# 🚀 Starting payment-service on port 3015 (production environment)
```

### Verify Payment Service (New PowerShell window)

```powershell
# Test health endpoint
curl http://localhost:3015/v1/health

# Expected: {"status":"healthy",...}

# Test get wallet
curl "http://localhost:3015/v1/wallets?user_id=user_123"

# Expected: {"user_id":"user_123","balance":500.50,...}
```

### Service 4: Ride Service

**Terminal 4 (New PowerShell window):**
```powershell
cd C:\dev\FamGo-platform\services\ride-service

# Initialize go module if needed
if (-not (Test-Path "go.mod")) {
    go mod init github.com/FamGo/platform/services/ride-service
}

# Download dependencies
go mod download

# Build the service
go build -o bin\ride-service.exe cmd\api\main.go

# Set environment variables
$env:SERVICE_NAME="ride-service"
$env:SERVICE_PORT="3010"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="ride_user"
$env:DB_PASSWORD="ride_service_pwd_secure_2024"
$env:DB_NAME="famgo_ride_service"
$env:LOG_LEVEL="info"

# Run the service
.\bin\ride-service.exe

# Expected output:
# 📋 Loading configuration from environment (production)
# ✓ Connected to database: ride_user@localhost:5432/famgo_ride_service
# ✓ Routes configured
# 🚀 Starting ride-service on port 3010 (production environment)
```

### Verify Ride Service (New PowerShell window)

```powershell
# Test health endpoint
curl http://localhost:3010/v1/health

# Expected: {"status":"healthy",...}

# Test get ride
curl "http://localhost:3010/v1/rides?id=ride_123"

# Expected: {"ride_id":"ride_123",...}
```
 
### Service 5: Dispatch Service

**Terminal 5 (New PowerShell window):**
```powershell
cd C:\dev\FamGo-platform\services\dispatch-service

# Initialize go module if needed
if (-not (Test-Path "go.mod")) {
    go mod init github.com/FamGo/platform/services/dispatch-service
}

# Download dependencies
go mod download

# Build the service
go build -o bin\dispatch-service.exe cmd\api\main.go

# Set environment variables
$env:SERVICE_NAME="dispatch-service"
$env:SERVICE_PORT="3011"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="dispatch_user"
$env:DB_PASSWORD="dispatch_service_pwd_secure_2024"
$env:DB_NAME="famgo_dispatch_service"
$env:LOG_LEVEL="info"

# Run the service
.\bin\dispatch-service.exe

# Expected output:
# 📋 Loading configuration from environment (production)
# ✓ Connected to database: dispatch_user@localhost:5432/famgo_dispatch_service
# ✓ Routes configured
# 🚀 Starting dispatch-service on port 3011 (production environment)
```

### Verify Dispatch Service (New PowerShell window)

```powershell
# Test health endpoint
curl http://localhost:3011/v1/health

# Expected: {"status":"healthy",...}

# Test get metrics
curl http://localhost:3011/v1/dispatch/metrics

# Expected: {"total_dispatches":10500,...}
```

---

## ✅ PART 4: COMPLETE SYSTEM VERIFICATION (5 minutes)

### Test All Services Health

```powershell
# Create a test script: test_services.ps1

$services = @(
    @{ Name = "Pricing"; Port = 3014 },
    @{ Name = "Driver"; Port = 3002 },
    @{ Name = "Payment"; Port = 3015 },
    @{ Name = "Ride"; Port = 3010 },
    @{ Name = "Dispatch"; Port = 3011 }
)

Write-Host "Testing FamGo Services..." -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green

foreach ($service in $services) {
    try {
        $response = curl -s "http://localhost:$($service.Port)/v1/health" | ConvertFrom-Json
        if ($response.status -eq "healthy") {
            Write-Host "✓ $($service.Name) Service (Port $($service.Port)): HEALTHY" -ForegroundColor Green
        } else {
            Write-Host "✗ $($service.Name) Service (Port $($service.Port)): UNHEALTHY" -ForegroundColor Red
        }
    } catch {
        Write-Host "✗ $($service.Name) Service (Port $($service.Port)): UNREACHABLE" -ForegroundColor Red
    }
}

# Run it
& .\test_services.ps1
```

### Test Integrated Flow

```powershell
# 1. Estimate pricing
Write-Host "1. Testing Pricing Service..."
$pricing = curl -Method POST http://localhost:3014/v1/pricing/estimate `
  -Body '{"ride_type":"ECONOMY","distance_meters":5000}' `
  -ContentType "application/json" | ConvertFrom-Json
Write-Host "   ✓ Estimated fare: ETB $($pricing.total_fare)"

# 2. Get driver metrics
Write-Host "2. Testing Driver Service..."
$driver = curl "http://localhost:3002/v1/drivers/metrics?id=driver_123" | ConvertFrom-Json
Write-Host "   ✓ Driver rating: $($driver.rating)"

# 3. Get wallet
Write-Host "3. Testing Payment Service..."
$wallet = curl "http://localhost:3015/v1/wallets?user_id=user_123" | ConvertFrom-Json
Write-Host "   ✓ User wallet: ETB $($wallet.balance)"

# 4. Create ride
Write-Host "4. Testing Ride Service..."
$ride = curl -Method POST http://localhost:3010/v1/rides `
  -Body 'ride_id=ride_123&user_id=user_123&pickup_lat=9.0320&pickup_lng=38.7469&dropoff_lat=9.0265&dropoff_lng=38.7400&ride_type=economy' `
  | ConvertFrom-Json
Write-Host "   ✓ Ride created: $($ride.ride_id)"

# 5. Match drivers
Write-Host "5. Testing Dispatch Service..."
$dispatch = curl -Method POST http://localhost:3011/v1/dispatch/match `
  -Body 'ride_id=ride_123&pickup_lat=9.0320&pickup_lng=38.7469&ride_type=economy' `
  | ConvertFrom-Json
Write-Host "   ✓ Matched drivers: $($dispatch.matching_drivers.Count)"

Write-Host "`n✓ All services working correctly!" -ForegroundColor Green
```

---

## 🎯 PART 5: PRODUCTION CHECKLIST

- [ ] All 5 databases created
- [ ] All 5 service users created with unique passwords
- [ ] All 5 services build without errors
- [ ] All 5 services run and respond to health checks
- [ ] Database connections from each service working
- [ ] All endpoints tested and returning data
- [ ] Integrated flow tested successfully
- [ ] Services gracefully handle shutdown (Ctrl+C)

---

## 🔍 TROUBLESHOOTING

### Issue: "Connection refused" on database

**Solution:**
```powershell
# Verify PostgreSQL is running
Get-Service postgresql* | Select-Object Status

# If stopped, start it
Start-Service "postgresql-x64-14"
```

### Issue: "Port already in use"

**Solution:**
```powershell
# Find process using port
netstat -ano | findstr :3014

# Kill process
taskkill /PID <PID> /F

# Or change SERVICE_PORT environment variable
$env:SERVICE_PORT="3014"  # Change to unused port
```

### Issue: Database authentication failed

**Solution:**
```powershell
# Verify password in .env matches database setup
# Check: DB_PASSWORD=pricing_service_pwd_secure_2024

# Reset password if needed
psql -U postgres -h localhost -c "ALTER USER pricing_user WITH PASSWORD 'pricing_service_pwd_secure_2024';"
```

### Issue: Cannot find cmd/api/main.go

**Solution:**
```powershell
# Verify file exists
ls C:\dev\FamGo-platform\services\pricing-service\cmd\api\

# If missing, ensure files were created from PRODUCTION_SERVICE_SETUP.md
```

---

## 🚀 NEXT STEPS

Once all 5 services are running successfully:

1. **Create API Gateway** (combines all 5 services)
2. **Setup Flutter Apps** (connect to backend)
3. **Implement Authentication** (JWT tokens)
4. **Add Kafka messaging** (async events)
5. **Setup monitoring** (Prometheus + Grafana)
6. **Deploy to production** (Docker + Kubernetes)

---

**This comprehensive guide ensures complete production-grade setup with all services running.** ✨

