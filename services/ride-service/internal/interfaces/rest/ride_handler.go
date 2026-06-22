package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/infrastructure/postgres"
	"github.com/gorilla/mux"
)

// RideHandler handles ride HTTP requests
type RideHandler struct {
	rideRepo *postgres.RideRepository
}

// NewRideHandler creates new ride handler
func NewRideHandler(rideRepo *postgres.RideRepository) *RideHandler {
	return &RideHandler{rideRepo: rideRepo}
}

// CreateRideRequest request DTO
type CreateRideRequestDTO struct {
	RiderID        string  `json:"rider_id"`
	PickupLat      float64 `json:"pickup_lat"`
	PickupLng      float64 `json:"pickup_lng"`
	PickupAddress  string  `json:"pickup_address"`
	DropoffLat     float64 `json:"dropoff_lat"`
	DropoffLng     float64 `json:"dropoff_lng"`
	DropoffAddress string  `json:"dropoff_address"`
	RideType       string  `json:"ride_type"`
}

// CreateRideRequest handler
func (h *RideHandler) CreateRideRequest(w http.ResponseWriter, r *http.Request) {
	var dto CreateRideRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req := entities.NewRideRequest(
		dto.RiderID,
		dto.PickupLat, dto.PickupLng, dto.PickupAddress,
		dto.DropoffLat, dto.DropoffLng, dto.DropoffAddress,
		entities.RideType(dto.RideType),
	)

	if err := h.rideRepo.CreateRideRequest(r.Context(), req); err != nil {
		http.Error(w, "Failed to create ride request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        req.ID,
		"status":    req.Status,
		"created_at": req.CreatedAt,
	})
}

// GetRideStatus handler
func (h *RideHandler) GetRideStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["rideID"]

	ride, err := h.rideRepo.GetRide(r.Context(), rideID)
	if err != nil {
		http.Error(w, "Ride not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":               ride.ID,
		"status":           ride.Status,
		"driver_id":        ride.DriverID,
		"pickup_lat":       ride.PickupLat,
		"pickup_lng":       ride.PickupLng,
		"dropoff_lat":      ride.DropoffLat,
		"dropoff_lng":      ride.DropoffLng,
		"estimated_fare":   ride.EstimatedFare,
		"actual_fare":      ride.ActualFare,
		"assigned_at":      ride.AssignedAt,
		"pickup_time":      ride.PickupTime,
		"dropoff_time":     ride.DropoffTime,
		"created_at":       ride.CreatedAt,
	})
}

// UpdateRideStatus request DTO
type UpdateRideStatusDTO struct {
	Status string `json:"status"`
}

// UpdateRideStatus handler
func (h *RideHandler) UpdateRideStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["rideID"]

	var dto UpdateRideStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ride, err := h.rideRepo.GetRide(r.Context(), rideID)
	if err != nil {
		http.Error(w, "Ride not found", http.StatusNotFound)
		return
	}

	newStatus := entities.RideStatus(dto.Status)
	if !ride.CanTransitionTo(newStatus) {
		http.Error(w, "Invalid status transition", http.StatusBadRequest)
		return
	}

	if err := h.rideRepo.UpdateRideStatus(r.Context(), rideID, newStatus); err != nil {
		http.Error(w, "Failed to update ride status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         rideID,
		"status":     newStatus,
		"updated_at": time.Now(),
	})
}

// GetRideHistory handler
func (h *RideHandler) GetRideHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	riderID := vars["riderID"]

	rides, err := h.rideRepo.GetRiderRideHistory(r.Context(), riderID, 50, 0)
	if err != nil {
		http.Error(w, "Failed to fetch ride history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total": len(rides),
		"rides": rides,
	})
}

// CancelRide handler
type CancelRideDTO struct {
	Reason   string `json:"reason"`
	CancelledBy string `json:"cancelled_by"`
}

func (h *RideHandler) CancelRide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["rideID"]

	var dto CancelRideDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.rideRepo.UpdateRideStatus(r.Context(), rideID, entities.StatusCancelled); err != nil {
		http.Error(w, "Failed to cancel ride", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     rideID,
		"status": "CANCELLED",
		"reason": dto.Reason,
	})
}

// Health handler
func (h *RideHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"service": "ride-service",
		"timestamp": time.Now(),
	})
}

// RegisterRoutes registers all ride handlers
func (h *RideHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/health", h.Health).Methods("GET")
	router.HandleFunc("/v1/rides/request", h.CreateRideRequest).Methods("POST")
	router.HandleFunc("/v1/rides/{rideID}/status", h.GetRideStatus).Methods("GET")
	router.HandleFunc("/v1/rides/{rideID}/status", h.UpdateRideStatus).Methods("PUT")
	router.HandleFunc("/v1/rides/{rideID}/cancel", h.CancelRide).Methods("POST")
	router.HandleFunc("/v1/riders/{riderID}/history", h.GetRideHistory).Methods("GET")
}
