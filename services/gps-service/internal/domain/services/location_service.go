// services/gps-service/internal/domain/services/location_service.go
// Location domain service for calculations and analysis

package services

import (
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

// LocationService provides location-related business logic
type LocationService struct{}

// NewLocationService creates a new location service
func NewLocationService() *LocationService {
	return &LocationService{}
}

// NearbyDrivers holds a driver with its distance
type NearbyDriver struct {
	Driver   *entities.DriverLocation
	Distance float64 // kilometers
	ETA      float64 // minutes
	Bearing  float64 // degrees (0-360)
}

// FindNearbyDrivers finds drivers within radius of a location
// Returns drivers sorted by distance, closest first
func (ls *LocationService) FindNearbyDrivers(
	referenceLocation *valueobjects.Geolocation,
	drivers []*entities.DriverLocation,
	radiusKm float64,
	limit int,
	baseSpeedKmH float64,
) []NearbyDriver {
	var nearby []NearbyDriver

	for _, driver := range drivers {
		if !driver.IsActive() {
			continue
		}

		distance := referenceLocation.DistanceToKm(driver.CurrentLocation)
		if distance <= radiusKm {
			eta := referenceLocation.EstimatedArrivalTime(driver.CurrentLocation, baseSpeedKmH)
			bearing := referenceLocation.BearingTo(driver.CurrentLocation)

			nearby = append(nearby, NearbyDriver{
				Driver:   driver,
				Distance: distance,
				ETA:      eta,
				Bearing:  bearing,
			})
		}
	}

	// Sort by distance (closest first)
	for i := 0; i < len(nearby)-1; i++ {
		for j := i + 1; j < len(nearby); j++ {
			if nearby[j].Distance < nearby[i].Distance {
				nearby[i], nearby[j] = nearby[j], nearby[i]
			}
		}
	}

	// Limit results
	if limit > 0 && len(nearby) > limit {
		return nearby[:limit]
	}

	return nearby
}

// CalculateRoute calculates route information between two locations
type RouteInfo struct {
	DistanceKm       float64
	DistanceMiles    float64
	ETAMinutes       float64
	Bearing          float64
	SpeedKmH         float64
	EstimatedArrival time.Time
}

func (ls *LocationService) CalculateRoute(
	from, to *valueobjects.Geolocation,
	baseSpeedKmH float64,
) *RouteInfo {
	distanceKm := from.DistanceToKm(to)
	distanceMiles := from.DistanceToMiles(to)
	etaMinutes := from.EstimatedArrivalTime(to, baseSpeedKmH)
	bearing := from.BearingTo(to)

	return &RouteInfo{
		DistanceKm:       distanceKm,
		DistanceMiles:    distanceMiles,
		ETAMinutes:       etaMinutes,
		Bearing:          bearing,
		SpeedKmH:         baseSpeedKmH,
		EstimatedArrival: time.Now().Add(time.Duration(etaMinutes) * time.Minute),
	}
}

// TrajectoryQuality holds trajectory analysis results
type TrajectoryQuality struct {
	IsValid         bool
	DistanceMoved   float64 // kilometers
	TimeElapsed     time.Duration
	CalculatedSpeed float64 // km/h
	ExpectedSpeed   float64 // km/h
	IsAnomaly       bool    // True if speed is unusually high
	Confidence      float64 // 0-1
}

// AnalyzeTrajectory analyzes movement between two locations
func (ls *LocationService) AnalyzeTrajectory(
	from, to *valueobjects.Geolocation,
	maxSpeedKmH float64,
) *TrajectoryQuality {
	if from == nil || to == nil {
		return &TrajectoryQuality{IsValid: false, Confidence: 0}
	}

	distanceKm := from.DistanceToKm(to)
	if distanceKm < 0.001 { // Less than 1 meter
		return &TrajectoryQuality{
			IsValid:         true,
			DistanceMoved:   distanceKm,
			TimeElapsed:     0,
			CalculatedSpeed: 0,
			ExpectedSpeed:   0,
			IsAnomaly:       false,
			Confidence:      1.0,
		}
	}

	timeElapsed := time.Duration(to.Timestamp-from.Timestamp) * time.Millisecond
	if timeElapsed <= 0 {
		return &TrajectoryQuality{
			IsValid:     false,
			Confidence:  0,
			IsAnomaly:   true,
		}
	}

	calculatedSpeedKmH := (distanceKm / timeElapsed.Hours())
	isAnomaly := calculatedSpeedKmH > maxSpeedKmH

	confidence := 1.0
	if isAnomaly {
		confidence = 0.1
	}

	return &TrajectoryQuality{
		IsValid:         true,
		DistanceMoved:   distanceKm,
		TimeElapsed:     timeElapsed,
		CalculatedSpeed: calculatedSpeedKmH,
		ExpectedSpeed:   maxSpeedKmH,
		IsAnomaly:       isAnomaly,
		Confidence:      confidence,
	}
}

// LocationCluster groups drivers by proximity
type LocationCluster struct {
	CenterLat float64
	CenterLng float64
	Drivers   []*entities.DriverLocation
	Count     int
	Radius    float64 // kilometers
}

// ClusterNearbyDrivers clusters drivers by proximity
func (ls *LocationService) ClusterNearbyDrivers(
	drivers []*entities.DriverLocation,
	clusterRadiusKm float64,
) []LocationCluster {
	if len(drivers) == 0 {
		return nil
	}

	visited := make(map[string]bool)
	var clusters []LocationCluster

	for i, driver := range drivers {
		if visited[driver.DriverID] {
			continue
		}

		cluster := LocationCluster{
			CenterLat: driver.CurrentLocation.Coordinates.Latitude,
			CenterLng: driver.CurrentLocation.Coordinates.Longitude,
			Radius:    clusterRadiusKm,
		}

		for j, other := range drivers {
			if !visited[other.DriverID] && driver.IsWithinRadius(other, clusterRadiusKm) {
				cluster.Drivers = append(cluster.Drivers, other)
				if i != j {
					visited[other.DriverID] = true
				}
			}
		}

		cluster.Count = len(cluster.Drivers)
		clusters = append(clusters, cluster)
	}

	return clusters
}

// UpdateLocationQuality checks location quality and returns any issues
func (ls *LocationService) UpdateLocationQuality(
	location *valueobjects.Geolocation,
	minAccuracyMeters float64,
) (bool, error) {
	if location == nil {
		return false, fmt.Errorf("location is nil")
	}

	if !location.IsValid() {
		return false, fmt.Errorf("location is invalid")
	}

	if location.Accuracy > minAccuracyMeters {
		return false, fmt.Errorf("location accuracy %f exceeds threshold %f meters", location.Accuracy, minAccuracyMeters)
	}

	return true, nil
}

// InterpolateTrajectory interpolates driver position between updates
func (ls *LocationService) InterpolateTrajectory(
	from, to *valueobjects.Geolocation,
	currentTime time.Time,
) *valueobjects.Geolocation {
	if from == nil || to == nil {
		return from
	}

	// Calculate fraction of time elapsed
	startTime := time.UnixMilli(from.Timestamp)
	endTime := time.UnixMilli(to.Timestamp)

	if !currentTime.After(startTime) {
		return from
	}
	if !currentTime.Before(endTime) {
		return to
	}

	totalDuration := endTime.Sub(startTime)
	elapsedDuration := currentTime.Sub(startTime)
	fraction := elapsedDuration.Seconds() / totalDuration.Seconds()

	return from.Interpolate(to, fraction)
}

// CalculateGeohash calculates geohash for spatial indexing
func (ls *LocationService) CalculateGeohash(
	lat, lng float64,
	precision int,
) (string, error) {
	coords, err := valueobjects.NewCoordinates(lat, lng)
	if err != nil {
		return "", err
	}

	// Simple geohash implementation
	// In production, use a proper geohashing library
	if precision < 1 || precision > 12 {
		precision = 8
	}

	hash := encodeGeohash(lat, lng, precision)
	return hash, nil
}

// Simple geohash encoding (base implementation)
func encodeGeohash(lat, lng float64, precision int) string {
	const (
		base32Alphabet = "0123456789bcdefghjkmnpqrstuvwxyz"
	)

	var hash string
	var minLat, maxLat = -90.0, 90.0
	var minLng, maxLng = -180.0, 180.0

	bit := 0
	ch := 0

	for len(hash) < precision {
		if len(hash)%2 == 0 {
			// Even bits: longitude
			mid := (minLng + maxLng) / 2
			if lng >= mid {
				ch = (ch << 1) | 1
				minLng = mid
			} else {
				ch = ch << 1
				maxLng = mid
			}
		} else {
			// Odd bits: latitude
			mid := (minLat + maxLat) / 2
			if lat >= mid {
				ch = (ch << 1) | 1
				minLat = mid
			} else {
				ch = ch << 1
				maxLat = mid
			}
		}

		bit++
		if bit == 5 {
			hash += string(base32Alphabet[ch])
			bit = 0
			ch = 0
		}
	}

	return hash
}

// ValidateLocationSequence validates a sequence of locations for consistency
func (ls *LocationService) ValidateLocationSequence(
	locations []*valueobjects.Geolocation,
	maxSpeedKmH float64,
) (bool, string) {
	if len(locations) < 2 {
		return true, ""
	}

	for i := 1; i < len(locations); i++ {
		trajectory := ls.AnalyzeTrajectory(locations[i-1], locations[i], maxSpeedKmH)
		if !trajectory.IsValid || trajectory.IsAnomaly {
			return false, fmt.Sprintf("anomaly detected between point %d and %d: speed %.2f km/h exceeds max %.2f km/h",
				i-1, i, trajectory.CalculatedSpeed, maxSpeedKmH)
		}
	}

	return true, ""
}
