// services/ride-service/internal/application/errors.go
// Application Errors

package application

import "errors"

var (
	errInvalidLocation         = errors.New("invalid location coordinates")
	errInvalidStateTransition  = errors.New("invalid state transition")
	errRideNotFound            = errors.New("ride not found")
	errInternalError           = errors.New("internal error")
)
