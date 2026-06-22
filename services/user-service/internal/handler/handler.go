package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"famgo/user-service/internal/model"
	"famgo/user-service/internal/service"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// Handler handles HTTP requests for user service
type Handler struct {
	userService *service.UserService
	logger      logger.Logger
}

// NewHandler creates a new handler
func NewHandler(userService *service.UserService, logger logger.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/v1/users", func(r chi.Router) {
		// Protected endpoints (auth middleware required)
		r.Get("/{userID}/profile", h.GetProfile)
		r.Put("/{userID}/profile", h.UpdateProfile)
		r.Get("/{userID}/preferences", h.GetPreferences)
		r.Put("/{userID}/preferences", h.UpdatePreferences)
		r.Get("/{userID}/addresses", h.GetAddresses)
		r.Post("/{userID}/addresses", h.CreateAddress)
		r.Delete("/addresses/{addressID}", h.DeleteAddress)
	})
}

// GetProfile retrieves user profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	profile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		h.logger.Warn("get profile failed", map[string]interface{}{"user_id": userID, "error": err})
		h.respondError(w, http.StatusNotFound, "PROFILE_NOT_FOUND", "User profile not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, profile)
}

// UpdateProfile updates user profile
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	var req model.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	profile, err := h.userService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		h.logger.Warn("update profile failed", map[string]interface{}{"user_id": userID, "error": err})
		h.respondError(w, http.StatusNotFound, "PROFILE_NOT_FOUND", "User profile not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, profile)
}

// GetPreferences retrieves user preferences
func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	prefs, err := h.userService.GetPreferences(r.Context(), userID)
	if err != nil {
		h.logger.Warn("get preferences failed", map[string]interface{}{"user_id": userID, "error": err})
		h.respondError(w, http.StatusNotFound, "PREFERENCES_NOT_FOUND", "User preferences not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, prefs)
}

// UpdatePreferences updates user preferences
func (h *Handler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	var req model.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	prefs, err := h.userService.UpdatePreferences(r.Context(), userID, &req)
	if err != nil {
		h.logger.Warn("update preferences failed", map[string]interface{}{"user_id": userID, "error": err})
		h.respondError(w, http.StatusNotFound, "PREFERENCES_NOT_FOUND", "User preferences not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, prefs)
}

// GetAddresses retrieves all saved addresses
func (h *Handler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	addresses, err := h.userService.GetAddresses(r.Context(), userID)
	if err != nil {
		h.logger.Warn("get addresses failed", map[string]interface{}{"user_id": userID, "error": err})
		h.respondError(w, http.StatusInternalServerError, "QUERY_FAILED", "Failed to retrieve addresses")
		return
	}

	h.respondSuccess(w, http.StatusOK, addresses)
}

// CreateAddress creates a new saved address
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_USER_ID", "User ID is required")
		return
	}

	var req model.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate required fields
	if req.Type == "" || req.AddressLine1 == "" || req.City == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Type, address_line_1, and city are required")
		return
	}

	address, err := h.userService.AddAddress(r.Context(), userID, &req)
	if err != nil {
		h.logger.Warn("create address failed", map[string]interface{}{"user_id": userID, "error": err})
		h.respondError(w, http.StatusBadRequest, "CREATE_ADDRESS_FAILED", err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, address)
}

// DeleteAddress deletes a saved address
func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressID := chi.URLParam(r, "addressID")
	if addressID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_ADDRESS_ID", "Address ID is required")
		return
	}

	if err := h.userService.DeleteAddress(r.Context(), addressID); err != nil {
		h.logger.Warn("delete address failed", map[string]interface{}{"address_id": addressID, "error": err})
		h.respondError(w, http.StatusNotFound, "ADDRESS_NOT_FOUND", "Address not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "Address deleted successfully",
	})
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
