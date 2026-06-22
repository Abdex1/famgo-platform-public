package valueobjects

import (
	"fmt"
	"time"
)

// Money represents a monetary amount in ETB (Ethiopian Birr)
type Money struct {
	Amount   float64
	Currency string
}

func NewMoney(amount float64) *Money {
	if amount < 0 {
		amount = 0
	}
	return &Money{
		Amount:   amount,
		Currency: "ETB",
	}
}

// Add returns a new Money value with added amounts
func (m *Money) Add(other *Money) *Money {
	return NewMoney(m.Amount + other.Amount)
}

// Multiply returns a new Money value multiplied by factor
func (m *Money) Multiply(factor float64) *Money {
	return NewMoney(m.Amount * factor)
}

// RideLocation represents pickup or dropoff location
type RideLocation struct {
	Latitude  float64
	Longitude float64
	Address   string
	Timestamp time.Time
}

func NewRideLocation(lat, lng float64, address string) (*RideLocation, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("invalid latitude: %f", lat)
	}
	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("invalid longitude: %f", lng)
	}

	return &RideLocation{
		Latitude:  lat,
		Longitude: lng,
		Address:   address,
		Timestamp: time.Now().UTC(),
	}, nil
}

// RideRoute represents the complete trip route
type RideRoute struct {
	PickupLocation    *RideLocation
	DropoffLocation   *RideLocation
	EstimatedDistance float64 // km
	EstimatedDuration int32   // seconds
	ActualDistance    float64 // km
	ActualDuration    int32   // seconds
}

func NewRideRoute(
	pickupLat, pickupLng float64,
	pickupAddr string,
	dropoffLat, dropoffLng float64,
	dropoffAddr string,
	estimatedDist float64,
	estimatedDur int32,
) (*RideRoute, error) {
	pickup, err := NewRideLocation(pickupLat, pickupLng, pickupAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid pickup location: %w", err)
	}

	dropoff, err := NewRideLocation(dropoffLat, dropoffLng, dropoffAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid dropoff location: %w", err)
	}

	return &RideRoute{
		PickupLocation:    pickup,
		DropoffLocation:   dropoff,
		EstimatedDistance: estimatedDist,
		EstimatedDuration: estimatedDur,
	}, nil
}

// FareSummary represents complete fare breakdown
type FareSummary struct {
	BaseFare        *Money
	DistanceFare    *Money
	TimeFare        *Money
	SurgeMultiplier float64
	SubTotal        *Money
	Total           *Money
}

func NewFareSummary(
	baseFare, distanceFare, timeFare *Money,
	surgeMultiplier float64,
) *FareSummary {
	subtotal := baseFare.Add(distanceFare).Add(timeFare)
	total := subtotal.Multiply(surgeMultiplier)

	return &FareSummary{
		BaseFare:        baseFare,
		DistanceFare:    distanceFare,
		TimeFare:        timeFare,
		SurgeMultiplier: surgeMultiplier,
		SubTotal:        subtotal,
		Total:           total,
	}
}

// RideRating represents driver or rider rating (1-5 stars with comment)
type RideRating struct {
	Stars      int32
	Comment    string
	CreatedAt  time.Time
}

func NewRideRating(stars int32, comment string) (*RideRating, error) {
	if stars < 1 || stars > 5 {
		return nil, fmt.Errorf("rating must be 1-5 stars, got %d", stars)
	}

	return &RideRating{
		Stars:     stars,
		Comment:   comment,
		CreatedAt: time.Now().UTC(),
	}, nil
}
