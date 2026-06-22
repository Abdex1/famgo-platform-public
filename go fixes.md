✅ How to Fix (Monorepo Submodule Approach)
Make sure your changes are committed and pushed:

powershell
git add packages/vault-sdk
git commit -m "Add vault-sdk module"
git push origin master
Create a submodule tag:

powershell
git tag packages/vault-sdk/v0.0.1
git push origin packages/vault-sdk/v0.0.1
⚠️ Notice the tag format: it must include the subfolder path (packages/vault-sdk/) followed by the version.

In auth-service/go.mod, require:

go
require github.com/Abdex1/FamGo-platform/packages/vault-sdk v0.0.1
Run:

powershell
go clean -modcache
go mod tidy






Your output reveals a few important things.

## Good news

This command:

```powershell
Get-ChildItem -Recurse -Filter *.go |
Select-String '"import "github.com/jinzhu/inflect"'
```

returned nothing.

That means:

✅ No Go source files import `sdk/resource`

✅ The bad dependency line was the only `sdk/resource` issue

✅ We do NOT need to modify imports

---

# Issue 1: Massive OTel Version Drift

You currently have:

| Version | Locations            |
| ------- | -------------------- |
| 1.20.0  | templates            |
| 1.21.0  | many services        |
| 1.35.0  | telemetry package    |
| 1.43.0  | event-bus, kafka-sdk |
| 1.44.0  | auth-service         |

This is exactly the kind of thing that breaks `go work sync`.

For a monorepo, all internal packages should generally align.

---

# Issue 2: Old Jaeger Packages

I see:

```go
go.opentelemetry.io/otel/exporters/jaeger
```

and

```go
go.opentelemetry.io/otel/exporters/jaeger/grpc
```

These are old APIs.

Recent OTel uses:

```go
go.opentelemetry.io/otel/exporters/otlp/otlptrace
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
```

Do not replace these yet until workspace is healthy.

---

# Safe Step 1 — Backup

From repo root:

```powershell
git add .
git commit -m "backup before otel standardization"
```

---

# Safe Step 2 — Standardize ALL OTel Core Versions

This updates:

```text
otel
otel/metric
otel/sdk
otel/sdk/metric
otel/sdk/trace
otel/trace
```

to 1.44.0.

Run from:

```powershell
cd C:\dev\FamGo-consolidated
```

```powershell
Get-ChildItem -Recurse -Filter go.mod | ForEach-Object {
    $content = Get-Content $_.FullName -Raw

    $content = $content -replace 'go\.opentelemetry\.io/otel v\d+\.\d+\.\d+','go.opentelemetry.io/otel v1.44.0'
    $content = $content -replace 'go\.opentelemetry\.io/otel/metric v\d+\.\d+\.\d+','go.opentelemetry.io/otel/metric v1.44.0'
    $content = $content -replace 'go\.opentelemetry\.io/otel/sdk v\d+\.\d+\.\d+','go.opentelemetry.io/otel/sdk v1.44.0'
    $content = $content -replace 'go\.opentelemetry\.io/otel/sdk/metric v\d+\.\d+\.\d+','go.opentelemetry.io/otel/sdk/metric v1.44.0'
    $content = $content -replace 'go\.opentelemetry\.io/otel/sdk/trace v\d+\.\d+\.\d+','go.opentelemetry.io/otel/sdk/trace v1.44.0'
    $content = $content -replace 'go\.opentelemetry\.io/otel/trace v\d+\.\d+\.\d+','go.opentelemetry.io/otel/trace v1.44.0'

    Set-Content $_.FullName $content
}
```

---

# Safe Step 3 — Verify

```powershell
Get-ChildItem -Recurse -Filter go.mod |
Select-String "go.opentelemetry.io/otel"
```

You should now see mostly:

```text
v1.44.0
```

---

# Safe Step 4 — Check Remaining Old Exporters

Run:

```powershell
Get-ChildItem -Recurse -Filter go.mod |
Select-String "exporters/jaeger"
```

and

```powershell
Get-ChildItem -Recurse -Filter *.go |
Select-String "exporters/jaeger"
```

Paste the output before changing anything.

---

# Safe Step 5 — Refresh Dependencies

After version alignment:

```powershell
go clean -cache
go clean -modcache
go clean -testcache
```

Then:

```powershell
Remove-Item "$env:USERPROFILE\go\pkg\mod\cache\vcs" `
-Recurse `
-Force `
-ErrorAction SilentlyContinue
```

---

# Safe Step 6 — Repair Every Module

From repo root:

```powershell
Get-ChildItem .\services -Directory |
Where-Object { Test-Path "$($_.FullName)\go.mod" } |
ForEach-Object {
    Write-Host "Fixing $($_.Name)"
    Push-Location $_.FullName
    go mod tidy
    Pop-Location
}
```

Then:

```powershell
Get-ChildItem .\packages -Directory |
Where-Object { Test-Path "$($_.FullName)\go.mod" } |
ForEach-Object {
    Write-Host "Fixing package $($_.Name)"
    Push-Location $_.FullName
    go mod tidy
    Pop-Location
}
```

---

# Safe Step 7 — Sync Workspace

Finally:

```powershell
go work sync
```

---

### One suspicious thing

Your output contains both:

```text
packages\telemetry\go.mod
```

and

```text
packages\telemetry\go\go.mod
```

That means you may have **nested Go modules**:

```text
packages/
 └ telemetry/
     go.mod
     go/
        go.mod
```

Nested modules frequently cause workspace and dependency problems.

Check:

```powershell
Get-ChildItem C:\dev\FamGo-consolidated\packages\telemetry -Recurse -Filter go.mod
```

If both are real modules, show me that output next. That may be the next blocker after the OTel version mismatch.














Required go.mod additions

After fixing telemetry.go:

cd C:\dev\FamGo-consolidated\services\auth-service

go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@v1.44.0

go get go.opentelemetry.io/otel/sdk@v1.44.0

go get go.opentelemetry.io/otel@v1.44.0

go mod tidy
Before running go work sync again

Run these checks:

Get-ChildItem -Recurse -Include *.go,go.mod |
Select-String "otel/exporters/jaeger"

Get-ChildItem -Recurse -Include go.mod,go.work |
Select-String "github.com/Abdex1/FamGo-platform v0.1.0"

Get-ChildItem -Recurse -Include go.mod,go.work |
Select-String "replace"

The current blocker for workspace synchronization is no longer Jaeger; it is the remaining reference to:

github.com/Abdex1/FamGo-platform v0.1.0

somewhere in the repository. Once that reference is found and corrected, go work sync should move much further.