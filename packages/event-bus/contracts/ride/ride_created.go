
// 15. CREATE CONTRACTS — RIDE DOMAIN

// packages/event-bus/contracts/ride/ride_created.go

package ride

type RideCreated struct {
	RideID        string  `json:"ride_id"`
	RiderID       string  `json:"rider_id"`

	PickupLat     float64 `json:"pickup_lat"`
	PickupLng     float64 `json:"pickup_lng"`

	DropoffLat    float64 `json:"dropoff_lat"`
	DropoffLng    float64 `json:"dropoff_lng"`

	VehicleType   string  `json:"vehicle_type"`

	RequestedAt   string  `json:"requested_at"`
}
