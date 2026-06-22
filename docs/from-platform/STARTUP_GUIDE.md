# 🚀 FAMGO PLATFORM - COMPLETE STARTUP GUIDE

**Last Updated**: 2024  
**Version**: 1.0.0  
**Status**: Production-Ready  

---

## 📋 QUICK START (2 MINUTES)

### LOCAL DEVELOPMENT
```powershell
# Option 1: Run all services at once
cd C:\dev\FamGo-platform
.\start_all_services.ps1 -Environment local
```

### PRODUCTION
```powershell
cd C:\dev\FamGo-platform
.\start_all_services.ps1 -Environment production
```

---

## 🔧 INDIVIDUAL SERVICE STARTUP

### Terminal 1: Pricing Service (Port 3014)
```powershell
cd C:\dev\FamGo-platform\services\pricing-service
.\start.ps1 -Environment local
```

### Terminal 2: Driver Service (Port 3002)
```powershell
cd C:\dev\FamGo-platform\services\driver-service
.\start.ps1 -Environment local
```

### Terminal 3: Payment Service (Port 3015)
```powershell
cd C:\dev\FamGo-platform\services\payment-service
.\start.ps1 -Environment local
```

### Terminal 4: Ride Service (Port 3010)
```powershell
cd C:\dev\FamGo-platform\services\ride-service
.\start.ps1 -Environment local
```

### Terminal 5: Dispatch Service (Port 3011)
```powershell
cd C:\dev\FamGo-platform\services\dispatch-service
.\start.ps1 -Environment local
```

---

## ✅ VERIFY ALL SERVICES ARE RUNNING

```powershell
# Test each service health endpoint
curl http://localhost:3014/v1/health  # Pricing Service
curl http://localhost:3002/v1/health  # Driver Service
curl http://localhost:3015/v1/health  # Payment Service
curl http://localhost:3010/v1/health  # Ride Service
curl http://localhost:3011/v1/health  # Dispatch Service
```

**Expected Response**:
```json
{
  "status": "healthy",
  "service": "service-name",
  "environment": "local",
  "version": "1.0.0",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## 📁 ENVIRONMENT CONFIGURATION FILES

### Structure
Each service has two environment files:

```
services/
├── pricing-service/
│   ├── .env.local           (Development)
│   ├── .env.production      (Production)
│   └── start.ps1            (Startup script)
├── driver-service/
│   ├── .env.local
│   ├── .env.production
│   └── start.ps1
├── payment-service/
│   ├── .env.local
│   ├── .env.production
│   └── start.ps1
├── ride-service/
│   ├── .env.local
│   ├── .env.production
│   └── start.ps1
└── dispatch-service/
    ├── .env.local
    ├── .env.production
    └── start.ps1
```

### Environment Differences

#### LOCAL (.env.local)
- **Database**: Shared `famgo_platform` (localhost)
- **User**: `postgres`
- **Logging**: DEBUG level
- **Redis**: No authentication
- **Features**: All disabled (metrics, tracing)

#### PRODUCTION (.env.production)
- **Database**: Dedicated per service (isolated security)
- **User**: Service-specific user (pricing_user, driver_user, etc.)
- **Logging**: INFO level
- **Redis**: Password protected
- **Features**: All enabled (metrics, tracing)
- **Kafka**: Production cluster
- **Security**: TLS enabled

---

## 🔐 PRODUCTION CREDENTIALS

### Database Users
```
Service          | User               | Password
-----------------+--------------------+--------------------------------
Pricing          | pricing_user       | Pr1c1ng@Secure2024!P@ssw0rd
Driver           | driver_user        | Driver@Secure2024!P@ssw0rd
Payment          | payment_user       | Payment@Secure2024!P@ssw0rd
Ride             | ride_user          | Ride@Secure2024!P@ssw0rd
Dispatch         | dispatch_user      | Dispatch@Secure2024!P@ssw0rd
```

### Shared Services
```
Redis Password   | RedisP@ss2024!Secure
Kafka Brokers    | kafka-prod-1:9092, kafka-prod-2:9092, kafka-prod-3:9092
Jaeger Host      | jaeger.production.internal:6831
```

---

## 🚨 TROUBLESHOOTING

### Issue: PowerShell `&&` operator not working
**Solution**: Use the provided `.ps1` scripts instead of inline commands
```powershell
# WRONG (causes error)
cd services/pricing-service && go build ... && run

# RIGHT (use provided scripts)
cd services/pricing-service
.\start.ps1
```

### Issue: Port already in use
**Solution**: Find and stop existing process
```powershell
# Find process using port 3014
netstat -ano | findstr :3014

# Kill process (replace PID with actual process ID)
taskkill /PID <PID> /F
```

### Issue: Database connection failed
**Solution**: Verify PostgreSQL is running
```powershell
# Check PostgreSQL status
& psql -U postgres -h localhost -c "SELECT version();"

# If failed, start PostgreSQL (Windows Service)
net start postgresql-x64-14
```

### Issue: Build fails with module errors
**Solution**: Clean and rebuild
```powershell
cd C:\dev\FamGo-platform\services\pricing-service
rm -Recurse -Force bin
rm go.sum
.\start.ps1
```

### Issue: Environment variables not loaded
**Solution**: Verify .env file exists and format is correct
```powershell
# Check file exists
ls -la C:\dev\FamGo-platform\services\pricing-service\.env.local

# Check format (should be KEY=VALUE)
cat C:\dev\FamGo-platform\services\pricing-service\.env.local
```

---

## 📊 SERVICE PORTS REFERENCE

| Service | Port | URL | Health |
|---------|------|-----|--------|
| Pricing | 3014 | http://localhost:3014 | /v1/health |
| Driver | 3002 | http://localhost:3002 | /v1/health |
| Payment | 3015 | http://localhost:3015 | /v1/health |
| Ride | 3010 | http://localhost:3010 | /v1/health |
| Dispatch | 3011 | http://localhost:3011 | /v1/health |

---

## 🧪 TEST API ENDPOINTS

### Pricing Service
```bash
# Estimate ride price
curl -X POST http://localhost:3014/v1/pricing/estimate \
  -d "ride_type=economy&distance_meters=5000&active_rides=10&available_drivers=15&is_pool=false"
```

### Driver Service
```bash
# Get driver details
curl http://localhost:3002/v1/drivers?id=driver-123

# Update driver location
curl -X POST http://localhost:3002/v1/drivers/location \
  -d "driver_id=driver-123&latitude=9.0320&longitude=38.7469&accuracy=10"

# Accept ride
curl -X POST http://localhost:3002/v1/drivers/accept-ride \
  -d "driver_id=driver-123&ride_id=ride-456"
```

### Payment Service
```bash
# Get wallet balance
curl http://localhost:3015/v1/wallets?user_id=user-123

# Process payment
curl -X POST http://localhost:3015/v1/payments/process \
  -d "ride_id=ride-456&user_id=user-123&amount=250.50&provider=telebirr"

# Get transaction history
curl http://localhost:3015/v1/transactions?user_id=user-123
```

### Ride Service
```bash
# Get ride details
curl http://localhost:3010/v1/rides?id=ride-456
```

### Dispatch Service
```bash
# Get matching metrics
curl http://localhost:3011/v1/dispatch/metrics
```

---

## 🎯 PRODUCTION DEPLOYMENT CHECKLIST

Before deploying to production:

- [ ] All `.env.production` files configured
- [ ] Database users created with secure passwords
- [ ] Redis cluster operational
- [ ] Kafka cluster running (3+ brokers)
- [ ] All services tested locally
- [ ] Health checks passing
- [ ] Logs monitored and configured
- [ ] Metrics collection enabled
- [ ] Tracing (Jaeger) configured
- [ ] TLS certificates installed
- [ ] Firewall rules configured
- [ ] Backup and recovery tested

---

## 📝 LOG LOCATIONS

### Service Logs (Local Development)
- Printed to console when running

### Service Logs (Production - Kubernetes)
```bash
kubectl logs <pod-name> -f
```

### Application Logs
All services output JSON-structured logs to stdout:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "service": "pricing-service",
  "message": "Starting service",
  "port": "3014"
}
```

---

## 🔄 GRACEFUL SHUTDOWN

Each service handles SIGINT and SIGTERM gracefully:

```powershell
# Press Ctrl+C in the service window
# Service will:
# 1. Stop accepting new connections
# 2. Wait for in-flight requests (30s timeout)
# 3. Close database connections
# 4. Exit cleanly
```

---

## 📈 MONITORING & HEALTH

### Service Health Endpoints
All services expose `/v1/health` endpoint that:
- Checks database connectivity
- Returns service version
- Returns environment info
- Returns timestamp

### Metrics Endpoints (Production)
Each service exposes Prometheus metrics:
- Pricing: `http://localhost:9001/metrics`
- Driver: `http://localhost:9002/metrics`
- Payment: `http://localhost:9003/metrics`
- Ride: `http://localhost:9004/metrics`
- Dispatch: `http://localhost:9005/metrics`

---

## 🆘 SUPPORT

**If services fail to start**:
1. Check all prerequisites are installed
2. Verify all `.env` files exist
3. Check PostgreSQL is running
4. Review error messages in console
5. Use troubleshooting section above

**For production issues**:
1. Check logs with `kubectl logs`
2. Verify database connectivity
3. Check Redis cluster status
4. Review Jaeger traces
5. Check Prometheus metrics

---

## ✨ SUCCESS INDICATORS

After startup, you should see:
- ✅ 5 service processes running
- ✅ All health endpoints responding with `"status": "healthy"`
- ✅ Database connections established
- ✅ No errors in logs
- ✅ All services listening on correct ports

---

**You're ready to go! Start services and begin development.** 🚀
