// services/safety-service/internal/domain/entities/sos_incident.go
package entities

import (
	"fmt"
	"time"
)

type IncidentStatus string

const (
	StatusInitiated  IncidentStatus = "initiated"
	StatusActive     IncidentStatus = "active"
	StatusEscalated  IncidentStatus = "escalated"
	StatusResolved   IncidentStatus = "resolved"
	StatusCancelled  IncidentStatus = "cancelled"
)

type EscalationLevel string

const (
	LevelEmergencyContact EscalationLevel = "emergency_contact"
	LevelPolice           EscalationLevel = "police"
	LevelAmbulance        EscalationLevel = "ambulance"
	LevelPlatformSupport  EscalationLevel = "platform_support"
)

type SOSIncident struct {
	ID                 string
	RideID             string
	UserID             string
	UserType           string
	Latitude           float64
	Longitude          float64
	Status             IncidentStatus
	IncidentType       string
	Description        string
	EscalationLevel    EscalationLevel
	EmergencyContacts  []string
	NotificationsSent  int
	PoliceNotified     bool
	PoliceNotifiedAt   *time.Time
	AmbulanceNotified  bool
	AmbulanceNotifiedAt *time.Time
	PlatformNotified   bool
	PlatformNotifiedAt *time.Time
	ResolvedAt         *time.Time
	ResolvedBy         *string
	ResolutionNotes    string
	InitiatedAt        time.Time
	LastEscalatedAt    *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}

func NewSOSIncident(rideID, userID, userType string, lat, lon float64, incidentType string) (*SOSIncident, error) {
	if rideID == "" || userID == "" {
		return nil, fmt.Errorf("ride ID and user ID required")
	}
	now := time.Now()
	return &SOSIncident{
		ID:              fmt.Sprintf("sos_%d", now.UnixNano()),
		RideID:          rideID,
		UserID:          userID,
		UserType:        userType,
		Latitude:        lat,
		Longitude:       lon,
		Status:          StatusInitiated,
		IncidentType:    incidentType,
		EscalationLevel: LevelEmergencyContact,
		InitiatedAt:     now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (s *SOSIncident) Activate() error {
	if s.Status != StatusInitiated {
		return fmt.Errorf("can only activate initiated incidents")
	}
	s.Status = StatusActive
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SOSIncident) Escalate(level EscalationLevel) error {
	if s.Status != StatusActive && s.Status != StatusEscalated {
		return fmt.Errorf("can only escalate active incidents")
	}
	s.Status = StatusEscalated
	s.EscalationLevel = level
	now := time.Now()
	s.LastEscalatedAt = &now
	s.UpdatedAt = now
	return nil
}

func (s *SOSIncident) Resolve(resolvedBy string, notes string) error {
	s.Status = StatusResolved
	s.ResolvedBy = &resolvedBy
	s.ResolutionNotes = notes
	now := time.Now()
	s.ResolvedAt = &now
	s.UpdatedAt = now
	return nil
}

func (s *SOSIncident) Cancel() error {
	if s.Status == StatusResolved {
		return fmt.Errorf("cannot cancel resolved incident")
	}
	s.Status = StatusCancelled
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SOSIncident) NotifyPolice() error {
	if s.PoliceNotified {
		return fmt.Errorf("police already notified")
	}
	now := time.Now()
	s.PoliceNotified = true
	s.PoliceNotifiedAt = &now
	s.NotificationsSent++
	s.UpdatedAt = now
	return nil
}

func (s *SOSIncident) NotifyAmbulance() error {
	if s.AmbulanceNotified {
		return fmt.Errorf("ambulance already notified")
	}
	now := time.Now()
	s.AmbulanceNotified = true
	s.AmbulanceNotifiedAt = &now
	s.NotificationsSent++
	s.UpdatedAt = now
	return nil
}

func (s *SOSIncident) NotifyPlatform() error {
	if s.PlatformNotified {
		return fmt.Errorf("platform already notified")
	}
	now := time.Now()
	s.PlatformNotified = true
	s.PlatformNotifiedAt = &now
	s.NotificationsSent++
	s.UpdatedAt = now
	return nil
}
