/*
# PHASE 2 — AUDIT LOGGING PIPELINE

# =========================================================

# STEP 6 — CREATE AUDIT LOGGER


internal/domain/services/audit_service.go
*/
package services

import (
	"context"
	"log"
)

type AuditService struct {
}

func NewAuditService() *AuditService {
	return &AuditService{}
}

func (s *AuditService) Log(
	ctx context.Context,
	action string,
	userID string,
	metadata map[string]any,
) {

	log.Printf(
		"[AUDIT] action=%s user=%s metadata=%v",
		action,
		userID,
		metadata,
	)
}
