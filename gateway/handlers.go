// gateway/handlers.go
package gateway

import (
	"encoding/json"
	"net/http"
)

// APIResponse standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

// HealthCheckHandler returns API health status
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{
		Success: true,
		Data: map[string]string{
			"status":  "ok",
			"service": "api-gateway",
		},
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

// RideHandler handles ride-related endpoints
func RideHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		response := APIResponse{
			Success: true,
			Data: map[string]string{
				"ride_id": "RIDE_001",
				"status":  "requested",
			},
			Code: http.StatusCreated,
		}
		json.NewEncoder(w).Encode(response)
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		response := APIResponse{
			Success: true,
			Data:    []interface{}{},
			Code:    http.StatusOK,
		}
		json.NewEncoder(w).Encode(response)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// PaymentHandler handles payment endpoints
func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := APIResponse{
		Success: true,
		Data: map[string]string{
			"payment_id": "PAY_001",
			"status":     "processing",
		},
		Code: http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}
