// services/ride-service/internal/application/grpc_clients.go
// gRPC Service Client Calls

package application

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/grpc-clients"
)

// RideGRPCClients manages cross-service gRPC calls
type RideGRPCClients struct {
	clientPool *grpcclient.GRPCClientPool
	logger     *zap.Logger
}

func NewRideGRPCClients(clientPool *grpcclient.GRPCClientPool, logger *zap.Logger) *RideGRPCClients {
	return &RideGRPCClients{
		clientPool: clientPool,
		logger:     logger,
	}
}

// CalculateFare calls Pricing Service to calculate fare
func (c *RideGRPCClients) CalculateFare(
	ctx context.Context,
	rideID string,
	pickupLat, pickupLon, dropoffLat, dropoffLon float64,
	vehicleType string,
) (float32, error) {
	c.logger.Info("calling pricing service",
		zap.String("ride_id", rideID),
		zap.String("operation", "CalculateFare"))

	// This would use the generated gRPC client stubs
	// For demonstration, showing the structure:
	// pricingClient := pricingpb.NewPricingServiceClient(c.clientPool.GetPricingClient())
	//
	// ctx, cancel := grpcclient.WithTimeout(ctx, 5*time.Second)
	// defer cancel()
	//
	// req := &pricingpb.CalculateFareRequest{
	//     RideId: rideID,
	//     PickupLat: pickupLat,
	//     PickupLon: pickupLon,
	//     DropoffLat: dropoffLat,
	//     DropoffLon: dropoffLon,
	//     VehicleType: vehicleType,
	//     ServiceType: "ECONOMY",
	// }
	//
	// resp, err := grpcclient.CallWithRetry(ctx, 3, time.Second, func(retryCtx context.Context) error {
	//     var err error
	//     resp, err = pricingClient.CalculateFare(retryCtx, req)
	//     return err
	// })
	//
	// if err != nil {
	//     c.logger.Error("pricing call failed", zap.Error(err))
	//     return 0, err
	// }
	//
	// c.logger.Info("fare calculated",
	//     zap.String("ride_id", rideID),
	//     zap.Float32("total_fare", resp.TotalFare))
	//
	// return resp.TotalFare, nil

	// For now, return a default value
	c.logger.Info("fare calculation mock (would use gRPC in production)")
	return 15.50, nil
}

// GetDriverLocation calls GPS Service to get driver location
func (c *RideGRPCClients) GetDriverLocation(
	ctx context.Context,
	driverID string,
) (lat, lon float64, err error) {
	c.logger.Info("calling gps service",
		zap.String("driver_id", driverID),
		zap.String("operation", "GetLocation"))

	// This would use the generated gRPC client stubs
	// Similar pattern to CalculateFare above

	// For now, return mock values
	c.logger.Info("driver location retrieved (mock)")
	return 37.7749, -122.4194, nil
}

// FindDrivers calls Dispatch Service to find available drivers
func (c *RideGRPCClients) FindDrivers(
	ctx context.Context,
	rideID string,
	pickupLat, pickupLon float64,
	requiredCount int,
) ([]string, error) {
	c.logger.Info("calling dispatch service",
		zap.String("ride_id", rideID),
		zap.String("operation", "FindDrivers"),
		zap.Int("required_count", requiredCount))

	// This would use the generated gRPC client stubs
	// Similar pattern to CalculateFare above

	// For now, return mock driver IDs
	c.logger.Info("drivers found (mock)")
	return []string{"driver1", "driver2", "driver3"}, nil
}

// SubscribeToLocationUpdates subscribes to driver location updates from GPS Service
func (c *RideGRPCClients) SubscribeToLocationUpdates(
	ctx context.Context,
	rideID string,
	driverID string,
) error {
	c.logger.Info("subscribing to location updates",
		zap.String("ride_id", rideID),
		zap.String("driver_id", driverID))

	// In production, this would set up a stream subscription
	// to receive location updates from GPS Service in real-time

	return nil
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxFailures     int
	FailureThreshold float32 // 0.5 = 50%
	OpenTimeout     time.Duration
	SuccessThreshold int
}

// ClientWithCircuitBreaker wraps a client with circuit breaker pattern
type ClientWithCircuitBreaker struct {
	conn              *grpc.ClientConn
	failureCount      int
	successCount      int
	isOpen            bool
	openedAt          time.Time
	config            CircuitBreakerConfig
	logger            *zap.Logger
}

// NewClientWithCircuitBreaker creates a new circuit breaker protected client
func NewClientWithCircuitBreaker(
	conn *grpc.ClientConn,
	config CircuitBreakerConfig,
	logger *zap.Logger,
) *ClientWithCircuitBreaker {
	return &ClientWithCircuitBreaker{
		conn:    conn,
		config:  config,
		logger:  logger,
		isOpen:  false,
	}
}

// Call executes a call with circuit breaker protection
func (cb *ClientWithCircuitBreaker) Call(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	// Check if circuit should be closed (timeout elapsed)
	if cb.isOpen {
		if time.Since(cb.openedAt) > cb.config.OpenTimeout {
			cb.logger.Info("circuit breaker half-open", zap.Duration("timeout", cb.config.OpenTimeout))
			cb.isOpen = false
			cb.failureCount = 0
			cb.successCount = 0
		} else {
			cb.logger.Warn("circuit breaker open, rejecting call")
			return fmt.Errorf("circuit breaker open")
		}
	}

	// Execute call
	err := fn(ctx)

	if err != nil {
		cb.failureCount++
		cb.successCount = 0

		failureRate := float32(cb.failureCount) / float32(cb.config.MaxFailures)
		if failureRate >= cb.config.FailureThreshold {
			cb.isOpen = true
			cb.openedAt = time.Now()
			cb.logger.Error("circuit breaker opened",
				zap.Int("failures", cb.failureCount),
				zap.Float32("rate", failureRate))
		}

		return err
	}

	// Success
	cb.successCount++
	if cb.successCount >= cb.config.SuccessThreshold {
		cb.failureCount = 0
		cb.successCount = 0
		cb.logger.Info("circuit breaker reset after successful calls",
			zap.Int("successes", cb.successCount))
	}

	return nil
}
