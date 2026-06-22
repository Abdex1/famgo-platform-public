package topics

const (
	MatchingStarted = "dispatch.matching.started.v1"
	DriverMatched   = "dispatch.driver.matched.v1"
	DriverAssigned  = "dispatch.driver.assigned.v1"
	MatchingFailed  = "dispatch.matching.failed.v1"
	MatchingExpired = "dispatch.matching.expired.v1"
)

var DispatchTopics = []string{
	MatchingStarted,
	DriverMatched,
	DriverAssigned,
	MatchingFailed,
	MatchingExpired,
}
