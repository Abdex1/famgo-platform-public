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
	"github.com/shopspring/decimal"
)

// EarningsService handles driver earnings tracking and settlements
// WEEK 3: Added alongside existing DriverService (NOT replacing it)
type EarningsService struct {
	earningsRepo  *repository.EarningsRepository
	settlementRepo *repository.SettlementRepository
	driverRepo    *repository.DriverRepository
	logger        logger.Logger
}

// NewEarningsService creates a new earnings service
func NewEarningsService(
	earningsRepo *repository.EarningsRepository,
	settlementRepo *repository.SettlementRepository,
	driverRepo *repository.DriverRepository,
	log logger.Logger,
) *EarningsService {
	return &EarningsService{
		earningsRepo:  earningsRepo,
		settlementRepo: settlementRepo,
		driverRepo:    driverRepo,
		logger:        log,
	}
}

// RecordEarning records earnings from a completed trip
func (s *EarningsService) RecordEarning(ctx context.Context, driverID string, tripID string, grossAmount decimal.Decimal) error {
	s.logger.Info("recording earning", map[string]interface{}{
		"driver_id":    driverID,
		"trip_id":      tripID,
		"gross_amount": grossAmount.String(),
	})

	// Calculate fees and taxes
	platformFee := calculatePlatformFee(grossAmount)
	taxAmount := calculateTax(grossAmount)
	netAmount := grossAmount.Sub(platformFee).Sub(taxAmount)

	earning := &model.DriverEarning{
		DriverID:      driverID,
		TripID:        tripID,
		GrossAmount:   grossAmount,
		PlatformFee:   platformFee,
		TaxAmount:     taxAmount,
		NetAmount:     netAmount,
		Currency:      "ETB",
		PaymentStatus: "pending",
		CreatedAt:     time.Now(),
	}

	if err := s.earningsRepo.RecordEarning(ctx, earning); err != nil {
		s.logger.Error("failed to record earning", map[string]interface{}{"error": err})
		return err
	}

	// Update driver total earnings
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err == nil {
		currentEarnings := driver.TotalEarnings
		driver.TotalEarnings = currentEarnings + netAmount.InexactFloat64()
		s.driverRepo.UpdateDriver(ctx, driver)
	}

	s.logger.Info("earning recorded", map[string]interface{}{
		"net_amount": netAmount.String(),
	})

	return nil
}

// GetTotalEarnings returns total earnings for a driver
func (s *EarningsService) GetTotalEarnings(ctx context.Context, driverID string) (decimal.Decimal, error) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return decimal.Zero, err
	}

	total := decimal.NewFromFloat(driver.TotalEarnings)
	return total, nil
}

// GetEarningsByPeriod returns earnings for a specific period
func (s *EarningsService) GetEarningsByPeriod(ctx context.Context, driverID string, startDate, endDate time.Time) ([]*model.DriverEarning, error) {
	earnings, err := s.earningsRepo.GetEarningsByPeriod(ctx, driverID, startDate, endDate)
	if err != nil {
		s.logger.Warn("failed to get earnings by period", map[string]interface{}{"error": err})
		return nil, err
	}

	return earnings, nil
}

// GenerateSettlement creates a settlement for a period
func (s *EarningsService) GenerateSettlement(ctx context.Context, driverID string, settlementPeriod string) (*model.DriverSettlement, error) {
	s.logger.Info("generating settlement", map[string]interface{}{
		"driver_id": driverID,
		"period":    settlementPeriod,
	})

	// Get all earnings for period
	earnings, err := s.earningsRepo.GetEarningsByDriver(ctx, driverID)
	if err != nil {
		return nil, err
	}

	// Calculate totals
	var totalGross, totalFees, totalTaxes decimal.Decimal
	var totalTrips int

	for _, earning := range earnings {
		totalGross = totalGross.Add(earning.GrossAmount)
		totalFees = totalFees.Add(earning.PlatformFee)
		totalTaxes = totalTaxes.Add(earning.TaxAmount)
		totalTrips++
	}

	totalNet := totalGross.Sub(totalFees).Sub(totalTaxes)

	settlement := &model.DriverSettlement{
		DriverID:         driverID,
		SettlementPeriod: settlementPeriod,
		TotalTrips:       totalTrips,
		TotalGross:       totalGross,
		TotalFees:        totalFees,
		TotalTaxes:       totalTaxes,
		TotalNet:         totalNet,
		PaymentMethod:    "bank_transfer",
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	if err := s.settlementRepo.CreateSettlement(ctx, settlement); err != nil {
		s.logger.Error("failed to create settlement", map[string]interface{}{"error": err})
		return nil, err
	}

	s.logger.Info("settlement generated", map[string]interface{}{
		"period":   settlementPeriod,
		"net":      totalNet.String(),
		"trips":    totalTrips,
	})

	return settlement, nil
}

// ProcessPayment processes settlement payment
func (s *EarningsService) ProcessPayment(ctx context.Context, settlementID string) error {
	s.logger.Info("processing payment", map[string]interface{}{"settlement_id": settlementID})

	// Mark settlement as completed
	if err := s.settlementRepo.ProcessPayment(ctx, settlementID); err != nil {
		s.logger.Error("failed to process payment", map[string]interface{}{"error": err})
		return err
	}

	s.logger.Info("payment processed", map[string]interface{}{"settlement_id": settlementID})
	return nil
}

// RatingService handles driver ratings and aggregation
type RatingService struct {
	ratingRepo *repository.RatingRepository
	driverRepo *repository.DriverRepository
	logger     logger.Logger
}

// NewRatingService creates a new rating service
func NewRatingService(ratingRepo *repository.RatingRepository, driverRepo *repository.DriverRepository, log logger.Logger) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		driverRepo: driverRepo,
		logger:     log,
	}
}

// AddRating adds a new rating for a driver
func (s *RatingService) AddRating(ctx context.Context, driverID string, tripID string, riderID string, rating int, comment string) error {
	s.logger.Info("adding rating", map[string]interface{}{
		"driver_id": driverID,
		"trip_id":   tripID,
		"rating":    rating,
	})

	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	ratingRecord := &model.DriverRating{
		DriverID:  driverID,
		TripID:    tripID,
		RiderID:   riderID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	if err := s.ratingRepo.CreateRating(ctx, ratingRecord); err != nil {
		s.logger.Error("failed to add rating", map[string]interface{}{"error": err})
		return err
	}

	// Update rating summary (would run every 5 minutes in production)
	_ = s.ratingRepo.UpdateRatingSummary(ctx, driverID)

	s.logger.Info("rating added", map[string]interface{}{"rating": rating})
	return nil
}

// GetAverageRating returns driver's average rating
func (s *RatingService) GetAverageRating(ctx context.Context, driverID string) (float64, error) {
	summary, err := s.ratingRepo.GetRatingSummary(ctx, driverID)
	if err != nil {
		return 0, err
	}

	return summary.AverageRating, nil
}

// GetRatingDistribution returns rating breakdown
func (s *RatingService) GetRatingDistribution(ctx context.Context, driverID string) (*model.RatingSummary, error) {
	summary, err := s.ratingRepo.GetRatingSummary(ctx, driverID)
	if err != nil {
		s.logger.Warn("failed to get rating distribution", map[string]interface{}{"error": err})
		return nil, err
	}

	return summary, nil
}

// FinancialReportService handles financial reporting
type FinancialReportService struct {
	earningsRepo  *repository.EarningsRepository
	settlementRepo *repository.SettlementRepository
	logger        logger.Logger
}

// NewFinancialReportService creates a new financial report service
func NewFinancialReportService(
	earningsRepo *repository.EarningsRepository,
	settlementRepo *repository.SettlementRepository,
	log logger.Logger,
) *FinancialReportService {
	return &FinancialReportService{
		earningsRepo:  earningsRepo,
		settlementRepo: settlementRepo,
		logger:        log,
	}
}

// GenerateMonthlyReport generates financial report for a month
func (s *FinancialReportService) GenerateMonthlyReport(ctx context.Context, driverID string, month string) (map[string]interface{}, error) {
	s.logger.Info("generating monthly report", map[string]interface{}{
		"driver_id": driverID,
		"month":     month,
	})

	// Get settlements for month
	settlements, err := s.settlementRepo.GetSettlements(ctx, driverID)
	if err != nil {
		return nil, err
	}

	// Filter by month and calculate totals
	var totalGross, totalFees, totalTaxes decimal.Decimal
	var totalTrips int

	for _, settlement := range settlements {
		if settlement.SettlementPeriod == month {
			totalGross = totalGross.Add(settlement.TotalGross)
			totalFees = totalFees.Add(settlement.TotalFees)
			totalTaxes = totalTaxes.Add(settlement.TotalTaxes)
			totalTrips += settlement.TotalTrips
		}
	}

	report := map[string]interface{}{
		"period":      month,
		"total_trips": totalTrips,
		"gross":       totalGross.String(),
		"fees":        totalFees.String(),
		"taxes":       totalTaxes.String(),
		"net":         totalGross.Sub(totalFees).Sub(totalTaxes).String(),
	}

	s.logger.Info("monthly report generated", map[string]interface{}{
		"month": month,
		"trips": totalTrips,
	})

	return report, nil
}

// CalculateTaxWithholding calculates taxes for a period
func (s *FinancialReportService) CalculateTaxWithholding(ctx context.Context, driverID string, period string) (decimal.Decimal, error) {
	settlements, err := s.settlementRepo.GetSettlements(ctx, driverID)
	if err != nil {
		return decimal.Zero, err
	}

	var totalTaxes decimal.Decimal

	for _, settlement := range settlements {
		if settlement.SettlementPeriod == period {
			totalTaxes = totalTaxes.Add(settlement.TotalTaxes)
		}
	}

	return totalTaxes, nil
}

// Helper functions for fee and tax calculations

func calculatePlatformFee(grossAmount decimal.Decimal) decimal.Decimal {
	// Platform takes 20% commission
	feePercentage := decimal.NewFromFloat(0.20)
	return grossAmount.Mul(feePercentage)
}

func calculateTax(grossAmount decimal.Decimal) decimal.Decimal {
	// Income tax calculation (simplified - 3% of platform fee)
	platformFee := calculatePlatformFee(grossAmount)
	taxPercentage := decimal.NewFromFloat(0.03)
	return platformFee.Mul(taxPercentage)
}
