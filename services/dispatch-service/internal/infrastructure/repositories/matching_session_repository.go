package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
)

type MatchingSessionRepository struct {
	pool *pgxpool.Pool
}

func NewMatchingSessionRepository(pool *pgxpool.Pool) *MatchingSessionRepository {
	return &MatchingSessionRepository{pool: pool}
}

func (r *MatchingSessionRepository) Create(ctx context.Context, session *entities.MatchingSession) error {
	query := `
		INSERT INTO matching_sessions (
			id, dispatch_request_id, ride_id, status, algorithm, search_radius_km,
			candidate_count, selected_driver_id, started_at, completed_at, expires_at,
			failure_reason, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`
	_, err := r.pool.Exec(ctx, query,
		session.ID, session.DispatchRequestID, session.RideID, string(session.Status),
		session.Algorithm, session.SearchRadiusKm, session.CandidateCount, session.SelectedDriverID,
		session.StartedAt, session.CompletedAt, session.ExpiresAt, session.FailureReason,
		session.CreatedAt, session.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create matching session: %w", err)
	}
	return nil
}

func (r *MatchingSessionRepository) Update(ctx context.Context, session *entities.MatchingSession) error {
	query := `
		UPDATE matching_sessions SET
			status = $1, candidate_count = $2, selected_driver_id = $3,
			completed_at = $4, failure_reason = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		string(session.Status), session.CandidateCount, session.SelectedDriverID,
		session.CompletedAt, session.FailureReason, session.UpdatedAt, session.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update matching session: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("matching session not found: %s", session.ID)
	}
	return nil
}

func (r *MatchingSessionRepository) GetByDispatchRequestID(
	ctx context.Context,
	dispatchRequestID string,
) (*entities.MatchingSession, error) {
	query := `
		SELECT id, dispatch_request_id, ride_id, status, algorithm, search_radius_km,
		       candidate_count, selected_driver_id, started_at, completed_at, expires_at,
		       failure_reason, created_at, updated_at
		FROM matching_sessions
		WHERE dispatch_request_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1
	`
	row := r.pool.QueryRow(ctx, query, dispatchRequestID)

	var session entities.MatchingSession
	var status string
	err := row.Scan(
		&session.ID, &session.DispatchRequestID, &session.RideID, &status, &session.Algorithm,
		&session.SearchRadiusKm, &session.CandidateCount, &session.SelectedDriverID,
		&session.StartedAt, &session.CompletedAt, &session.ExpiresAt, &session.FailureReason,
		&session.CreatedAt, &session.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("matching session not found")
		}
		return nil, fmt.Errorf("failed to scan matching session: %w", err)
	}
	session.Status = entities.MatchingSessionStatus(status)
	return &session, nil
}
