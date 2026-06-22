// services/ride-service/tests/integration/event_workflow_test.go
// Integration Tests for Event-Driven Workflows

package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/shared/contracts/events"
)

// TestEventWorkflow_RideRequested tests the complete ride requested workflow
func TestEventWorkflow_RideRequested(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create command handler
	mockRepo := NewMockRideRepository()
	mockEventPub := NewMockEventPublisher()
	createHandler := application.NewCreateRideHandler(
		mockRepo,
		nil, // cache
		nil, // domain service
		mockEventPub,
		logger,
	)

	// Execute: Create ride
	cmd := application.CreateRideCommand{
		PassengerID: "passenger123",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	rideID, err := createHandler.Handle(ctx, cmd)

	// Verify: Ride created
	assert.NoError(t, err)
	assert.NotEmpty(t, rideID)

	// Verify: Event published
	require.Equal(t, 1, len(mockEventPub.PublishedEvents))
	event := mockEventPub.PublishedEvents[0]
	assert.Equal(t, "RideRequested", event.EventType)
	assert.Equal(t, rideID, event.AggregateID)

	// Verify: Event contains correct data
	var payload map[string]interface{}
	err = json.Unmarshal(event.Data, &payload)
	assert.NoError(t, err)
	assert.Equal(t, "passenger123", payload["passenger_id"])
}

// TestEventWorkflow_DriverAssigned tests the driver assignment workflow
func TestEventWorkflow_DriverAssigned(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	mockRepo := NewMockRideRepository()
	subscriber := application.NewRideEventSubscriber(mockRepo, logger)

	// Simulate DriverAssigned event from dispatch service
	event := &events.Event{
		EventID:     "evt-123",
		EventType:   "DriverAssigned",
		AggregateID: "ride-456",
		Data: func() []byte {
			payload := map[string]interface{}{
				"ride_id":        "ride-456",
				"driver_id":      "driver-789",
				"estimated_fare": 15.50,
			}
			data, _ := json.Marshal(payload)
			return data
		}(),
	}

	// Execute: Handle event
	err := subscriber.HandleDriverAssigned(ctx, event)

	// Verify: No error
	assert.NoError(t, err)

	// Verify: Ride was updated
	assert.Equal(t, 1, len(mockRepo.SavedRides))
	savedRide := mockRepo.SavedRides[0]
	assert.Equal(t, "driver-789", savedRide.DriverID)
	assert.Equal(t, float32(15.50), savedRide.EstimatedFare)
}

// TestEventWorkflow_PaymentProcessed tests payment processing workflow
func TestEventWorkflow_PaymentProcessed(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	mockRepo := NewMockRideRepository()
	subscriber := application.NewRideEventSubscriber(mockRepo, logger)

	// Simulate PaymentProcessed event from payment service
	event := &events.Event{
		EventID:     "evt-payment-1",
		EventType:   "PaymentProcessed",
		AggregateID: "ride-456",
		Data: func() []byte {
			payload := map[string]interface{}{
				"ride_id": "ride-456",
				"amount":  18.75,
				"status":  "SUCCESS",
			}
			data, _ := json.Marshal(payload)
			return data
		}(),
	}

	// Execute: Handle event
	err := subscriber.HandlePaymentProcessed(ctx, event)

	// Verify: No error (payment successful)
	assert.NoError(t, err)
}

// TestEventWorkflow_MultipleEventsSequence tests a sequence of related events
func TestEventWorkflow_MultipleEventsSequence(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	mockRepo := NewMockRideRepository()
	mockEventPub := NewMockEventPublisher()
	subscriber := application.NewRideEventSubscriber(mockRepo, logger)

	// Step 1: Create ride
	createHandler := application.NewCreateRideHandler(
		mockRepo,
		nil,
		nil,
		mockEventPub,
		logger,
	)

	cmd := application.CreateRideCommand{
		PassengerID: "pass-001",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	rideID, err := createHandler.Handle(ctx, cmd)
	assert.NoError(t, err)

	// Verify: RideRequested event published
	assert.Equal(t, 1, len(mockEventPub.PublishedEvents))
	assert.Equal(t, "RideRequested", mockEventPub.PublishedEvents[0].EventType)

	// Step 2: Dispatch assigns driver (simulated event)
	driverEvent := &events.Event{
		EventID:     "evt-driver-1",
		EventType:   "DriverAssigned",
		AggregateID: rideID,
		Data: func() []byte {
			payload := map[string]interface{}{
				"ride_id":        rideID,
				"driver_id":      "driver-001",
				"estimated_fare": 16.00,
			}
			data, _ := json.Marshal(payload)
			return data
		}(),
	}

	err = subscriber.HandleDriverAssigned(ctx, driverEvent)
	assert.NoError(t, err)

	// Verify: Ride updated with driver info
	ride, err := mockRepo.GetByID(ctx, rideID)
	assert.NoError(t, err)
	assert.Equal(t, "driver-001", ride.DriverID)

	// Step 3: Payment processed (simulated event)
	paymentEvent := &events.Event{
		EventID:     "evt-payment-1",
		EventType:   "PaymentProcessed",
		AggregateID: rideID,
		Data: func() []byte {
			payload := map[string]interface{}{
				"ride_id": rideID,
				"amount":  16.00,
				"status":  "SUCCESS",
			}
			data, _ := json.Marshal(payload)
			return data
		}(),
	}

	err = subscriber.HandlePaymentProcessed(ctx, paymentEvent)
	assert.NoError(t, err)

	// Verify: Complete workflow sequence executed
	assert.Equal(t, 1, len(mockEventPub.PublishedEvents)) // Only RideRequested published
	assert.Equal(t, 1, len(mockRepo.SavedRides))          // Ride updated once
}

// TestEventWorkflow_ConcurrentRides tests multiple concurrent rides
func TestEventWorkflow_ConcurrentRides(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	mockRepo := NewMockRideRepository()
	mockEventPub := NewMockEventPublisher()

	createHandler := application.NewCreateRideHandler(
		mockRepo,
		nil,
		nil,
		mockEventPub,
		logger,
	)

	// Create multiple rides concurrently
	numRides := 10
	rideIDs := make([]string, 0, numRides)

	done := make(chan error, numRides)
	for i := 0; i < numRides; i++ {
		go func(index int) {
			cmd := application.CreateRideCommand{
				PassengerID: "passenger-" + string(rune(index)),
				PickupLat:   37.7749 + float64(index)*0.001,
				PickupLon:   -122.4194 + float64(index)*0.001,
				DropoffLat:  37.8044,
				DropoffLon:  -122.2712,
			}

			id, err := createHandler.Handle(ctx, cmd)
			if err == nil {
				rideIDs = append(rideIDs, id)
			}
			done <- err
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < numRides; i++ {
		err := <-done
		assert.NoError(t, err)
	}

	// Verify: All rides created
	assert.Equal(t, numRides, len(mockEventPub.PublishedEvents))
	for _, event := range mockEventPub.PublishedEvents {
		assert.Equal(t, "RideRequested", event.EventType)
	}
}

// Mock implementations for testing

type MockRideRepository struct {
	SavedRides map[string]interface{}
	saved      []interface{}
}

func NewMockRideRepository() *MockRideRepository {
	return &MockRideRepository{
		SavedRides: make(map[string]interface{}),
		saved:      make([]interface{}, 0),
	}
}

func (m *MockRideRepository) Save(ctx context.Context, ride interface{}) error {
	m.saved = append(m.saved, ride)
	return nil
}

func (m *MockRideRepository) GetByID(ctx context.Context, id string) (interface{}, error) {
	if ride, ok := m.SavedRides[id]; ok {
		return ride, nil
	}
	return nil, nil
}

type MockEventPublisher struct {
	PublishedEvents []*events.Event
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		PublishedEvents: make([]*events.Event, 0),
	}
}

func (m *MockEventPublisher) Publish(ctx context.Context, event *events.Event) error {
	m.PublishedEvents = append(m.PublishedEvents, event)
	return nil
}
