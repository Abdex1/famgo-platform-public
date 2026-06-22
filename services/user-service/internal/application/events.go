// services/user-service/internal/application/events.go
// User Service Event Publishing

package application

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/event-bus"
	"github.com/Abdex1/FamGo-platform/shared/contracts/events"
)

// UserEventPublisher publishes user service events
type UserEventPublisher struct {
	eventBus event_bus.EventBus
	logger   *zap.Logger
}

func NewUserEventPublisher(eventBus event_bus.EventBus, logger *zap.Logger) *UserEventPublisher {
	return &UserEventPublisher{
		eventBus: eventBus,
		logger:   logger,
	}
}

// PublishUserRegistered publishes user registration event
func (p *UserEventPublisher) PublishUserRegistered(
	ctx context.Context,
	userID string,
	phone string,
	email string,
	userType string, // PASSENGER, DRIVER
) error {
	payload := map[string]interface{}{
		"user_id":    userID,
		"phone":      phone,
		"email":      email,
		"user_type":  userType,
		"timestamp":  time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "UserRegistered",
		AggregateID: userID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published UserRegistered event",
		zap.String("user_id", userID),
		zap.String("user_type", userType))

	return nil
}

// PublishUserProfileUpdated publishes user profile update event
func (p *UserEventPublisher) PublishUserProfileUpdated(
	ctx context.Context,
	userID string,
	firstName string,
	lastName string,
	email string,
) error {
	payload := map[string]interface{}{
		"user_id":     userID,
		"first_name":  firstName,
		"last_name":   lastName,
		"email":       email,
		"timestamp":   time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "UserProfileUpdated",
		AggregateID: userID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published UserProfileUpdated event", zap.String("user_id", userID))
	return nil
}

// PublishDriverVerified publishes driver verification event
func (p *UserEventPublisher) PublishDriverVerified(ctx context.Context, userID string) error {
	payload := map[string]interface{}{
		"user_id":    userID,
		"status":     "VERIFIED",
		"timestamp":  time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "DriverVerified",
		AggregateID: userID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published DriverVerified event", zap.String("user_id", userID))
	return nil
}

// PublishDriverSuspended publishes driver suspension event
func (p *UserEventPublisher) PublishDriverSuspended(ctx context.Context, userID string, reason string) error {
	payload := map[string]interface{}{
		"user_id":    userID,
		"status":     "SUSPENDED",
		"reason":     reason,
		"timestamp":  time.Now().UTC(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	event := &events.Event{
		EventType:   "DriverSuspended",
		AggregateID: userID,
		Data:        data,
	}

	if err := p.eventBus.Publish(ctx, event); err != nil {
		p.logger.Error("failed to publish event", zap.Error(err))
		return err
	}

	p.logger.Info("published DriverSuspended event",
		zap.String("user_id", userID),
		zap.String("reason", reason))

	return nil
}
