// services/wallet-service/internal/domain/entities/wallet.go
package entities

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID              string
	WalletID        string
	TransactionType string
	Amount          float64
	BalanceBefore   float64
	BalanceAfter    float64
	Status          string
	ReferenceID     *string
	Description     string
	CreatedAt       time.Time
	CreatedBy       string
}

type WalletSnapshot struct {
	ID        string
	WalletID  string
	Balance   float64
	TotalIn   float64
	TotalOut  float64
	CreatedAt time.Time
}

type Wallet struct {
	ID                string
	UserID            string
	UserType          string
	Balance           float64
	Currency          string
	TotalDeposited    float64
	TotalWithdrawn    float64
	TotalEarned       float64
	IsActive          bool
	KYCVerified       bool
	LastTransactionAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

func NewWallet(userID string) (*Wallet, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID required")
	}
	now := time.Now()
	return &Wallet{
		ID:             fmt.Sprintf("wal_%d", now.UnixNano()),
		UserID:         userID,
		Balance:        0,
		Currency:       "ETB",
		TotalDeposited: 0,
		TotalWithdrawn: 0,
		TotalEarned:    0,
		IsActive:       true,
		KYCVerified:    false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (w *Wallet) RecordTransaction(txType string, amount float64, ref *string, desc string, createdBy string) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	balanceBefore := w.Balance

	switch txType {
	case "deposit":
		w.Balance += amount
		w.TotalDeposited += amount
	case "withdrawal":
		if w.Balance < amount {
			return nil, fmt.Errorf("insufficient balance")
		}
		w.Balance -= amount
		w.TotalWithdrawn += amount
	case "earning":
		w.Balance += amount
		w.TotalEarned += amount
	case "refund":
		w.Balance += amount
	case "charge":
		if w.Balance < amount {
			return nil, fmt.Errorf("insufficient balance for charge")
		}
		w.Balance -= amount
	default:
		return nil, fmt.Errorf("unknown transaction type: %s", txType)
	}

	now := time.Now()
	w.LastTransactionAt = &now
	w.UpdatedAt = now

	return &Transaction{
		ID:              fmt.Sprintf("txn_%d", now.UnixNano()),
		WalletID:        w.ID,
		TransactionType: txType,
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    w.Balance,
		Status:          "completed",
		ReferenceID:     ref,
		Description:     desc,
		CreatedAt:       now,
		CreatedBy:       createdBy,
	}, nil
}

func (w *Wallet) Transfer(amount float64, recipientWallet *Wallet, reason string) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}
	if w.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	if recipientWallet == nil {
		return fmt.Errorf("recipient wallet required")
	}

	w.Balance -= amount
	recipientWallet.Balance += amount

	now := time.Now()
	w.UpdatedAt = now
	recipientWallet.UpdatedAt = now
	w.LastTransactionAt = &now
	recipientWallet.LastTransactionAt = &now

	return nil
}

func (w *Wallet) Reconcile(actualBalance float64) error {
	if actualBalance < 0 {
		return fmt.Errorf("balance cannot be negative")
	}
	w.Balance = actualBalance
	w.UpdatedAt = time.Now()
	return nil
}
