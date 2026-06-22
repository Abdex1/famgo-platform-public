// services/payment-service/internal/domain/entities/payment.go
package entities

import (
	"fmt"
	"time"
)

type PaymentStatus string

const (
	StatusInitiated  PaymentStatus = "initiated"
	StatusPending    PaymentStatus = "pending"
	StatusCompleted  PaymentStatus = "completed"
	StatusFailed     PaymentStatus = "failed"
	StatusRefunded   PaymentStatus = "refunded"
	StatusReversed   PaymentStatus = "reversed"
)

type PaymentMethod string

const (
	MethodTelebirr PaymentMethod = "telebirr"
	MethodCBEBirr  PaymentMethod = "cbe_birr"
	MethodChapa    PaymentMethod = "chapa"
	MethodWallet   PaymentMethod = "wallet"
)

type Payment struct {
	ID                string
	RideID            string
	RiderID           string
	DriverID          string
	Amount            float64
	Currency          string
	Method            PaymentMethod
	Provider          string
	Status            PaymentStatus
	ProviderRef       *string
	ProviderTxnID     *string
	InitiatedAt       time.Time
	CompletedAt       *time.Time
	FailedAt          *time.Time
	RefundedAt        *time.Time
	ReversedAt        *time.Time
	FailureReason     string
	RefundAmount      *float64
	RefundReason      *string
	RetryCount        int
	MaxRetries        int
	WebhookVerified   bool
	WebhookVerifiedAt *time.Time
	Metadata          map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
	CreatedBy         string
	UpdatedBy         string
}

func NewPayment(rideID, riderID, driverID string, amount float64, method PaymentMethod) (*Payment, error) {
	if rideID == "" || riderID == "" {
		return nil, fmt.Errorf("ride ID and rider ID required")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	now := time.Now()
	return &Payment{
		ID:         fmt.Sprintf("pay_%d", now.UnixNano()),
		RideID:     rideID,
		RiderID:    riderID,
		DriverID:   driverID,
		Amount:     amount,
		Currency:   "ETB",
		Method:     method,
		Status:     StatusInitiated,
		InitiatedAt: now,
		RetryCount: 0,
		MaxRetries: 3,
		Metadata:   make(map[string]interface{}),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (p *Payment) IsValid() bool {
	return p.ID != "" && p.RideID != "" && p.RiderID != "" && p.Amount > 0 && p.DeletedAt == nil
}

func (p *Payment) CanRetry() bool {
	return p.RetryCount < p.MaxRetries && (p.Status == StatusFailed || p.Status == StatusInitiated)
}

func (p *Payment) SetPending(providerRef string) error {
	if p.Status != StatusInitiated {
		return fmt.Errorf("can only move to pending from initiated, current: %s", p.Status)
	}
	p.Status = StatusPending
	p.ProviderRef = &providerRef
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) Complete(txnID string) error {
	if p.Status != StatusPending {
		return fmt.Errorf("can only complete pending payments, current: %s", p.Status)
	}
	now := time.Now()
	p.Status = StatusCompleted
	p.ProviderTxnID = &txnID
	p.CompletedAt = &now
	p.UpdatedAt = now
	return nil
}

func (p *Payment) Fail(reason string) error {
	if p.Status == StatusCompleted || p.Status == StatusRefunded {
		return fmt.Errorf("cannot fail completed or refunded payment")
	}
	now := time.Now()
	p.Status = StatusFailed
	p.FailureReason = reason
	p.FailedAt = &now
	p.RetryCount++
	p.UpdatedAt = now
	return nil
}

func (p *Payment) Refund(amount float64, reason string) error {
	if p.Status != StatusCompleted {
		return fmt.Errorf("can only refund completed payments")
	}
	if amount > p.Amount {
		return fmt.Errorf("refund amount exceeds payment amount")
	}
	now := time.Now()
	p.Status = StatusRefunded
	p.RefundAmount = &amount
	p.RefundReason = &reason
	p.RefundedAt = &now
	p.UpdatedAt = now
	return nil
}

func (p *Payment) Reverse(reason string) error {
	if p.Status != StatusCompleted && p.Status != StatusRefunded {
		return fmt.Errorf("can only reverse completed or refunded payments")
	}
	now := time.Now()
	p.Status = StatusReversed
	p.FailureReason = reason
	p.ReversedAt = &now
	p.UpdatedAt = now
	return nil
}

func (p *Payment) VerifyWebhook() {
	now := time.Now()
	p.WebhookVerified = true
	p.WebhookVerifiedAt = &now
	p.UpdatedAt = now
}
