# MANUAL STARTUP - NO SCRIPTS REQUIRED

## Open 5 PowerShell Windows (ONE FOR EACH SERVICE)

### Window 1: PRICING SERVICE (Port 3014)
```powershell
cd C:\dev\FamGo-platform\services\pricing-service
.\start.ps1 -Environment local
```

### Window 2: DRIVER SERVICE (Port 3002)
```powershell
cd C:\dev\FamGo-platform\services\driver-service
.\start.ps1 -Environment local
```

### Window 3: PAYMENT SERVICE (Port 3015)
```powershell
cd C:\dev\FamGo-platform\services\payment-service
.\start.ps1 -Environment local
```

### Window 4: RIDE SERVICE (Port 3010)
```powershell
cd C:\dev\FamGo-platform\services\ride-service
.\start.ps1 -Environment local
```

### Window 5: DISPATCH SERVICE (Port 3011)
```powershell
cd C:\dev\FamGo-platform\services\dispatch-service
.\start.ps1 -Environment local
```

## Verify Services (in NEW window)
```powershell
cd C:\dev\FamGo-platform
.\test_services.ps1
```

## Expected Output
```
Testing Pricing (Port 3014)... OK
Testing Driver (Port 3002)... OK
Testing Payment (Port 3015)... OK
Testing Ride (Port 3010)... OK
Testing Dispatch (Port 3011)... OK
```

## If Any Shows FAILED
1. Wait 5 more seconds
2. Check the service window for error messages
3. Look at the error logs in that window
