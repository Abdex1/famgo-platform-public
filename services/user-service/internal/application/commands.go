// services/user-service/internal/application/commands.go
// User Service Commands and Handlers

package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
	"go.uber.org/zap"
)

// ===== REGISTER USER COMMAND =====

type RegisterUserCommand struct {
	Phone     string
	Email     string
	FirstName string
	LastName  string
	UserType  string // "driver" or "passenger"
}

type RegisterUserHandler struct {
	userRepo    domain.UserRepository
	userCache   domain.UserCache
	userService *domain.UserService
	logger      *zap.Logger
}

func NewRegisterUserHandler(
	userRepo domain.UserRepository,
	userCache domain.UserCache,
	userService *domain.UserService,
	logger *zap.Logger,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo:    userRepo,
		userCache:   userCache,
		userService: userService,
		logger:      logger,
	}
}

func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (string, error) {
	// Validate input
	if !h.userService.ValidatePhoneNumber(cmd.Phone) {
		h.logger.Warn("invalid phone number", zap.String("phone", cmd.Phone))
		return "", errInvalidPhoneNumber
	}

	if !h.userService.ValidateEmail(cmd.Email) {
		h.logger.Warn("invalid email", zap.String("email", cmd.Email))
		return "", errInvalidEmail
	}

	// Check if user already exists
	existingUser, _ := h.userRepo.GetUserByPhone(ctx, cmd.Phone)
	if existingUser != nil {
		h.logger.Info("user already exists", zap.String("phone", cmd.Phone))
		return "", errUserAlreadyExists
	}

	// Create new user
	user := domain.NewUser(cmd.Phone, cmd.Email, cmd.FirstName, cmd.LastName)

	// Persist user
	if err := h.userRepo.CreateUser(ctx, user); err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		return "", err
	}

	// Cache user
	h.userCache.SetUser(ctx, user, 3600)

	h.logger.Info("user registered", zap.String("userID", user.ID), zap.String("phone", cmd.Phone))

	return user.ID, nil
}

// ===== UPDATE PROFILE COMMAND =====

type UpdateProfileCommand struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
}

type UpdateProfileHandler struct {
	userRepo    domain.UserRepository
	userCache   domain.UserCache
	userService *domain.UserService
	logger      *zap.Logger
}

func NewUpdateProfileHandler(
	userRepo domain.UserRepository,
	userCache domain.UserCache,
	userService *domain.UserService,
	logger *zap.Logger,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		userRepo:    userRepo,
		userCache:   userCache,
		userService: userService,
		logger:      logger,
	}
}

func (h *UpdateProfileHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) error {
	// Get user
	user, err := h.userRepo.GetUser(ctx, cmd.UserID)
	if err != nil {
		h.logger.Error("user not found", zap.String("userID", cmd.UserID))
		return err
	}

	// Update fields
	user.FirstName = cmd.FirstName
	user.LastName = cmd.LastName
	user.Email = cmd.Email
	user.UpdatedAt = time.Now()

	// Persist
	if err := h.userRepo.UpdateUser(ctx, user); err != nil {
		h.logger.Error("failed to update user", zap.Error(err))
		return err
	}

	// Invalidate cache
	h.userCache.DeleteUser(ctx, cmd.UserID)

	h.logger.Info("user profile updated", zap.String("userID", cmd.UserID))

	return nil
}

// ===== ACTIVATE USER COMMAND =====

type ActivateUserCommand struct {
	UserID string
}

type ActivateUserHandler struct {
	userRepo    domain.UserRepository
	userCache   domain.UserCache
	userService *domain.UserService
	logger      *zap.Logger
}

func NewActivateUserHandler(
	userRepo domain.UserRepository,
	userCache domain.UserCache,
	userService *domain.UserService,
	logger *zap.Logger,
) *ActivateUserHandler {
	return &ActivateUserHandler{
		userRepo:    userRepo,
		userCache:   userCache,
		userService: userService,
		logger:      logger,
	}
}

func (h *ActivateUserHandler) Handle(ctx context.Context, cmd ActivateUserCommand) error {
	// Get user
	user, err := h.userRepo.GetUser(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	// Check if can activate
	if !h.userService.CanActivateUser(user) {
		h.logger.Warn("cannot activate user", zap.String("userID", cmd.UserID), zap.String("status", string(user.Status)))
		return errCannotActivateUser
	}

	// Activate
	user.Activate()

	// Persist
	if err := h.userRepo.UpdateUser(ctx, user); err != nil {
		h.logger.Error("failed to activate user", zap.Error(err))
		return err
	}

	// Invalidate cache
	h.userCache.DeleteUser(ctx, cmd.UserID)

	h.logger.Info("user activated", zap.String("userID", cmd.UserID))

	return nil
}

// ===== VERIFY DRIVER COMMAND =====

type VerifyDriverCommand struct {
	DriverID string
}

type VerifyDriverHandler struct {
	driverRepo  domain.DriverProfileRepository
	driverCache domain.UserCache
	userService *domain.UserService
	logger      *zap.Logger
}

func NewVerifyDriverHandler(
	driverRepo domain.DriverProfileRepository,
	driverCache domain.UserCache,
	userService *domain.UserService,
	logger *zap.Logger,
) *VerifyDriverHandler {
	return &VerifyDriverHandler{
		driverRepo:  driverRepo,
		driverCache: driverCache,
		userService: userService,
		logger:      logger,
	}
}

func (h *VerifyDriverHandler) Handle(ctx context.Context, cmd VerifyDriverCommand) error {
	// Get driver profile
	profile, err := h.driverRepo.GetProfile(ctx, cmd.DriverID)
	if err != nil {
		h.logger.Error("driver profile not found", zap.String("driverID", cmd.DriverID))
		return err
	}

	// Check if can verify
	if !h.userService.CanVerifyDriver(profile) {
		h.logger.Warn("cannot verify driver", zap.String("driverID", cmd.DriverID))
		return errCannotVerifyDriver
	}

	// Verify
	profile.Verify()

	// Persist
	if err := h.driverRepo.UpdateProfile(ctx, profile); err != nil {
		h.logger.Error("failed to verify driver", zap.Error(err))
		return err
	}

	// Invalidate cache
	h.driverCache.DeleteDriverProfile(ctx, cmd.DriverID)

	h.logger.Info("driver verified", zap.String("driverID", cmd.DriverID))

	return nil
}

// ===== CREATE DRIVER PROFILE COMMAND =====

type CreateDriverProfileCommand struct {
	UserID         string
	LicenseNumber  string
	LicenseExpiry  time.Time
	VehicleNumber  string
	VehicleType    string
}

type CreateDriverProfileHandler struct {
	driverRepo  domain.DriverProfileRepository
	driverCache domain.UserCache
	logger      *zap.Logger
}

func NewCreateDriverProfileHandler(
	driverRepo domain.DriverProfileRepository,
	driverCache domain.UserCache,
	logger *zap.Logger,
) *CreateDriverProfileHandler {
	return &CreateDriverProfileHandler{
		driverRepo:  driverRepo,
		driverCache: driverCache,
		logger:      logger,
	}
}

func (h *CreateDriverProfileHandler) Handle(ctx context.Context, cmd CreateDriverProfileCommand) (string, error) {
	// Create profile
	profile := domain.NewDriverProfile(cmd.UserID, cmd.LicenseNumber, cmd.VehicleNumber, cmd.VehicleType)
	profile.LicenseExpiry = cmd.LicenseExpiry

	// Persist
	if err := h.driverRepo.CreateProfile(ctx, profile); err != nil {
		h.logger.Error("failed to create driver profile", zap.Error(err))
		return "", err
	}

	// Cache
	h.driverCache.SetDriverProfile(ctx, profile, 3600)

	h.logger.Info("driver profile created", zap.String("profileID", profile.ID), zap.String("userID", cmd.UserID))

	return profile.ID, nil
}
