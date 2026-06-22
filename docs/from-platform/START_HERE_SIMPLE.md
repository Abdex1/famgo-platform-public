# 🚀 START HERE - SUPER SIMPLE

## Step 1: Open PowerShell

Press `Win + R`, type `powershell`, press Enter

## Step 2: Run This Command

```powershell
cd C:\dev\FamGo-platform; .\start_all_services.ps1 -Environment local
```

## Step 3: Wait

Wait 30-60 seconds. 5 new windows will open (one per service).

## Step 4: Verify (Optional)

In a NEW PowerShell window:
```powershell
cd C:\dev\FamGo-platform
.\test_services.ps1
```

Should show all 5 services as "OK"

---

## 🆘 IF SCRIPT ERRORS

Use MANUAL startup instead:

Open 5 separate PowerShell windows and run ONE in each:

**Window 1:**
```powershell
cd C:\dev\FamGo-platform\services\pricing-service; .\start.ps1
```

**Window 2:**
```powershell
cd C:\dev\FamGo-platform\services\driver-service; .\start.ps1
```

**Window 3:**
```powershell
cd C:\dev\FamGo-platform\services\payment-service; .\start.ps1
```

**Window 4:**
```powershell
cd C:\dev\FamGo-platform\services\ride-service; .\start.ps1
```

**Window 5:**
```powershell
cd C:\dev\FamGo-platform\services\dispatch-service; .\start.ps1
```

---

## ✅ SUCCESS

When you see in each window:
```
[12:34:56] Starting service-name on port XXXX
[12:34:57] ===== SERVICE STARTED =====
[12:34:57] 🚀 Starting service-name
```

You're done! All services running.

---

## 🧪 TEST

New PowerShell window:
```powershell
Invoke-WebRequest http://localhost:3014/v1/health
Invoke-WebRequest http://localhost:3002/v1/health
Invoke-WebRequest http://localhost:3015/v1/health
Invoke-WebRequest http://localhost:3010/v1/health
Invoke-WebRequest http://localhost:3011/v1/health
```

All should return `"status":"healthy"`

---

**That's it! System is ready to use.** 🎉
