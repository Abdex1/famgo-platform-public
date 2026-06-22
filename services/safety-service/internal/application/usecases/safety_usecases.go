// services/safety-service/internal/application/usecases/safety_usecases.go
package usecases

import (
	"context"
	"fmt"

	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/infrastructure/repositories"
)

type SafetyUseCases struct {
	repo *repositories.SOSRepository
}

func NewSafetyUseCases(repo *repositories.SOSRepository) *SafetyUseCases {
	return &SafetyUseCases{repo: repo}
}

type InitiateSOSInput struct {
	RideID       string
	UserID       string
	UserType     string
	Latitude     float64
	Longitude    float64
	IncidentType string
	Description  string
}

type SOSOutput struct {
	IncidentID      string
	Status          string
	EscalationLevel string
}

func (uc *SafetyUseCases) InitiateSOS(ctx context.Context, input *InitiateSOSInput) (*SOSOutput, error) {
	if input == nil || input.RideID == "" {
		return nil, fmt.Errorf("ride ID required")
	}

	incident, err := entities.NewSOSIncident(input.RideID, input.UserID, input.UserType,
		input.Latitude, input.Longitude, input.IncidentType)
	if err != nil {
		return nil, err
	}
	incident.Description = input.Description

	if err := incident.Activate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, incident); err != nil {
		return nil, err
	}

	return &SOSOutput{
		IncidentID:      incident.ID,
		Status:          string(incident.Status),
		EscalationLevel: string(incident.EscalationLevel),
	}, nil
}

type GetIncidentInput struct {
	IncidentID string
}

func (uc *SafetyUseCases) GetIncident(ctx context.Context, input *GetIncidentInput) (*SOSOutput, error) {
	if input == nil || input.IncidentID == "" {
		return nil, fmt.Errorf("incident ID required")
	}

	incident, err := uc.repo.GetByID(ctx, input.IncidentID)
	if err != nil {
		return nil, err
	}

	return &SOSOutput{
		IncidentID:      incident.ID,
		Status:          string(incident.Status),
		EscalationLevel: string(incident.EscalationLevel),
	}, nil
}

type EscalateIncidentInput struct {
	IncidentID      string
	EscalationLevel string
}

func (uc *SafetyUseCases) EscalateIncident(ctx context.Context, input *EscalateIncidentInput) error {
	if input == nil || input.IncidentID == "" {
		return fmt.Errorf("incident ID required")
	}

	incident, err := uc.repo.GetByID(ctx, input.IncidentID)
	if err != nil {
		return err
	}

	if err := incident.Escalate(entities.EscalationLevel(input.EscalationLevel)); err != nil {
		return err
	}

	return uc.repo.Update(ctx, incident)
}

type ResolveIncidentInput struct {
	IncidentID     string
	ResolvedBy     string
	ResolutionNote string
}

func (uc *SafetyUseCases) ResolveIncident(ctx context.Context, input *ResolveIncidentInput) error {
	if input == nil || input.IncidentID == "" {
		return fmt.Errorf("incident ID required")
	}

	incident, err := uc.repo.GetByID(ctx, input.IncidentID)
	if err != nil {
		return err
	}

	if err := incident.Resolve(input.ResolvedBy, input.ResolutionNote); err != nil {
		return err
	}

	return uc.repo.Update(ctx, incident)
}
