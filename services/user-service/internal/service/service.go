package service

import (
	"context"
	"fmt"

	"famgo/user-service/internal/model"
	"famgo/user-service/internal/repository"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// UserService handles user profile operations
type UserService struct {
	userRepo      *repository.UserRepository
	prefsRepo     *repository.PreferencesRepository
	addressRepo   *repository.AddressRepository
	logger        logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, log logger.Logger) *UserService {
	db := userRepo.db
	return &UserService{
		userRepo:    userRepo,
		prefsRepo:   repository.NewPreferencesRepository(db),
		addressRepo: repository.NewAddressRepository(db),
		logger:      log,
	}
}

// GetProfile retrieves user profile
func (s *UserService) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	profile, err := s.userRepo.GetProfileByID(ctx, userID)
	if err != nil {
		s.logger.Warn("get profile failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	return profile, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.UserProfile, error) {
	// Get existing profile
	profile, err := s.userRepo.GetProfileByID(ctx, userID)
	if err != nil {
		s.logger.Warn("profile not found", map[string]interface{}{"user_id": userID})
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		profile.FirstName = req.FirstName
	}
	if req.LastName != "" {
		profile.LastName = req.LastName
	}
	if req.ProfilePictureURL != "" {
		profile.ProfilePictureURL = req.ProfilePictureURL
	}

	// Save to database
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		s.logger.Error("update profile failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	s.logger.Info("profile updated", map[string]interface{}{"user_id": userID})
	return profile, nil
}

// GetPreferences retrieves user preferences
func (s *UserService) GetPreferences(ctx context.Context, userID string) (*model.UserPreferences, error) {
	prefs, err := s.prefsRepo.GetPreferencesByUserID(ctx, userID)
	if err != nil {
		s.logger.Warn("get preferences failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	return prefs, nil
}

// UpdatePreferences updates user preferences
func (s *UserService) UpdatePreferences(ctx context.Context, userID string, req *model.UpdatePreferencesRequest) (*model.UserPreferences, error) {
	// Get existing preferences
	prefs, err := s.prefsRepo.GetPreferencesByUserID(ctx, userID)
	if err != nil {
		s.logger.Warn("preferences not found", map[string]interface{}{"user_id": userID})
		return nil, err
	}

	// Update fields
	prefs.NotificationEmail = req.NotificationEmail
	prefs.NotificationSMS = req.NotificationSMS
	if req.Language != "" {
		prefs.Language = req.Language
	}

	// Save to database
	if err := s.prefsRepo.UpdatePreferences(ctx, prefs); err != nil {
		s.logger.Error("update preferences failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	s.logger.Info("preferences updated", map[string]interface{}{"user_id": userID})
	return prefs, nil
}

// CreatePreferences creates default preferences for new user
func (s *UserService) CreatePreferences(ctx context.Context, userID string) (*model.UserPreferences, error) {
	prefs, err := s.prefsRepo.CreatePreferences(ctx, userID)
	if err != nil {
		s.logger.Error("create preferences failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	return prefs, nil
}

// AddAddress adds a new saved address
func (s *UserService) AddAddress(ctx context.Context, userID string, req *model.AddressRequest) (*model.UserAddress, error) {
	// Validate address type
	if req.Type != "home" && req.Type != "work" {
		return nil, fmt.Errorf("invalid address type: %s", req.Type)
	}

	address, err := s.addressRepo.CreateAddress(ctx, userID, req)
	if err != nil {
		s.logger.Error("add address failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	s.logger.Info("address added", map[string]interface{}{"user_id": userID, "address_id": address.ID})
	return address, nil
}

// GetAddresses retrieves all saved addresses
func (s *UserService) GetAddresses(ctx context.Context, userID string) ([]*model.UserAddress, error) {
	addresses, err := s.addressRepo.GetAddressesByUserID(ctx, userID)
	if err != nil {
		s.logger.Warn("get addresses failed", map[string]interface{}{"user_id": userID, "error": err})
		return nil, err
	}

	return addresses, nil
}

// DeleteAddress deletes a saved address
func (s *UserService) DeleteAddress(ctx context.Context, addressID string) error {
	if err := s.addressRepo.DeleteAddress(ctx, addressID); err != nil {
		s.logger.Error("delete address failed", map[string]interface{}{"address_id": addressID, "error": err})
		return err
	}

	s.logger.Info("address deleted", map[string]interface{}{"address_id": addressID})
	return nil
}
