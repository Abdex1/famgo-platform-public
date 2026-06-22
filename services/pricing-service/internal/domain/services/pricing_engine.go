package services

import (
	"time"

	"github.com/Abdex1/FamGo-platform/services/pricing-service/internal/domain/entities"
	"github.com/google/uuid"
)

// PricingEngine calculates fares with all components
type PricingEngine struct {
	repo entities.PricingRepository
}

// NewPricingEngine creates pricing engine
func NewPricingEngine(repo entities.PricingRepository) *PricingEngine {
	return &PricingEngine{repo: repo}
}

// CalculateFare calculates complete fare with all components
// Formula: Fare = BaseFare + (Distance × Rate) + (Duration × Rate) + Surge + Tax - Discount
func (pe *PricingEngine) CalculateFare(
	rideID, rideType string,
	distanceMeters, durationSeconds int,
	pickupLat, pickupLng, dropoffLat, dropoffLng float64,
	isPool bool,
	activeRides, availableDrivers int,
	discountCode *entities.DiscountCode,
) *entities.FareCalculation {

	// City detection (in production, use geolocation service)
	city := "Addis Ababa"

	// Get pricing rule
	rule := pe.getPricingRule(rideType, city)
	if rule == nil {
		rule = pe.getDefaultRule(rideType)
	}

	// Apply pool discount to rates if pooling
	var distanceRate, timeRate = rule.DistanceRate, rule.TimeRate
	if isPool {
		poolDiscountFactor := 1.0 - (rule.PoolDiscount / 100.0)
		distanceRate *= poolDiscountFactor
		timeRate *= poolDiscountFactor
	}

	// Calculate fare components
	baseFare := rule.BaseFare
	distanceFare := calculateDistanceFare(distanceMeters, distanceRate)
	timeFare := calculateTimeFare(durationSeconds, timeRate)

	// Subtotal before surge
	subtotalBeforeSurge := baseFare + distanceFare + timeFare
	if subtotalBeforeSurge < rule.MinimumFare {
		subtotalBeforeSurge = rule.MinimumFare
	}

	// Calculate surge
	surgeMultiplier := pe.CalculateSurgeMultiplier(city, pickupLat, pickupLng, activeRides, availableDrivers)
	surgeAmount := subtotalBeforeSurge * (surgeMultiplier - 1.0)
	if surgeMultiplier > rule.SurgeFactorMax {
		surgeAmount = subtotalBeforeSurge * (rule.SurgeFactorMax - 1.0)
		surgeMultiplier = rule.SurgeFactorMax
	}

	// Total before tax
	totalBeforeTax := subtotalBeforeSurge + surgeAmount

	// Calculate taxes
	taxes := totalBeforeTax * (rule.TaxPercentage / 100.0)

	// Apply discount
	var discountAmount float64
	var discountCodeID *string
	if discountCode != nil {
		discountAmount = pe.calculateDiscount(totalBeforeTax+taxes, discountCode)
		discountCodeID = &discountCode.ID
	}

	// Final fare
	finalFare := totalBeforeTax + taxes - discountAmount
	if finalFare < rule.MinimumFare {
		finalFare = rule.MinimumFare
	}

	return &entities.FareCalculation{
		ID:                  uuid.New().String(),
		RideID:              rideID,
		RideType:            rideType,
		DistanceMeters:      distanceMeters,
		DurationSeconds:     durationSeconds,
		PickupLat:           pickupLat,
		PickupLng:           pickupLng,
		DropoffLat:          dropoffLat,
		DropoffLng:          dropoffLng,
		BaseFare:            baseFare,
		DistanceFare:        distanceFare,
		TimeFare:            timeFare,
		SubtotalBeforeSurge: subtotalBeforeSurge,
		SurgeMultiplier:     surgeMultiplier,
		SurgeAmount:         surgeAmount,
		Taxes:               taxes,
		DiscountCodeID:      discountCodeID,
		DiscountAmount:      discountAmount,
		FinalFare:           finalFare,
		IsPool:              isPool,
		City:                city,
		CalculatedAt:        time.Now(),
		CreatedAt:           time.Now(),
	}
}

// EstimateFare returns a quick fare estimate without persistence
func (pe *PricingEngine) EstimateFare(
	rideType string,
	distanceMeters int,
	activeRides, availableDrivers int,
	isPool bool,
) map[string]interface{} {

	city := "Addis Ababa"
	rule := pe.getPricingRule(rideType, city)
	if rule == nil {
		rule = pe.getDefaultRule(rideType)
	}

	// Apply pool discount
	var distanceRate = rule.DistanceRate
	if isPool {
		poolDiscountFactor := 1.0 - (rule.PoolDiscount / 100.0)
		distanceRate *= poolDiscountFactor
	}

	// Calculate
	baseFare := rule.BaseFare
	distanceFare := calculateDistanceFare(distanceMeters, distanceRate)
	subtotal := baseFare + distanceFare
	if subtotal < rule.MinimumFare {
		subtotal = rule.MinimumFare
	}

	surgeMultiplier := pe.CalculateSurgeMultiplier(city, 0, 0, activeRides, availableDrivers)
	surgeAmount := subtotal * (surgeMultiplier - 1.0)

	taxes := (subtotal + surgeAmount) * (rule.TaxPercentage / 100.0)
	finalFare := subtotal + surgeAmount + taxes

	return map[string]interface{}{
		"base_fare":        baseFare,
		"distance_fare":    distanceFare,
		"subtotal":         subtotal,
		"surge_multiplier": surgeMultiplier,
		"surge_amount":     surgeAmount,
		"taxes":            taxes,
		"final_fare":       finalFare,
	}
}

// CalculateSurgeMultiplier calculates surge based on supply/demand + time
func (pe *PricingEngine) CalculateSurgeMultiplier(
	city string,
	latitude, longitude float64,
	activeRides, availableDrivers int,
) float64 {

	// Time-based surge
	currentHour := time.Now().Hour()
	timeMultiplier := 1.0

	// Peak hours: 6-9 AM, 5-8 PM
	if (currentHour >= 6 && currentHour < 9) || (currentHour >= 17 && currentHour < 20) {
		timeMultiplier = 1.5
	}

	// Supply-demand surge
	demandSupplyRatio := calculateDemandSupplyRatio(activeRides, availableDrivers)
	supplyDemandMultiplier := 1.0

	if demandSupplyRatio > 0 {
		// Formula: 1.0 + (ratio - 1.0) * 0.5 with diminishing returns
		supplyDemandMultiplier = 1.0 + (demandSupplyRatio - 1.0) * 0.5
	}

	// Combined with weights: Time (40%), Supply-Demand (60%)
	combinedMultiplier := (timeMultiplier * 0.4) + (supplyDemandMultiplier * 0.6)

	// Clamp between 1.0 and 5.0
	if combinedMultiplier < 1.0 {
		combinedMultiplier = 1.0
	}
	if combinedMultiplier > 5.0 {
		combinedMultiplier = 5.0
	}

	return combinedMultiplier
}

// GetSurgeMultiplierAtLocation calculates surge for specific location
func (pe *PricingEngine) GetSurgeMultiplierAtLocation(
	latitude, longitude float64,
	activeRides, availableDrivers int,
) float64 {
	return pe.CalculateSurgeMultiplier("Addis Ababa", latitude, longitude, activeRides, availableDrivers)
}

// Helper functions

func calculateDistanceFare(distanceMeters int, ratePerKm float64) float64 {
	distanceKm := float64(distanceMeters) / 1000.0
	return distanceKm * ratePerKm
}

func calculateTimeFare(durationSeconds int, ratePerMinute float64) float64 {
	durationMinutes := float64(durationSeconds) / 60.0
	return durationMinutes * ratePerMinute
}

func calculateDemandSupplyRatio(activeRides, availableDrivers int) float64 {
	if availableDrivers == 0 {
		return 5.0 // Maximum surge if no drivers
	}
	return float64(activeRides) / float64(availableDrivers)
}

func (pe *PricingEngine) calculateDiscount(totalAmount float64, code *entities.DiscountCode) float64 {
	if code == nil {
		return 0
	}

	var discount float64

	if code.DiscountType == "FIXED" {
		discount = code.DiscountValue
	} else if code.DiscountType == "PERCENTAGE" {
		discount = totalAmount * (code.DiscountValue / 100.0)
		if code.MaxDiscount != nil && discount > *code.MaxDiscount {
			discount = *code.MaxDiscount
		}
	}

	return discount
}

// getPricingRule retrieves active pricing rule (placeholder - would use context)
func (pe *PricingEngine) getPricingRule(rideType, city string) *entities.PricingRule {
	// In production, use context and repository with caching
	return nil
}

// getDefaultRule returns default pricing rule
func (pe *PricingEngine) getDefaultRule(rideType string) *entities.PricingRule {
	rules := map[string]*entities.PricingRule{
		"ECONOMY": {
			BaseFare:       20.0,
			DistanceRate:   10.0,
			TimeRate:       0.33,
			MinimumFare:    15.0,
			SurgeFactorMax: 5.0,
			TaxPercentage:  2.0,
			PoolDiscount:   25.0,
		},
		"COMFORT": {
			BaseFare:       30.0,
			DistanceRate:   13.0,
			TimeRate:       0.43,
			MinimumFare:    25.0,
			SurgeFactorMax: 5.0,
			TaxPercentage:  2.0,
			PoolDiscount:   20.0,
		},
		"BUSINESS": {
			BaseFare:       40.0,
			DistanceRate:   18.0,
			TimeRate:       0.60,
			MinimumFare:    35.0,
			SurgeFactorMax: 5.0,
			TaxPercentage:  2.0,
			PoolDiscount:   15.0,
		},
		"POOL": {
			BaseFare:       15.0,
			DistanceRate:   8.0,
			TimeRate:       0.25,
			MinimumFare:    10.0,
			SurgeFactorMax: 5.0,
			TaxPercentage:  2.0,
			PoolDiscount:   0.0,
		},
	}

	if rule, exists := rules[rideType]; exists {
		return rule
	}

	// Fallback to ECONOMY
	return rules["ECONOMY"]
}
