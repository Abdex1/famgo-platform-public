
// 9. CREATE TOPIC OWNERSHIP GOVERNANCE

// packages/event-bus/governance/ownership.go

package governance

var TopicOwners = map[string]string{
	"ride.created.v1":            "ride-service",
	"ride.accepted.v1":           "dispatch-service",
	"driver.location.updated.v1": "gps-service",
	"payment.completed.v1":       "payment-service",
	"auth.login.succeeded.v1":    "auth-service",
}
