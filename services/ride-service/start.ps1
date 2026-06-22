# RIDE SERVICE - STARTUP SCRIPT
param([string]$Environment = "local")

$ServiceName = "ride-service"
$ServicePort = "3010"
$ServicePath = "$PSScriptRoot"
$EnvFile = "$ServicePath\.env.$Environment"

if (-not (Test-Path $EnvFile)) {
    Write-Host "ERROR: Environment file not found: $EnvFile" -ForegroundColor Red
    exit 1
}

Write-Host "[$(Get-Date -Format 'HH:mm:ss')] Loading environment from: $EnvFile" -ForegroundColor Green

$envContent = Get-Content $EnvFile | Where-Object { $_ -notmatch "^\s*#" -and $_ -match "=" }
foreach ($line in $envContent) {
    $parts = $line -split "=", 2
    if ($parts.Count -eq 2) {
        $key = $parts[0].Trim()
        $value = $parts[1].Trim()
        [Environment]::SetEnvironmentVariable($key, $value, "Process")
    }
}

Write-Host "[$(Get-Date -Format 'HH:mm:ss')] Starting $ServiceName on port $ServicePort ($Environment)" -ForegroundColor Green

Push-Location $ServicePath

if (-not (Test-Path "bin\$ServiceName.exe")) {
    Write-Host "[$(Get-Date -Format 'HH:mm:ss')] Binary not found. Building..." -ForegroundColor Yellow
    
    if (-not (Test-Path "go.mod")) {
        & go mod init "github.com/FamGo/platform/services/$ServiceName" 2>&1 | Out-Null
    }
    
    & go mod download 2>&1 | Out-Null
    & go build -o "bin\$ServiceName.exe" "cmd\api\main.go" 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Build failed!" -ForegroundColor Red
        Pop-Location
        exit 1
    }
}

Write-Host "[$(Get-Date -Format 'HH:mm:ss')] ===== SERVICE STARTED =====" -ForegroundColor Green
Write-Host ""

& ".\bin\$ServiceName.exe"
