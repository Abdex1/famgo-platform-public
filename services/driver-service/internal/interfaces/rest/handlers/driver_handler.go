package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// DriverProfileRequest represents driver profile update
type DriverProfileRequest struct {
	LicenseNumber  string `json:"license_number"`
	LicenseExpiry  string `json:"license_expiry"`
	BankAccount    map[string]string `json:"bank_account"`
	EmergencyContact map[string]string `json:"emergency_contact"`
}

// VehicleRequest represents vehicle creation/update
type VehicleRequest struct {
	LicensePlate       string `json:"license_plate"`
	VehicleType        string `json:"vehicle_type"` // economy, comfort, premium, xl
	Make               string `json:"make"`
	Model              string `json:"model"`
	Year               int    `json:"year"`
	Color              string `json:"color"`
	RegistrationExpiry string `json:"registration_expiry"`
	InsuranceExpiry    string `json:"insurance_expiry"`
}

// LocationRequest represents location update
type LocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   int     `json:"heading"`
	Speed     int     `json:"speed"`
	Accuracy  int     `json:"accuracy"`
}

// DriverHandler handles driver-related endpoints
type DriverHandler struct {
	db *gorm.DB
}

// NewDriverHandler creates a new driver handler
func NewDriverHandler(db *gorm.DB) *DriverHandler {
	return &DriverHandler{db: db}
}

// GetProfile handles GET /v1/drivers/:id
func (h *DriverHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	// TODO: Fetch driver profile from database

	writeSuccessResponse(w, http.StatusOK, "Driver profile retrieved", map[string]interface{}{
		"id":              driverID,
		"user_id":         uuid.New().String(),
		"status":          "online",
		"rating":          4.9,
		"total_rides":     150,
		"total_earnings":  25000.00,
		"acceptance_rate": 98.5,
		"is_verified":     true,
	})
}

// UpdateProfile handles PUT /v1/drivers/:id
func (h *DriverHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	var req DriverProfileRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: Update driver profile
	// TODO: Emit driver.profile.updated event

	writeSuccessResponse(w, http.StatusOK, "Driver profile updated", map[string]string{
		"driver_id": driverID,
	})
}

// GoOnline handles POST /v1/drivers/:id/go-online
func (h *DriverHandler) GoOnline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	// TODO: Update driver status to online
	// TODO: Emit driver.online event
	// TODO: Save location from request

	writeSuccessResponse(w, http.StatusOK, "Driver is now online", map[string]string{
		"driver_id": driverID,
		"status":    "online",
	})
}

// GoOffline handles POST /v1/drivers/:id/go-offline
func (h *DriverHandler) GoOffline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	// TODO: Update driver status to offline
	// TODO: Emit driver.offline event

	writeSuccessResponse(w, http.StatusOK, "Driver is now offline", map[string]string{
		"driver_id": driverID,
		"status":    "offline",
	})
}

// GetVehicles handles GET /v1/drivers/:id/vehicles
func (h *DriverHandler) GetVehicles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	// TODO: Fetch driver vehicles

	writeSuccessResponse(w, http.StatusOK, "Vehicles retrieved", map[string]interface{}{
		"driver_id": driverID,
		"vehicles": []map[string]interface{}{
			{
				"id":              uuid.New().String(),
				"license_plate":   "AA1234",
				"vehicle_type":    "economy",
				"make":            "Toyota",
				"model":           "Corolla",
				"color":           "Silver",
				"registration_expiry": "2025-12-31",
			},
		},
	})
}

// AddVehicle handles POST /v1/drivers/:id/vehicles
func (h *DriverHandler) AddVehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	var req VehicleRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: Create vehicle in database
	// TODO: Emit driver.vehicle.added event

	writeSuccessResponse(w, http.StatusCreated, "Vehicle added successfully", map[string]string{
		"driver_id":     driverID,
		"vehicle_id":    uuid.New().String(),
		"license_plate": req.LicensePlate,
	})
}

// GetDocuments handles GET /v1/drivers/:id/documents
func (h *DriverHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	// TODO: Fetch driver documents

	writeSuccessResponse(w, http.StatusOK, "Documents retrieved", map[string]interface{}{
		"driver_id": driverID,
		"documents": []map[string]interface{}{
			{
				"type":    "license",
				"status":  "approved",
				"expiry":  "2026-12-31",
			},
			{
				"type":    "insurance",
				"status":  "approved",
				"expiry":  "2025-06-30",
			},
		},
	})
}

// GetEarnings handles GET /v1/drivers/:id/earnings
func (h *DriverHandler) GetEarnings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	// TODO: Fetch driver earnings for period

	writeSuccessResponse(w, http.StatusOK, "Earnings retrieved", map[string]interface{}{
		"driver_id":        driverID,
		"total_earnings":   5000.00,
		"rides_completed":  45,
		"rides_cancelled":  2,
		"average_fare":     111.11,
		"period":           "weekly",
		"start_date":       "2024-01-01",
		"end_date":         "2024-01-07",
	})
}

// UpdateLocation handles POST /v1/drivers/:id/location
func (h *DriverHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	driverID := vars["id"]

	var req LocationRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: Update driver location in Redis GEO
	// TODO: Emit driver.location.updated event

	writeSuccessResponse(w, http.StatusOK, "Location updated", map[string]interface{}{
		"driver_id": driverID,
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
	})
}

// Health handles health check
func (h *DriverHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeSuccessResponse(w, http.StatusOK, "Driver service is healthy", map[string]string{
		"status": "ok",
	})
}

// Helper functions

func writeSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}

	json.NewEncoder(w).Encode(response)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error":   http.StatusText(statusCode),
		"message": message,
		"status":  statusCode,
	}

	json.NewEncoder(w).Encode(response)
}
