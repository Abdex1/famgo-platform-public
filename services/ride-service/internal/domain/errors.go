// services/ride-service/internal/domain/errors.go
// Domain-level errors

package domain

import "errors"

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrRideNotFound           = errors.New("ride not found")
	ErrRideAlreadyAssigned    = errors.New("ride already assigned")
	ErrRideAlreadyCancelled   = errors.New("ride already cancelled")
	ErrRideNotActive          = errors.New("ride not active")
	ErrInvalidLocation        = errors.New("invalid location")
	ErrMissingDriver          = errors.New("missing driver assignment")
)
