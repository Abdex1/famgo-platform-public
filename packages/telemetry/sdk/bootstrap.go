
// STEP 16 — SDK BOOTSTRAP

 
// packages/telemetry/go/sdk/bootstrap.go
  
package sdk

import (
    "context"

    "github.com/Abdex1/FamGo-platform/packages/telemetry/config"
    "github.com/Abdex1/FamGo-platform/packages/telemetry/logger"
    "github.com/Abdex1/FamGo-platform/packages/telemetry/tracing"
)

type SDK struct {
    Logger *logger.Logger
}

func Bootstrap(ctx context.Context, cfg config.Config) (*SDK, error) {
    _, err := tracing.InitProvider(ctx, cfg.OTLPEndpoint)
    if err != nil {
        return nil, err
    }

    log := logger.New(cfg.ServiceName, cfg.Environment)

    return &SDK{
        Logger: log,
    }, nil
}
