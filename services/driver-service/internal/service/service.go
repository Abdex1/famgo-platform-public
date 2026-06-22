package service

import (
	"context"
	"fmt"

	"famgo/driver-service/internal/model"
	"famgo/driver-service/internal/repository"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// DriverService handles driver operations (WEEK 1 FOUNDATION)
// Full verification, documents, location in WEEK 3
type DriverService struct {
	driverRepo *repository.DriverRepository
	stateRepo  *repository.DriverStateRepository
	logger     logger.Logger
}

// NewDriverService creates a new driver service
func NewDriverService(driverRepo *repository.DriverRepository, log logger.Logger) *DriverService {
	return &DriverService{
		driverRepo: driverRepo,
		stateRepo:  repository.NewDriverStateRepository(driverRepo.DB()),
		logger:     log,
	}
}

// RegisterDriver creates a new driver registration (2-step process)
func (s *DriverService) RegisterDriver(ctx context.Context, driver *model.Driver) error {
	if driver.AuthID == "" {
		return fmt.Errorf("auth_id is required")
	}

	if err := s.driverRepo.CreateDriver(ctx, driver); err != nil {
		s.logger.Error("failed to create driver", map[string]interface{}{"error": err})
		return err
	}

	// Initialize state: pending
	_, err := s.stateRepo.TransitionState(ctx, driver.ID, "", model.DriverStatePending, "Initial registration")
	if err != nil {
		s.logger.Error("failed to initialize driver state", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("driver registered", map[string]interface{}{"driver_id": driver.ID})
	return nil
}

// GetProfile retrieves driver profile
func (s *DriverService) GetProfile(ctx context.Context, driverID string) (*model.Driver, error) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		s.logger.Warn("get profile failed", map[string]interface{}{"driver_id": driverID, "error": err})
		return nil, err
	}

	return driver, nil
}

// UpdateProfile updates driver information (foundation - license only)
func (s *DriverService) UpdateProfile(ctx context.Context, driverID string, licenseNumber string) (*model.Driver, error) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		s.logger.Warn("driver not found", map[string]interface{}{"driver_id": driverID})
		return nil, err
	}

	if licenseNumber != "" {
		driver.LicenseNumber = licenseNumber
	}

	if err := s.driverRepo.UpdateDriver(ctx, driver); err != nil {
		s.logger.Error("update profile failed", map[string]interface{}{"driver_id": driverID, "error": err})
		return nil, err
	}

	s.logger.Info("driver profile updated", map[string]interface{}{"driver_id": driverID})
	return driver, nil
}

// TransitionState transitions driver to a new state (Pattern 4: State Machine)
func (s *DriverService) TransitionState(ctx context.Context, driverID string, newState, reason string) (*model.DriverState, error) {
	// Get current driver
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		s.logger.Warn("driver not found", map[string]interface{}{"driver_id": driverID})
		return nil, err
	}

	// Validate transition
	if !model.IsValidTransition(driver.Status, newState) {
		return nil, fmt.Errorf("invalid state transition: %s -> %s", driver.Status, newState)
	}

	// Record state transition
	state, err := s.stateRepo.TransitionState(ctx, driverID, driver.Status, newState, reason)
	if err != nil {
		s.logger.Error("state transition failed", map[string]interface{}{"driver_id": driverID, "error": err})
		return nil, err
	}

	// Update driver status
	driver.Status = newState
	if err := s.driverRepo.UpdateDriver(ctx, driver); err != nil {
		s.logger.Error("failed to update driver status", map[string]interface{}{"driver_id": driverID, "error": err})
		return nil, err
	}

	s.logger.Info("driver status transitioned", map[string]interface{}{
		"driver_id":   driverID,
		"from_state":  driver.Status,
		"to_state":    newState,
	})

	return state, nil
}

// GetStateHistory retrieves state transition history
func (s *DriverService) GetStateHistory(ctx context.Context, driverID string, limit int) ([]*model.DriverState, error) {
	states, err := s.stateRepo.GetStateHistory(ctx, driverID, limit)
	if err != nil {
		s.logger.Warn("get state history failed", map[string]interface{}{"driver_id": driverID, "error": err})
		return nil, err
	}

	return states, nil
}

// IsDriverActive checks if driver is active and available
func (s *DriverService) IsDriverActive(ctx context.Context, driverID string) (bool, error) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return false, err
	}

	return driver.Status == model.DriverStateActive, nil
}
