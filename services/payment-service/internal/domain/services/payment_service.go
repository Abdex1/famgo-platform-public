// services/payment-service/internal/domain/services/payment_service.go
package services

import (
	"fmt"
	"time"

	"github.com/Abdex1/FamGo-platform/services/payment-service/internal/domain/entities"
)

type PaymentService struct {
	minAmount     float64
	maxAmount     float64
	maxRetries    int
	retryDelayMs  int
	timeoutSec    int
}

func NewPaymentService(minAmount, maxAmount float64, maxRetries, retryDelayMs, timeoutSec int) *PaymentService {
	return &PaymentService{
		minAmount:    minAmount,
		maxAmount:    maxAmount,
		maxRetries:   maxRetries,
		retryDelayMs: retryDelayMs,
		timeoutSec:   timeoutSec,
	}
}

func (ps *PaymentService) ValidatePayment(amount float64, method string) (bool, string) {
	if amount < ps.minAmount {
		return false, fmt.Sprintf("amount %.2f below minimum %.2f", amount, ps.minAmount)
	}
	if amount > ps.maxAmount {
		return false, fmt.Sprintf("amount %.2f exceeds maximum %.2f", amount, ps.maxAmount)
	}
	if method != "telebirr" && method != "cbe_birr" && method != "chapa" && method != "wallet" {
		return false, fmt.Sprintf("unsupported payment method: %s", method)
	}
	return true, ""
}

func (ps *PaymentService) SelectProvider(method string) string {
	switch method {
	case "telebirr":
		return "telebirr"
	case "cbe_birr":
		return "cbe_birr"
	case "chapa":
		return "chapa"
	case "wallet":
		return "wallet"
	default:
		return "telebirr" // Default fallback
	}
}

func (ps *PaymentService) GetRetryDelay(attempt int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s, 8s, etc.
	delayMs := ps.retryDelayMs * (1 << uint(attempt))
	return time.Duration(delayMs) * time.Millisecond
}

func (ps *PaymentService) ShouldRetry(payment *entities.Payment) bool {
	return payment.CanRetry()
}

func (ps *PaymentService) GetPaymentTimeout() time.Duration {
	return time.Duration(ps.timeoutSec) * time.Second
}

type ProviderResponse struct {
	TransactionID string
	Status        string
	ErrorCode     string
	ErrorMessage  string
}

func (ps *PaymentService) ValidateProviderResponse(resp *ProviderResponse) (bool, string) {
	if resp == nil {
		return false, "nil response from provider"
	}
	if resp.Status == "failed" || resp.Status == "error" {
		return false, fmt.Sprintf("%s: %s", resp.ErrorCode, resp.ErrorMessage)
	}
	if resp.Status != "success" && resp.Status != "pending" {
		return false, fmt.Sprintf("unknown status: %s", resp.Status)
	}
	if resp.TransactionID == "" {
		return false, "empty transaction ID"
	}
	return true, ""
}

type ChargebackReason struct {
	Code        string
	Description string
	ResolutionDays int
}

func (ps *PaymentService) HandleChargeback(reason string) *ChargebackReason {
	chargebackMap := map[string]ChargebackReason{
		"duplicate": {
			Code:             "duplicate",
			Description:      "Customer claims duplicate charge",
			ResolutionDays:   30,
		},
		"fraud": {
			Code:             "fraud",
			Description:      "Customer reports fraudulent transaction",
			ResolutionDays:   45,
		},
		"unauthorized": {
			Code:             "unauthorized",
			Description:      "Customer did not authorize transaction",
			ResolutionDays:   60,
		},
		"not_received": {
			Code:             "not_received",
			Description:      "Customer claims service not received",
			ResolutionDays:   30,
		},
	}

	if cb, exists := chargebackMap[reason]; exists {
		return &cb
	}

	return &ChargebackReason{
		Code:             "other",
		Description:      "Miscellaneous chargeback",
		ResolutionDays:   45,
	}
}

func (ps *PaymentService) CalculateServiceFee(amount float64, provider string) float64 {
	feePercentage := 0.03 // 3% base fee
	switch provider {
	case "telebirr":
		feePercentage = 0.025 // 2.5%
	case "cbe_birr":
		feePercentage = 0.035 // 3.5%
	case "chapa":
		feePercentage = 0.03 // 3%
	}
	return amount * feePercentage
}
