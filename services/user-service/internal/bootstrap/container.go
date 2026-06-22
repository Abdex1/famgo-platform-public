// services/user-service/internal/bootstrap/container.go
// Dependency Injection Container

package bootstrap

import (
	"database/sql"
	"log"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/infrastructure"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/transport"
)

// Container holds all dependencies
type Container struct {
	DB     *sql.DB
	Redis  *redis.Client
	Logger *zap.Logger

	// Repositories
	UserRepo             domain.UserRepository
	DriverProfileRepo    domain.DriverProfileRepository
	PassengerProfileRepo domain.PassengerProfileRepository

	// Cache
	UserCache domain.UserCache

	// Domain Services
	UserService *domain.UserService

	// Application Handlers
	RegisterUserHandler         *application.RegisterUserHandler
	UpdateProfileHandler        *application.UpdateProfileHandler
	ActivateUserHandler         *application.ActivateUserHandler
	VerifyDriverHandler         *application.VerifyDriverHandler
	CreateDriverProfileHandler  *application.CreateDriverProfileHandler
	GetUserHandler              *application.GetUserHandler
	GetDriverProfileHandler     *application.GetDriverProfileHandler

	// Transport
	HTTPHandler *transport.HTTPHandler
}

// NewContainer creates and initializes all dependencies
func NewContainer(db *sql.DB, redisClient *redis.Client, logger *zap.Logger) *Container {
	c := &Container{
		DB:     db,
		Redis:  redisClient,
		Logger: logger,
	}

	// Initialize repositories
	c.UserRepo = infrastructure.NewPostgresUserRepository(db)
	c.DriverProfileRepo = infrastructure.NewPostgresDriverProfileRepository(db)
	c.PassengerProfileRepo = infrastructure.NewPostgresPassengerProfileRepository(db)

	// Initialize cache
	c.UserCache = infrastructure.NewRedisUserCache(redisClient, 3600)

	// Initialize domain services
	c.UserService = domain.NewUserService()

	// Initialize command handlers
	c.RegisterUserHandler = application.NewRegisterUserHandler(
		c.UserRepo,
		c.UserCache,
		c.UserService,
		logger,
	)

	c.UpdateProfileHandler = application.NewUpdateProfileHandler(
		c.UserRepo,
		c.UserCache,
		c.UserService,
		logger,
	)

	c.ActivateUserHandler = application.NewActivateUserHandler(
		c.UserRepo,
		c.UserCache,
		c.UserService,
		logger,
	)

	c.VerifyDriverHandler = application.NewVerifyDriverHandler(
		c.DriverProfileRepo,
		c.UserCache,
		c.UserService,
		logger,
	)

	c.CreateDriverProfileHandler = application.NewCreateDriverProfileHandler(
		c.DriverProfileRepo,
		c.UserCache,
		logger,
	)

	// Initialize query handlers
	c.GetUserHandler = application.NewGetUserHandler(
		c.UserRepo,
		c.UserCache,
		logger,
	)

	c.GetDriverProfileHandler = application.NewGetDriverProfileHandler(
		c.DriverProfileRepo,
		c.UserCache,
		logger,
	)

	// Initialize transport
	c.HTTPHandler = transport.NewHTTPHandler(
		c.RegisterUserHandler,
		c.UpdateProfileHandler,
		c.ActivateUserHandler,
		c.VerifyDriverHandler,
		c.CreateDriverProfileHandler,
		c.GetUserHandler,
		c.GetDriverProfileHandler,
		logger,
	)

	logger.Info("dependency injection container initialized")

	return c
}
