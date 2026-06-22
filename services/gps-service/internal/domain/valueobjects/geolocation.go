// services/gps-service/internal/domain/valueobjects/geolocation.go
// Geolocation value object with distance and bearing calculations

package valueobjects

import (
	"fmt"
	"math"
)

const (
	// EarthRadiusKm is Earth's radius in kilometers
	EarthRadiusKm = 6371.0
	// EarthRadiusMiles is Earth's radius in miles
	EarthRadiusMiles = 3959.0
)

// Coordinates represents latitude and longitude
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NewCoordinates creates a new coordinates value object with validation
func NewCoordinates(lat, lng float64) (*Coordinates, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90, got %f", lat)
	}
	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180, got %f", lng)
	}
	return &Coordinates{
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

// IsValid checks if coordinates are valid
func (c *Coordinates) IsValid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 &&
		c.Longitude >= -180 && c.Longitude <= 180
}

// String returns string representation
func (c *Coordinates) String() string {
	return fmt.Sprintf("%.8f,%.8f", c.Latitude, c.Longitude)
}

// Geolocation represents a location with coordinates, altitude, and accuracy
type Geolocation struct {
	Coordinates *Coordinates `json:"coordinates"`
	Altitude    float64      `json:"altitude"`      // meters
	Accuracy    float64      `json:"accuracy"`      // meters
	Speed       float64      `json:"speed"`         // m/s
	Heading     float64      `json:"heading"`       // degrees (0-360)
	Timestamp   int64        `json:"timestamp"`     // Unix timestamp in milliseconds
}

// NewGeolocation creates a new geolocation value object
func NewGeolocation(lat, lng, altitude, accuracy, speed, heading float64, timestamp int64) (*Geolocation, error) {
	coords, err := NewCoordinates(lat, lng)
	if err != nil {
		return nil, err
	}

	if accuracy < 0 {
		return nil, fmt.Errorf("accuracy must be non-negative, got %f", accuracy)
	}
	if speed < 0 {
		return nil, fmt.Errorf("speed must be non-negative, got %f", speed)
	}
	if heading < 0 || heading > 360 {
		return nil, fmt.Errorf("heading must be between 0 and 360, got %f", heading)
	}
	if timestamp <= 0 {
		return nil, fmt.Errorf("timestamp must be positive, got %d", timestamp)
	}

	return &Geolocation{
		Coordinates: coords,
		Altitude:    altitude,
		Accuracy:    accuracy,
		Speed:       speed,
		Heading:     heading,
		Timestamp:   timestamp,
	}, nil
}

// IsValid checks if geolocation is valid
func (g *Geolocation) IsValid() bool {
	return g.Coordinates != nil &&
		g.Coordinates.IsValid() &&
		g.Accuracy >= 0 &&
		g.Speed >= 0 &&
		g.Heading >= 0 && g.Heading <= 360 &&
		g.Timestamp > 0
}

// DistanceToKm calculates distance to another location in kilometers using Haversine formula
func (g *Geolocation) DistanceToKm(other *Geolocation) float64 {
	return haversine(
		g.Coordinates.Latitude,
		g.Coordinates.Longitude,
		other.Coordinates.Latitude,
		other.Coordinates.Longitude,
		EarthRadiusKm,
	)
}

// DistanceToMiles calculates distance to another location in miles
func (g *Geolocation) DistanceToMiles(other *Geolocation) float64 {
	return haversine(
		g.Coordinates.Latitude,
		g.Coordinates.Longitude,
		other.Coordinates.Latitude,
		other.Coordinates.Longitude,
		EarthRadiusMiles,
	)
}

// BearingTo calculates bearing to another location in degrees (0-360)
func (g *Geolocation) BearingTo(other *Geolocation) float64 {
	lat1 := degreesToRadians(g.Coordinates.Latitude)
	lat2 := degreesToRadians(other.Coordinates.Latitude)
	dLng := degreesToRadians(other.Coordinates.Longitude - g.Coordinates.Longitude)

	y := math.Sin(dLng) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLng)

	bearing := radiansToDegrees(math.Atan2(y, x))
	// Normalize to 0-360
	if bearing < 0 {
		bearing += 360
	}
	return math.Mod(bearing, 360)
}

// EstimatedArrivalTime calculates ETA to another location in minutes
// baseSpeed is in km/h
func (g *Geolocation) EstimatedArrivalTime(other *Geolocation, baseSpeedKmH float64) float64 {
	distanceKm := g.DistanceToKm(other)
	if baseSpeedKmH <= 0 {
		return 0
	}
	return (distanceKm / baseSpeedKmH) * 60 // Convert to minutes
}

// IsWithinRadius checks if another location is within a radius in kilometers
func (g *Geolocation) IsWithinRadius(other *Geolocation, radiusKm float64) bool {
	distance := g.DistanceToKm(other)
	return distance <= radiusKm
}

// Interpolate calculates an intermediate point between two locations at time fraction
// fraction should be 0-1 (0 = g, 1 = other)
func (g *Geolocation) Interpolate(other *Geolocation, fraction float64) *Geolocation {
	if fraction < 0 || fraction > 1 {
		return g
	}

	lat1 := degreesToRadians(g.Coordinates.Latitude)
	lng1 := degreesToRadians(g.Coordinates.Longitude)
	lat2 := degreesToRadians(other.Coordinates.Latitude)
	lng2 := degreesToRadians(other.Coordinates.Longitude)

	a := math.Sin((1-fraction)*angularDistance(lat1, lng1, lat2, lng2)) / math.Sin(angularDistance(lat1, lng1, lat2, lng2))
	b := math.Sin(fraction*angularDistance(lat1, lng1, lat2, lng2)) / math.Sin(angularDistance(lat1, lng1, lat2, lng2))

	x := a*math.Cos(lat1)*math.Cos(lng1) + b*math.Cos(lat2)*math.Cos(lng2)
	y := a*math.Cos(lat1)*math.Sin(lng1) + b*math.Cos(lat2)*math.Sin(lng2)
	z := a*math.Sin(lat1) + b*math.Sin(lat2)

	newLat := radiansToDegrees(math.Atan2(z, math.Sqrt(x*x+y*y)))
	newLng := radiansToDegrees(math.Atan2(y, x))

	// Interpolate altitude and accuracy linearly
	newAltitude := g.Altitude + (other.Altitude-g.Altitude)*fraction
	newAccuracy := g.Accuracy + (other.Accuracy-g.Accuracy)*fraction

	result, _ := NewGeolocation(
		newLat, newLng,
		newAltitude,
		newAccuracy,
		g.Speed,
		g.Heading,
		int64(float64(g.Timestamp)+float64(other.Timestamp-g.Timestamp)*fraction),
	)
	return result
}

// Helper functions

// haversine calculates the great-circle distance between two points on a sphere
func haversine(lat1, lng1, lat2, lng2, radius float64) float64 {
	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)
	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Asin(math.Sqrt(a))
	return radius * c
}

// angularDistance calculates angular distance between two points
func angularDistance(lat1, lng1, lat2, lng2 float64) float64 {
	return 2 * math.Asin(math.Sqrt(
		math.Sin((lat2-lat1)/2)*math.Sin((lat2-lat1)/2)+
			math.Cos(lat1)*math.Cos(lat2)*
				math.Sin((lng2-lng1)/2)*math.Sin((lng2-lng1)/2),
	))
}

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// radiansToDegrees converts radians to degrees
func radiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}
