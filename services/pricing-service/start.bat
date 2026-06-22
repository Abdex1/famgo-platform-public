@echo off
REM Pricing Service Launcher
REM Sets environment variables and starts the service

echo.
echo ============================================
echo  FamGo Pricing Service Startup
echo ============================================
echo.

REM Set environment variables
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=famgo_user
set DB_PASSWORD=famgo_secure
set DB_NAME=famgo_platform
set SERVICE_PORT=3014

echo [CONFIG] Loading environment variables...
echo   - DB_HOST: %DB_HOST%
echo   - DB_PORT: %DB_PORT%
echo   - DB_USER: %DB_USER%
echo   - DB_NAME: %DB_NAME%
echo   - SERVICE_PORT: %SERVICE_PORT%
echo.

REM Check if binary exists
if not exist "bin\pricing-service.exe" (
    echo [ERROR] Binary not found! Building...
    go build -o bin/pricing-service cmd/api/main.go
    if %ERRORLEVEL% neq 0 (
        echo [ERROR] Build failed!
        exit /b 1
    )
    echo [SUCCESS] Build completed
    echo.
)

echo [INFO] Starting Pricing Service...
echo.
bin\pricing-service.exe

if %ERRORLEVEL% neq 0 (
    echo [ERROR] Service failed to start
    exit /b 1
)
