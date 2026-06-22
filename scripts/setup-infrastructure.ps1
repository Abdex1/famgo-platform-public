# scripts/setup-infrastructure.ps1
# ============================================================================
# FamGo Platform Infrastructure Setup Script (PowerShell)
# ============================================================================
# This script automates the setup and verification of the entire infrastructure
#
# Usage: .\scripts\setup-infrastructure.ps1 [command]
# Commands: setup, verify, start, stop, clean, logs
# ============================================================================

param (
    [string]$Command = "setup",
    [string]$Service = ""
)

# ============================================================================
# CONFIGURATION
# ============================================================================
$ScriptDir   = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = Resolve-Path "$ScriptDir\.."
$DockerComposeFile = "$ProjectRoot\infra\docker\docker-compose.yml"
$EnvFile     = "$ProjectRoot\.env.local"
$LogDir      = "$ProjectRoot\logs"
$Timestamp   = Get-Date -Format "yyyyMMdd_HHmmss"

# ============================================================================
# UTILITY FUNCTIONS
# ============================================================================

function Log-Info    { Write-Host "[INFO]    $args" -ForegroundColor Blue }
function Log-Success { Write-Host "[SUCCESS] $args" -ForegroundColor Green }
function Log-Warning { Write-Host "[WARNING] $args" -ForegroundColor Yellow }
function Log-Error   { Write-Host "[ERROR]   $args" -ForegroundColor Red }

function Check-Prerequisites {
    Log-Info "Checking prerequisites..."

    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Log-Error "Docker is not installed"
        exit 1
    }

    if (-not (Get-Command docker-compose -ErrorAction SilentlyContinue)) {
        Log-Error "Docker Compose is not installed"
        exit 1
    }

    if (-not (Test-Path $EnvFile)) {
        Log-Error ".env.local file not found"
        Log-Info "Please create .env.local from .env.example"
        exit 1
    }

    if (-not (Test-Path $DockerComposeFile)) {
        Log-Error "Docker Compose file not found: $DockerComposeFile"
        exit 1
    }

    Log-Success "Prerequisites check passed"
}

function Load-EnvFile {
    if (Test-Path $EnvFile) {
        Get-Content $EnvFile | ForEach-Object {
            if ($_ -match "^(.*?)=(.*)$") {
                [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2])
            }
        }
        Log-Success "Environment variables loaded from .env.local"
    } else {
        Log-Warning "No .env.local file found"
    }
}

function Create-LogDirectory {
    if (-not (Test-Path $LogDir)) {
        New-Item -ItemType Directory -Force -Path $LogDir | Out-Null
        Log-Success "Created logs directory: $LogDir"
    }
}

# ============================================================================
# SETUP FUNCTION
# ============================================================================

function Setup {
    Log-Info "Starting infrastructure setup..."
    Check-Prerequisites
    Load-EnvFile
    Create-LogDirectory

    Log-Info "Building Docker images..."
    docker-compose -f $DockerComposeFile --env-file $EnvFile build

    Log-Info "Creating volumes and networks..."
    docker-compose -f $DockerComposeFile --env-file $EnvFile create

    Log-Success "Infrastructure setup completed!"
}

# ============================================================================
# START FUNCTION
# ============================================================================

function Start-Infrastructure {
    Log-Info "Starting containers..."
    docker-compose -f $DockerComposeFile --env-file $EnvFile up -d

    Log-Info "Waiting for services to start (30 seconds)..."
    Start-Sleep -Seconds 30

    Log-Success "Containers started!"
    Verify
}

# ============================================================================
# STOP FUNCTION
# ============================================================================

function Stop-Infrastructure {
    Log-Info "Stopping containers..."
    docker-compose -f $DockerComposeFile --env-file $EnvFile down
    Log-Success "Containers stopped!"
}

# ============================================================================
# VERIFY FUNCTION
# ============================================================================

function Verify {
    Log-Info "Verifying infrastructure..."
    $failed = 0

    # PostgreSQL
    Log-Info "Checking PostgreSQL..."
    if (docker-compose -f $DockerComposeFile --env-file $EnvFile exec -T postgres pg_isready -U $env:POSTGRES_USER) {
        Log-Success "PostgreSQL: OK"
    } else { Log-Error "PostgreSQL: FAILED"; $failed++ }

    # Redis
    Log-Info "Checking Redis..."
    if ((docker-compose -f $DockerComposeFile --env-file $EnvFile exec -T redis redis-cli ping) -match "PONG") {
        Log-Success "Redis: OK"
    } else { Log-Error "Redis: FAILED"; $failed++ }

    # Kafka
    Log-Info "Checking Kafka..."
    if (docker ps | Select-String "famgo-kafka") {
        Log-Success "Kafka: OK"
    } else { Log-Error "Kafka: FAILED"; $failed++ }

    # MinIO
    Log-Info "Checking MinIO..."
    try { Invoke-WebRequest "http://localhost:9000/minio/health/live" -UseBasicParsing | Out-Null; Log-Success "MinIO: OK" }
    catch { Log-Error "MinIO: FAILED"; $failed++ }

    # ClickHouse
    Log-Info "Checking ClickHouse..."
    try { Invoke-WebRequest "http://localhost:8123/ping" -UseBasicParsing | Out-Null; Log-Success "ClickHouse: OK" }
    catch { Log-Error "ClickHouse: FAILED"; $failed++ }

    # Prometheus
    Log-Info "Checking Prometheus..."
    try { Invoke-WebRequest "http://localhost:9090/-/healthy" -UseBasicParsing | Out-Null; Log-Success "Prometheus: OK" }
    catch { Log-Error "Prometheus: FAILED"; $failed++ }

    # Grafana
    Log-Info "Checking Grafana..."
    try { Invoke-WebRequest "http://localhost:3001/api/health" -UseBasicParsing | Out-Null; Log-Success "Grafana: OK" }
    catch { Log-Error "Grafana: FAILED"; $failed++ }

    # Loki
    Log-Info "Checking Loki..."
    try { Invoke-WebRequest "http://localhost:3100/ready" -UseBasicParsing | Out-Null; Log-Success "Loki: OK" }
    catch { Log-Error "Loki: FAILED"; $failed++ }

    # Jaeger
    Log-Info "Checking Jaeger..."
    try { Invoke-WebRequest "http://localhost:16686/api/services" -UseBasicParsing | Out-Null; Log-Success "Jaeger: OK" }
    catch { Log-Error "Jaeger: FAILED"; $failed++ }

    if ($failed -eq 0) {
        Log-Success "All services verified successfully!"
        Write-Host "`nAccess URLs:"
        Write-Host "  PostgreSQL:    localhost:5432"
        Write-Host "  Redis:         localhost:6379"
        Write-Host "  Kafka:         localhost:9092"
        Write-Host "  MinIO:         http://localhost:9000 (Console: http://localhost:9001)"
        Write-Host "  ClickHouse:    http://localhost:8123"
        Write-Host "  Prometheus:    http://localhost:9090"
        Write-Host "  Grafana:       http://localhost:3001 (admin/admin_dev_password)"
        Write-Host "  Loki:          http://localhost:3100"
        Write-Host "  Jaeger:        http://localhost:16686"
        Write-Host "  Nginx:         http://localhost:80"
    } else {
        Log-Error "$failed service(s) failed verification"
    }
}

# ============================================================================
# CLEAN FUNCTION
# ============================================================================

function Clean {
    Log-Warning "This will remove all containers, volumes, and data!"
    $confirm = Read-Host "Are you sure? (yes/no)"
    if ($confirm -match "^(yes|y)$") {
        Log-Info "Cleaning up infrastructure..."
        docker-compose -f $DockerComposeFile --env-file $EnvFile down -v
        Log-Success "Infrastructure cleaned!"
    } else {
        Log-Info "Cleanup cancelled"
    }
}

# ============================================================================
# LOGS FUNCTION
# ============================================================================

function Logs {
    if ([string]::IsNullOrEmpty($Service)) {
        docker-compose -f $DockerComposeFile --env-file $EnvFile logs -f
    } else {
        docker-compose -f $DockerComposeFile --env-file $EnvFile logs -f $Service
    }
}

# ============================================================================
# MAIN
# ============================================================================

switch ($Command) {
    "setup"  { Setup }
    "start"  { Start-Infrastructure }
    "stop"   { Stop-Infrastructure }
    "verify" { Verify }
    "clean"  { Clean }
    "logs"   { Logs }
    default  { Write-Host "Usage: .\setup-infrastructure.ps1 {setup|start|stop|verify|clean|logs [service]}"; exit 1 }
}
