// services/fraud-service/internal/infrastructure/repositories/fraud_repository.go
package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/fraud-service/internal/domain/entities"
)

type FraudRepository struct {
	pool *pgxpool.Pool
}

func NewFraudRepository(pool *pgxpool.Pool) *FraudRepository {
	return &FraudRepository{pool: pool}
}

func (r *FraudRepository) Create(ctx context.Context, check *entities.FraudCheck) error {
	query := `
		INSERT INTO fraud_checks 
		(id, ride_id, user_id, user_type, risk_score, risk_level, flags_triggered,
		 location_anomalies, velocity_anomalies, payment_anomalies, behavior_anomalies,
		 is_blacklisted, is_review, action, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.pool.Exec(ctx, query,
		check.ID, check.RideID, check.UserID, check.UserType, check.RiskScore, 
		string(check.RiskLevel), check.FlagsTriggered, check.LocationAnomalies,
		check.VelocityAnomalies, check.PaymentAnomalies, check.BehaviorAnomalies,
		check.IsBlacklisted, check.IsReview, check.Action, check.CreatedAt, check.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create fraud check: %w", err)
	}
	return nil
}

func (r *FraudRepository) GetByID(ctx context.Context, id string) (*entities.FraudCheck, error) {
	query := `
		SELECT id, ride_id, user_id, user_type, risk_score, risk_level, flags_triggered,
		       location_anomalies, velocity_anomalies, payment_anomalies, behavior_anomalies,
		       is_blacklisted, is_review, review_reason, review_by, reviewed_at,
		       action, is_manual_override, manual_override_reason, created_at, updated_at
		FROM fraud_checks WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var check entities.FraudCheck
	var flags []string

	err := row.Scan(
		&check.ID, &check.RideID, &check.UserID, &check.UserType, &check.RiskScore,
		&check.RiskLevel, &flags, &check.LocationAnomalies, &check.VelocityAnomalies,
		&check.PaymentAnomalies, &check.BehaviorAnomalies, &check.IsBlacklisted,
		&check.IsReview, &check.ReviewReason, &check.ReviewBy, &check.ReviewedAt,
		&check.Action, &check.IsManualOverride, &check.ManualOverrideReason,
		&check.CreatedAt, &check.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("fraud check not found")
		}
		return nil, fmt.Errorf("failed to query fraud check: %w", err)
	}

	check.FlagsTriggered = flags
	return &check, nil
}

func (r *FraudRepository) Update(ctx context.Context, check *entities.FraudCheck) error {
	query := `
		UPDATE fraud_checks SET
			risk_score = $1, risk_level = $2, flags_triggered = $3,
			location_anomalies = $4, velocity_anomalies = $5, payment_anomalies = $6,
			behavior_anomalies = $7, is_review = $8, review_reason = $9, review_by = $10,
			reviewed_at = $11, action = $12, is_manual_override = $13,
			manual_override_reason = $14, updated_at = $15
		WHERE id = $16
	`
	result, err := r.pool.Exec(ctx, query,
		check.RiskScore, string(check.RiskLevel), check.FlagsTriggered,
		check.LocationAnomalies, check.VelocityAnomalies, check.PaymentAnomalies,
		check.BehaviorAnomalies, check.IsReview, check.ReviewReason, check.ReviewBy,
		check.ReviewedAt, check.Action, check.IsManualOverride,
		check.ManualOverrideReason, check.UpdatedAt, check.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update fraud check: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("fraud check not found")
	}
	return nil
}

func (r *FraudRepository) GetHighRiskChecks(ctx context.Context, limit int) ([]*entities.FraudCheck, error) {
	query := `
		SELECT id, ride_id, user_id, user_type, risk_score, risk_level, flags_triggered,
		       location_anomalies, velocity_anomalies, payment_anomalies, behavior_anomalies,
		       is_blacklisted, is_review, review_reason, review_by, reviewed_at,
		       action, is_manual_override, manual_override_reason, created_at, updated_at
		FROM fraud_checks WHERE risk_level = 'high' AND is_review = true
		ORDER BY created_at DESC LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query high risk checks: %w", err)
	}
	defer rows.Close()

	var checks []*entities.FraudCheck

	for rows.Next() {
		var check entities.FraudCheck
		var flags []string

		err := rows.Scan(
			&check.ID, &check.RideID, &check.UserID, &check.UserType, &check.RiskScore,
			&check.RiskLevel, &flags, &check.LocationAnomalies, &check.VelocityAnomalies,
			&check.PaymentAnomalies, &check.BehaviorAnomalies, &check.IsBlacklisted,
			&check.IsReview, &check.ReviewReason, &check.ReviewBy, &check.ReviewedAt,
			&check.Action, &check.IsManualOverride, &check.ManualOverrideReason,
			&check.CreatedAt, &check.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fraud check: %w", err)
		}

		check.FlagsTriggered = flags
		checks = append(checks, &check)
	}

	return checks, rows.Err()
}
