# 🚀 SIMPLE STEP-BY-STEP EXECUTION

## Step 1: Build All Services (5 minutes)

```powershell
cd C:\dev\FamGo-platform
.\build_all_services.ps1
```

Wait for all 5 to say "SUCCESS"

## Step 2: Start Services (5 separate windows)

Open 5 new PowerShell windows and run ONE command in each:

**Window 1 - Pricing Service (Port 3014)**
```powershell
cd C:\dev\FamGo-platform\services\pricing-service
.\start.ps1
```

**Window 2 - Driver Service (Port 3002)**
```powershell
cd C:\dev\FamGo-platform\services\driver-service
.\start.ps1
```

**Window 3 - Payment Service (Port 3015)**
```powershell
cd C:\dev\FamGo-platform\services\payment-service
.\start.ps1
```

**Window 4 - Ride Service (Port 3010)**
```powershell
cd C:\dev\FamGo-platform\services\ride-service
.\start.ps1
```

**Window 5 - Dispatch Service (Port 3011)**
```powershell
cd C:\dev\FamGo-platform\services\dispatch-service
.\start.ps1
```

## Step 3: Wait for Startup

Each window should show:
```
[HH:MM:SS] Starting service-name on port XXXX (local)
[HH:MM:SS] ===== SERVICE STARTED =====
[HH:MM:SS] 🚀 Starting service-name
```

Wait 30-60 seconds for all to start.

## Step 4: Verify Services (in NEW window)

```powershell
cd C:\dev\FamGo-platform
.\test_services.ps1
```

Should show:
```
Testing Pricing (Port 3014)... OK
Testing Driver (Port 3002)... OK
Testing Payment (Port 3015)... OK
Testing Ride (Port 3010)... OK
Testing Dispatch (Port 3011)... OK
```

## Step 5: Build Flutter Apps (when ready)

```powershell
# Driver App
cd C:\dev\FamGo-platform\mobile\flutter-driver-app
flutter pub get
flutter build apk --debug

# Passenger App
cd C:\dev\FamGo-platform\mobile\flutter-passenger-app
flutter pub get
flutter build apk --debug
```

---

## ✅ SUCCESS

When all services show OK and Flutter APKs are built, your system is complete and ready!

Each service runs independently in its own window. To stop any service, press Ctrl+C in that window.
