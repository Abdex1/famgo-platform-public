// services/user-service/internal/application/queries.go
// User Service Queries and Handlers

package application

import (
	"context"

	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
	"go.uber.org/zap"
)

// ===== GET USER QUERY =====

type GetUserQuery struct {
	UserID string
}

type GetUserHandler struct {
	userRepo  domain.UserRepository
	userCache domain.UserCache
	logger    *zap.Logger
}

func NewGetUserHandler(
	userRepo domain.UserRepository,
	userCache domain.UserCache,
	logger *zap.Logger,
) *GetUserHandler {
	return &GetUserHandler{
		userRepo:  userRepo,
		userCache: userCache,
		logger:    logger,
	}
}

func (h *GetUserHandler) Handle(ctx context.Context, q GetUserQuery) (*domain.User, error) {
	// Try cache first
	if user, err := h.userCache.GetUser(ctx, q.UserID); err == nil {
		h.logger.Debug("user found in cache", zap.String("userID", q.UserID))
		return user, nil
	}

	// Query from database
	user, err := h.userRepo.GetUser(ctx, q.UserID)
	if err != nil {
		h.logger.Error("user not found", zap.String("userID", q.UserID))
		return nil, errUserNotFound
	}

	// Cache result
	h.userCache.SetUser(ctx, user, 3600)

	return user, nil
}

// ===== GET DRIVER PROFILE QUERY =====

type GetDriverProfileQuery struct {
	ProfileID string
}

type GetDriverProfileHandler struct {
	driverRepo  domain.DriverProfileRepository
	driverCache domain.UserCache
	logger      *zap.Logger
}

func NewGetDriverProfileHandler(
	driverRepo domain.DriverProfileRepository,
	driverCache domain.UserCache,
	logger *zap.Logger,
) *GetDriverProfileHandler {
	return &GetDriverProfileHandler{
		driverRepo:  driverRepo,
		driverCache: driverCache,
		logger:      logger,
	}
}

func (h *GetDriverProfileHandler) Handle(ctx context.Context, q GetDriverProfileQuery) (*domain.DriverProfile, error) {
	// Try cache
	if profile, err := h.driverCache.GetDriverProfile(ctx, q.ProfileID); err == nil {
		h.logger.Debug("driver profile found in cache", zap.String("profileID", q.ProfileID))
		return profile, nil
	}

	// Query from database
	profile, err := h.driverRepo.GetProfile(ctx, q.ProfileID)
	if err != nil {
		h.logger.Error("driver profile not found", zap.String("profileID", q.ProfileID))
		return nil, errDriverProfileNotFound
	}

	// Cache result
	h.driverCache.SetDriverProfile(ctx, profile, 3600)

	return profile, nil
}

// ===== GET DRIVER PROFILE BY USER ID QUERY =====

type GetDriverProfileByUserIDQuery struct {
	UserID string
}

type GetDriverProfileByUserIDHandler struct {
	driverRepo  domain.DriverProfileRepository
	driverCache domain.UserCache
	logger      *zap.Logger
}

func NewGetDriverProfileByUserIDHandler(
	driverRepo domain.DriverProfileRepository,
	driverCache domain.UserCache,
	logger *zap.Logger,
) *GetDriverProfileByUserIDHandler {
	return &GetDriverProfileByUserIDHandler{
		driverRepo:  driverRepo,
		driverCache: driverCache,
		logger:      logger,
	}
}

func (h *GetDriverProfileByUserIDHandler) Handle(ctx context.Context, q GetDriverProfileByUserIDQuery) (*domain.DriverProfile, error) {
	// Query from database
	profile, err := h.driverRepo.GetByUserID(ctx, q.UserID)
	if err != nil {
		h.logger.Error("driver profile not found for user", zap.String("userID", q.UserID))
		return nil, errDriverProfileNotFound
	}

	// Cache result
	h.driverCache.SetDriverProfile(ctx, profile, 3600)

	return profile, nil
}

// ===== GET PASSENGER PROFILE QUERY =====

type GetPassengerProfileQuery struct {
	ProfileID string
}

type GetPassengerProfileHandler struct {
	passengerRepo   domain.PassengerProfileRepository
	passengerCache  domain.UserCache
	logger          *zap.Logger
}

func NewGetPassengerProfileHandler(
	passengerRepo domain.PassengerProfileRepository,
	passengerCache domain.UserCache,
	logger *zap.Logger,
) *GetPassengerProfileHandler {
	return &GetPassengerProfileHandler{
		passengerRepo:   passengerRepo,
		passengerCache:  passengerCache,
		logger:          logger,
	}
}

func (h *GetPassengerProfileHandler) Handle(ctx context.Context, q GetPassengerProfileQuery) (*domain.PassengerProfile, error) {
	// Try cache
	if profile, err := h.passengerCache.GetPassengerProfile(ctx, q.ProfileID); err == nil {
		h.logger.Debug("passenger profile found in cache", zap.String("profileID", q.ProfileID))
		return profile, nil
	}

	// Query from database
	profile, err := h.passengerRepo.GetProfile(ctx, q.ProfileID)
	if err != nil {
		h.logger.Error("passenger profile not found", zap.String("profileID", q.ProfileID))
		return nil, errPassengerProfileNotFound
	}

	// Cache result
	h.passengerCache.SetPassengerProfile(ctx, profile, 3600)

	return profile, nil
}

// ===== LIST ACTIVE DRIVERS QUERY =====

type ListActiveDriversQuery struct {
	Limit  int
	Offset int
}

type ListActiveDriversHandler struct {
	driverRepo domain.DriverProfileRepository
	logger     *zap.Logger
}

func NewListActiveDriversHandler(
	driverRepo domain.DriverProfileRepository,
	logger *zap.Logger,
) *ListActiveDriversHandler {
	return &ListActiveDriversHandler{
		driverRepo: driverRepo,
		logger:     logger,
	}
}

func (h *ListActiveDriversHandler) Handle(ctx context.Context, q ListActiveDriversQuery) ([]domain.DriverProfile, error) {
	drivers, err := h.driverRepo.ListActiveDrivers(ctx, q.Limit, q.Offset)
	if err != nil {
		h.logger.Error("failed to list active drivers", zap.Error(err))
		return nil, err
	}

	return drivers, nil
}
