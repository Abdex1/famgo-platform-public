package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Abdex1/FamGo-platform/packages/telemetry"
	"github.com/Abdex1/FamGo-consolidated/services/gps-service/internal/application"
)

// HTTPHandler handles HTTP requests for GPS service
type HTTPHandler struct {
	updateLocationHandler    *application.UpdateDriverLocationHandler
	getLocationHandler       *application.GetDriverLocationHandler
	getNearbyDriversHandler  *application.GetNearbyDriversHandler
	metrics                  telemetry.Metrics
	logger                   telemetry.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	updateLocationHandler *application.UpdateDriverLocationHandler,
	getLocationHandler *application.GetDriverLocationHandler,
	getNearbyDriversHandler *application.GetNearbyDriversHandler,
	metrics telemetry.Metrics,
	logger telemetry.Logger,
) *HTTPHandler {
	return &HTTPHandler{
		updateLocationHandler:   updateLocationHandler,
		getLocationHandler:      getLocationHandler,
		getNearbyDriversHandler: getNearbyDriversHandler,
		metrics:                 metrics,
		logger:                  logger,
	}
}

// UpdateLocation handles POST /api/gps/location
func (h *HTTPHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var cmd application.UpdateDriverLocationCommand
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		h.logger.Error("Failed to parse request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "invalid request body", http.StatusBadRequest)
		h.metrics.RecordError("UpdateLocation", err)
		return
	}

	// Handle command
	err = h.updateLocationHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to update location", map[string]interface{}{
			"driver_id": cmd.DriverID,
			"error":     err.Error(),
		})
		http.Error(w, fmt.Sprintf("failed to update location: %v", err), http.StatusInternalServerError)
		h.metrics.RecordError("UpdateLocation", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// GetLocation handles GET /api/gps/location/{driver_id}
func (h *HTTPHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	// Extract driver ID from query parameter
	driverID := r.URL.Query().Get("driver_id")
	if driverID == "" {
		http.Error(w, "driver_id query parameter required", http.StatusBadRequest)
		return
	}

	// Handle query
	q := application.GetDriverLocationQuery{DriverID: driverID}
	location, err := h.getLocationHandler.Handle(r.Context(), q)
	if err != nil {
		h.logger.Error("Failed to get location", map[string]interface{}{
			"driver_id": driverID,
			"error":     err.Error(),
		})
		http.Error(w, fmt.Sprintf("failed to get location: %v", err), http.StatusInternalServerError)
		h.metrics.RecordError("GetLocation", err)
		return
	}

	if location == nil {
		http.Error(w, "location not found", http.StatusNotFound)
		return
	}

	// Return location
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}

// GetNearbyDrivers handles GET /api/gps/nearby
func (h *HTTPHandler) GetNearbyDrivers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	var latitude, longitude, radiusM float64
	fmt.Sscanf(r.URL.Query().Get("latitude"), "%f", &latitude)
	fmt.Sscanf(r.URL.Query().Get("longitude"), "%f", &longitude)
	fmt.Sscanf(r.URL.Query().Get("radius_m"), "%f", &radiusM)

	if latitude == 0 || longitude == 0 || radiusM == 0 {
		http.Error(w, "latitude, longitude, and radius_m query parameters required", http.StatusBadRequest)
		return
	}

	// Handle query
	q := application.GetNearbyDriversQuery{
		Latitude:  latitude,
		Longitude: longitude,
		RadiusM:   radiusM,
	}
	drivers, err := h.getNearbyDriversHandler.Handle(r.Context(), q)
	if err != nil {
		h.logger.Error("Failed to get nearby drivers", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, fmt.Sprintf("failed to get nearby drivers: %v", err), http.StatusInternalServerError)
		h.metrics.RecordError("GetNearbyDrivers", err)
		return
	}

	// Return drivers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"drivers": drivers,
		"count":   len(drivers),
	})
}

// ========================================
// Health Checks

// Live handles GET /health (liveness probe)
func (h *HTTPHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

// Ready handles GET /ready (readiness probe)
func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Check dependencies (database, redis, etc.)
	// For now, just return OK
	// In production, check actual dependencies

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// Startup handles GET /startup (startup probe)
func (h *HTTPHandler) Startup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// Metrics handles GET /metrics (Prometheus metrics)
func (h *HTTPHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	// Return Prometheus metrics
	// This is handled by the telemetry package
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# GPS Service Metrics\n"))
	w.Write([]byte("gps_service_requests_total{operation=\"UpdateLocation\"} 0\n"))
}
