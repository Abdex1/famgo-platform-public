package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Abdex1/FamGo-platform/services/pricing-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/pricing-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/pricing-service/internal/infrastructure/postgres"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PricingHandler handles pricing HTTP requests
type PricingHandler struct {
	pricingRepo    *postgres.PricingRuleRepository
	pricingEngine  *services.PricingEngine
}

// NewPricingHandler creates new pricing handler
func NewPricingHandler(pricingRepo *postgres.PricingRuleRepository, engine *services.PricingEngine) *PricingHandler {
	return &PricingHandler{
		pricingRepo:   pricingRepo,
		pricingEngine: engine,
	}
}

// CalculateFareRequest DTO
type CalculateFareRequest struct {
	RideID          string  `json:"ride_id"`
	RideType        string  `json:"ride_type"`
	DistanceMeters  int     `json:"distance_meters"`
	DurationSeconds int     `json:"duration_seconds"`
	PickupLat       float64 `json:"pickup_lat"`
	PickupLng       float64 `json:"pickup_lng"`
	DropoffLat      float64 `json:"dropoff_lat"`
	DropoffLng      float64 `json:"dropoff_lng"`
	IsPool          bool    `json:"is_pool"`
	ActiveRides     int     `json:"active_rides"`
	AvailableDrivers int    `json:"available_drivers"`
	DiscountCode    string  `json:"discount_code,omitempty"`
}

// CalculateFare handler - calculates complete fare with all components
func (h *PricingHandler) CalculateFare(w http.ResponseWriter, r *http.Request) {
	var req CalculateFareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get discount code if provided
	var discountCode *entities.DiscountCode
	if req.DiscountCode != "" {
		var err error
		discountCode, err = h.pricingRepo.ValidateDiscountCode(r.Context(), req.DiscountCode)
		if err != nil {
			// Discount not found or expired - proceed without it
			discountCode = nil
		}
	}

	// Calculate fare
	fare := h.pricingEngine.CalculateFare(
		req.RideID, req.RideType,
		req.DistanceMeters, req.DurationSeconds,
		req.PickupLat, req.PickupLng, req.DropoffLat, req.DropoffLng,
		req.IsPool,
		req.ActiveRides, req.AvailableDrivers,
		discountCode,
	)

	// Save to database
	if err := h.pricingRepo.SaveFareCalculation(r.Context(), fare); err != nil {
		http.Error(w, "Failed to save fare calculation", http.StatusInternalServerError)
		return
	}

	// Decrement discount code usage if applied
	if discountCode != nil {
		_ = h.pricingRepo.DecrementDiscountCodeUsage(r.Context(), discountCode.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"calculation_id":     fare.ID,
		"ride_id":            fare.RideID,
		"base_fare":          fare.BaseFare,
		"distance_fare":      fare.DistanceFare,
		"time_fare":          fare.TimeFare,
		"subtotal":           fare.SubtotalBeforeSurge,
		"surge_multiplier":   fare.SurgeMultiplier,
		"surge_amount":       fare.SurgeAmount,
		"taxes":              fare.Taxes,
		"discount_amount":    fare.DiscountAmount,
		"final_fare":         fare.FinalFare,
		"is_pool":            fare.IsPool,
		"calculated_at":      fare.CalculatedAt,
	})
}

// EstimateFareRequest DTO
type EstimateFareRequest struct {
	RideType        string  `json:"ride_type"`
	DistanceMeters  int     `json:"distance_meters"`
	ActiveRides     int     `json:"active_rides"`
	AvailableDrivers int    `json:"available_drivers"`
	IsPool          bool    `json:"is_pool"`
}

// EstimateFare handler - quick estimate without persistence
func (h *PricingHandler) EstimateFare(w http.ResponseWriter, r *http.Request) {
	var req EstimateFareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	estimate := h.pricingEngine.EstimateFare(
		req.RideType,
		req.DistanceMeters,
		req.ActiveRides,
		req.AvailableDrivers,
		req.IsPool,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estimate)
}

// GetSurgeMultiplier handler - returns current surge for location
type SurgeRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ActiveRides int    `json:"active_rides"`
	AvailableDrivers int `json:"available_drivers"`
}

func (h *PricingHandler) GetSurgeMultiplier(w http.ResponseWriter, r *http.Request) {
	var req SurgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	surgeMultiplier := h.pricingEngine.GetSurgeMultiplierAtLocation(
		req.Latitude, req.Longitude,
		req.ActiveRides, req.AvailableDrivers,
	)

	// Save surge history
	surge := &entities.SurgeHistory{
		ID:               uuid.New().String(),
		Timestamp:        time.Now(),
		City:             "Addis Ababa", // In production, use geolocation service
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		SurgeMultiplier:  surgeMultiplier,
		ActiveRides:      req.ActiveRides,
		AvailableDrivers: req.AvailableDrivers,
		Reason:           "SUPPLY_DEMAND",
		CreatedAt:        time.Now(),
	}
	_ = h.pricingRepo.SaveSurgeMultiplier(r.Context(), surge)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"surge_multiplier": surgeMultiplier,
		"timestamp":        time.Now(),
		"active_rides":     req.ActiveRides,
		"available_drivers": req.AvailableDrivers,
	})
}

// ApplyDiscountRequest DTO
type ApplyDiscountRequest struct {
	DiscountCode string  `json:"discount_code"`
	FareAmount   float64 `json:"fare_amount"`
}

// ApplyDiscount handler - validates and applies discount code
func (h *PricingHandler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	var req ApplyDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	discount, err := h.pricingRepo.ValidateDiscountCode(r.Context(), req.DiscountCode)
	if err != nil {
		http.Error(w, "Invalid or expired discount code", http.StatusBadRequest)
		return
	}

	// Check minimum fare
	minimumFare := 0.0
	if discount.MinimumFareAmount != nil {
		minimumFare = *discount.MinimumFareAmount
	}
	if req.FareAmount < minimumFare {
		http.Error(w, "Fare does not meet minimum for this discount", http.StatusBadRequest)
		return
	}

	// Calculate discount amount
	var discountAmount float64
	if discount.DiscountType == "FIXED" {
		discountAmount = discount.DiscountValue
	} else if discount.DiscountType == "PERCENTAGE" {
		discountAmount = req.FareAmount * (discount.DiscountValue / 100.0)
		if discount.MaxDiscount != nil && discountAmount > *discount.MaxDiscount {
			discountAmount = *discount.MaxDiscount
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"discount_code":   req.DiscountCode,
		"original_fare":   req.FareAmount,
		"discount_amount": discountAmount,
		"final_fare":      req.FareAmount - discountAmount,
		"discount_type":   discount.DiscountType,
	})
}

// GetPricingStats handler - returns pricing statistics
func (h *PricingHandler) GetPricingStats(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Addis Ababa"
	}

	avgFares, err := h.pricingRepo.GetAverageFareByRideType(r.Context(), city, 7)
	if err != nil {
		http.Error(w, "Failed to fetch statistics", http.StatusInternalServerError)
		return
	}

	surgeHistory, err := h.pricingRepo.GetSurgeHistory(r.Context(), city, 24)
	if err != nil {
		surgeHistory = []entities.SurgeHistory{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"city":            city,
		"average_fares":   avgFares,
		"surge_history":   surgeHistory,
		"period_days":     7,
		"timestamp":       time.Now(),
	})
}

// Health handler
func (h *PricingHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "pricing-service",
		"timestamp": time.Now(),
	})
}

// RegisterRoutes registers all pricing handlers
func (h *PricingHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/health", h.Health).Methods("GET")
	router.HandleFunc("/v1/pricing/calculate", h.CalculateFare).Methods("POST")
	router.HandleFunc("/v1/pricing/estimate", h.EstimateFare).Methods("POST")
	router.HandleFunc("/v1/pricing/surge", h.GetSurgeMultiplier).Methods("POST")
	router.HandleFunc("/v1/pricing/apply-discount", h.ApplyDiscount).Methods("POST")
	router.HandleFunc("/v1/pricing/statistics", h.GetPricingStats).Methods("GET")
}
