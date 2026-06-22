// services/wallet-service/internal/application/usecases/wallet_usecases.go
package usecases

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/wallet-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/wallet-service/internal/infrastructure/repositories"
)

type WalletUseCases struct {
	repo *repositories.WalletRepository
}

func NewWalletUseCases(repo *repositories.WalletRepository) *WalletUseCases {
	return &WalletUseCases{repo: repo}
}

type CreateWalletInput struct {
	UserID   string
	UserType string
}

type WalletOutput struct {
	WalletID string
	UserID   string
	Balance  float64
	Currency string
}

func (uc *WalletUseCases) CreateWallet(ctx context.Context, input *CreateWalletInput) (*WalletOutput, error) {
	if input == nil || input.UserID == "" {
		return nil, fmt.Errorf("user ID required")
	}

	wallet, err := entities.NewWallet(input.UserID)
	if err != nil {
		return nil, err
	}
	wallet.UserType = input.UserType

	if err := uc.repo.CreateWallet(ctx, wallet); err != nil {
		return nil, err
	}

	return &WalletOutput{
		WalletID: wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}

type GetWalletInput struct {
	WalletID string
}

func (uc *WalletUseCases) GetWallet(ctx context.Context, input *GetWalletInput) (*WalletOutput, error) {
	if input == nil || input.WalletID == "" {
		return nil, fmt.Errorf("wallet ID required")
	}

	wallet, err := uc.repo.GetWallet(ctx, input.WalletID)
	if err != nil {
		return nil, err
	}

	return &WalletOutput{
		WalletID: wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}

type RecordTransactionInput struct {
	WalletID   string
	TxType     string
	Amount     float64
	Description string
	CreatedBy  string
}

func (uc *WalletUseCases) RecordTransaction(ctx context.Context, input *RecordTransactionInput) (*WalletOutput, error) {
	if input == nil || input.WalletID == "" {
		return nil, fmt.Errorf("wallet ID required")
	}

	wallet, err := uc.repo.GetWallet(ctx, input.WalletID)
	if err != nil {
		return nil, err
	}

	txn, err := wallet.RecordTransaction(input.TxType, input.Amount, nil, input.Description, input.CreatedBy)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.RecordTransaction(ctx, txn); err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateWallet(ctx, wallet); err != nil {
		return nil, err
	}

	return &WalletOutput{
		WalletID: wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}
