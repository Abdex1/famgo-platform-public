package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
)

// DriverCandidateEngine discovers and normalizes driver candidates.
type DriverCandidateEngine struct {
	discovery ports.DriverDiscovery
}

func NewDriverCandidateEngine(discovery ports.DriverDiscovery) *DriverCandidateEngine {
	return &DriverCandidateEngine{discovery: discovery}
}

// FindWithinRadius returns online drivers within the search radius.
func (e *DriverCandidateEngine) FindWithinRadius(
	ctx context.Context,
	latitude, longitude, radiusKm float64,
	limit int,
) ([]ports.DriverCandidate, error) {
	if radiusKm <= 0 {
		return nil, fmt.Errorf("search radius must be positive")
	}
	if limit <= 0 {
		limit = 50
	}

	candidates, err := e.discovery.FindDriversWithinRadius(ctx, latitude, longitude, radiusKm, limit)
	if err != nil {
		return nil, fmt.Errorf("driver discovery failed: %w", err)
	}

	filtered := make([]ports.DriverCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.IsOnline {
			filtered = append(filtered, candidate)
		}
	}

	return filtered, nil
}

// SortByDistance orders candidates nearest-first (Sprint 1 baseline algorithm).
func (e *DriverCandidateEngine) SortByDistance(candidates []ports.DriverCandidate) []ports.DriverCandidate {
	sorted := append([]ports.DriverCandidate(nil), candidates...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].DistanceKm == sorted[j].DistanceKm {
			return sorted[i].ETAMinutes < sorted[j].ETAMinutes
		}
		return sorted[i].DistanceKm < sorted[j].DistanceKm
	})
	return sorted
}
