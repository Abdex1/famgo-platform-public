cd C:\dev\FamGo-platform

After updating all apps, run:

pnpm turbo run build

Then test locally:

pnpm turbo run dev



docker compose --env-file .\.env up -d






That `github.com/your-org/...` is only a placeholder import path.

Since you are currently local on Windows at:

```text
C:\dev\FamGo-platform
```
pnpm --filter @famgo/analytics-dashboard build

pnpm --filter admin-dashboard build
pnpm --filter support-dashboard build
pnpm --filter operator-dashboard build
pnpm --filter rider-web build
pnpm --filter driver-web build





you have 3 production-grade choices.

For enterprise repos, the best practice is to use the future Git repository path NOW, even before pushing.

If your future GitHub username is:

```text
abdisha
```

then use:

```bash
go mod init github.com/abdisha/famgo-auth-service
```

OR, for monorepo-style naming (better for your architecture):

```bash
go mod init github.com/abdisha/FamGo-platform/services/auth-service
```

I strongly recommend the SECOND one because your repo is a monorepo.

So inside:

```text
C:\dev\FamGo-platform\services\auth-service
```

run:

```bash
go mod init github.com/abdisha/FamGo-platform/services/auth-service
```

Later when you create the GitHub repository, keep the SAME name:

```text
FamGo-platform
```

under your GitHub account:

```text
abdisha
```

Then the import path stays valid.

Enterprise recommendation for ALL Go services:

```text
github.com/abdisha/FamGo-platform/services/auth-service
github.com/abdisha/FamGo-platform/services/dispatch-service
github.com/abdisha/FamGo-platform/services/gps-service
github.com/abdisha/FamGo-platform/services/payment-service
```

etc.

This gives:

* clean internal imports
* proper Go module resolution
* future CI/CD compatibility
* Kubernetes build consistency
* proper dependency tracing
* monorepo-safe architecture

Example internal imports later:

```go
import "github.com/abdisha/FamGo-platform/shared/contracts/grpc"
```

or

```go
import "github.com/abdisha/FamGo-platform/shared/utilities/logger"
```

Do NOT use:

* local Windows paths
* `C:\dev\...`
* fake placeholders
* temporary names

because Go modules are designed around repository import paths, not filesystem paths.



to list a dir and files"
powershell
Get-ChildItem -Recurse |
Where-Object {
    $_.FullName -notmatch 'node_modules|\.next|dist|build|coverage|\.turbo|\.git|\.pnpm-store'
} |
Select-Object FullName


Better enterprise version (much cleaner):

Get-ChildItem -Recurse |
Where-Object {
    $_.FullName -notmatch 'node_modules|\.next|dist|build|coverage|\.turbo|\.git|\.pnpm-store'
} |
Select-Object FullName

Best option for exporting to a file:

tree /F /A `
| findstr /V "node_modules .next dist build coverage .turbo .git" `
> project-structure.txt

Then open:

notepad project-structure.txt

If you want the MOST enterprise-readable architecture output, use this instead:

Get-ChildItem -Recurse |
Where-Object {
    $_.FullName -notmatch 'node_modules|\.next|dist|build|coverage|\.turbo|\.git|\.pnpm-store'
} |
ForEach-Object {
    $_.FullName.Replace("C:\dev\FamGo-platform\", "")
} |
Out-File project-structure-latest.txt

Then:

notepad project-structure-clean.txt

This produces a much cleaner monorepo inventory for:

architecture reviews
CI/CD audits
documentation
dependency planning
production readiness reviews
onboarding documentation

You can also generate ONLY directories:

Get-ChildItem -Recurse -Directory |
Where-Object {
    $_.FullName -notmatch 'node_modules|\.next|dist|build|coverage|\.turbo|\.git'
} |
ForEach-Object {
    $_.FullName.Replace("C:\dev\FamGo-platform\", "")
}

Or ONLY files:

Get-ChildItem -Recurse -File |
Where-Object {
    $_.FullName -notmatch 'node_modules|\.next|dist|build|coverage|\.turbo|\.git'
} |
ForEach-Object {
    $_.FullName.Replace("C:\dev\FamGo-platform\", "")
}