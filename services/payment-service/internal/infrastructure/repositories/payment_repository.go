// services/payment-service/internal/infrastructure/repositories/payment_repository.go
package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/entities"
)

type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Create(ctx context.Context, p *entities.Payment) error {
	query := `
		INSERT INTO payments (
			id, ride_id, rider_id, driver_id, amount, currency, method, provider,
			status, initiated_at, retry_count, max_retries, metadata,
			created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	err := r.pool.QueryRow(ctx, query,
		p.ID, p.RideID, p.RiderID, p.DriverID, p.Amount, p.Currency, string(p.Method), p.Provider,
		string(p.Status), p.InitiatedAt, p.RetryCount, p.MaxRetries, p.Metadata,
		p.CreatedAt, p.UpdatedAt, p.CreatedBy,
	).Scan()

	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *PaymentRepository) Update(ctx context.Context, p *entities.Payment) error {
	query := `
		UPDATE payments SET
			status = $1, provider_ref = $2, provider_txn_id = $3,
			completed_at = $4, failed_at = $5, refunded_at = $6,
			reversed_at = $7, failure_reason = $8, refund_amount = $9,
			refund_reason = $10, retry_count = $11, webhook_verified = $12,
			webhook_verified_at = $13, updated_at = $14, updated_by = $15
		WHERE id = $16 AND deleted_at IS NULL
	`

	result, err := r.pool.Exec(ctx, query,
		string(p.Status), p.ProviderRef, p.ProviderTxnID,
		p.CompletedAt, p.FailedAt, p.RefundedAt,
		p.ReversedAt, p.FailureReason, p.RefundAmount,
		p.RefundReason, p.RetryCount, p.WebhookVerified,
		p.WebhookVerifiedAt, p.UpdatedAt, p.UpdatedBy, p.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found: %s", p.ID)
	}
	return nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*entities.Payment, error) {
	query := `
		SELECT id, ride_id, rider_id, driver_id, amount, currency, method, provider,
		       status, provider_ref, provider_txn_id, initiated_at, completed_at,
		       failed_at, refunded_at, reversed_at, failure_reason, refund_amount,
		       refund_reason, retry_count, max_retries, webhook_verified,
		       webhook_verified_at, metadata, created_at, updated_at, created_by, updated_by
		FROM payments WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.pool.QueryRow(ctx, query, id)

	var p entities.Payment
	var metadata string

	err := row.Scan(
		&p.ID, &p.RideID, &p.RiderID, &p.DriverID, &p.Amount, &p.Currency, &p.Method, &p.Provider,
		&p.Status, &p.ProviderRef, &p.ProviderTxnID, &p.InitiatedAt, &p.CompletedAt,
		&p.FailedAt, &p.RefundedAt, &p.ReversedAt, &p.FailureReason, &p.RefundAmount,
		&p.RefundReason, &p.RetryCount, &p.MaxRetries, &p.WebhookVerified,
		&p.WebhookVerifiedAt, &metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to query payment: %w", err)
	}

	return &p, nil
}

func (r *PaymentRepository) GetByRideID(ctx context.Context, rideID string) (*entities.Payment, error) {
	query := `
		SELECT id, ride_id, rider_id, driver_id, amount, currency, method, provider,
		       status, provider_ref, provider_txn_id, initiated_at, completed_at,
		       failed_at, refunded_at, reversed_at, failure_reason, refund_amount,
		       refund_reason, retry_count, max_retries, webhook_verified,
		       webhook_verified_at, metadata, created_at, updated_at, created_by, updated_by
		FROM payments WHERE ride_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1
	`

	row := r.pool.QueryRow(ctx, query, rideID)

	var p entities.Payment
	var metadata string

	err := row.Scan(
		&p.ID, &p.RideID, &p.RiderID, &p.DriverID, &p.Amount, &p.Currency, &p.Method, &p.Provider,
		&p.Status, &p.ProviderRef, &p.ProviderTxnID, &p.InitiatedAt, &p.CompletedAt,
		&p.FailedAt, &p.RefundedAt, &p.ReversedAt, &p.FailureReason, &p.RefundAmount,
		&p.RefundReason, &p.RetryCount, &p.MaxRetries, &p.WebhookVerified,
		&p.WebhookVerifiedAt, &metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to query payment: %w", err)
	}

	return &p, nil
}

func (r *PaymentRepository) GetPendingPayments(ctx context.Context, limit int) ([]*entities.Payment, error) {
	query := `
		SELECT id, ride_id, rider_id, driver_id, amount, currency, method, provider,
		       status, provider_ref, provider_txn_id, initiated_at, completed_at,
		       failed_at, refunded_at, reversed_at, failure_reason, refund_amount,
		       refund_reason, retry_count, max_retries, webhook_verified,
		       webhook_verified_at, metadata, created_at, updated_at, created_by, updated_by
		FROM payments
		WHERE status IN ('initiated', 'pending') AND retry_count < max_retries
		  AND deleted_at IS NULL
		ORDER BY created_at ASC LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending payments: %w", err)
	}
	defer rows.Close()

	var payments []*entities.Payment

	for rows.Next() {
		var p entities.Payment
		var metadata string

		err := rows.Scan(
			&p.ID, &p.RideID, &p.RiderID, &p.DriverID, &p.Amount, &p.Currency, &p.Method, &p.Provider,
			&p.Status, &p.ProviderRef, &p.ProviderTxnID, &p.InitiatedAt, &p.CompletedAt,
			&p.FailedAt, &p.RefundedAt, &p.ReversedAt, &p.FailureReason, &p.RefundAmount,
			&p.RefundReason, &p.RetryCount, &p.MaxRetries, &p.WebhookVerified,
			&p.WebhookVerifiedAt, &metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		payments = append(payments, &p)
	}

	return payments, rows.Err()
}
