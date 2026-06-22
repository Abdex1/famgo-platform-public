// services/wallet-service/internal/infrastructure/repositories/wallet_repository.go
package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/wallet-service/internal/domain/entities"
)

type WalletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{pool: pool}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, w *entities.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, user_type, balance, currency, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query, w.ID, w.UserID, w.UserType, w.Balance, w.Currency, w.IsActive, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	return nil
}

func (r *WalletRepository) GetWallet(ctx context.Context, walletID string) (*entities.Wallet, error) {
	query := `
		SELECT id, user_id, user_type, balance, currency, total_deposited, total_withdrawn,
		       total_earned, is_active, kyc_verified, last_transaction_at, created_at, updated_at
		FROM wallets WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, walletID)

	var w entities.Wallet
	err := row.Scan(
		&w.ID, &w.UserID, &w.UserType, &w.Balance, &w.Currency, &w.TotalDeposited,
		&w.TotalWithdrawn, &w.TotalEarned, &w.IsActive, &w.KYCVerified, &w.LastTransactionAt,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to query wallet: %w", err)
	}
	return &w, nil
}

func (r *WalletRepository) UpdateWallet(ctx context.Context, w *entities.Wallet) error {
	query := `
		UPDATE wallets SET
			balance = $1, total_deposited = $2, total_withdrawn = $3, total_earned = $4,
			last_transaction_at = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		w.Balance, w.TotalDeposited, w.TotalWithdrawn, w.TotalEarned,
		w.LastTransactionAt, w.UpdatedAt, w.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("wallet not found")
	}
	return nil
}

func (r *WalletRepository) RecordTransaction(ctx context.Context, t *entities.Transaction) error {
	query := `
		INSERT INTO wallet_transactions 
		(id, wallet_id, transaction_type, amount, balance_before, balance_after,
		 status, reference_id, description, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.pool.Exec(ctx, query,
		t.ID, t.WalletID, t.TransactionType, t.Amount, t.BalanceBefore, t.BalanceAfter,
		t.Status, t.ReferenceID, t.Description, t.CreatedAt, t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}
	return nil
}

func (r *WalletRepository) CreateSnapshot(ctx context.Context, snap *entities.WalletSnapshot) error {
	query := `
		INSERT INTO wallet_snapshots (id, wallet_id, balance, total_in, total_out, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query, snap.ID, snap.WalletID, snap.Balance, snap.TotalIn, snap.TotalOut, snap.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}
	return nil
}
