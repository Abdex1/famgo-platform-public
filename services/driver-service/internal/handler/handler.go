package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"famgo/driver-service/internal/model"
	"famgo/driver-service/internal/service"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// Handler handles HTTP requests for driver service
type Handler struct {
	driverService *service.DriverService
	logger        logger.Logger
}

// NewHandler creates a new handler
func NewHandler(driverService *service.DriverService, logger logger.Logger) *Handler {
	return &Handler{
		driverService: driverService,
		logger:        logger,
	}
}

// RegisterRoutes registers all HTTP routes (WEEK 1 FOUNDATION)
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/v1/drivers", func(r chi.Router) {
		// Registration endpoints (2-step)
		r.Post("/register", h.Register)
		r.Post("/verify-register", h.VerifyRegister)

		// Profile endpoints
		r.Get("/{driverID}/profile", h.GetProfile)
		r.Put("/{driverID}/profile", h.UpdateProfile)

		// State management (foundation)
		r.Get("/{driverID}/state", h.GetCurrentState)
		r.Get("/{driverID}/state-history", h.GetStateHistory)
		r.Post("/{driverID}/state-transition", h.TransitionState)
	})
}

// Register - Step 1: Send OTP for driver registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Phone == "" || req.LicenseNumber == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email, phone, and license_number are required")
		return
	}

	// In WEEK 1: Just confirm registration started
	// Full OTP logic in auth-service integration (WEEK 3)
	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "OTP sent to email",
		"email":   req.Email,
	})
}

// VerifyRegister - Step 2: Verify OTP and create driver account
func (h *Handler) VerifyRegister(w http.ResponseWriter, r *http.Request) {
	var req model.VerifyRegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.OTP == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and OTP are required")
		return
	}

	// In WEEK 1: Foundation structure only
	// Full verification in WEEK 3
	h.respondSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "Driver registered successfully",
		"status":  "pending",
	})
}

// GetProfile retrieves driver profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "driverID")
	if driverID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DRIVER_ID", "Driver ID is required")
		return
	}

	driver, err := h.driverService.GetProfile(r.Context(), driverID)
	if err != nil {
		h.logger.Warn("get profile failed", map[string]interface{}{"driver_id": driverID, "error": err})
		h.respondError(w, http.StatusNotFound, "DRIVER_NOT_FOUND", "Driver not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, driver)
}

// UpdateProfile updates driver profile (foundation - license only)
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "driverID")
	if driverID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DRIVER_ID", "Driver ID is required")
		return
	}

	var req struct {
		LicenseNumber string `json:"license_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	driver, err := h.driverService.UpdateProfile(r.Context(), driverID, req.LicenseNumber)
	if err != nil {
		h.logger.Warn("update profile failed", map[string]interface{}{"driver_id": driverID, "error": err})
		h.respondError(w, http.StatusNotFound, "DRIVER_NOT_FOUND", "Driver not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, driver)
}

// GetCurrentState retrieves driver's current state (Pattern 4)
func (h *Handler) GetCurrentState(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "driverID")
	if driverID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DRIVER_ID", "Driver ID is required")
		return
	}

	driver, err := h.driverService.GetProfile(r.Context(), driverID)
	if err != nil {
		h.logger.Warn("get state failed", map[string]interface{}{"driver_id": driverID, "error": err})
		h.respondError(w, http.StatusNotFound, "DRIVER_NOT_FOUND", "Driver not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"driver_id": driver.ID,
		"status":    driver.Status,
	})
}

// GetStateHistory retrieves state transition history (Pattern 4)
func (h *Handler) GetStateHistory(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "driverID")
	if driverID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DRIVER_ID", "Driver ID is required")
		return
	}

	limit := 10 // Default limit
	queryLimit := r.URL.Query().Get("limit")
	if queryLimit != "" {
		fmt.Sscanf(queryLimit, "%d", &limit)
	}

	states, err := h.driverService.GetStateHistory(r.Context(), driverID, limit)
	if err != nil {
		h.logger.Warn("get state history failed", map[string]interface{}{"driver_id": driverID, "error": err})
		// Return empty array instead of error
		h.respondSuccess(w, http.StatusOK, []interface{}{})
		return
	}

	h.respondSuccess(w, http.StatusOK, states)
}

// TransitionState transitions driver to a new state (Pattern 4: State Machine)
func (h *Handler) TransitionState(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "driverID")
	if driverID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DRIVER_ID", "Driver ID is required")
		return
	}

	var req model.StateTransitionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.NewState == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "new_state is required")
		return
	}

	// Transition state
	state, err := h.driverService.TransitionState(r.Context(), driverID, req.NewState, req.Reason)
	if err != nil {
		h.logger.Warn("state transition failed", map[string]interface{}{"driver_id": driverID, "error": err})
		h.respondError(w, http.StatusBadRequest, "TRANSITION_FAILED", err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, state)
}

// Response helpers

// respondSuccess sends a success response
func (h *Handler) respondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := map[string]interface{}{
		"code": "OK",
		"data": data,
	}

	json.NewEncoder(w).Encode(resp)
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := map[string]interface{}{
		"code":    code,
		"message": message,
	}

	json.NewEncoder(w).Encode(resp)
}

