package services

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/valueobjects"
)

// PaymentProcessor defines payment processing interface
type PaymentProcessor interface {
	Process(ctx context.Context, payment *entities.Payment) error
	Verify(ctx context.Context, externalRef string) (bool, error)
	Refund(ctx context.Context, payment *entities.Payment, amount float64) error
}

// PaymentProcessorService routes payments to correct provider
type PaymentProcessorService struct {
	providers map[valueobjects.PaymentProvider]PaymentProcessor
}

func NewPaymentProcessorService() *PaymentProcessorService {
	return &PaymentProcessorService{
		providers: make(map[valueobjects.PaymentProvider]PaymentProcessor),
	}
}

// RegisterProvider registers a payment provider
func (pps *PaymentProcessorService) RegisterProvider(
	provider valueobjects.PaymentProvider,
	processor PaymentProcessor,
) error {
	if !provider.IsValid() {
		return fmt.Errorf("invalid provider: %s", provider)
	}

	pps.providers[provider] = processor
	return nil
}

// Process routes payment to appropriate provider
func (pps *PaymentProcessorService) Process(ctx context.Context, payment *entities.Payment) error {
	provider := valueobjects.PaymentProvider(payment.PaymentMethod)

	proc, exists := pps.providers[provider]
	if !exists {
		return fmt.Errorf("no processor for provider: %s", provider)
	}

	return proc.Process(ctx, payment)
}

// Verify checks payment status with provider
func (pps *PaymentProcessorService) Verify(ctx context.Context, payment *entities.Payment) (bool, error) {
	provider := valueobjects.PaymentProvider(payment.PaymentMethod)

	proc, exists := pps.providers[provider]
	if !exists {
		return false, fmt.Errorf("no processor for provider: %s", provider)
	}

	return proc.Verify(ctx, payment.ExternalReference)
}

// RefundOrchestrator handles refund operations
type RefundOrchestrator struct {
	paymentProcessor *PaymentProcessorService
}

func NewRefundOrchestrator(processor *PaymentProcessorService) *RefundOrchestrator {
	return &RefundOrchestrator{
		paymentProcessor: processor,
	}
}

// ProcessRefund handles refund request
func (ro *RefundOrchestrator) ProcessRefund(
	ctx context.Context,
	payment *entities.Payment,
	amount float64,
	reason string,
) error {
	// Validate refund eligibility
	if payment.Status != entities.PaymentStatusCompleted {
		return fmt.Errorf("can only refund completed payments, current: %s", payment.Status)
	}

	if amount <= 0 || amount > payment.Amount {
		return fmt.Errorf("invalid refund amount: %f (payment: %f)", amount, payment.Amount)
	}

	// Call provider refund
	if err := ro.paymentProcessor.providers[valueobjects.PaymentProvider(payment.PaymentMethod)].
		Refund(ctx, payment, amount); err != nil {
		return fmt.Errorf("provider refund failed: %w", err)
	}

	// Update payment state
	if err := payment.Refund(amount, reason); err != nil {
		return err
	}

	return nil
}

// PaymentValidationService validates payments
type PaymentValidationService struct {
	minAmount float64
	maxAmount float64
}

func NewPaymentValidationService(minAmount, maxAmount float64) *PaymentValidationService {
	return &PaymentValidationService{
		minAmount: minAmount,
		maxAmount: maxAmount,
	}
}

// ValidatePayment performs business rule validation
func (pvs *PaymentValidationService) ValidatePayment(payment *entities.Payment) error {
	if payment.Amount < pvs.minAmount {
		return fmt.Errorf("amount below minimum: %f < %f", payment.Amount, pvs.minAmount)
	}

	if payment.Amount > pvs.maxAmount {
		return fmt.Errorf("amount exceeds maximum: %f > %f", payment.Amount, pvs.maxAmount)
	}

	if !valueobjects.PaymentProvider(payment.PaymentMethod).IsValid() {
		return fmt.Errorf("invalid payment method: %s", payment.PaymentMethod)
	}

	if payment.RideID == "" || payment.UserID == "" {
		return fmt.Errorf("missing required fields")
	}

	return nil
}

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	maxRetries   int32
	backoffBase  int32
	backoffMax   int32
}

func NewRetryPolicy(maxRetries, backoffBase, backoffMax int32) *RetryPolicy {
	return &RetryPolicy{
		maxRetries:  maxRetries,
		backoffBase: backoffBase,
		backoffMax:  backoffMax,
	}
}

// CanRetry determines if payment can be retried
func (rp *RetryPolicy) CanRetry(payment *entities.Payment) bool {
	return payment.Status == entities.PaymentStatusFailed &&
		payment.RetryCount < rp.maxRetries
}

// GetBackoffDuration calculates exponential backoff
func (rp *RetryPolicy) GetBackoffDuration(retryCount int32) int32 {
	backoff := rp.backoffBase
	for i := int32(0); i < retryCount; i++ {
		backoff = backoff * 2
		if backoff > rp.backoffMax {
			backoff = rp.backoffMax
			break
		}
	}
	return backoff
}
