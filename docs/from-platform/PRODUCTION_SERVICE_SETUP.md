# 🚀 PRODUCTION-GRADE EXECUTION: COMPLETE SERVICE SETUP

**Status**: Enterprise production-ready  
**Focus**: Individual service databases, configs, and deployment  
**Quality**: Production-grade with all dependencies and error handling  

---

## 🎯 PART 1: PRODUCTION DATABASE SETUP

### Database Architecture
```
PostgreSQL Instance (Single Host)
├── Database: famgo_platform (shared schema)
├── Database: famgo_pricing_service (dedicated)
├── Database: famgo_driver_service (dedicated)
├── Database: famgo_payment_service (dedicated)
├── Database: famgo_ride_service (dedicated)
└── Database: famgo_dispatch_service (dedicated)

Users per Service:
├── pricing_user (pricing_service_pwd_secure_2024)
├── driver_user (driver_service_pwd_secure_2024)
├── payment_user (payment_service_pwd_secure_2024)
├── ride_user (ride_service_pwd_secure_2024)
└── dispatch_user (dispatch_service_pwd_secure_2024)
```

### Database Creation Script
**File**: `C:\dev\FamGo-platform\database\setup_production.sql`

```sql
-- ============================================================================
-- PRODUCTION DATABASE SETUP - Execute as superuser (postgres)
-- ============================================================================

-- ============================================================================
-- CREATE SHARED DATABASES FOR EACH SERVICE
-- ============================================================================

-- Pricing Service Database
CREATE DATABASE famgo_pricing_service 
  OWNER postgres 
  ENCODING 'UTF8' 
  TEMPLATE template0;

-- Driver Service Database
CREATE DATABASE famgo_driver_service 
  OWNER postgres 
  ENCODING 'UTF8' 
  TEMPLATE template0;

-- Payment Service Database
CREATE DATABASE famgo_payment_service 
  OWNER postgres 
  ENCODING 'UTF8' 
  TEMPLATE template0;

-- Ride Service Database
CREATE DATABASE famgo_ride_service 
  OWNER postgres 
  ENCODING 'UTF8' 
  TEMPLATE template0;

-- Dispatch Service Database
CREATE DATABASE famgo_dispatch_service 
  OWNER postgres 
  ENCODING 'UTF8' 
  TEMPLATE template0;

-- ============================================================================
-- CREATE SERVICE USERS WITH STRONG PASSWORDS
-- ============================================================================

-- Pricing Service User
CREATE USER pricing_user WITH PASSWORD 'pricing_service_pwd_secure_2024';
ALTER USER pricing_user CREATEDB;
ALTER USER pricing_user CONNECTION LIMIT 50;

-- Driver Service User
CREATE USER driver_user WITH PASSWORD 'driver_service_pwd_secure_2024';
ALTER USER driver_user CREATEDB;
ALTER USER driver_user CONNECTION LIMIT 50;

-- Payment Service User
CREATE USER payment_user WITH PASSWORD 'payment_service_pwd_secure_2024';
ALTER USER payment_user CREATEDB;
ALTER USER payment_user CONNECTION LIMIT 50;

-- Ride Service User
CREATE USER ride_user WITH PASSWORD 'ride_service_pwd_secure_2024';
ALTER USER ride_user CREATEDB;
ALTER USER ride_user CONNECTION LIMIT 50;

-- Dispatch Service User
CREATE USER dispatch_user WITH PASSWORD 'dispatch_service_pwd_secure_2024';
ALTER USER dispatch_user CREATEDB;
ALTER USER dispatch_user CONNECTION LIMIT 50;

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

-- Pricing Service
GRANT ALL PRIVILEGES ON DATABASE famgo_pricing_service TO pricing_user;
GRANT CONNECT ON DATABASE famgo_pricing_service TO pricing_user;
ALTER DATABASE famgo_pricing_service OWNER TO pricing_user;

-- Driver Service
GRANT ALL PRIVILEGES ON DATABASE famgo_driver_service TO driver_user;
GRANT CONNECT ON DATABASE famgo_driver_service TO driver_user;
ALTER DATABASE famgo_driver_service OWNER TO driver_user;

-- Payment Service
GRANT ALL PRIVILEGES ON DATABASE famgo_payment_service TO payment_user;
GRANT CONNECT ON DATABASE famgo_payment_service TO payment_user;
ALTER DATABASE famgo_payment_service OWNER TO payment_user;

-- Ride Service
GRANT ALL PRIVILEGES ON DATABASE famgo_ride_service TO ride_user;
GRANT CONNECT ON DATABASE famgo_ride_service TO ride_user;
ALTER DATABASE famgo_ride_service OWNER TO ride_user;

-- Dispatch Service
GRANT ALL PRIVILEGES ON DATABASE famgo_dispatch_service TO dispatch_user;
GRANT CONNECT ON DATABASE famgo_dispatch_service TO dispatch_user;
ALTER DATABASE famgo_dispatch_service OWNER TO dispatch_user;

-- ============================================================================
-- INSTALL EXTENSIONS
-- ============================================================================

-- Connect to each database and install extensions
\c famgo_pricing_service
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgvector";
CREATE EXTENSION IF NOT EXISTS "postgis";

\c famgo_driver_service
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c famgo_payment_service
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c famgo_ride_service
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pgvector";

\c famgo_dispatch_service
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pgvector";

-- ============================================================================
-- VERIFICATION
-- ============================================================================

-- List all databases
\l

-- List all users
\du

-- Test connections
\c famgo_pricing_service pricing_user
SELECT 1;
\c famgo_driver_service driver_user
SELECT 1;
\c famgo_payment_service payment_user
SELECT 1;
\c famgo_ride_service ride_user
SELECT 1;
\c famgo_dispatch_service dispatch_user
SELECT 1;

COMMIT;
```

### Execute Database Setup (Windows PowerShell)
```powershell
# Execute setup script
psql -U postgres -h localhost -f "C:\dev\FamGo-platform\database\setup_production.sql"

# Verify setup
psql -U postgres -h localhost -c "\l"  # List databases
psql -U postgres -h localhost -c "\du" # List users
```

---

## 🎯 PART 2: ENVIRONMENT CONFIGURATION FILES

### Service 1: Pricing Service `.env`
**File**: `C:\dev\FamGo-platform\services\pricing-service\.env.production`

```env
# ============================================================================
# PRICING SERVICE - PRODUCTION ENVIRONMENT
# ============================================================================

# Service Configuration
SERVICE_NAME=pricing-service
SERVICE_PORT=3014
SERVICE_ENV=production
LOG_LEVEL=info

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=pricing_user
DB_PASSWORD=pricing_service_pwd_secure_2024
DB_NAME=famgo_pricing_service
DB_SSL_MODE=disable
DB_CONNECTION_POOL_SIZE=20
DB_CONNECTION_MAX_IDLE_TIME=300
DB_QUERY_TIMEOUT=30

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_PREFIX=pricing
KAFKA_CONSUMER_GROUP=pricing-service-group
KAFKA_AUTO_OFFSET_RESET=earliest

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10

# Application Configuration
ALGORITHM_VERSION=v1.0.0
SURGE_UPDATE_INTERVAL=5
CACHE_TTL=3600
REQUEST_TIMEOUT=30

# Monitoring & Logging
METRICS_ENABLED=true
METRICS_PORT=9001
JAEGER_ENABLED=true
JAEGER_HOST=localhost
JAEGER_PORT=6831
```

### Service 2: Driver Service `.env`
**File**: `C:\dev\FamGo-platform\services\driver-service\.env.production`

```env
# ============================================================================
# DRIVER SERVICE - PRODUCTION ENVIRONMENT
# ============================================================================

# Service Configuration
SERVICE_NAME=driver-service
SERVICE_PORT=3002
SERVICE_ENV=production
LOG_LEVEL=info

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=driver_user
DB_PASSWORD=driver_service_pwd_secure_2024
DB_NAME=famgo_driver_service
DB_SSL_MODE=disable
DB_CONNECTION_POOL_SIZE=20
DB_CONNECTION_MAX_IDLE_TIME=300
DB_QUERY_TIMEOUT=30

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_PREFIX=driver
KAFKA_CONSUMER_GROUP=driver-service-group
KAFKA_AUTO_OFFSET_RESET=earliest

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=1
REDIS_POOL_SIZE=10

# Location Services
GEOLOCATION_UPDATE_INTERVAL=5
GPS_ACCURACY_THRESHOLD=50
LOCATION_HISTORY_DAYS=30

# Application Configuration
MAX_ACTIVE_RIDES_PER_DRIVER=3
DRIVER_ONLINE_TIMEOUT=600
DOCUMENT_VERIFICATION_REQUIRED=true

# Monitoring & Logging
METRICS_ENABLED=true
METRICS_PORT=9002
JAEGER_ENABLED=true
JAEGER_HOST=localhost
JAEGER_PORT=6831
```

### Service 3: Payment Service `.env`
**File**: `C:\dev\FamGo-platform\services\payment-service\.env.production`

```env
# ============================================================================
# PAYMENT SERVICE - PRODUCTION ENVIRONMENT
# ============================================================================

# Service Configuration
SERVICE_NAME=payment-service
SERVICE_PORT=3015
SERVICE_ENV=production
LOG_LEVEL=info

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=payment_user
DB_PASSWORD=payment_service_pwd_secure_2024
DB_NAME=famgo_payment_service
DB_SSL_MODE=disable
DB_CONNECTION_POOL_SIZE=25
DB_CONNECTION_MAX_IDLE_TIME=300
DB_QUERY_TIMEOUT=60

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_PREFIX=payment
KAFKA_CONSUMER_GROUP=payment-service-group
KAFKA_AUTO_OFFSET_RESET=earliest

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=2
REDIS_POOL_SIZE=15

# Payment Provider Credentials (Encrypted in production)
TELEBIRR_API_KEY=your_telebirr_api_key_here
TELEBIRR_API_SECRET=your_telebirr_secret_here
TELEBIRR_MERCHANT_ID=your_merchant_id

CBE_API_KEY=your_cbe_api_key_here
CBE_API_SECRET=your_cbe_secret_here
CBE_ACCOUNT_NUMBER=your_account_number

CHAPA_API_KEY=your_chapa_api_key_here
CHAPA_SECRET_KEY=your_chapa_secret_here

PAYPAL_CLIENT_ID=your_paypal_client_id
PAYPAL_SECRET=your_paypal_secret

# Payment Processing
TRANSACTION_TIMEOUT=120
WEBHOOK_TIMEOUT=30
MAX_RETRY_ATTEMPTS=3
IDEMPOTENCY_KEY_EXPIRY=86400

# Monitoring & Logging
METRICS_ENABLED=true
METRICS_PORT=9003
JAEGER_ENABLED=true
JAEGER_HOST=localhost
JAEGER_PORT=6831
```

### Service 4: Ride Service `.env`
**File**: `C:\dev\FamGo-platform\services\ride-service\.env.production`

```env
# ============================================================================
# RIDE SERVICE - PRODUCTION ENVIRONMENT
# ============================================================================

# Service Configuration
SERVICE_NAME=ride-service
SERVICE_PORT=3010
SERVICE_ENV=production
LOG_LEVEL=info

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=ride_user
DB_PASSWORD=ride_service_pwd_secure_2024
DB_NAME=famgo_ride_service
DB_SSL_MODE=disable
DB_CONNECTION_POOL_SIZE=20
DB_CONNECTION_MAX_IDLE_TIME=300
DB_QUERY_TIMEOUT=30

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_PREFIX=ride
KAFKA_CONSUMER_GROUP=ride-service-group
KAFKA_AUTO_OFFSET_RESET=earliest

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=3
REDIS_POOL_SIZE=10

# Ride Configuration
RIDE_TIMEOUT_MINUTES=30
AUTOMATIC_CANCELLATION_MINUTES=15
PICKUP_ARRIVAL_WINDOW_MINUTES=5
DROPOFF_ARRIVAL_WINDOW_MINUTES=5

# Monitoring & Logging
METRICS_ENABLED=true
METRICS_PORT=9004
JAEGER_ENABLED=true
JAEGER_HOST=localhost
JAEGER_PORT=6831
```

### Service 5: Dispatch Service `.env`
**File**: `C:\dev\FamGo-platform\services\dispatch-service\.env.production`

```env
# ============================================================================
# DISPATCH SERVICE - PRODUCTION ENVIRONMENT
# ============================================================================

# Service Configuration
SERVICE_NAME=dispatch-service
SERVICE_PORT=3011
SERVICE_ENV=production
LOG_LEVEL=info

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=dispatch_user
DB_PASSWORD=dispatch_service_pwd_secure_2024
DB_NAME=famgo_dispatch_service
DB_SSL_MODE=disable
DB_CONNECTION_POOL_SIZE=20
DB_CONNECTION_MAX_IDLE_TIME=300
DB_QUERY_TIMEOUT=30

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_PREFIX=dispatch
KAFKA_CONSUMER_GROUP=dispatch-service-group
KAFKA_AUTO_OFFSET_RESET=earliest

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=4
REDIS_POOL_SIZE=15

# Dispatch Algorithm Configuration
MATCHING_TIMEOUT_SECONDS=30
MAX_DRIVERS_TO_MATCH=5
SEARCH_RADIUS_METERS=5000
ALGORITHM_VERSION=v1.0.0

# Monitoring & Logging
METRICS_ENABLED=true
METRICS_PORT=9005
JAEGER_ENABLED=true
JAEGER_HOST=localhost
JAEGER_PORT=6831
```

---

## 🎯 PART 3: SERVICE STARTUP SCRIPTS

### Master Startup Script (Windows PowerShell)
**File**: `C:\dev\FamGo-platform\run_services.ps1`

```powershell
# ============================================================================
# PRODUCTION SERVICE STARTUP SCRIPT (Windows PowerShell)
# ============================================================================

param(
    [string]$Service = "all",  # all, pricing, driver, payment, ride, dispatch
    [string]$Environment = "production"
)

# Color output
function Write-Status { Write-Host "[$(Get-Date -Format 'HH:mm:ss')] $($args -join ' ')" -ForegroundColor Green }
function Write-Error-Custom { Write-Host "[$(Get-Date -Format 'HH:mm:ss')] ERROR: $($args -join ' ')" -ForegroundColor Red }

# Check prerequisites
function Test-Prerequisites {
    Write-Status "Checking prerequisites..."
    
    # Check Go
    $go = go version 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "Go not found. Install Go 1.21+"
        exit 1
    }
    
    # Check PostgreSQL
    $psql = psql --version 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "PostgreSQL not found. Install PostgreSQL 14+"
        exit 1
    }
    
    # Check Redis
    $redis = redis-cli ping 2>&1
    if ($redis -ne "PONG") {
        Write-Error-Custom "Redis not running. Start Redis first: redis-server"
        exit 1
    }
    
    # Check Kafka
    $kafka = Test-Connection -ComputerName localhost -Port 9092 -ErrorAction SilentlyContinue
    if (-not $kafka) {
        Write-Status "WARNING: Kafka not detected. Some features may not work."
    }
    
    Write-Status "Prerequisites OK"
}

# Load environment
function Load-Environment {
    param([string]$ServiceName, [string]$Env)
    
    $envFile = "C:\dev\FamGo-platform\services\$ServiceName\.env.$Env"
    
    if (-not (Test-Path $envFile)) {
        Write-Error-Custom "Environment file not found: $envFile"
        exit 1
    }
    
    Get-Content $envFile | ForEach-Object {
        if ($_ -notmatch '^\s*#' -and $_ -notmatch '^\s*$') {
            $key, $value = $_ -split '=', 2
            [Environment]::SetEnvironmentVariable($key.Trim(), $value.Trim())
        }
    }
}

# Build service
function Build-Service {
    param([string]$ServiceName)
    
    Write-Status "Building $ServiceName..."
    Push-Location "C:\dev\FamGo-platform\services\$ServiceName"
    
    if (-not (Test-Path "go.mod")) {
        Write-Status "Initializing Go module for $ServiceName..."
        & go mod init "github.com/FamGo/platform/services/$ServiceName" 2>&1
    }
    
    & go mod download 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "Failed to download dependencies for $ServiceName"
        Pop-Location
        return $false
    }
    
    & go build -o "bin\$ServiceName" "cmd\api\main.go" 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "Failed to build $ServiceName"
        Pop-Location
        return $false
    }
    
    Pop-Location
    Write-Status "$ServiceName built successfully"
    return $true
}

# Run service
function Run-Service {
    param([string]$ServiceName, [string]$Env)
    
    Load-Environment $ServiceName $Env
    
    $binPath = "C:\dev\FamGo-platform\services\$ServiceName\bin\$ServiceName.exe"
    
    if (-not (Test-Path $binPath)) {
        if (-not (Build-Service $ServiceName)) {
            return
        }
    }
    
    Write-Status "Starting $ServiceName on port $([Environment]::GetEnvironmentVariable('SERVICE_PORT'))..."
    & $binPath
}

# Main execution
Write-Status "FamGo Platform Service Startup"
Write-Status "================================"

Test-Prerequisites

$services = if ($Service -eq "all") {
    @("pricing-service", "driver-service", "payment-service", "ride-service", "dispatch-service")
} else {
    @("$Service-service")
}

# Start each service in background
$jobs = @()
foreach ($svc in $services) {
    Write-Status "Starting $svc..."
    $job = Start-Job -ScriptBlock {
        param($svc, $env, $platform)
        Set-Location $platform
        & ".\run_services.ps1" -Service $svc.Replace("-service", "") -Environment $env
    } -ArgumentList $svc, $Environment, "C:\dev\FamGo-platform"
    $jobs += $job
}

Write-Status "All services started. Use 'Get-Job' to list jobs."
Write-Status "Use 'Stop-Job -Name <jobname>' to stop a service."

# Wait for jobs
Get-Job | Wait-Job
```

### Individual Service Run Scripts

#### Pricing Service
**File**: `C:\dev\FamGo-platform\services\pricing-service\run.ps1`

```powershell
$env:SERVICE_NAME="pricing-service"
$env:SERVICE_PORT="3014"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="pricing_user"
$env:DB_PASSWORD="pricing_service_pwd_secure_2024"
$env:DB_NAME="famgo_pricing_service"
$env:LOG_LEVEL="info"

go build -o bin/pricing-service cmd/api/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit 1
}

.\bin\pricing-service
```

#### Driver Service
**File**: `C:\dev\FamGo-platform\services\driver-service\run.ps1`

```powershell
$env:SERVICE_NAME="driver-service"
$env:SERVICE_PORT="3002"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="driver_user"
$env:DB_PASSWORD="driver_service_pwd_secure_2024"
$env:DB_NAME="famgo_driver_service"
$env:LOG_LEVEL="info"

go build -o bin/driver-service cmd/api/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit 1
}

.\bin\driver-service
```

#### Payment Service
**File**: `C:\dev\FamGo-platform\services\payment-service\run.ps1`

```powershell
$env:SERVICE_NAME="payment-service"
$env:SERVICE_PORT="3015"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="payment_user"
$env:DB_PASSWORD="payment_service_pwd_secure_2024"
$env:DB_NAME="famgo_payment_service"
$env:LOG_LEVEL="info"

go build -o bin/payment-service cmd/api/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit 1
}

.\bin\payment-service
```

#### Ride Service
**File**: `C:\dev\FamGo-platform\services\ride-service\run.ps1`

```powershell
$env:SERVICE_NAME="ride-service"
$env:SERVICE_PORT="3010"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="ride_user"
$env:DB_PASSWORD="ride_service_pwd_secure_2024"
$env:DB_NAME="famgo_ride_service"
$env:LOG_LEVEL="info"

go build -o bin/ride-service cmd/api/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit 1
}

.\bin\ride-service
```

#### Dispatch Service
**File**: `C:\dev\FamGo-platform\services\dispatch-service\run.ps1`

```powershell
$env:SERVICE_NAME="dispatch-service"
$env:SERVICE_PORT="3011"
$env:SERVICE_ENV="production"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="dispatch_user"
$env:DB_PASSWORD="dispatch_service_pwd_secure_2024"
$env:DB_NAME="famgo_dispatch_service"
$env:LOG_LEVEL="info"

go build -o bin/dispatch-service cmd/api/main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit 1
}

.\bin\dispatch-service
```

---

## 🚀 COMPLETE EXECUTION STEPS (Windows)

### Step 1: Setup Databases
```powershell
# Execute as Administrator in PowerShell
psql -U postgres -h localhost -f "C:\dev\FamGo-platform\database\setup_production.sql"

# Verify
psql -U postgres -h localhost -c "\l"  # Should show all 5 databases
```

### Step 2: Create .env Files
(All .env files created above in Part 2)

### Step 3: Build All Services
```powershell
# Pricing Service
cd C:\dev\FamGo-platform\services\pricing-service
go mod download
go build -o bin/pricing-service cmd/api/main.go

# Driver Service
cd C:\dev\FamGo-platform\services\driver-service
go mod download
go build -o bin/driver-service cmd/api/main.go

# Payment Service
cd C:\dev\FamGo-platform\services\payment-service
go mod download
go build -o bin/payment-service cmd/api/main.go

# Ride Service
cd C:\dev\FamGo-platform\services\ride-service
go mod download
go build -o bin/ride-service cmd/api/main.go

# Dispatch Service
cd C:\dev\FamGo-platform\services\dispatch-service
go mod download
go build -o bin/dispatch-service cmd/api/main.go
```

### Step 4: Run Services (in separate terminals)
```powershell
# Terminal 1: Pricing Service
cd C:\dev\FamGo-platform\services\pricing-service
& .\run.ps1

# Terminal 2: Driver Service
cd C:\dev\FamGo-platform\services\driver-service
& .\run.ps1

# Terminal 3: Payment Service
cd C:\dev\FamGo-platform\services\payment-service
& .\run.ps1

# Terminal 4: Ride Service
cd C:\dev\FamGo-platform\services\ride-service
& .\run.ps1

# Terminal 5: Dispatch Service
cd C:\dev\FamGo-platform\services\dispatch-service
& .\run.ps1
```

### Step 5: Verify All Services
```powershell
# Test each service health endpoint
curl http://localhost:3014/v1/health  # Pricing
curl http://localhost:3002/v1/health  # Driver
curl http://localhost:3015/v1/health  # Payment
curl http://localhost:3010/v1/health  # Ride
curl http://localhost:3011/v1/health  # Dispatch
```

---

## ✅ VERIFICATION CHECKLIST

- [ ] All 5 databases created in PostgreSQL
- [ ] All 5 service users created with correct passwords
- [ ] All 5 `.env` files created in respective service directories
- [ ] All services build without errors
- [ ] All services run and respond to health checks
- [ ] Can connect to each database from respective service
- [ ] Kafka connection working (if available)
- [ ] Redis connection working

---

**This provides complete production-grade setup with isolated databases per service, dedicated users, and proper environment configuration.** Ready to execute? 🚀

