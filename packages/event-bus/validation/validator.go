// 13. CREATE VALIDATION ENGINE

// packages/event-bus/validation/validator.go

package validation

import (
	"github.com/go-playground/validator/v10"

	"github.com/Abdex1/FamGo-platform/packages/event-bus/envelope"
)

var validate = validator.New()

func ValidateEnvelope(event envelope.EventEnvelope) error {
	return validate.Struct(event)
}
