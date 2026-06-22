// services/payment-service/internal/application/usecases/payment_usecases.go
package usecases

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/infrastructure/repositories"
)

type PaymentUseCases struct {
	repo    *repositories.PaymentRepository
	service *services.PaymentService
}

func NewPaymentUseCases(repo *repositories.PaymentRepository, service *services.PaymentService) *PaymentUseCases {
	return &PaymentUseCases{repo: repo, service: service}
}

type InitiatePaymentInput struct {
	RideID   string
	RiderID  string
	DriverID string
	Amount   float64
	Method   string
}

type PaymentOutput struct {
	PaymentID string
	Status    string
	Amount    float64
}

func (uc *PaymentUseCases) InitiatePayment(ctx context.Context, input *InitiatePaymentInput) (*PaymentOutput, error) {
	if input == nil || input.RideID == "" {
		return nil, fmt.Errorf("invalid payment input")
	}

	valid, errMsg := uc.service.ValidatePayment(input.Amount, input.Method)
	if !valid {
		return nil, fmt.Errorf("invalid payment: %s", errMsg)
	}

	payment, err := entities.NewPayment(input.RideID, input.RiderID, input.DriverID, input.Amount, entities.PaymentMethod(input.Method))
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	if err := uc.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to persist payment: %w", err)
	}

	return &PaymentOutput{
		PaymentID: payment.ID,
		Status:    string(payment.Status),
		Amount:    payment.Amount,
	}, nil
}

type CompletePaymentInput struct {
	PaymentID     string
	TransactionID string
}

func (uc *PaymentUseCases) CompletePayment(ctx context.Context, input *CompletePaymentInput) error {
	if input == nil || input.PaymentID == "" {
		return fmt.Errorf("invalid payment input")
	}

	payment, err := uc.repo.GetByID(ctx, input.PaymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if err := payment.Complete(input.TransactionID); err != nil {
		return fmt.Errorf("failed to complete payment: %w", err)
	}

	if err := uc.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

type RefundPaymentInput struct {
	PaymentID string
	Amount    float64
	Reason    string
}

func (uc *PaymentUseCases) RefundPayment(ctx context.Context, input *RefundPaymentInput) error {
	if input == nil || input.PaymentID == "" {
		return fmt.Errorf("invalid refund input")
	}

	payment, err := uc.repo.GetByID(ctx, input.PaymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if err := payment.Refund(input.Amount, input.Reason); err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	if err := uc.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

type GetPaymentInput struct {
	PaymentID string
}

func (uc *PaymentUseCases) GetPayment(ctx context.Context, input *GetPaymentInput) (*PaymentOutput, error) {
	if input == nil || input.PaymentID == "" {
		return nil, fmt.Errorf("invalid payment input")
	}

	payment, err := uc.repo.GetByID(ctx, input.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	return &PaymentOutput{
		PaymentID: payment.ID,
		Status:    string(payment.Status),
		Amount:    payment.Amount,
	}, nil
}

func (uc *PaymentUseCases) HandleWebhook(ctx context.Context, paymentID string) error {
	payment, err := uc.repo.GetByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	payment.VerifyWebhook()
	if err := uc.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}
