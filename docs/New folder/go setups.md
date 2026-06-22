Create the Correct go.work

After fixing old module paths:

From root:

cd C:\dev\FamGo-consolidated

go work init

Add services:

go work use ./services/auth-service
go work use ./services/dispatch-service
go work use ./services/driver-service
go work use ./services/fraud-service
go work use ./services/gps-service
go work use ./services/notification-service
go work use ./services/payment-service
go work use ./services/pooling-service
go work use ./services/pricing-service
go work use ./services/ride-service
go work use ./services/safety-service
go work use ./services/user-service
go work use ./services/wallet-service

Add shared packages:

go work use ./packages/event-bus
go work use ./packages/kafka-sdk
go work use ./packages/redis-platform
go work use ./packages/telemetry
go work use ./packages/telemetry/go
go work use ./packages/vault-sdk

Then:

go work sync
Final Cleanup

After all replacements:

go clean -cache
go clean -modcache
go clean -testcache

Delete stale VCS cache:

Remove-Item `
"$env:USERPROFILE\go\pkg\mod\cache\vcs" `
-Recurse `
-Force