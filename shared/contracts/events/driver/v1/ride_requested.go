
//shared/contracts/events/ride/v1/ride_requested.go

package v1

import "time"

type RideRequested struct {
    TripID         string    `json:"trip_id"`
    RiderID        string    `json:"rider_id"`
    PickupLat      float64   `json:"pickup_lat"`
    PickupLng      float64   `json:"pickup_lng"`
    DestinationLat float64   `json:"destination_lat"`
    DestinationLng float64   `json:"destination_lng"`
    VehicleType    string    `json:"vehicle_type"`
    RequestedAt    time.Time `json:"requested_at"`
}
