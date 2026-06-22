// services/fraud-service/internal/domain/entities/fraud_check.go
package entities

import (
	"fmt"
	"time"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type FraudCheck struct {
	ID                  string
	RideID              string
	UserID              string
	UserType            string
	RiskScore           float64
	RiskLevel           RiskLevel
	FlagsTriggered      []string
	LocationAnomalies   bool
	VelocityAnomalies   bool
	PaymentAnomalies    bool
	BehaviorAnomalies   bool
	IsBlacklisted       bool
	IsReview            bool
	ReviewReason        string
	ReviewBy            *string
	ReviewedAt          *time.Time
	Action              string
	IsManualOverride    bool
	ManualOverrideReason *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewFraudCheck(rideID, userID string, userType string) (*FraudCheck, error) {
	if rideID == "" || userID == "" {
		return nil, fmt.Errorf("ride ID and user ID required")
	}

	now := time.Now()
	return &FraudCheck{
		ID:             fmt.Sprintf("chk_%d", now.UnixNano()),
		RideID:         rideID,
		UserID:         userID,
		UserType:       userType,
		RiskScore:      0.0,
		RiskLevel:      RiskLow,
		FlagsTriggered: []string{},
		Action:         "allow",
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (f *FraudCheck) AddFlag(flag string) {
	for _, existing := range f.FlagsTriggered {
		if existing == flag {
			return // Already exists
		}
	}
	f.FlagsTriggered = append(f.FlagsTriggered, flag)
}

func (f *FraudCheck) CalculateRiskScore(highThreshold, mediumThreshold float64) {
	baseScore := 0.0

	// Location anomaly: +0.2
	if f.LocationAnomalies {
		baseScore += 0.2
		f.AddFlag("location_anomaly")
	}

	// Velocity anomaly: +0.25
	if f.VelocityAnomalies {
		baseScore += 0.25
		f.AddFlag("velocity_anomaly")
	}

	// Payment anomaly: +0.15
	if f.PaymentAnomalies {
		baseScore += 0.15
		f.AddFlag("payment_anomaly")
	}

	// Behavior anomaly: +0.1
	if f.BehaviorAnomalies {
		baseScore += 0.1
		f.AddFlag("behavior_anomaly")
	}

	// Blacklisted: +0.3
	if f.IsBlacklisted {
		baseScore += 0.3
		f.AddFlag("blacklisted_user")
	}

	f.RiskScore = baseScore

	if f.RiskScore >= highThreshold {
		f.RiskLevel = RiskHigh
		f.Action = "block"
		f.IsReview = true
	} else if f.RiskScore >= mediumThreshold {
		f.RiskLevel = RiskMedium
		f.Action = "review"
		f.IsReview = true
	} else {
		f.RiskLevel = RiskLow
		f.Action = "allow"
	}

	f.UpdatedAt = time.Now()
}

func (f *FraudCheck) Override(reason string) error {
	if f.IsManualOverride {
		return fmt.Errorf("already overridden")
	}
	f.IsManualOverride = true
	f.ManualOverrideReason = &reason
	f.UpdatedAt = time.Now()
	return nil
}

func (f *FraudCheck) Review(reviewedBy string, reason string) error {
	if !f.IsReview {
		return fmt.Errorf("no review needed")
	}
	now := time.Now()
	f.ReviewBy = &reviewedBy
	f.ReviewedAt = &now
	f.ReviewReason = reason
	f.UpdatedAt = now
	return nil
}
