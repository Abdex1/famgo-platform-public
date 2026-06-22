/*
# STEP 5 — CREATE DOMAIN ENTITIES

 USER ENTITY

services/auth-service/internal/domain/entities/user.go


# FILE: user.go
*/
package entities

import "time"

type UserRole string

const (
	RoleRider      UserRole = "rider"
	RoleDriver     UserRole = "driver"
	RoleSupport    UserRole = "support"
	RoleAdmin      UserRole = "admin"
	RoleOps        UserRole = "ops"
	RoleFraudAgent UserRole = "fraud-agent"
	RoleSuperAdmin UserRole = "super-admin"
)

type User struct {
	ID string

	Email string
	Phone string

	PasswordHash string

	Role UserRole

	EmailVerified bool
	PhoneVerified bool

	Status string

	CreatedAt time.Time
	UpdatedAt time.Time
}
