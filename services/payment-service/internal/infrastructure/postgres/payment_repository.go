package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/entities"
)

// PaymentRepository handles payment persistence
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository creates new payment repository
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// SavePaymentTransaction saves payment transaction
func (r *PaymentRepository) SavePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error {
	query := `
		INSERT INTO payment_transactions 
		(id, ride_id, user_id, driver_id, amount, currency, payment_method, provider, 
		 provider_charge_id, status, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.RideID, tx.UserID, tx.DriverID, tx.Amount, tx.Currency,
		tx.PaymentMethod, tx.Provider, tx.ProviderChargeID, tx.Status,
		tx.ErrorMessage, tx.CreatedAt, tx.UpdatedAt,
	)
	return err
}

// GetPaymentTransaction retrieves transaction by ID
func (r *PaymentRepository) GetPaymentTransaction(ctx context.Context, txID string) (*entities.PaymentTransaction, error) {
	query := `
		SELECT id, ride_id, user_id, driver_id, amount, currency, payment_method, provider,
		       provider_charge_id, status, error_message, created_at, updated_at, completed_at
		FROM payment_transactions WHERE id = $1
	`

	tx := &entities.PaymentTransaction{}
	err := r.db.QueryRowContext(ctx, query, txID).Scan(
		&tx.ID, &tx.RideID, &tx.UserID, &tx.DriverID, &tx.Amount, &tx.Currency,
		&tx.PaymentMethod, &tx.Provider, &tx.ProviderChargeID, &tx.Status,
		&tx.ErrorMessage, &tx.CreatedAt, &tx.UpdatedAt, &tx.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("payment transaction not found")
	}
	return tx, err
}

// UpdatePaymentStatus updates transaction status
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, txID, status string, errMsg *string) error {
	query := `
		UPDATE payment_transactions 
		SET status = $1, error_message = $2, updated_at = NOW(), completed_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, errMsg, txID)
	return err
}

// SavePaymentMethod saves user payment method
func (r *PaymentRepository) SavePaymentMethod(ctx context.Context, method *entities.PaymentMethod) error {
	query := `
		INSERT INTO payment_methods 
		(id, user_id, method_type, provider, last4_digits, expiry_month, expiry_year, is_default, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		method.ID, method.UserID, method.MethodType, method.Provider,
		method.Last4Digits, method.ExpiryMonth, method.ExpiryYear,
		method.IsDefault, method.Status, method.CreatedAt,
	)
	return err
}

// GetUserPaymentMethods retrieves all active payment methods for user
func (r *PaymentRepository) GetUserPaymentMethods(ctx context.Context, userID string) ([]entities.PaymentMethod, error) {
	query := `
		SELECT id, user_id, method_type, provider, last4_digits, expiry_month, expiry_year, is_default, status
		FROM payment_methods 
		WHERE user_id = $1 AND status = 'ACTIVE'
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []entities.PaymentMethod
	for rows.Next() {
		var method entities.PaymentMethod
		if err := rows.Scan(
			&method.ID, &method.UserID, &method.MethodType, &method.Provider,
			&method.Last4Digits, &method.ExpiryMonth, &method.ExpiryYear,
			&method.IsDefault, &method.Status,
		); err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	return methods, rows.Err()
}

// GetDefaultPaymentMethod retrieves user's default payment method
func (r *PaymentRepository) GetDefaultPaymentMethod(ctx context.Context, userID string) (*entities.PaymentMethod, error) {
	query := `
		SELECT id, user_id, method_type, provider, last4_digits, expiry_month, expiry_year, is_default, status
		FROM payment_methods 
		WHERE user_id = $1 AND is_default = true AND status = 'ACTIVE'
		LIMIT 1
	`

	method := &entities.PaymentMethod{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&method.ID, &method.UserID, &method.MethodType, &method.Provider,
		&method.Last4Digits, &method.ExpiryMonth, &method.ExpiryYear,
		&method.IsDefault, &method.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No default method
	}
	return method, err
}

// GetTransactionsByUser retrieves all transactions for user
func (r *PaymentRepository) GetTransactionsByUser(ctx context.Context, userID string, limit int, offset int) ([]entities.PaymentTransaction, error) {
	query := `
		SELECT id, ride_id, user_id, driver_id, amount, currency, payment_method, provider,
		       provider_charge_id, status, created_at
		FROM payment_transactions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entities.PaymentTransaction
	for rows.Next() {
		var tx entities.PaymentTransaction
		if err := rows.Scan(
			&tx.ID, &tx.RideID, &tx.UserID, &tx.DriverID, &tx.Amount, &tx.Currency,
			&tx.PaymentMethod, &tx.Provider, &tx.ProviderChargeID, &tx.Status, &tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// CountUserTransactions counts user transactions
func (r *PaymentRepository) CountUserTransactions(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM payment_transactions WHERE user_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// GetTransactionsByStatus retrieves transactions by status
func (r *PaymentRepository) GetTransactionsByStatus(ctx context.Context, status string, limit int) ([]entities.PaymentTransaction, error) {
	query := `
		SELECT id, ride_id, user_id, driver_id, amount, currency, payment_method, provider,
		       provider_charge_id, status, created_at
		FROM payment_transactions 
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entities.PaymentTransaction
	for rows.Next() {
		var tx entities.PaymentTransaction
		if err := rows.Scan(
			&tx.ID, &tx.RideID, &tx.UserID, &tx.DriverID, &tx.Amount, &tx.Currency,
			&tx.PaymentMethod, &tx.Provider, &tx.ProviderChargeID, &tx.Status, &tx.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// SaveRefundRequest saves refund request
func (r *PaymentRepository) SaveRefundRequest(ctx context.Context, req *entities.RefundRequest) error {
	query := `
		INSERT INTO refund_requests 
		(id, payment_id, ride_id, user_id, amount, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		req.ID, req.PaymentID, req.RideID, req.UserID, req.Amount,
		req.Reason, req.Status, req.CreatedAt, req.UpdatedAt,
	)
	return err
}

// GetRefundRequest retrieves refund request
func (r *PaymentRepository) GetRefundRequest(ctx context.Context, refundID string) (*entities.RefundRequest, error) {
	query := `
		SELECT id, payment_id, ride_id, user_id, amount, reason, status, approved_by, created_at, updated_at
		FROM refund_requests WHERE id = $1
	`

	req := &entities.RefundRequest{}
	err := r.db.QueryRowContext(ctx, query, refundID).Scan(
		&req.ID, &req.PaymentID, &req.RideID, &req.UserID, &req.Amount,
		&req.Reason, &req.Status, &req.ApprovedBy, &req.CreatedAt, &req.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("refund request not found")
	}
	return req, err
}

// GetPaymentStatistics returns payment statistics
func (r *PaymentRepository) GetPaymentStatistics(ctx context.Context, days int) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_transactions,
			SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END) as successful,
			SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed,
			AVG(amount) as avg_amount,
			MAX(amount) as max_amount,
			SUM(CASE WHEN status = 'SUCCESS' THEN amount ELSE 0 END) as total_revenue
		FROM payment_transactions 
		WHERE created_at > NOW() - INTERVAL '1 day' * $1
	`

	var totalTx, successful, failed sql.NullInt64
	var avgAmount, maxAmount, totalRevenue sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, days).Scan(
		&totalTx, &successful, &failed, &avgAmount, &maxAmount, &totalRevenue,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_transactions": totalTx.Int64,
		"successful":         successful.Int64,
		"failed":             failed.Int64,
		"average_amount":     avgAmount.Float64,
		"max_amount":         maxAmount.Float64,
		"total_revenue":      totalRevenue.Float64,
	}, nil
}
