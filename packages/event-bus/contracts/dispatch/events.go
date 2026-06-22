package dispatch

type MatchingStarted struct {
	DispatchRequestID string  `json:"dispatch_request_id"`
	RideID            string  `json:"ride_id"`
	RiderID           string  `json:"rider_id"`
	SearchRadiusKm    float64 `json:"search_radius_km"`
}

type DriverMatched struct {
	DispatchRequestID string   `json:"dispatch_request_id"`
	RideID            string   `json:"ride_id"`
	DriverID          string   `json:"driver_id"`
	ProposedDrivers   []string `json:"proposed_drivers"`
	MatchScore        float64  `json:"match_score"`
}

type DriverAssigned struct {
	DispatchRequestID string `json:"dispatch_request_id"`
	RideID            string `json:"ride_id"`
	DriverID          string `json:"driver_id"`
}

type MatchingFailed struct {
	DispatchRequestID string `json:"dispatch_request_id"`
	RideID            string `json:"ride_id"`
	Reason            string `json:"reason"`
	AttemptCount      int    `json:"attempt_count"`
}

type MatchingExpired struct {
	DispatchRequestID string `json:"dispatch_request_id"`
	RideID            string `json:"ride_id"`
}
