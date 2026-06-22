package valueobjects

import (
	"fmt"
	"time"
)

// PaymentAmount represents a monetary amount
type PaymentAmount struct {
	Amount   float64
	Currency string
}

func NewPaymentAmount(amount float64, currency string) (*PaymentAmount, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive: %f", amount)
	}

	return &PaymentAmount{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Add returns sum of two amounts
func (pa *PaymentAmount) Add(other *PaymentAmount) (*PaymentAmount, error) {
	if pa.Currency != other.Currency {
		return nil, fmt.Errorf("currency mismatch: %s vs %s", pa.Currency, other.Currency)
	}

	return NewPaymentAmount(pa.Amount+other.Amount, pa.Currency)
}

// Subtract returns difference
func (pa *PaymentAmount) Subtract(other *PaymentAmount) (*PaymentAmount, error) {
	if pa.Currency != other.Currency {
		return nil, fmt.Errorf("currency mismatch: %s vs %s", pa.Currency, other.Currency)
	}

	if pa.Amount < other.Amount {
		return nil, fmt.Errorf("insufficient funds: %f < %f", pa.Amount, other.Amount)
	}

	return NewPaymentAmount(pa.Amount-other.Amount, pa.Currency)
}

// PaymentProvider represents a payment provider abstraction
type PaymentProvider string

const (
	PaymentProviderTelebirr PaymentProvider = "TELEBIRR"
	PaymentProviderCbeBirr  PaymentProvider = "CBE_BIRR"
	PaymentProviderChapa    PaymentProvider = "CHAPA"
)

// IsValid checks if provider is valid
func (pp PaymentProvider) IsValid() bool {
	return pp == PaymentProviderTelebirr ||
		pp == PaymentProviderCbeBirr ||
		pp == PaymentProviderChapa
}

// WebhookSignature represents webhook authentication
type WebhookSignature struct {
	Signature  string
	Timestamp  time.Time
	ProviderID string
}

func NewWebhookSignature(signature string, providerID string) *WebhookSignature {
	return &WebhookSignature{
		Signature:  signature,
		Timestamp:  time.Now().UTC(),
		ProviderID: providerID,
	}
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID     string
	Amount        float64
	Reason        string
	InitiatedBy   string // RIDER, DRIVER, ADMIN
	RequestedAt   time.Time
}

func NewRefundRequest(paymentID string, amount float64, reason string, initiatedBy string) *RefundRequest {
	return &RefundRequest{
		PaymentID:   paymentID,
		Amount:      amount,
		Reason:      reason,
		InitiatedBy: initiatedBy,
		RequestedAt: time.Now().UTC(),
	}
}

// PaymentReconciliation represents a reconciliation record
type PaymentReconciliation struct {
	ID                 string
	PaymentID          string
	ProviderReference  string
	ProviderAmount     float64
	SystemAmount       float64
	ReconciliationTime time.Time
	Status             string // MATCHED, MISMATCH, PENDING
	Notes              string
}

func NewPaymentReconciliation(
	id string,
	paymentID string,
	providerRef string,
	providerAmount float64,
	systemAmount float64,
) *PaymentReconciliation {
	status := "MATCHED"
	if providerAmount != systemAmount {
		status = "MISMATCH"
	}

	return &PaymentReconciliation{
		ID:                 id,
		PaymentID:          paymentID,
		ProviderReference:  providerRef,
		ProviderAmount:     providerAmount,
		SystemAmount:       systemAmount,
		ReconciliationTime: time.Now().UTC(),
		Status:             status,
	}
}
