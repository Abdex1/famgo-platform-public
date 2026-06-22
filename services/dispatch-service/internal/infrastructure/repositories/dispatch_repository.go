// services/dispatch-service/internal/infrastructure/repositories/dispatch_repository.go
package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
)

type DispatchRepository struct {
	pool *pgxpool.Pool
}

func NewDispatchRepository(pool *pgxpool.Pool) *DispatchRepository {
	return &DispatchRepository{pool: pool}
}

func (r *DispatchRepository) Create(ctx context.Context, dr *entities.DispatchRequest) error {
	query := `
		INSERT INTO dispatch_requests (
			id, ride_id, rider_id, pickup_latitude, pickup_longitude,
			dropoff_latitude, dropoff_longitude, status, search_radius,
			max_search_radius, attempt_count, max_attempts, expiry_time,
			requested_at, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.pool.Exec(ctx, query,
		dr.ID, dr.RideID, dr.RiderID,
		dr.PickupLatitude, dr.PickupLongitude,
		dr.DropoffLatitude, dr.DropoffLongitude,
		string(dr.Status), dr.SearchRadius, dr.MaxSearchRadius,
		dr.AttemptCount, dr.MaxAttempts, dr.ExpiryTime,
		dr.RequestedAt, dr.CreatedAt, dr.UpdatedAt, dr.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create dispatch request: %w", err)
	}

	return nil
}

func (r *DispatchRepository) Update(ctx context.Context, dr *entities.DispatchRequest) error {
	query := `
		UPDATE dispatch_requests SET
			status = $1, matched_driver_id = $2, proposed_driver_ids = $3,
			matching_started_at = $4, matched_at = $5, accepted_at = $6,
			rejected_at = $7, expired_at = $8, rejection_reason = $9,
			search_radius = $10, attempt_count = $11, expiry_time = $12,
			updated_at = $13, updated_by = $14
		WHERE id = $15 AND deleted_at IS NULL
	`

	result, err := r.pool.Exec(ctx, query,
		string(dr.Status), dr.MatchedDriverID, dr.ProposedDriverIDs,
		dr.MatchingStartedAt, dr.MatchedAt, dr.AcceptedAt,
		dr.RejectedAt, dr.ExpiredAt, dr.RejectionReason,
		dr.SearchRadius, dr.AttemptCount, dr.ExpiryTime,
		dr.UpdatedAt, dr.UpdatedBy, dr.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update dispatch request: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("dispatch request not found: %s", dr.ID)
	}

	return nil
}

func (r *DispatchRepository) GetByID(ctx context.Context, id string) (*entities.DispatchRequest, error) {
	query := `
		SELECT id, ride_id, rider_id, pickup_latitude, pickup_longitude,
		       dropoff_latitude, dropoff_longitude, status, matched_driver_id,
		       proposed_driver_ids, requested_at, matching_started_at, matched_at,
		       accepted_at, rejected_at, expired_at, rejection_reason,
		       search_radius, max_search_radius, attempt_count, max_attempts,
		       expiry_time, created_at, updated_at, created_by, updated_by
		FROM dispatch_requests
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanDispatchRequest(ctx, query, id)
}

func (r *DispatchRepository) GetByRideID(ctx context.Context, rideID string) (*entities.DispatchRequest, error) {
	query := `
		SELECT id, ride_id, rider_id, pickup_latitude, pickup_longitude,
		       dropoff_latitude, dropoff_longitude, status, matched_driver_id,
		       proposed_driver_ids, requested_at, matching_started_at, matched_at,
		       accepted_at, rejected_at, expired_at, rejection_reason,
		       search_radius, max_search_radius, attempt_count, max_attempts,
		       expiry_time, created_at, updated_at, created_by, updated_by
		FROM dispatch_requests
		WHERE ride_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1
	`

	return r.scanDispatchRequest(ctx, query, rideID)
}

func (r *DispatchRepository) GetPendingRequests(ctx context.Context, limit int) ([]*entities.DispatchRequest, error) {
	query := `
		SELECT id, ride_id, rider_id, pickup_latitude, pickup_longitude,
		       dropoff_latitude, dropoff_longitude, status, matched_driver_id,
		       proposed_driver_ids, requested_at, matching_started_at, matched_at,
		       accepted_at, rejected_at, expired_at, rejection_reason,
		       search_radius, max_search_radius, attempt_count, max_attempts,
		       expiry_time, created_at, updated_at, created_by, updated_by
		FROM dispatch_requests
		WHERE status = 'pending' AND expiry_time > NOW() AND deleted_at IS NULL
		ORDER BY requested_at ASC LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending requests: %w", err)
	}
	defer rows.Close()

	return r.scanDispatchRequests(rows)
}

func (r *DispatchRepository) scanDispatchRequest(ctx context.Context, query string, args ...interface{}) (*entities.DispatchRequest, error) {
	row := r.pool.QueryRow(ctx, query, args...)

	var (
		id, rideID, riderID, status, rejectionReason string
		matchedDriverID *string
		proposedDriverIDs []string
		pickupLat, pickupLng, dropoffLat, dropoffLng, searchRadius, maxSearchRadius float64
		attemptCount, maxAttempts int
		requestedAt time.Time
		matchingStartedAtPtr, matchedAtPtr, acceptedAtPtr, rejectedAtPtr, expiredAtPtr *time.Time
		expiryTime, createdAt, updatedAt time.Time
		createdBy, updatedBy string
	)

	err := row.Scan(
		&id, &rideID, &riderID, &pickupLat, &pickupLng,
		&dropoffLat, &dropoffLng, &status, &matchedDriverID,
		&proposedDriverIDs, &requestedAt, &matchingStartedAtPtr, &matchedAtPtr,
		&acceptedAtPtr, &rejectedAtPtr, &expiredAtPtr, &rejectionReason,
		&searchRadius, &maxSearchRadius, &attemptCount, &maxAttempts,
		&expiryTime, &createdAt, &updatedAt, &createdBy, &updatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("dispatch request not found")
		}
		return nil, fmt.Errorf("failed to scan dispatch request: %w", err)
	}

	return &entities.DispatchRequest{
		ID:                 id,
		RideID:             rideID,
		RiderID:            riderID,
		PickupLatitude:     pickupLat,
		PickupLongitude:    pickupLng,
		DropoffLatitude:    dropoffLat,
		DropoffLongitude:   dropoffLng,
		Status:             entities.MatchStatus(status),
		MatchedDriverID:    matchedDriverID,
		ProposedDriverIDs:  proposedDriverIDs,
		RequestedAt:        requestedAt,
		MatchingStartedAt:  matchingStartedAtPtr,
		MatchedAt:          matchedAtPtr,
		AcceptedAt:         acceptedAtPtr,
		RejectedAt:         rejectedAtPtr,
		ExpiredAt:          expiredAtPtr,
		RejectionReason:    rejectionReason,
		SearchRadius:       searchRadius,
		MaxSearchRadius:    maxSearchRadius,
		AttemptCount:       attemptCount,
		MaxAttempts:        maxAttempts,
		ExpiryTime:         expiryTime,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		CreatedBy:          createdBy,
		UpdatedBy:          updatedBy,
	}, nil
}

func (r *DispatchRepository) scanDispatchRequests(rows pgx.Rows) ([]*entities.DispatchRequest, error) {
	var requests []*entities.DispatchRequest

	for rows.Next() {
		var (
			id, rideID, riderID, status, rejectionReason string
			matchedDriverID *string
			proposedDriverIDs []string
			pickupLat, pickupLng, dropoffLat, dropoffLng, searchRadius, maxSearchRadius float64
			attemptCount, maxAttempts int
			requestedAt time.Time
			matchingStartedAtPtr, matchedAtPtr, acceptedAtPtr, rejectedAtPtr, expiredAtPtr *time.Time
			expiryTime, createdAt, updatedAt time.Time
			createdBy, updatedBy string
		)

		err := rows.Scan(
			&id, &rideID, &riderID, &pickupLat, &pickupLng,
			&dropoffLat, &dropoffLng, &status, &matchedDriverID,
			&proposedDriverIDs, &requestedAt, &matchingStartedAtPtr, &matchedAtPtr,
			&acceptedAtPtr, &rejectedAtPtr, &expiredAtPtr, &rejectionReason,
			&searchRadius, &maxSearchRadius, &attemptCount, &maxAttempts,
			&expiryTime, &createdAt, &updatedAt, &createdBy, &updatedBy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan dispatch request: %w", err)
		}

		requests = append(requests, &entities.DispatchRequest{
			ID:                id,
			RideID:            rideID,
			RiderID:           riderID,
			PickupLatitude:    pickupLat,
			PickupLongitude:   pickupLng,
			DropoffLatitude:   dropoffLat,
			DropoffLongitude:  dropoffLng,
			Status:            entities.MatchStatus(status),
			MatchedDriverID:   matchedDriverID,
			ProposedDriverIDs: proposedDriverIDs,
			RequestedAt:       requestedAt,
			MatchingStartedAt: matchingStartedAtPtr,
			MatchedAt:         matchedAtPtr,
			AcceptedAt:        acceptedAtPtr,
			RejectedAt:        rejectedAtPtr,
			ExpiredAt:         expiredAtPtr,
			RejectionReason:   rejectionReason,
			SearchRadius:      searchRadius,
			MaxSearchRadius:   maxSearchRadius,
			AttemptCount:      attemptCount,
			MaxAttempts:       maxAttempts,
			ExpiryTime:        expiryTime,
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
			CreatedBy:         createdBy,
			UpdatedBy:         updatedBy,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate dispatch requests: %w", err)
	}

	return requests, nil
}

type DispatchStats struct {
	TotalMatches       int
	SuccessfulMatches  int
	FailedMatches      int
	SuccessRate        float64
	AverageTimeToMatch float64
}

func (r *DispatchRepository) GetDispatchStats(ctx context.Context, startDate, endDate time.Time) (*DispatchStats, error) {
	if startDate.IsZero() {
		startDate = time.Now().UTC().Add(-24 * time.Hour)
	}
	if endDate.IsZero() {
		endDate = time.Now().UTC()
	}

	query := `
		SELECT
			COUNT(*) FILTER (WHERE status IN ('matched', 'accepted', 'completed')) AS successful_matches,
			COUNT(*) FILTER (WHERE status IN ('failed', 'expired')) AS failed_matches,
			COUNT(*) AS total_matches,
			COALESCE(AVG(EXTRACT(EPOCH FROM (matched_at - matching_started_at))), 0) AS avg_match_seconds
		FROM dispatch_requests
		WHERE requested_at BETWEEN $1 AND $2 AND deleted_at IS NULL
	`

	var successful, failed, total int
	var avgSeconds float64
	if err := r.pool.QueryRow(ctx, query, startDate, endDate).Scan(&successful, &failed, &total, &avgSeconds); err != nil {
		return nil, fmt.Errorf("failed to query dispatch stats: %w", err)
	}

	successRate := 0.0
	if total > 0 {
		successRate = float64(successful) / float64(total) * 100
	}

	return &DispatchStats{
		TotalMatches:       total,
		SuccessfulMatches:  successful,
		FailedMatches:      failed,
		SuccessRate:        successRate,
		AverageTimeToMatch: avgSeconds,
	}, nil
}
