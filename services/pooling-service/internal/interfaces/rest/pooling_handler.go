package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/domain/services"
	"github.com/Abdex1/FamGo-platform/services/pooling-service/internal/infrastructure/postgres"
	"github.com/gorilla/mux"
)

// PoolingHandler handles pooling HTTP requests
type PoolingHandler struct {
	poolRepo       *postgres.PoolRepository
	poolingEngine  *services.PoolingEngine
}

// NewPoolingHandler creates new pooling handler
func NewPoolingHandler(poolRepo *postgres.PoolRepository, engine *services.PoolingEngine) *PoolingHandler {
	return &PoolingHandler{
		poolRepo:      poolRepo,
		poolingEngine: engine,
	}
}

// FindPoolMatches request DTO
type FindPoolMatchesDTO struct {
	RideID       string  `json:"ride_id"`
	DriverID     string  `json:"driver_id"`
	PickupLat    float64 `json:"pickup_lat"`
	PickupLng    float64 `json:"pickup_lng"`
	DropoffLat   float64 `json:"dropoff_lat"`
	DropoffLng   float64 `json:"dropoff_lng"`
	PickupAddr   string  `json:"pickup_address"`
	DropoffAddr  string  `json:"dropoff_address"`
	EstDistance  int     `json:"estimated_distance_meters"`
	EstDuration  int     `json:"estimated_duration_seconds"`
	EstFare      float64 `json:"estimated_fare"`
	FemaleOnly   bool    `json:"female_only"`
	MaxDetour    int     `json:"max_detour_minutes"`
	MaxWait      int     `json:"max_wait_minutes"`
	MinOverlap   float64 `json:"min_route_overlap"`
}

// FindPoolMatches handler - finds compatible rides for pooling
func (h *PoolingHandler) FindPoolMatches(w http.ResponseWriter, r *http.Request) {
	var dto FindPoolMatchesDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create target ride request
	targetRide := &entities.PoolRequest{
		RideID:             dto.RideID,
		DriverID:           dto.DriverID,
		PickupLat:          dto.PickupLat,
		PickupLng:          dto.PickupLng,
		DropoffLat:         dto.DropoffLat,
		DropoffLng:         dto.DropoffLng,
		PickupAddress:      dto.PickupAddr,
		DropoffAddress:     dto.DropoffAddr,
		EstimatedDistance:  dto.EstDistance,
		EstimatedDuration:  dto.EstDuration,
		EstimatedFare:      dto.EstFare,
		FemaleOnly:         dto.FemaleOnly,
		MaxDetourMinutes:   dto.MaxDetour,
		MaxWaitMinutes:     dto.MaxWait,
		MinRouteOverlap:    dto.MinOverlap,
		CreatedAt:          time.Now(),
	}

	// Get available pool requests
	availableRides, err := h.poolRepo.GetActivePoolRequests(r.Context(), 100)
	if err != nil {
		http.Error(w, "Failed to fetch available rides", http.StatusInternalServerError)
		return
	}

	// Convert to pointers
	ridePointers := make([]*entities.PoolRequest, len(availableRides))
	for i := range availableRides {
		ridePointers[i] = &availableRides[i]
	}

	// Find compatible candidates
	candidates := h.poolingEngine.FindPoolCandidates(targetRide, ridePointers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ride_id":     dto.RideID,
		"total":       len(candidates),
		"candidates":  candidates,
		"timestamp":   time.Now(),
	})
}

// CreatePoolRequest handler - creates new pool
type CreatePoolDTO struct {
	DriverID    string   `json:"driver_id"`
	RideIDs     []string `json:"ride_ids"`
	MaxSize     int      `json:"max_size"`
}

func (h *PoolingHandler) CreatePool(w http.ResponseWriter, r *http.Request) {
	var dto CreatePoolDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create pool
	pool := entities.NewPoolGroup(dto.DriverID, dto.MaxSize)

	// Add rides
	for _, rideID := range dto.RideIDs {
		pool.AddRide(rideID)
	}

	// Save pool
	if err := h.poolRepo.CreatePool(r.Context(), pool); err != nil {
		http.Error(w, "Failed to create pool", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pool_id":    pool.ID,
		"driver_id":  pool.DriverID,
		"rides":      pool.RideIDs,
		"size":       pool.CurrentSize,
		"status":     pool.Status,
		"created_at": pool.CreatedAt,
	})
}

// ActivatePool handler
func (h *PoolingHandler) ActivatePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	poolID := vars["poolID"]

	pool, err := h.poolRepo.GetPool(r.Context(), poolID)
	if err != nil {
		http.Error(w, "Pool not found", http.StatusNotFound)
		return
	}

	if !pool.CanBeActivated() {
		http.Error(w, "Pool cannot be activated (insufficient rides)", http.StatusBadRequest)
		return
	}

	pool.Activate()
	if err := h.poolRepo.UpdatePoolStatus(r.Context(), poolID, pool.Status); err != nil {
		http.Error(w, "Failed to activate pool", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pool_id": poolID,
		"status":  "ACTIVE",
		"rides":   pool.RideIDs,
	})
}

// CompletePool handler
func (h *PoolingHandler) CompletePool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	poolID := vars["poolID"]

	pool, err := h.poolRepo.GetPool(r.Context(), poolID)
	if err != nil {
		http.Error(w, "Pool not found", http.StatusNotFound)
		return
	}

	pool.Complete()
	if err := h.poolRepo.UpdatePoolStatus(r.Context(), poolID, pool.Status); err != nil {
		http.Error(w, "Failed to complete pool", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pool_id":      poolID,
		"status":       "COMPLETED",
		"completed_at": time.Now(),
	})
}

// GetPoolStatistics handler
func (h *PoolingHandler) GetPoolStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.poolRepo.GetPoolStatistics(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Health handler
func (h *PoolingHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "pooling-service",
		"timestamp": time.Now(),
	})
}

// RegisterRoutes registers all pooling handlers
func (h *PoolingHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/health", h.Health).Methods("GET")
	router.HandleFunc("/v1/pooling/find-matches", h.FindPoolMatches).Methods("POST")
	router.HandleFunc("/v1/pooling/pools", h.CreatePool).Methods("POST")
	router.HandleFunc("/v1/pooling/pools/{poolID}/activate", h.ActivatePool).Methods("POST")
	router.HandleFunc("/v1/pooling/pools/{poolID}/complete", h.CompletePool).Methods("POST")
	router.HandleFunc("/v1/pooling/statistics", h.GetPoolStatistics).Methods("GET")
}
