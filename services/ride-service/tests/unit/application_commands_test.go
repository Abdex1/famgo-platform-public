// services/ride-service/tests/unit/application_commands_test.go
// Application Layer Command Handler Tests

package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// MockRideRepository for testing
type MockRideRepository struct {
	mock.Mock
}

func (m *MockRideRepository) Save(ctx context.Context, ride *domain.Ride) error {
	args := m.Called(ctx, ride)
	return args.Error(0)
}

func (m *MockRideRepository) GetByID(ctx context.Context, id string) (*domain.Ride, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ride), args.Error(1)
}

// MockEventPublisher for testing
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(ctx context.Context, eventType string, aggregateID string, data interface{}) error {
	args := m.Called(ctx, eventType, aggregateID, data)
	return args.Error(0)
}

// Test CreateRideHandler
func TestCreateRideHandler_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockRideRepository)
	mockEventPub := new(MockEventPublisher)
	logger, _ := zap.NewDevelopment()

	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(r *domain.Ride) bool {
		return r.PassengerID == "p123" &&
			r.PickupLat == 37.7749 &&
			r.DropoffLat == 37.8044
	})).Return(nil)

	mockEventPub.On("Publish", mock.Anything, "RideRequested", mock.Anything, mock.Anything).Return(nil)

	handler := application.NewCreateRideHandler(mockRepo, mockEventPub, logger)
	cmd := application.CreateRideCommand{
		PassengerID: "p123",
		PickupLat:   37.7749,
		PickupLon:   -122.4194,
		DropoffLat:  37.8044,
		DropoffLon:  -122.2712,
	}

	// Execute
	rideID, err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, rideID)
	mockRepo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
	mockEventPub.AssertCalled(t, "Publish", mock.Anything, "RideRequested", mock.Anything, mock.Anything)
}

// Test AssignDriverHandler
func TestAssignDriverHandler_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockRideRepository)
	mockEventPub := new(MockEventPublisher)
	logger, _ := zap.NewDevelopment()

	ride := &domain.Ride{
		ID:          "ride123",
		PassengerID: "p123",
		Status:      domain.RideStatusRequested,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, "ride123").Return(ride, nil)
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(r *domain.Ride) bool {
		return r.DriverID == "d456"
	})).Return(nil)
	mockEventPub.On("Publish", mock.Anything, "DriverAssigned", "ride123", mock.Anything).Return(nil)

	handler := application.NewAssignDriverHandler(mockRepo, mockEventPub, logger)
	cmd := application.AssignDriverCommand{
		RideID:   "ride123",
		DriverID: "d456",
	}

	// Execute
	err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "GetByID", mock.Anything, "ride123")
	mockRepo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
	mockEventPub.AssertCalled(t, "Publish", mock.Anything, "DriverAssigned", "ride123", mock.Anything)
}

// Test StartRideHandler
func TestStartRideHandler_Success(t *testing.T) {
	mockRepo := new(MockRideRepository)
	mockEventPub := new(MockEventPublisher)
	logger, _ := zap.NewDevelopment()

	ride := &domain.Ride{
		ID:        "ride123",
		Status:    domain.RideStatusAssigned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, "ride123").Return(ride, nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockEventPub.On("Publish", mock.Anything, "RideStarted", "ride123", mock.Anything).Return(nil)

	handler := application.NewStartRideHandler(mockRepo, mockEventPub, logger)
	cmd := application.StartRideCommand{RideID: "ride123"}

	err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	mockEventPub.AssertCalled(t, "Publish", mock.Anything, "RideStarted", "ride123", mock.Anything)
}

// Test CompleteRideHandler
func TestCompleteRideHandler_Success(t *testing.T) {
	mockRepo := new(MockRideRepository)
	mockEventPub := new(MockEventPublisher)
	logger, _ := zap.NewDevelopment()

	ride := &domain.Ride{
		ID:        "ride123",
		Status:    domain.RideStatusStarted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, "ride123").Return(ride, nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockEventPub.On("Publish", mock.Anything, "RideCompleted", "ride123", mock.Anything).Return(nil)

	handler := application.NewCompleteRideHandler(mockRepo, mockEventPub, logger)
	cmd := application.CompleteRideCommand{
		RideID:     "ride123",
		ActualFare: 18.75,
	}

	err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	mockEventPub.AssertCalled(t, "Publish", mock.Anything, "RideCompleted", "ride123", mock.Anything)
}

// Test CancelRideHandler
func TestCancelRideHandler_Success(t *testing.T) {
	mockRepo := new(MockRideRepository)
	mockEventPub := new(MockEventPublisher)
	logger, _ := zap.NewDevelopment()

	ride := &domain.Ride{
		ID:        "ride123",
		Status:    domain.RideStatusAssigned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, "ride123").Return(ride, nil)
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockEventPub.On("Publish", mock.Anything, "RideCancelled", "ride123", mock.Anything).Return(nil)

	handler := application.NewCancelRideHandler(mockRepo, mockEventPub, logger)
	cmd := application.CancelRideCommand{
		RideID: "ride123",
		Reason: "passenger request",
	}

	err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	mockEventPub.AssertCalled(t, "Publish", mock.Anything, "RideCancelled", "ride123", mock.Anything)
}
