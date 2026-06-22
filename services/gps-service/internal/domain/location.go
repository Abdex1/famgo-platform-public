package domain

import "time"

// LocationUpdate represents GPS location data
type LocationUpdate struct {
	ID        string    `json:"id"`
	DriverID  string    `json:"driver_id"`
	RideID    string    `json:"ride_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Heading   float64   `json:"heading"`
	Speed     float64   `json:"speed"`
	Accuracy  float64   `json:"accuracy"`
	Timestamp time.Time `json:"timestamp"`
}

// DriverLocation represents driver's current location
type DriverLocation struct {
	DriverID  string    `json:"driver_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Heading   float64   `json:"heading"`
	Speed     float64   `json:"speed"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NearbyDriver represents driver in proximity search
type NearbyDriver struct {
	DriverID  string  `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"` // meters
}
