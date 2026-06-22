// services/user-service/internal/transport/http_handler.go
// User Service HTTP Handlers

package transport

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain"
)

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	registerUserHandler         *application.RegisterUserHandler
	updateProfileHandler        *application.UpdateProfileHandler
	activateUserHandler         *application.ActivateUserHandler
	verifyDriverHandler         *application.VerifyDriverHandler
	createDriverProfileHandler  *application.CreateDriverProfileHandler
	getUserHandler              *application.GetUserHandler
	getDriverProfileHandler     *application.GetDriverProfileHandler
	logger                      *zap.Logger
}

func NewHTTPHandler(
	registerUserHandler *application.RegisterUserHandler,
	updateProfileHandler *application.UpdateProfileHandler,
	activateUserHandler *application.ActivateUserHandler,
	verifyDriverHandler *application.VerifyDriverHandler,
	createDriverProfileHandler *application.CreateDriverProfileHandler,
	getUserHandler *application.GetUserHandler,
	getDriverProfileHandler *application.GetDriverProfileHandler,
	logger *zap.Logger,
) *HTTPHandler {
	return &HTTPHandler{
		registerUserHandler:        registerUserHandler,
		updateProfileHandler:       updateProfileHandler,
		activateUserHandler:        activateUserHandler,
		verifyDriverHandler:        verifyDriverHandler,
		createDriverProfileHandler: createDriverProfileHandler,
		getUserHandler:             getUserHandler,
		getDriverProfileHandler:    getDriverProfileHandler,
		logger:                     logger,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// User endpoints
	router.HandleFunc("/api/user/register", h.RegisterUser).Methods("POST")
	router.HandleFunc("/api/user/{userID}", h.GetUser).Methods("GET")
	router.HandleFunc("/api/user/profile", h.UpdateProfile).Methods("PUT")
	router.HandleFunc("/api/user/activate", h.ActivateUser).Methods("POST")

	// Driver endpoints
	router.HandleFunc("/api/driver/profile", h.CreateDriverProfile).Methods("POST")
	router.HandleFunc("/api/driver/{driverID}", h.GetDriverProfile).Methods("GET")
	router.HandleFunc("/api/driver/verify", h.VerifyDriver).Methods("POST")

	// Health endpoints
	router.HandleFunc("/health", h.Health).Methods("GET")
	router.HandleFunc("/ready", h.Ready).Methods("GET")
	router.HandleFunc("/startup", h.Startup).Methods("GET")
}

// ===== USER ENDPOINTS =====

// POST /api/user/register
func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		UserType  string `json:"user_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	cmd := application.RegisterUserCommand{
		Phone:     req.Phone,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserType:  req.UserType,
	}

	userID, err := h.registerUserHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.logger.Error("register user failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

// GET /api/user/{userID}
func (h *HTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	q := application.GetUserQuery{UserID: userID}
	user, err := h.getUserHandler.Handle(r.Context(), q)
	if err != nil {
		h.logger.Error("get user failed", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// PUT /api/user/profile
func (h *HTTPHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd := application.UpdateProfileCommand{
		UserID:    req.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if err := h.updateProfileHandler.Handle(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// POST /api/user/activate
func (h *HTTPHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd := application.ActivateUserCommand{UserID: req.UserID}

	if err := h.activateUserHandler.Handle(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// ===== DRIVER ENDPOINTS =====

// POST /api/driver/profile
func (h *HTTPHandler) CreateDriverProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID        string `json:"user_id"`
		LicenseNumber string `json:"license_number"`
		LicenseExpiry string `json:"license_expiry"`
		VehicleNumber string `json:"vehicle_number"`
		VehicleType   string `json:"vehicle_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expiry, _ := time.Parse("2006-01-02", req.LicenseExpiry)

	cmd := application.CreateDriverProfileCommand{
		UserID:        req.UserID,
		LicenseNumber: req.LicenseNumber,
		LicenseExpiry: expiry,
		VehicleNumber: req.VehicleNumber,
		VehicleType:   req.VehicleType,
	}

	profileID, err := h.createDriverProfileHandler.Handle(r.Context(), cmd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"profile_id": profileID})
}

// GET /api/driver/{driverID}
func (h *HTTPHandler) GetDriverProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["driverID"]

	q := application.GetDriverProfileQuery{ProfileID: driverID}
	profile, err := h.getDriverProfileHandler.Handle(r.Context(), q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// POST /api/driver/verify
func (h *HTTPHandler) VerifyDriver(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DriverID string `json:"driver_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd := application.VerifyDriverCommand{DriverID: req.DriverID}

	if err := h.verifyDriverHandler.Handle(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// ===== HEALTH ENDPOINTS =====

// GET /health - Liveness probe
func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

// GET /ready - Readiness probe
func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// GET /startup - Startup probe
func (h *HTTPHandler) Startup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}
