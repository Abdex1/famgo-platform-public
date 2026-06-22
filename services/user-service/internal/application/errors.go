// services/user-service/internal/application/errors.go
// Application-level errors

package application

import "errors"

var (
	errInvalidPhoneNumber  = errors.New("invalid phone number")
	errInvalidEmail        = errors.New("invalid email")
	errUserAlreadyExists   = errors.New("user already exists")
	errCannotActivateUser  = errors.New("cannot activate user in current status")
	errCannotVerifyDriver  = errors.New("cannot verify driver in current status")
	errUserNotFound        = errors.New("user not found")
	errDriverProfileNotFound = errors.New("driver profile not found")
	errPassengerProfileNotFound = errors.New("passenger profile not found")
	errInternalError       = errors.New("internal error")
)
