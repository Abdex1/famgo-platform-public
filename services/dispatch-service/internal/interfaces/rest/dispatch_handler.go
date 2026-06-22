package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/usecases"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/services"
	"github.com/gorilla/mux"
)

// DispatchHandler handles dispatch REST requests.
type DispatchHandler struct {
	useCases        *usecases.DispatchUseCases
	candidateEngine *services.DriverCandidateEngine
	searchRadiusKm  float64
	nearbyLimit     int
}

// NewDispatchHandler creates a REST handler wired to application use cases.
func NewDispatchHandler(
	useCases *usecases.DispatchUseCases,
	candidateEngine *services.DriverCandidateEngine,
	searchRadiusKm float64,
	nearbyLimit int,
) *DispatchHandler {
	if searchRadiusKm <= 0 {
		searchRadiusKm = 5.0
	}
	if nearbyLimit <= 0 {
		nearbyLimit = 50
	}
	return &DispatchHandler{
		useCases:        useCases,
		candidateEngine: candidateEngine,
		searchRadiusKm:  searchRadiusKm,
		nearbyLimit:     nearbyLimit,
	}
}

// MatchRideRequestDTO is the REST payload for manual match requests.
type MatchRideRequestDTO struct {
	RideID     string  `json:"ride_id"`
	RiderID    string  `json:"rider_id"`
	PickupLat  float64 `json:"pickup_lat"`
	PickupLng  float64 `json:"pickup_lng"`
	DropoffLat float64 `json:"dropoff_lat"`
	DropoffLng float64 `json:"dropoff_lng"`
	RideType   string  `json:"ride_type"`
}

// MatchRide starts dispatch matching through the application layer.
func (h *DispatchHandler) MatchRide(w http.ResponseWriter, r *http.Request) {
	var dto MatchRideRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.useCases.MatchRide(r.Context(), &usecases.MatchRideInput{
		RideID:     dto.RideID,
		RiderID:    dto.RiderID,
		PickupLat:  dto.PickupLat,
		PickupLng:  dto.PickupLng,
		DropoffLat: dto.DropoffLat,
		DropoffLng: dto.DropoffLng,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dispatch_request_id": result.DispatchRequestID,
		"ride_id":             dto.RideID,
		"assigned_driver_id":  result.MatchedDriverID,
		"status":              result.Status,
		"proposed_drivers":    result.ProposedDrivers,
		"matched_at":          result.MatchedAt,
	})
}

// GetNearbyDriversDTO is the REST payload for nearby driver lookup.
type GetNearbyDriversDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetNearbyDrivers returns online drivers near a pickup point.
func (h *DispatchHandler) GetNearbyDrivers(w http.ResponseWriter, r *http.Request) {
	var dto GetNearbyDriversDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	drivers, err := h.candidateEngine.FindWithinRadius(
		r.Context(),
		dto.Latitude,
		dto.Longitude,
		h.searchRadiusKm,
		h.nearbyLimit,
	)
	if err != nil {
		http.Error(w, "Failed to find drivers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":   len(drivers),
		"drivers": drivers,
	})
}

// Health returns service health metadata.
func (h *DispatchHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "dispatch-service",
		"timestamp": time.Now(),
	})
}

// RegisterRoutes registers dispatch REST routes.
func (h *DispatchHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/health", h.Health).Methods("GET")
	router.HandleFunc("/v1/dispatch/match", h.MatchRide).Methods("POST")
	router.HandleFunc("/v1/dispatch/nearby-drivers", h.GetNearbyDrivers).Methods("POST")
}
