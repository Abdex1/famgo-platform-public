package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/domain/entities"
	"github.com/lib/pq"
)

// PoolRepository handles pool data persistence
type PoolRepository struct {
	db *sql.DB
}

// NewPoolRepository creates new pool repository
func NewPoolRepository(db *sql.DB) *PoolRepository {
	return &PoolRepository{db: db}
}

// CreatePool saves new pool group
func (r *PoolRepository) CreatePool(ctx context.Context, pool *entities.PoolGroup) error {
	rideIDsJSON, _ := json.Marshal(pool.RideIDs)
	sequenceJSON, _ := json.Marshal(pool.PassengerSequence)
	faresJSON, _ := json.Marshal(pool.IndividualFares)

	query := `
		INSERT INTO pool_groups 
		(id, driver_id, status, max_size, current_size, ride_ids, passenger_sequence, 
		 dropoff_sequence, total_distance_meters, total_duration_seconds, pooled_fare, 
		 individual_fares, compatibility_score, estimated_profit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := r.db.ExecContext(ctx, query,
		pool.ID, pool.DriverID, pool.Status, pool.MaxSize, pool.CurrentSize,
		rideIDsJSON, sequenceJSON, sequenceJSON, pool.TotalDistance, pool.TotalDuration,
		pool.PooledFare, faresJSON, pool.CompatibilityScore, pool.EstimatedProfit,
		pool.CreatedAt, pool.UpdatedAt,
	)
	return err
}

// GetPool retrieves pool by ID
func (r *PoolRepository) GetPool(ctx context.Context, poolID string) (*entities.PoolGroup, error) {
	query := `
		SELECT id, driver_id, status, max_size, current_size, ride_ids, passenger_sequence,
		       total_distance_meters, total_duration_seconds, pooled_fare, individual_fares,
		       compatibility_score, estimated_profit, created_at, completed_at, updated_at
		FROM pool_groups WHERE id = $1
	`

	pool := &entities.PoolGroup{}
	var rideIDsJSON, sequenceJSON, faresJSON []byte

	err := r.db.QueryRowContext(ctx, query, poolID).Scan(
		&pool.ID, &pool.DriverID, &pool.Status, &pool.MaxSize, &pool.CurrentSize,
		&rideIDsJSON, &sequenceJSON, &pool.TotalDistance, &pool.TotalDuration,
		&pool.PooledFare, &faresJSON, &pool.CompatibilityScore, &pool.EstimatedProfit,
		&pool.CreatedAt, &pool.CompletedAt, &pool.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("pool not found")
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	json.Unmarshal(rideIDsJSON, &pool.RideIDs)
	json.Unmarshal(sequenceJSON, &pool.PassengerSequence)
	json.Unmarshal(faresJSON, &pool.IndividualFares)

	return pool, nil
}

// UpdatePoolStatus updates pool status
func (r *PoolRepository) UpdatePoolStatus(ctx context.Context, poolID string, status entities.PoolStatus) error {
	query := `UPDATE pool_groups SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), poolID)
	return err
}

// AddRideToPool adds ride to existing pool
func (r *PoolRepository) AddRideToPool(ctx context.Context, poolID, rideID string) error {
	pool, err := r.GetPool(ctx, poolID)
	if err != nil {
		return err
	}

	if pool.AddRide(rideID) {
		rideIDsJSON, _ := json.Marshal(pool.RideIDs)
		sequenceJSON, _ := json.Marshal(pool.PassengerSequence)

		query := `
			UPDATE pool_groups 
			SET current_size = $1, ride_ids = $2, passenger_sequence = $3, updated_at = $4
			WHERE id = $5
		`
		_, err = r.db.ExecContext(ctx, query, pool.CurrentSize, rideIDsJSON, sequenceJSON, time.Now(), poolID)
		return err
	}

	return errors.New("pool is full")
}

// GetActivePoolsByDriver retrieves active pools for driver
func (r *PoolRepository) GetActivePoolsByDriver(ctx context.Context, driverID string) ([]entities.PoolGroup, error) {
	query := `
		SELECT id, driver_id, status, max_size, current_size, ride_ids, passenger_sequence,
		       total_distance_meters, total_duration_seconds, pooled_fare, compatibility_score,
		       estimated_profit, created_at, updated_at
		FROM pool_groups 
		WHERE driver_id = $1 AND status IN ('FORMING', 'ACTIVE')
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []entities.PoolGroup
	for rows.Next() {
		var pool entities.PoolGroup
		var rideIDsJSON, sequenceJSON []byte

		err := rows.Scan(
			&pool.ID, &pool.DriverID, &pool.Status, &pool.MaxSize, &pool.CurrentSize,
			&rideIDsJSON, &sequenceJSON, &pool.TotalDistance, &pool.TotalDuration,
			&pool.PooledFare, &pool.CompatibilityScore, &pool.EstimatedProfit,
			&pool.CreatedAt, &pool.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		json.Unmarshal(rideIDsJSON, &pool.RideIDs)
		json.Unmarshal(sequenceJSON, &pool.PassengerSequence)
		pools = append(pools, pool)
	}

	return pools, rows.Err()
}

// GetPoolStatistics retrieves pool metrics
func (r *PoolRepository) GetPoolStatistics(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_pools,
			SUM(CASE WHEN status = 'ACTIVE' THEN 1 ELSE 0 END) as active_pools,
			SUM(CASE WHEN status = 'COMPLETED' THEN 1 ELSE 0 END) as completed_pools,
			AVG(CAST(current_size AS FLOAT)) as avg_pool_size,
			AVG(compatibility_score) as avg_compatibility_score,
			SUM(estimated_profit) as total_estimated_profit
		FROM pool_groups
	`

	stats := make(map[string]interface{})
	var totalPools, activePools, completedPools sql.NullInt64
	var avgPoolSize, avgScore, totalProfit sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query).Scan(
		&totalPools, &activePools, &completedPools, &avgPoolSize, &avgScore, &totalProfit,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if totalPools.Valid {
		stats["total_pools"] = totalPools.Int64
	}
	if activePools.Valid {
		stats["active_pools"] = activePools.Int64
	}
	if completedPools.Valid {
		stats["completed_pools"] = completedPools.Int64
	}
	if avgPoolSize.Valid {
		stats["avg_pool_size"] = avgPoolSize.Float64
	}
	if avgScore.Valid {
		stats["avg_compatibility_score"] = avgScore.Float64
	}
	if totalProfit.Valid {
		stats["total_estimated_profit"] = totalProfit.Float64
	}

	return stats, nil
}

// CreatePoolRequest saves pool request
func (r *PoolRepository) CreatePoolRequest(ctx context.Context, req *entities.PoolRequest) error {
	query := `
		INSERT INTO pool_requests 
		(ride_id, driver_id, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, 
		 pickup_address, dropoff_address, estimated_distance_meters, estimated_duration_seconds,
		 estimated_fare, female_only, max_detour_minutes, max_wait_minutes, min_route_overlap,
		 created_at, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.ExecContext(ctx, query,
		req.RideID, req.DriverID, req.PickupLat, req.PickupLng, req.DropoffLat, req.DropoffLng,
		req.PickupAddress, req.DropoffAddress, req.EstimatedDistance, req.EstimatedDuration,
		req.EstimatedFare, req.FemaleOnly, req.MaxDetourMinutes, req.MaxWaitMinutes,
		req.MinRouteOverlap, req.CreatedAt, req.CreatedAt.Add(5*time.Minute), "PENDING",
	)
	return err
}

// GetActivePoolRequests retrieves eligible rides for pooling
func (r *PoolRepository) GetActivePoolRequests(ctx context.Context, limit int) ([]entities.PoolRequest, error) {
	query := `
		SELECT ride_id, driver_id, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng,
		       pickup_address, dropoff_address, estimated_distance_meters, estimated_duration_seconds,
		       estimated_fare, female_only, max_detour_minutes, max_wait_minutes, min_route_overlap,
		       created_at
		FROM pool_requests 
		WHERE status = 'PENDING' AND expires_at > NOW()
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entities.PoolRequest
	for rows.Next() {
		var req entities.PoolRequest
		err := rows.Scan(
			&req.RideID, &req.DriverID, &req.PickupLat, &req.PickupLng, &req.DropoffLat, &req.DropoffLng,
			&req.PickupAddress, &req.DropoffAddress, &req.EstimatedDistance, &req.EstimatedDuration,
			&req.EstimatedFare, &req.FemaleOnly, &req.MaxDetourMinutes, &req.MaxWaitMinutes,
			&req.MinRouteOverlap, &req.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, rows.Err()
}
