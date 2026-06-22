
// 16. CREATE EVENT REGISTRY

// packages/event-bus/registry/topics.go

package registry

var RegisteredTopics = []string{
	"ride.created.v1",
	"ride.accepted.v1",
	"ride.cancelled.v1",
	"driver.location.updated.v1",
	"payment.completed.v1",
	"auth.login.succeeded.v1",
	"dispatch.matching.started.v1",
	"dispatch.driver.matched.v1",
	"dispatch.driver.assigned.v1",
	"dispatch.matching.failed.v1",
	"dispatch.matching.expired.v1",
}
