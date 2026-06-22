// services/ride-service/tests/integration/full_workflow_test.go
// Full Integration Tests for Complete Ride Workflows

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// TestFullWorkflow_CreateToComplete tests complete ride lifecycle
func TestFullWorkflow_CreateToComplete(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Setup repositories and handlers
	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()

	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)
	assignHandler := application.NewAssignDriverHandler(rideRepo, nil, nil, eventPub, logger)
	startHandler := application.NewStartRideHandler(rideRepo, nil, nil, eventPub, logger)
	completeHandler := application.NewCompleteRideHandler(rideRepo, nil, nil, eventPub, logger)

	// Step 1: Create ride
	createCmd := application.CreateRideCommand{
		PassengerID: "passenger-001",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	rideID, err := createHandler.Handle(ctx, createCmd)
	require.NoError(t, err)
	require.NotEmpty(t, rideID)
	assert.Equal(t, 1, len(eventPub.PublishedEvents))
	assert.Equal(t, "RideRequested", eventPub.PublishedEvents[0].EventType)

	// Step 2: Assign driver
	assignCmd := application.AssignDriverCommand{
		RideID:   rideID,
		DriverID: "driver-001",
	}

	err = assignHandler.Handle(ctx, assignCmd)
	require.NoError(t, err)
	assert.Equal(t, 2, len(eventPub.PublishedEvents))
	assert.Equal(t, "DriverAssigned", eventPub.PublishedEvents[1].EventType)

	// Verify ride state
	ride, err := rideRepo.GetByID(ctx, rideID)
	require.NoError(t, err)
	assert.NotNil(t, ride)

	// Step 3: Start ride
	startCmd := application.StartRideCommand{RideID: rideID}
	err = startHandler.Handle(ctx, startCmd)
	require.NoError(t, err)
	assert.Equal(t, 3, len(eventPub.PublishedEvents))
	assert.Equal(t, "RideStarted", eventPub.PublishedEvents[2].EventType)

	// Step 4: Complete ride
	completeCmd := application.CompleteRideCommand{
		RideID:     rideID,
		ActualFare: 18.75,
	}

	err = completeHandler.Handle(ctx, completeCmd)
	require.NoError(t, err)
	assert.Equal(t, 4, len(eventPub.PublishedEvents))
	assert.Equal(t, "RideCompleted", eventPub.PublishedEvents[3].EventType)

	t.Logf("✅ Complete workflow successful: %s → %s", rideID, "COMPLETED")
}

// TestFullWorkflow_CreateToCancel tests ride cancellation flow
func TestFullWorkflow_CreateToCancel(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()

	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)
	cancelHandler := application.NewCancelRideHandler(rideRepo, nil, nil, eventPub, logger)

	// Create ride
	createCmd := application.CreateRideCommand{
		PassengerID: "passenger-002",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	rideID, err := createHandler.Handle(ctx, createCmd)
	require.NoError(t, err)

	// Cancel ride
	cancelCmd := application.CancelRideCommand{
		RideID: rideID,
		Reason: "passenger changed mind",
	}

	err = cancelHandler.Handle(ctx, cancelCmd)
	require.NoError(t, err)

	// Verify events
	assert.Equal(t, 2, len(eventPub.PublishedEvents))
	assert.Equal(t, "RideRequested", eventPub.PublishedEvents[0].EventType)
	assert.Equal(t, "RideCancelled", eventPub.PublishedEvents[1].EventType)

	t.Logf("✅ Cancellation workflow successful: %s cancelled", rideID)
}

// TestFullWorkflow_MultipleRidesConcurrent tests concurrent ride creation
func TestFullWorkflow_MultipleRidesConcurrent(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()

	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)

	// Create 20 rides concurrently
	numRides := 20
	done := make(chan error, numRides)
	rideIDs := make([]string, numRides)

	for i := 0; i < numRides; i++ {
		go func(index int) {
			cmd := application.CreateRideCommand{
				PassengerID: "passenger-" + string(rune(index)),
				PickupLat:   37.7749 + float64(index)*0.001,
				PickupLon:   -122.4194,
				DropoffLat:  37.8044,
				DropoffLon:  -122.2712,
			}

			id, err := createHandler.Handle(ctx, cmd)
			rideIDs[index] = id
			done <- err
		}(i)
	}

	// Wait for all to complete
	errorCount := 0
	for i := 0; i < numRides; i++ {
		err := <-done
		if err != nil {
			errorCount++
		}
	}

	assert.Equal(t, 0, errorCount)
	assert.Equal(t, numRides, len(eventPub.PublishedEvents))

	t.Logf("✅ Concurrent workflow successful: %d rides created", numRides)
}

// TestFullWorkflow_HighThroughput tests ride creation under load
func TestFullWorkflow_HighThroughput(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()
	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)

	// Create 100 rides sequentially
	start := time.Now()
	successCount := 0
	failureCount := 0

	for i := 0; i < 100; i++ {
		cmd := application.CreateRideCommand{
			PassengerID: "passenger-load-" + string(rune(i%10)),
			PickupLat:   37.7749 + float64(i%10)*0.001,
			PickupLon:   -122.4194,
			DropoffLat:  37.8044 + float64(i%10)*0.001,
			DropoffLon:  -122.2712,
		}

		_, err := createHandler.Handle(ctx, cmd)
		if err != nil {
			failureCount++
		} else {
			successCount++
		}
	}

	duration := time.Since(start)

	assert.Equal(t, 100, successCount)
	assert.Equal(t, 0, failureCount)
	assert.Equal(t, 100, len(eventPub.PublishedEvents))

	throughput := float64(successCount) / duration.Seconds()
	t.Logf("✅ High throughput test: %d rides in %.2fs = %.0f rides/sec", successCount, duration.Seconds(), throughput)
}

// TestFullWorkflow_StateTransitionValidation tests invalid state transitions
func TestFullWorkflow_StateTransitionValidation(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Test that invalid state transitions are rejected
	ride := &domain.Ride{
		ID:          "test-ride",
		PassengerID: "passenger-001",
		Status:      domain.RideStatusRequested,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Valid transition: REQUESTED → SEARCHING
	assert.True(t, ride.CanTransitionTo(domain.RideStatusSearching))

	// Invalid transition: REQUESTED → COMPLETED (must go through other states)
	assert.False(t, ride.CanTransitionTo(domain.RideStatusCompleted))

	// Invalid transition: REQUESTED → ASSIGNED (must go through SEARCHING first)
	assert.False(t, ride.CanTransitionTo(domain.RideStatusAssigned))

	// Transition to CANCELLED is always allowed
	assert.True(t, ride.CanTransitionTo(domain.RideStatusCancelled))

	t.Log("✅ State transition validation successful")
}

// TestFullWorkflow_RapidAssignmentAndCompletion tests rapid assignment after creation
func TestFullWorkflow_RapidAssignmentAndCompletion(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()

	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)
	assignHandler := application.NewAssignDriverHandler(rideRepo, nil, nil, eventPub, logger)
	completeHandler := application.NewCompleteRideHandler(rideRepo, nil, nil, eventPub, logger)

	// Create ride
	createCmd := application.CreateRideCommand{
		PassengerID: "passenger-rapid",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	rideID, err := createHandler.Handle(ctx, createCmd)
	require.NoError(t, err)

	// Immediately assign driver
	assignCmd := application.AssignDriverCommand{
		RideID:   rideID,
		DriverID: "driver-rapid",
	}

	err = assignHandler.Handle(ctx, assignCmd)
	require.NoError(t, err)

	// Immediately complete ride
	completeCmd := application.CompleteRideCommand{
		RideID:     rideID,
		ActualFare: 12.50,
	}

	err = completeHandler.Handle(ctx, completeCmd)
	require.NoError(t, err)

	// Verify 3 events published
	assert.Equal(t, 3, len(eventPub.PublishedEvents))

	t.Logf("✅ Rapid assignment workflow successful: %s", rideID)
}

// ============= PERFORMANCE TESTS =============

// BenchmarkRideCreation benchmarks ride creation performance
func BenchmarkRideCreation(b *testing.B) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()
	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)

	cmd := application.CreateRideCommand{
		PassengerID: "passenger-bench",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = createHandler.Handle(ctx, cmd)
	}
}

// BenchmarkRideAssignment benchmarks driver assignment performance
func BenchmarkRideAssignment(b *testing.B) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	rideRepo := NewMockRideRepository()
	eventPub := NewMockEventPublisher()

	// Create a ride first
	createHandler := application.NewCreateRideHandler(rideRepo, nil, nil, eventPub, logger)
	createCmd := application.CreateRideCommand{
		PassengerID: "passenger-bench",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	rideID, _ := createHandler.Handle(ctx, createCmd)

	assignHandler := application.NewAssignDriverHandler(rideRepo, nil, nil, eventPub, logger)
	assignCmd := application.AssignDriverCommand{
		RideID:   rideID,
		DriverID: "driver-bench",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = assignHandler.Handle(ctx, assignCmd)
	}
}
