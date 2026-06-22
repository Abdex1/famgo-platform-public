# dev.ps1 — streamlined workflow for FamGo monorepo

param(
    [switch]$Sync,
    [switch]$Test,
    [switch]$Quick
)

Write-Host "🚀 FamGo Dev Workflow Starting..."

# Pre-cache dependencies once
function PreCache {
    Write-Host "📦 Pre-caching dependencies..."
    go mod download all
}

# Sync workspace only if requested
if ($Sync) {
    Write-Host "🔄 Syncing go.work..."
    go work sync
}

# Quick test run (no coverage, short mode)
if ($Quick) {
    Write-Host "⚡ Running quick tests..."
    go test ./internal/... -short
}

# Full test run with coverage
if ($Test) {
    Write-Host "🧪 Running full tests with coverage..."
    go test ./... -cover -p 8
}

Write-Host "✅ Workflow complete."
