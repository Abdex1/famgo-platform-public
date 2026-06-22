package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Abdex1/FamGo-platform/services/pricing-service/internal/domain/entities"
)

// PricingRuleRepository handles pricing rule persistence
type PricingRuleRepository struct {
	db *sql.DB
}

// NewPricingRuleRepository creates new pricing rule repository
func NewPricingRuleRepository(db *sql.DB) *PricingRuleRepository {
	return &PricingRuleRepository{db: db}
}

// CreatePricingRule saves pricing rule
func (r *PricingRuleRepository) CreatePricingRule(ctx context.Context, rule *entities.PricingRule) error {
	query := `
		INSERT INTO pricing_rules 
		(id, ride_type, city, base_fare, distance_rate, time_rate, minimum_fare, 
		 surge_factor_max, tax_percentage, pool_discount, active_from, active_until, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.ExecContext(ctx, query,
		rule.ID, rule.RideType, rule.City, rule.BaseFare, rule.DistanceRate, rule.TimeRate,
		rule.MinimumFare, rule.SurgeFactorMax, rule.TaxPercentage, rule.PoolDiscount,
		rule.ActiveFrom, rule.ActiveUntil, rule.Status, rule.CreatedAt, rule.UpdatedAt,
	)
	return err
}

// GetPricingRule gets rule by ID
func (r *PricingRuleRepository) GetPricingRule(ctx context.Context, ruleID string) (*entities.PricingRule, error) {
	query := `
		SELECT id, ride_type, city, base_fare, distance_rate, time_rate, minimum_fare,
		       surge_factor_max, tax_percentage, pool_discount, active_from, active_until, status,
		       created_at, updated_at
		FROM pricing_rules WHERE id = $1 AND status = 'ACTIVE'
	`

	rule := &entities.PricingRule{}
	err := r.db.QueryRowContext(ctx, query, ruleID).Scan(
		&rule.ID, &rule.RideType, &rule.City, &rule.BaseFare, &rule.DistanceRate, &rule.TimeRate,
		&rule.MinimumFare, &rule.SurgeFactorMax, &rule.TaxPercentage, &rule.PoolDiscount,
		&rule.ActiveFrom, &rule.ActiveUntil, &rule.Status,
		&rule.CreatedAt, &rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("pricing rule not found")
	}
	return rule, err
}

// GetActiveRuleForRideType gets active rule for ride type and city
func (r *PricingRuleRepository) GetActiveRuleForRideType(ctx context.Context, rideType, city string) (*entities.PricingRule, error) {
	query := `
		SELECT id, ride_type, city, base_fare, distance_rate, time_rate, minimum_fare,
		       surge_factor_max, tax_percentage, pool_discount, active_from, active_until, status,
		       created_at, updated_at
		FROM pricing_rules 
		WHERE ride_type = $1 AND city = $2 AND status = 'ACTIVE'
		AND active_from <= NOW() AND (active_until IS NULL OR active_until > NOW())
		LIMIT 1
	`

	rule := &entities.PricingRule{}
	err := r.db.QueryRowContext(ctx, query, rideType, city).Scan(
		&rule.ID, &rule.RideType, &rule.City, &rule.BaseFare, &rule.DistanceRate, &rule.TimeRate,
		&rule.MinimumFare, &rule.SurgeFactorMax, &rule.TaxPercentage, &rule.PoolDiscount,
		&rule.ActiveFrom, &rule.ActiveUntil, &rule.Status,
		&rule.CreatedAt, &rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no active pricing rule found")
	}
	return rule, err
}

// SaveFareCalculation persists fare calculation
func (r *PricingRuleRepository) SaveFareCalculation(ctx context.Context, fare *entities.FareCalculation) error {
	query := `
		INSERT INTO fare_calculations 
		(id, ride_id, ride_type, distance_meters, duration_seconds, pickup_lat, pickup_lng, 
		 dropoff_lat, dropoff_lng, base_fare, distance_fare, time_fare, subtotal_before_surge,
		 surge_multiplier, surge_amount, taxes, discount_code_id, discount_amount, final_fare, is_pool, city, calculated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	`

	_, err := r.db.ExecContext(ctx, query,
		fare.ID, fare.RideID, fare.RideType, fare.DistanceMeters, fare.DurationSeconds,
		fare.PickupLat, fare.PickupLng, fare.DropoffLat, fare.DropoffLng,
		fare.BaseFare, fare.DistanceFare, fare.TimeFare, fare.SubtotalBeforeSurge,
		fare.SurgeMultiplier, fare.SurgeAmount, fare.Taxes,
		fare.DiscountCodeID, fare.DiscountAmount, fare.FinalFare,
		fare.IsPool, fare.City, fare.CalculatedAt, fare.CreatedAt,
	)
	return err
}

// SaveSurgeMultiplier saves surge history for analytics
func (r *PricingRuleRepository) SaveSurgeMultiplier(ctx context.Context, surge *entities.SurgeHistory) error {
	query := `
		INSERT INTO surge_history 
		(id, timestamp, city, latitude, longitude, surge_multiplier, active_rides, available_drivers, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		surge.ID, surge.Timestamp, surge.City, surge.Latitude, surge.Longitude,
		surge.SurgeMultiplier, surge.ActiveRides, surge.AvailableDrivers,
		surge.Reason, surge.CreatedAt,
	)
	return err
}

// ValidateDiscountCode validates discount code
func (r *PricingRuleRepository) ValidateDiscountCode(ctx context.Context, code string) (*entities.DiscountCode, error) {
	query := `
		SELECT id, code, discount_type, discount_value, max_discount, minimum_fare_amount,
		       max_uses, uses_remaining, valid_from, valid_until, applicable_to_pooling, status
		FROM discount_codes 
		WHERE code = $1 AND status = 'ACTIVE' AND uses_remaining > 0
		AND valid_from <= NOW() AND valid_until > NOW()
		LIMIT 1
	`

	discount := &entities.DiscountCode{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&discount.ID, &discount.Code, &discount.DiscountType, &discount.DiscountValue,
		&discount.MaxDiscount, &discount.MinimumFareAmount,
		&discount.MaxUses, &discount.UsesRemaining, &discount.ValidFrom, &discount.ValidUntil,
		&discount.ApplicableToPooling, &discount.Status,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("discount code not found or expired")
	}
	return discount, err
}

// DecrementDiscountCodeUsage decreases remaining uses
func (r *PricingRuleRepository) DecrementDiscountCodeUsage(ctx context.Context, codeID string) error {
	query := `UPDATE discount_codes SET uses_remaining = uses_remaining - 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, codeID)
	return err
}

// GetFareCalculationHistory gets all calculations for ride
func (r *PricingRuleRepository) GetFareCalculationHistory(ctx context.Context, rideID string) ([]entities.FareCalculation, error) {
	query := `
		SELECT id, ride_id, ride_type, distance_meters, duration_seconds, base_fare, distance_fare,
		       time_fare, surge_multiplier, final_fare, created_at
		FROM fare_calculations WHERE ride_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fares []entities.FareCalculation
	for rows.Next() {
		var fare entities.FareCalculation
		err := rows.Scan(
			&fare.ID, &fare.RideID, &fare.RideType, &fare.DistanceMeters, &fare.DurationSeconds,
			&fare.BaseFare, &fare.DistanceFare, &fare.TimeFare, &fare.SurgeMultiplier,
			&fare.FinalFare, &fare.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		fares = append(fares, fare)
	}

	return fares, rows.Err()
}

// GetSurgeHistory retrieves surge history for analytics
func (r *PricingRuleRepository) GetSurgeHistory(ctx context.Context, city string, hours int) ([]entities.SurgeHistory, error) {
	query := `
		SELECT id, timestamp, city, surge_multiplier, active_rides, available_drivers, reason
		FROM surge_history 
		WHERE city = $1 AND timestamp > NOW() - INTERVAL '1 hour' * $2
		ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, city, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []entities.SurgeHistory
	for rows.Next() {
		var surge entities.SurgeHistory
		err := rows.Scan(
			&surge.ID, &surge.Timestamp, &surge.City, &surge.SurgeMultiplier,
			&surge.ActiveRides, &surge.AvailableDrivers, &surge.Reason,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, surge)
	}

	return history, rows.Err()
}

// GetAverageFareByRideType gets average fare stats
func (r *PricingRuleRepository) GetAverageFareByRideType(ctx context.Context, city string, days int) (map[string]float64, error) {
	query := `
		SELECT ride_type, AVG(final_fare) as avg_fare
		FROM fare_calculations 
		WHERE city = $1 AND created_at > NOW() - INTERVAL '1 day' * $2
		GROUP BY ride_type
	`

	rows, err := r.db.QueryContext(ctx, query, city, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var rideType string
		var avgFare float64
		if err := rows.Scan(&rideType, &avgFare); err != nil {
			return nil, err
		}
		result[rideType] = avgFare
	}

	return result, rows.Err()
}
