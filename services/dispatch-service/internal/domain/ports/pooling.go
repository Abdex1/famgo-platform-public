package ports

import "context"

// PoolingStrategyHook allows future ride-pooling without changing assignment APIs.
type PoolingStrategyHook interface {
	CanPool(ctx context.Context, rideID string, driverID string) (bool, error)
}

// NoOpPoolingHook is the default until pooling-service is integrated.
type NoOpPoolingHook struct{}

func (NoOpPoolingHook) CanPool(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}
