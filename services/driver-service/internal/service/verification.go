//go:build week3

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"famgo/driver-service/internal/model"
	"famgo/driver-service/internal/repository"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// VerificationService handles driver KYC, training, and compliance verification
// WEEK 3: Added alongside existing DriverService (NOT replacing it)
type VerificationService struct {
	verificationRepo *repository.VerificationRepository
	documentRepo     *repository.DocumentRepository
	trainingRepo     *repository.TrainingRepository
	backgroundRepo   *repository.BackgroundCheckRepository
	driverRepo       *repository.DriverRepository
	logger           logger.Logger
}

// NewVerificationService creates a new verification service
func NewVerificationService(
	verificationRepo *repository.VerificationRepository,
	documentRepo *repository.DocumentRepository,
	trainingRepo *repository.TrainingRepository,
	backgroundRepo *repository.BackgroundCheckRepository,
	driverRepo *repository.DriverRepository,
	log logger.Logger,
) *VerificationService {
	return &VerificationService{
		verificationRepo: verificationRepo,
		documentRepo:     documentRepo,
		trainingRepo:     trainingRepo,
		backgroundRepo:   backgroundRepo,
		driverRepo:       driverRepo,
		logger:           log,
	}
}

// InitiateKYC starts the KYC process for a driver
func (s *VerificationService) InitiateKYC(ctx context.Context, driverID string) error {
	s.logger.Info("initiating KYC", map[string]interface{}{"driver_id": driverID})

	// Check if driver exists
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		s.logger.Error("driver not found", map[string]interface{}{"driver_id": driverID})
		return err
	}

	// Create verification record if not exists
	verification := &model.DriverVerification{
		DriverID:           driverID,
		KYCStatus:          "pending",
		VerificationStatus: "pending",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.verificationRepo.CreateVerification(ctx, verification); err != nil {
		s.logger.Error("failed to create verification record", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("KYC initiated", map[string]interface{}{"driver_id": driverID, "kyc_status": "pending"})
	return nil
}

// CheckVerificationStatus returns current verification status
func (s *VerificationService) CheckVerificationStatus(ctx context.Context, driverID string) (*model.DriverVerification, error) {
	verification, err := s.verificationRepo.GetVerificationByDriverID(ctx, driverID)
	if err != nil {
		s.logger.Warn("verification not found", map[string]interface{}{"driver_id": driverID})
		return nil, err
	}

	return verification, nil
}

// ApproveDriver approves a driver (moves from pending to approved)
func (s *VerificationService) ApproveDriver(ctx context.Context, driverID string, reason string) error {
	s.logger.Info("approving driver", map[string]interface{}{"driver_id": driverID})

	// Update verification status
	if err := s.verificationRepo.UpdateVerificationStatus(ctx, driverID, "approved"); err != nil {
		s.logger.Error("failed to update verification status", map[string]interface{}{"error": err})
		return err
	}

	// Update driver status (use state machine from existing DriverService)
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return err
	}

	// Transition to approved state
	driver.Status = model.DriverStateApproved
	if err := s.driverRepo.UpdateDriver(ctx, driver); err != nil {
		s.logger.Error("failed to update driver status", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("driver approved", map[string]interface{}{"driver_id": driverID})
	return nil
}

// RejectDriver rejects a driver with reason
func (s *VerificationService) RejectDriver(ctx context.Context, driverID string, reason string) error {
	s.logger.Info("rejecting driver", map[string]interface{}{"driver_id": driverID, "reason": reason})

	verification := &model.DriverVerification{
		DriverID:       driverID,
		KYCStatus:      "rejected",
		RejectedReason: reason,
		UpdatedAt:      time.Now(),
	}

	if err := s.verificationRepo.UpdateVerification(ctx, verification); err != nil {
		s.logger.Error("failed to update verification", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("driver rejected", map[string]interface{}{"driver_id": driverID})
	return nil
}

// CheckComplianceStatus returns compliance checklist status
func (s *VerificationService) CheckComplianceStatus(ctx context.Context, driverID string) (*model.ComplianceChecklist, error) {
	verification, err := s.verificationRepo.GetVerificationByDriverID(ctx, driverID)
	if err != nil {
		return nil, err
	}

	return &verification.ComplianceChecklist, nil
}

// EnforceCompliance validates all compliance items
func (s *VerificationService) EnforceCompliance(ctx context.Context, driverID string) error {
	s.logger.Info("enforcing compliance", map[string]interface{}{"driver_id": driverID})

	verification, err := s.verificationRepo.GetVerificationByDriverID(ctx, driverID)
	if err != nil {
		return err
	}

	// Check each compliance item
	if !verification.ComplianceChecklist.BackgroundCheckPassed {
		return errors.New("background check not passed")
	}

	if !verification.ComplianceChecklist.VehicleInspectionPassed {
		return errors.New("vehicle inspection not passed")
	}

	if !verification.ComplianceChecklist.InsuranceVerified {
		return errors.New("insurance not verified")
	}

	if !verification.ComplianceChecklist.TrainingCompleted {
		return errors.New("training not completed")
	}

	if !verification.ComplianceChecklist.TermsAccepted {
		return errors.New("terms not accepted")
	}

	s.logger.Info("compliance check passed", map[string]interface{}{"driver_id": driverID})
	return nil
}

// DocumentService handles document upload and verification
type DocumentService struct {
	documentRepo *repository.DocumentRepository
	logger       logger.Logger
}

// NewDocumentService creates a new document service
func NewDocumentService(documentRepo *repository.DocumentRepository, log logger.Logger) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		logger:       log,
	}
}

// UploadDocument uploads a new document for a driver
func (s *DocumentService) UploadDocument(ctx context.Context, driverID string, docType string, documentURL string) error {
	s.logger.Info("uploading document", map[string]interface{}{
		"driver_id":   driverID,
		"doc_type":    docType,
		"document_url": documentURL,
	})

	doc := &model.DriverDocument{
		DriverID:           driverID,
		DocumentType:       docType,
		DocumentURL:        documentURL,
		UploadDate:         time.Now(),
		VerificationStatus: "pending",
		CreatedAt:          time.Now(),
	}

	if err := s.documentRepo.CreateDocument(ctx, doc); err != nil {
		s.logger.Error("failed to upload document", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("document uploaded", map[string]interface{}{"doc_type": docType})
	return nil
}

// GetDocuments retrieves all documents for a driver
func (s *DocumentService) GetDocuments(ctx context.Context, driverID string) ([]*model.DriverDocument, error) {
	docs, err := s.documentRepo.GetDocumentsByDriverID(ctx, driverID)
	if err != nil {
		s.logger.Warn("failed to get documents", map[string]interface{}{"driver_id": driverID})
		return nil, err
	}

	return docs, nil
}

// VerifyDocument marks a document as verified or rejected
func (s *DocumentService) VerifyDocument(ctx context.Context, docID string, approved bool, reason string) error {
	s.logger.Info("verifying document", map[string]interface{}{
		"doc_id":   docID,
		"approved": approved,
	})

	status := "approved"
	if !approved {
		status = "rejected"
	}

	if err := s.documentRepo.VerifyDocument(ctx, docID, status, reason); err != nil {
		s.logger.Error("failed to verify document", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("document verified", map[string]interface{}{"status": status})
	return nil
}

// TrainingService handles driver training completion
type TrainingService struct {
	trainingRepo *repository.TrainingRepository
	logger       logger.Logger
}

// NewTrainingService creates a new training service
func NewTrainingService(trainingRepo *repository.TrainingRepository, log logger.Logger) *TrainingService {
	return &TrainingService{
		trainingRepo: trainingRepo,
		logger:       log,
	}
}

// GetTrainingProgress retrieves training progress for a driver
func (s *TrainingService) GetTrainingProgress(ctx context.Context, driverID string) ([]*model.DriverTraining, error) {
	trainings, err := s.trainingRepo.GetTrainingByDriverID(ctx, driverID)
	if err != nil {
		s.logger.Warn("failed to get training progress", map[string]interface{}{"driver_id": driverID})
		return nil, err
	}

	return trainings, nil
}

// CompleteModule marks a training module as complete
func (s *TrainingService) CompleteModule(ctx context.Context, trainingID string, moduleNum int) error {
	s.logger.Info("completing training module", map[string]interface{}{
		"training_id": trainingID,
		"module":      moduleNum,
	})

	if err := s.trainingRepo.CompleteModule(ctx, trainingID, moduleNum); err != nil {
		s.logger.Error("failed to complete module", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("module completed", map[string]interface{}{"module": moduleNum})
	return nil
}

// ComplianceService handles compliance checking
type ComplianceService struct {
	verificationRepo *repository.VerificationRepository
	logger           logger.Logger
}

// NewComplianceService creates a new compliance service
func NewComplianceService(verificationRepo *repository.VerificationRepository, log logger.Logger) *ComplianceService {
	return &ComplianceService{
		verificationRepo: verificationRepo,
		logger:           log,
	}
}

// CheckAllCompliance validates all compliance items for a driver
func (s *ComplianceService) CheckAllCompliance(ctx context.Context, driverID string) (bool, error) {
	s.logger.Info("checking all compliance", map[string]interface{}{"driver_id": driverID})

	verification, err := s.verificationRepo.GetVerificationByDriverID(ctx, driverID)
	if err != nil {
		return false, err
	}

	checklist := verification.ComplianceChecklist

	// All items must be true
	allCompliant := checklist.BackgroundCheckPassed &&
		checklist.VehicleInspectionPassed &&
		checklist.InsuranceVerified &&
		checklist.TrainingCompleted &&
		checklist.TermsAccepted

	s.logger.Info("compliance check complete", map[string]interface{}{
		"driver_id":   driverID,
		"compliant":   allCompliant,
		"bg_check":    checklist.BackgroundCheckPassed,
		"inspection":  checklist.VehicleInspectionPassed,
		"insurance":   checklist.InsuranceVerified,
		"training":    checklist.TrainingCompleted,
		"terms":       checklist.TermsAccepted,
	})

	return allCompliant, nil
}

// AcceptTerms marks terms as accepted in compliance checklist
func (s *ComplianceService) AcceptTerms(ctx context.Context, driverID string) error {
	s.logger.Info("accepting terms", map[string]interface{}{"driver_id": driverID})

	verification, err := s.verificationRepo.GetVerificationByDriverID(ctx, driverID)
	if err != nil {
		return err
	}

	verification.ComplianceChecklist.TermsAccepted = true
	verification.ComplianceChecklist.ComplianceSignedAt = &[]time.Time{time.Now()}[0]

	if err := s.verificationRepo.UpdateVerification(ctx, verification); err != nil {
		s.logger.Error("failed to update compliance", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("terms accepted", map[string]interface{}{"driver_id": driverID})
	return nil
}
