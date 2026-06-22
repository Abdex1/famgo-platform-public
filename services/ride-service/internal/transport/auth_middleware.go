// services/ride-service/internal/transport/auth_middleware.go
// Authentication & Authorization Middleware

package transport

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/auth-client"
)

// AuthContext holds authenticated user information
type AuthContext struct {
	UserID string
	Email  string
	Roles  []string
}

// AuthMiddleware handles JWT validation and authorization
type AuthMiddleware struct {
	authClient authclient.AuthClient
	logger     *zap.Logger
	secLogger  *StructuredLogger
}

func NewAuthMiddleware(authClient authclient.AuthClient, logger *zap.Logger, secLogger *StructuredLogger) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
		logger:     logger,
		secLogger:  secLogger,
	}
}

// ValidateToken validates JWT token and extracts user information
func (am *AuthMiddleware) ValidateToken(ctx context.Context, token string) (*AuthContext, error) {
	start := time.Now()

	if token == "" {
		am.secLogger.LogSecurityEvent(
			ctx,
			"AUTH",
			"unknown",
			"ride-service",
			"validate_token",
			"DENY",
			map[string]string{"reason": "missing_token"},
		)
		return nil, fmt.Errorf("missing authorization token")
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// Validate token using auth service
	userInfo, err := am.authClient.ValidateToken(ctx, token)
	if err != nil {
		am.secLogger.LogSecurityEvent(
			ctx,
			"AUTH",
			"unknown",
			"ride-service",
			"validate_token",
			"DENY",
			map[string]string{"reason": "invalid_token", "error": err.Error()},
		)
		am.logger.Warn("token validation failed", zap.Error(err), zap.Duration("duration_ms", time.Since(start)))
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	authCtx := &AuthContext{
		UserID: userInfo.UserID,
		Email:  userInfo.Email,
		Roles:  userInfo.Roles,
	}

	am.secLogger.LogSecurityEvent(
		ctx,
		"AUTH",
		authCtx.UserID,
		"ride-service",
		"validate_token",
		"ALLOW",
		map[string]string{"roles": strings.Join(authCtx.Roles, ",")},
	)

	am.logger.Debug("token validated successfully",
		zap.String("user_id", authCtx.UserID),
		zap.Strings("roles", authCtx.Roles),
		zap.Duration("duration_ms", time.Since(start)))

	return authCtx, nil
}

// CheckAuthorization checks if user has required roles for an operation
func (am *AuthMiddleware) CheckAuthorization(
	ctx context.Context,
	authCtx *AuthContext,
	resource string,
	action string,
	requiredRoles []string,
) (bool, error) {
	// Check if user has any of the required roles
	hasRole := false
	for _, userRole := range authCtx.Roles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				hasRole = true
				break
			}
		}
		if hasRole {
			break
		}
	}

	if !hasRole {
		am.secLogger.LogSecurityEvent(
			ctx,
			"AUTHZ",
			authCtx.UserID,
			resource,
			action,
			"DENY",
			map[string]string{
				"required_roles": strings.Join(requiredRoles, ","),
				"user_roles":     strings.Join(authCtx.Roles, ","),
			},
		)
		return false, fmt.Errorf("insufficient permissions")
	}

	am.secLogger.LogSecurityEvent(
		ctx,
		"AUTHZ",
		authCtx.UserID,
		resource,
		action,
		"ALLOW",
		map[string]string{
			"user_roles": strings.Join(authCtx.Roles, ","),
		},
	)

	return true, nil
}

// ==================== RBAC RULES ====================

// RBACRules defines role-based access control for operations
var RBACRules = map[string]map[string][]string{
	// Resource: Ride Operations
	"rides:create": {
		"POST /rides": {"PASSENGER", "ADMIN"},
	},
	"rides:read": {
		"GET /rides/{rideID}": {"PASSENGER", "DRIVER", "ADMIN"},
	},
	"rides:list": {
		"GET /passengers/{passengerID}/rides": {"PASSENGER", "ADMIN"},
		"GET /drivers/{driverID}/rides":       {"DRIVER", "ADMIN"},
	},
	"rides:assign": {
		"POST /rides/{rideID}/assign": {"DISPATCHER", "ADMIN"},
	},
	"rides:start": {
		"POST /rides/{rideID}/start": {"DRIVER", "ADMIN"},
	},
	"rides:complete": {
		"POST /rides/{rideID}/complete": {"DRIVER", "ADMIN"},
	},
	"rides:cancel": {
		"POST /rides/{rideID}/cancel": {"PASSENGER", "DRIVER", "ADMIN"},
	},
}

// GetRequiredRoles returns required roles for an endpoint
func GetRequiredRoles(resource, endpoint string) []string {
	if rules, ok := RBACRules[resource]; ok {
		if roles, ok := rules[endpoint]; ok {
			return roles
		}
	}
	return []string{} // Default: no specific role required (public endpoint)
}

// ==================== INPUT VALIDATION ====================

// InputValidator validates request input
type InputValidator struct {
	logger *StructuredLogger
}

func NewInputValidator(logger *StructuredLogger) *InputValidator {
	return &InputValidator{logger: logger}
}

// ValidateRideCreation validates CreateRide request
func (iv *InputValidator) ValidateRideCreation(ctx context.Context, passengerID string, pickupLat, pickupLon, dropoffLat, dropoffLon float64) error {
	// Validate passenger ID
	if passengerID == "" {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", "unknown", "ride", "create", "INVALID", map[string]string{
			"field":  "passenger_id",
			"reason": "empty",
		})
		return fmt.Errorf("passenger_id cannot be empty")
	}

	// Validate coordinates (valid latitude/longitude ranges)
	if !isValidLatitude(pickupLat) || !isValidLongitude(pickupLon) {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", passengerID, "ride", "create", "INVALID", map[string]string{
			"field":  "pickup_location",
			"reason": "invalid_coordinates",
		})
		return fmt.Errorf("invalid pickup coordinates")
	}

	if !isValidLatitude(dropoffLat) || !isValidLongitude(dropoffLon) {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", passengerID, "ride", "create", "INVALID", map[string]string{
			"field":  "dropoff_location",
			"reason": "invalid_coordinates",
		})
		return fmt.Errorf("invalid dropoff coordinates")
	}

	// Validate locations are different
	if pickupLat == dropoffLat && pickupLon == dropoffLon {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", passengerID, "ride", "create", "INVALID", map[string]string{
			"field":  "locations",
			"reason": "pickup_and_dropoff_same",
		})
		return fmt.Errorf("pickup and dropoff locations must be different")
	}

	return nil
}

// ValidateFareAmount validates fare amount is reasonable
func (iv *InputValidator) ValidateFareAmount(ctx context.Context, fare float32) error {
	const minFare = 1.0
	const maxFare = 10000.0

	if fare < minFare {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", "system", "fare", "validate", "INVALID", map[string]string{
			"field":  "fare",
			"reason": "below_minimum",
		})
		return fmt.Errorf("fare below minimum: %.2f", minFare)
	}

	if fare > maxFare {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", "system", "fare", "validate", "INVALID", map[string]string{
			"field":  "fare",
			"reason": "above_maximum",
		})
		return fmt.Errorf("fare above maximum: %.2f", maxFare)
	}

	return nil
}

// ValidateRideID validates ride ID format
func (iv *InputValidator) ValidateRideID(ctx context.Context, rideID string) error {
	if rideID == "" {
		return fmt.Errorf("ride_id cannot be empty")
	}

	if len(rideID) > 100 {
		iv.logger.LogSecurityEvent(ctx, "VALIDATION", "unknown", "ride", "validate_id", "INVALID", map[string]string{
			"field":  "ride_id",
			"reason": "too_long",
		})
		return fmt.Errorf("ride_id too long")
	}

	return nil
}

// Helper functions

func isValidLatitude(lat float64) bool {
	return lat >= -90.0 && lat <= 90.0
}

func isValidLongitude(lon float64) bool {
	return lon >= -180.0 && lon <= 180.0
}

// ==================== AUDIT LOGGING ====================

// AuditLogger logs security-sensitive operations
type AuditLogger struct {
	logger *zap.Logger
}

func NewAuditLogger(logger *zap.Logger) *AuditLogger {
	return &AuditLogger{logger: logger}
}

// LogAuditEvent logs an audit event
func (al *AuditLogger) LogAuditEvent(
	ctx context.Context,
	action string, // CREATE, UPDATE, DELETE, etc.
	resource string,
	resourceID string,
	userID string,
	status string, // SUCCESS, FAILURE
	details map[string]string,
) {
	fields := []zap.Field{
		zap.String("audit_event", "true"),
		zap.String("action", action),
		zap.String("resource", resource),
		zap.String("resource_id", resourceID),
		zap.String("user_id", userID),
		zap.String("status", status),
		zap.Time("timestamp", time.Now().UTC()),
	}

	for key, value := range details {
		fields = append(fields, zap.String(key, value))
	}

	al.logger.Info("AUDIT", fields...)
}

// LogRideCreation logs ride creation
func (al *AuditLogger) LogRideCreation(ctx context.Context, rideID, userID string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILURE"
	}

	al.LogAuditEvent(ctx, "CREATE", "ride", rideID, userID, status, map[string]string{
		"operation": "create_ride",
	})
}

// LogRideCompletion logs ride completion
func (al *AuditLogger) LogRideCompletion(ctx context.Context, rideID, userID string, fare float32, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILURE"
	}

	al.LogAuditEvent(ctx, "UPDATE", "ride", rideID, userID, status, map[string]string{
		"operation": "complete_ride",
		"fare":      fmt.Sprintf("%.2f", fare),
	})
}

// LogRideCancellation logs ride cancellation
func (al *AuditLogger) LogRideCancellation(ctx context.Context, rideID, userID, reason string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILURE"
	}

	al.LogAuditEvent(ctx, "UPDATE", "ride", rideID, userID, status, map[string]string{
		"operation": "cancel_ride",
		"reason":    reason,
	})
}
