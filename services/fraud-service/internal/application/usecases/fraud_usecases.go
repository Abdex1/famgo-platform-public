// services/fraud-service/internal/application/usecases/fraud_usecases.go
package usecases

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/fraud-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/fraud-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/fraud-service/internal/infrastructure/repositories"
)

type FraudUseCases struct {
	repo           *repositories.FraudRepository
	fraudService   *services.FraudService
	highThreshold  float64
	mediumThreshold float64
}

func NewFraudUseCases(repo *repositories.FraudRepository, svc *services.FraudService, 
	highThreshold, mediumThreshold float64) *FraudUseCases {
	return &FraudUseCases{
		repo:            repo,
		fraudService:    svc,
		highThreshold:   highThreshold,
		mediumThreshold: mediumThreshold,
	}
}

type CheckRideInput struct {
	RideID    string
	UserID    string
	UserType  string
	Amount    float64
	Latitude  float64
	Longitude float64
}

type CheckRideOutput struct {
	CheckID     string
	RiskScore   float64
	RiskLevel   string
	Action      string
	FlagsTriggered []string
}

func (uc *FraudUseCases) CheckRide(ctx context.Context, input *CheckRideInput) (*CheckRideOutput, error) {
	if input == nil || input.RideID == "" {
		return nil, fmt.Errorf("ride ID required")
	}

	check, err := entities.NewFraudCheck(input.RideID, input.UserID, input.UserType)
	if err != nil {
		return nil, err
	}

	// Simulate anomaly detection (in production, would query historical data)
	// For demo: random detection based on amount
	if input.Amount > 10000 {
		check.PaymentAnomalies = true
	}

	check.CalculateRiskScore(uc.highThreshold, uc.mediumThreshold)

	if err := uc.repo.Create(ctx, check); err != nil {
		return nil, err
	}

	return &CheckRideOutput{
		CheckID:        check.ID,
		RiskScore:      check.RiskScore,
		RiskLevel:      string(check.RiskLevel),
		Action:         check.Action,
		FlagsTriggered: check.FlagsTriggered,
	}, nil
}

type GetCheckInput struct {
	CheckID string
}

func (uc *FraudUseCases) GetCheck(ctx context.Context, input *GetCheckInput) (*CheckRideOutput, error) {
	if input == nil || input.CheckID == "" {
		return nil, fmt.Errorf("check ID required")
	}

	check, err := uc.repo.GetByID(ctx, input.CheckID)
	if err != nil {
		return nil, err
	}

	return &CheckRideOutput{
		CheckID:        check.ID,
		RiskScore:      check.RiskScore,
		RiskLevel:      string(check.RiskLevel),
		Action:         check.Action,
		FlagsTriggered: check.FlagsTriggered,
	}, nil
}

type ReviewCheckInput struct {
	CheckID    string
	ReviewedBy string
	Reason     string
}

func (uc *FraudUseCases) ReviewCheck(ctx context.Context, input *ReviewCheckInput) error {
	if input == nil || input.CheckID == "" {
		return fmt.Errorf("check ID required")
	}

	check, err := uc.repo.GetByID(ctx, input.CheckID)
	if err != nil {
		return err
	}

	if err := check.Review(input.ReviewedBy, input.Reason); err != nil {
		return err
	}

	return uc.repo.Update(ctx, check)
}

type OverrideCheckInput struct {
	CheckID string
	Reason  string
}

func (uc *FraudUseCases) OverrideCheck(ctx context.Context, input *OverrideCheckInput) error {
	if input == nil || input.CheckID == "" {
		return fmt.Errorf("check ID required")
	}

	check, err := uc.repo.GetByID(ctx, input.CheckID)
	if err != nil {
		return err
	}

	if err := check.Override(input.Reason); err != nil {
		return err
	}

	return uc.repo.Update(ctx, check)
}
