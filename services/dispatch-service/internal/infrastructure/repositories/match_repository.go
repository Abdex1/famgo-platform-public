package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
)

// MatchRepository handles match result persistence
type MatchRepository struct {
	pool *pgxpool.Pool
}

func NewMatchRepository(pool *pgxpool.Pool) *MatchRepository {
	return &MatchRepository{pool: pool}
}

// Create inserts a match result
func (repo *MatchRepository) Create(ctx context.Context, result *entities.MatchResult) error {
	query := `
		INSERT INTO match_results (match_request_id, ride_id, selected_driver_id, score, eta, distance, confidence, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := repo.pool.Exec(ctx, query,
		result.MatchRequestID,
		result.RideID,
		result.SelectedDriverID,
		result.Score,
		result.ETA,
		result.Distance,
		result.Confidence,
		result.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create match result: %w", err)
	}

	return nil
}

// GetByRideID retrieves the match result for a ride
func (repo *MatchRepository) GetByRideID(ctx context.Context, rideID string) (*entities.MatchResult, error) {
	query := `
		SELECT match_request_id, ride_id, selected_driver_id, score, eta, distance, confidence, created_at
		FROM match_results
		WHERE ride_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := repo.pool.QueryRow(ctx, query, rideID)

	result := &entities.MatchResult{}
	err := row.Scan(
		&result.MatchRequestID,
		&result.RideID,
		&result.SelectedDriverID,
		&result.Score,
		&result.ETA,
		&result.Distance,
		&result.Confidence,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get match result: %w", err)
	}

	return result, nil
}

// GetStats returns dispatch statistics
func (repo *MatchRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total_matches,
			AVG(score) as avg_score,
			AVG(confidence) as avg_confidence,
			AVG(eta) as avg_eta
		FROM match_results
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`

	row := repo.pool.QueryRow(ctx, query)

	var totalMatches int
	var avgScore, avgConfidence, avgETA float64

	err := row.Scan(&totalMatches, &avgScore, &avgConfidence, &avgETA)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return map[string]interface{}{
		"total_matches":   totalMatches,
		"avg_score":       avgScore,
		"avg_confidence":  avgConfidence,
		"avg_eta":         avgETA,
	}, nil
}
