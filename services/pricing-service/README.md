# Pricing Service - FamGo Platform

Microservice for calculating ride fares, managing pricing rules, and handling surge pricing for the FamGo transportation platform.

## Features

- ✅ **Complete Fare Calculation**: Base fare + distance + time + surge + tax - discount
- ✅ **Dynamic Surge Pricing**: Time-based (peak hours) + supply-demand based
- ✅ **Discount Code Management**: Fixed/percentage discounts with validation
- ✅ **Multiple Ride Types**: ECONOMY, COMFORT, BUSINESS, POOL
- ✅ **Pool Discounts**: Automatic discounts for shared rides
- ✅ **Pricing Rules**: City-specific, time-based pricing configuration
- ✅ **Analytics**: Surge history and fare statistics
- ✅ **Minimal Fare Enforcement**: Prevents unrealistic low fares

## Architecture

```
pricing-service/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/
│   │   ├── entities/            # Domain models & interfaces
│   │   └── services/            # Business logic (pricing engine)
│   ├── infrastructure/
│   │   └── postgres/            # Database layer
│   └── interfaces/
│       └── rest/                # HTTP handlers
├── bin/
│   └── pricing-service          # Compiled binary
└── start.bat / start.ps1        # Startup scripts
```

## Quick Start

### Prerequisites
- Go 1.20+
- PostgreSQL 12+
- Environment variables set (see below)

### Installation

1. **Clone and navigate**:
```bash
cd services/pricing-service
```

2. **Install dependencies**:
```bash
go mod download
```

3. **Build**:
```bash
go build -o bin/pricing-service cmd/api/main.go
```

4. **Set environment variables** (Windows CMD):
```batch
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=famgo_user
set DB_PASSWORD=famgo_secure
set DB_NAME=famgo_platform
set SERVICE_PORT=3014
```

5. **Run using startup script**:
```batch
# Batch script
.\start.bat

# Or PowerShell
.\start.ps1

# Or directly
.\bin\pricing-service.exe
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | famgo_user | Database user |
| `DB_PASSWORD` | famgo_secure | Database password |
| `DB_NAME` | famgo_platform | Database name |
| `SERVICE_PORT` | 3014 | Service HTTP port |

## API Endpoints

### Health Check
```http
GET /v1/health
```
Returns service status.

### Calculate Fare (Complete)
```http
POST /v1/pricing/calculate
Content-Type: application/json

{
  "ride_id": "ride_123",
  "ride_type": "ECONOMY",
  "distance_meters": 5000,
  "duration_seconds": 900,
  "pickup_lat": 9.0,
  "pickup_lng": 38.7,
  "dropoff_lat": 9.5,
  "dropoff_lng": 38.9,
  "is_pool": false,
  "active_rides": 50,
  "available_drivers": 20,
  "discount_code": "FAMGO10"
}
```

**Response**:
```json
{
  "calculation_id": "calc_abc123",
  "ride_id": "ride_123",
  "base_fare": 20.0,
  "distance_fare": 50.0,
  "time_fare": 4.95,
  "subtotal": 74.95,
  "surge_multiplier": 1.25,
  "surge_amount": 18.74,
  "taxes": 1.87,
  "discount_amount": 9.45,
  "final_fare": 85.11,
  "is_pool": false,
  "calculated_at": "2024-01-20T10:30:00Z"
}
```

### Estimate Fare (Quick)
```http
POST /v1/pricing/estimate
Content-Type: application/json

{
  "ride_type": "ECONOMY",
  "distance_meters": 5000,
  "active_rides": 50,
  "available_drivers": 20,
  "is_pool": false
}
```

**Response**:
```json
{
  "base_fare": 20.0,
  "distance_fare": 50.0,
  "subtotal": 70.0,
  "surge_multiplier": 1.25,
  "surge_amount": 17.5,
  "taxes": 1.75,
  "final_fare": 89.25
}
```

### Get Surge Multiplier
```http
POST /v1/pricing/surge
Content-Type: application/json

{
  "latitude": 9.0,
  "longitude": 38.7,
  "active_rides": 50,
  "available_drivers": 20
}
```

**Response**:
```json
{
  "surge_multiplier": 1.25,
  "timestamp": "2024-01-20T10:30:00Z",
  "active_rides": 50,
  "available_drivers": 20
}
```

### Apply Discount Code
```http
POST /v1/pricing/apply-discount
Content-Type: application/json

{
  "discount_code": "FAMGO10",
  "fare_amount": 100.0
}
```

**Response**:
```json
{
  "discount_code": "FAMGO10",
  "original_fare": 100.0,
  "discount_amount": 10.0,
  "final_fare": 90.0,
  "discount_type": "PERCENTAGE"
}
```

### Get Pricing Statistics
```http
GET /v1/pricing/statistics?city=Addis+Ababa
```

**Response**:
```json
{
  "city": "Addis Ababa",
  "average_fares": {
    "ECONOMY": 85.50,
    "COMFORT": 125.30,
    "BUSINESS": 165.75,
    "POOL": 65.20
  },
  "surge_history": [...],
  "period_days": 7,
  "timestamp": "2024-01-20T10:30:00Z"
}
```

## Fare Calculation Formula

```
Final Fare = (BaseFare + DistanceFare + TimeFare) × SurgeMultiplier + Taxes - Discount

Where:
  BaseFare          = Ride type minimum charge
  DistanceFare      = Distance(km) × RatePerKm
  TimeFare          = Duration(min) × RatePerMinute
  SurgeMultiplier   = Time-based(40%) + Supply-Demand(60%)
  Taxes             = (Subtotal + Surge) × TaxPercentage
  Discount          = Discount code value
  MinimumFare       = Enforced floor price
```

## Pricing Rules

### Default Pricing (ETB)

| Ride Type | Base Fare | Distance/km | Time/min | Min Fare | Max Surge | Pool Discount |
|-----------|-----------|-------------|----------|----------|-----------|---------------|
| ECONOMY | 20.0 | 10.0 | 0.33 | 15.0 | 5.0x | 25% |
| COMFORT | 30.0 | 13.0 | 0.43 | 25.0 | 5.0x | 20% |
| BUSINESS | 40.0 | 18.0 | 0.60 | 35.0 | 5.0x | 15% |
| POOL | 15.0 | 8.0 | 0.25 | 10.0 | 5.0x | 0% |

### Surge Multiplier Calculation

**Time-based**:
- Peak hours (6-9 AM, 5-8 PM): 1.5x
- Off-peak: 1.0x

**Supply-Demand**:
- Ratio = ActiveRides / AvailableDrivers
- Multiplier = 1.0 + (Ratio - 1.0) × 0.5

**Combined** (clamped 1.0-5.0):
- Final = (Time × 0.4) + (Supply-Demand × 0.6)

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/domain/services -v

# Benchmark tests
go test -bench=. -benchmem ./internal/domain/services
```

## Troubleshooting

### Service won't start

**Error**: "Failed to connect to database"

**Solution**:
1. Verify PostgreSQL is running
2. Check credentials: `DB_USER`, `DB_PASSWORD`, `DB_NAME`
3. Verify host and port: `DB_HOST`, `DB_PORT`
4. Test connection: `psql -h localhost -U famgo_user -d famgo_platform`

### Port already in use

**Error**: "listen tcp :3014: bind: An attempt was made to use a socket address that belongs to a process..."

**Solution**:
```bash
# Find process using port
netstat -ano | findstr :3014

# Kill process (replace PID)
taskkill /PID <PID> /F

# Or use different port
set SERVICE_PORT=3015
```

### Build issues

**Error**: "command not found: go"

**Solution**: Ensure Go is installed and in PATH
```bash
go version
```

**Error**: "package not found"

**Solution**: Download dependencies
```bash
go mod download
go mod tidy
```

## Development

### Project Structure

- **`entities`**: Domain models (PricingRule, FareCalculation, DiscountCode)
- **`services`**: Business logic (PricingEngine with surge calculation)
- **`postgres`**: Database operations (PricingRuleRepository)
- **`rest`**: HTTP handlers and routes

### Adding New Features

1. Define entity in `internal/domain/entities/`
2. Add repository methods to interface in `entities/repository.go`
3. Implement in `internal/infrastructure/postgres/`
4. Add business logic in `internal/domain/services/`
5. Create HTTP handler in `internal/interfaces/rest/`
6. Register route in `RegisterRoutes()`

### Code Style

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Vet code
go vet ./...
```

## Performance

**Benchmarks** (on typical hardware):
- Fare Calculation: ~10μs
- Surge Calculation: ~2μs
- Estimate: ~8μs

## Deployment

### Docker

```dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/pricing-service cmd/api/main.go

FROM alpine:latest
COPY --from=builder /app/bin/pricing-service /usr/local/bin/
EXPOSE 3014
CMD ["pricing-service"]
```

### Environment Setup (Production)

```bash
export DB_HOST=prod-postgres.internal
export DB_PORT=5432
export DB_USER=prod_famgo_user
export DB_PASSWORD=$(aws secretsmanager get-secret-value ...)
export DB_NAME=famgo_prod
export SERVICE_PORT=3014
```

## Monitoring

Health check endpoint for Kubernetes/monitoring:
```bash
curl http://localhost:3014/v1/health
```

Response indicates service is running and connected to database.

## Contributing

1. Create feature branch: `git checkout -b feature/new-pricing-logic`
2. Write tests first
3. Implement feature
4. Run tests: `go test ./...`
5. Format: `go fmt ./...`
6. Submit PR

## License

Proprietary - FamGo Platform
