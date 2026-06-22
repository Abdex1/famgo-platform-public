# ⚡ QUICK START REFERENCE CARD

## 🚀 START ALL SERVICES (30 SECONDS)

### Development
```powershell
cd C:\dev\FamGo-platform
.\start_all_services.ps1 -Environment local
```

### Production
```powershell
cd C:\dev\FamGo-platform
.\start_all_services.ps1 -Environment production
```

---

## 📊 VERIFY SERVICES (COPY & PASTE)

```powershell
# All in one PowerShell command
$ports = @(3014, 3002, 3015, 3010, 3011)
$services = @('pricing', 'driver', 'payment', 'ride', 'dispatch')
for ($i = 0; $i -lt 5; $i++) {
    $response = curl -s "http://localhost:$($ports[$i])/v1/health"
    Write-Host "$($services[$i]): $response"
}
```

---

## 🔧 INDIVIDUAL SERVICE STARTUP

```powershell
# Pricing (3014)
cd C:\dev\FamGo-platform\services\pricing-service
.\start.ps1 -Environment local

# Driver (3002)
cd C:\dev\FamGo-platform\services\driver-service
.\start.ps1 -Environment local

# Payment (3015)
cd C:\dev\FamGo-platform\services\payment-service
.\start.ps1 -Environment local

# Ride (3010)
cd C:\dev\FamGo-platform\services\ride-service
.\start.ps1 -Environment local

# Dispatch (3011)
cd C:\dev\FamGo-platform\services\dispatch-service
.\start.ps1 -Environment local
```

---

## 🐛 TROUBLESHOOTING

### Port Already in Use
```powershell
# Find process
netstat -ano | findstr :3014

# Kill it
taskkill /PID <PID> /F
```

### Database Connection Failed
```powershell
# Test PostgreSQL
psql -U postgres -h localhost -c "SELECT 1"

# Start PostgreSQL (Windows)
net start postgresql-x64-14
```

### Build Failed
```powershell
cd C:\dev\FamGo-platform\services\pricing-service
rm -Recurse -Force bin
rm go.sum
.\start.ps1
```

---

## 📝 ENVIRONMENT FILES

| Service | Local | Production |
|---------|-------|------------|
| Pricing | `.env.local` | `.env.production` |
| Driver | `.env.local` | `.env.production` |
| Payment | `.env.local` | `.env.production` |
| Ride | `.env.local` | `.env.production` |
| Dispatch | `.env.local` | `.env.production` |

---

## 🌐 API HEALTH CHECKS

```bash
curl http://localhost:3014/v1/health  # Pricing
curl http://localhost:3002/v1/health  # Driver
curl http://localhost:3015/v1/health  # Payment
curl http://localhost:3010/v1/health  # Ride
curl http://localhost:3011/v1/health  # Dispatch
```

---

## 🔐 PRODUCTION CREDENTIALS

```
Pricing:  pricing_user    / Pr1c1ng@Secure2024!P@ssw0rd
Driver:   driver_user     / Driver@Secure2024!P@ssw0rd
Payment:  payment_user    / Payment@Secure2024!P@ssw0rd
Ride:     ride_user       / Ride@Secure2024!P@ssw0rd
Dispatch: dispatch_user   / Dispatch@Secure2024!P@ssw0rd
Redis:    password        / RedisP@ss2024!Secure
```

---

## 📊 SERVICE PORTS

| Service | Port |
|---------|------|
| Pricing | 3014 |
| Driver | 3002 |
| Payment | 3015 |
| Ride | 3010 |
| Dispatch | 3011 |

---

## 📁 KEY FILES

```
C:\dev\FamGo-platform\
├── start_all_services.ps1       ← Master startup (PowerShell)
├── start_all_services.bat       ← Batch alternative
├── STARTUP_GUIDE.md             ← Complete guide
└── services\
    ├── pricing-service\
    │   ├── .env.local
    │   ├── .env.production
    │   └── start.ps1
    ├── driver-service\
    ├── payment-service\
    ├── ride-service\
    └── dispatch-service\
```

---

## ✅ SUCCESS INDICATORS

After running `start_all_services.ps1`:

- [ ] 5 windows opened (one per service)
- [ ] No error messages
- [ ] Each shows "🚀 Starting service_name on port XXXX"
- [ ] Health checks return `"status":"healthy"`

---

## 📞 HELP

| Issue | Solution |
|-------|----------|
| && operator error | Use `.ps1` scripts (fixed!) |
| Port in use | `taskkill /PID <PID> /F` |
| DB connection failed | Check PostgreSQL running |
| Service won't build | Delete `bin/` and `go.sum` |
| Health check fails | Wait 5 seconds, service initializing |

---

**That's it! You're ready to go.** 🚀

Full documentation: `STARTUP_GUIDE.md`  
Configuration details: `CONFIGURATION_DELIVERY.md`
