/*
# PHASE 4 — SUSPICIOUS LOGIN DETECTION

# =========================================================

# STEP 13 — CREATE RISK ENGINE


internal/domain/services/risk_service.go
*/
package services

type RiskService struct {
}

func NewRiskService() *RiskService {
	return &RiskService{}
}

func (r *RiskService) DetectImpossibleTravel(
	oldIP string,
	newIP string,
) bool {

	if oldIP != newIP {
		return true
	}

	return false
}
