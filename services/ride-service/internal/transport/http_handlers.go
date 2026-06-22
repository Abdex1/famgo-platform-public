// services/ride-service/internal/transport/http_handlers.go
// HTTP Request Handlers

package transport

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// HTTPServer handles HTTP requests
type HTTPServer struct {
	createRideHandler     *application.CreateRideHandler
	assignDriverHandler   *application.AssignDriverHandler
	startRideHandler      *application.StartRideHandler
	completeRideHandler   *application.CompleteRideHandler
	cancelRideHandler     *application.CancelRideHandler
	getRideHandler        *application.GetRideHandler
	getPassengerRidesHandler *application.GetPassengerRidesHandler
	getDriverRidesHandler *application.GetDriverRidesHandler
	logger                *zap.Logger
}

func NewHTTPServer(
	createRideHandler *application.CreateRideHandler,
	assignDriverHandler *application.AssignDriverHandler,
	startRideHandler *application.StartRideHandler,
	completeRideHandler *application.CompleteRideHandler,
	cancelRideHandler *application.CancelRideHandler,
	getRideHandler *application.GetRideHandler,
	getPassengerRidesHandler *application.GetPassengerRidesHandler,
	getDriverRidesHandler *application.GetDriverRidesHandler,
	logger *zap.Logger,
) *HTTPServer {
	return &HTTPServer{
		createRideHandler:     createRideHandler,
		assignDriverHandler:   assignDriverHandler,
		startRideHandler:      startRideHandler,
		completeRideHandler:   completeRideHandler,
		cancelRideHandler:     cancelRideHandler,
		getRideHandler:        getRideHandler,
		getPassengerRidesHandler: getPassengerRidesHandler,
		getDriverRidesHandler: getDriverRidesHandler,
		logger:                logger,
	}
}

// ===== REQUEST/RESPONSE TYPES =====

type CreateRideRequest struct {
	PassengerID string  `json:"passenger_id" validate:"required"`
	PickupLat   float64 `json:"pickup_lat" validate:"required,min=-90,max=90"`
	PickupLon   float64 `json:"pickup_lon" validate:"required,min=-180,max=180"`
	DropoffLat  float64 `json:"dropoff_lat" validate:"required,min=-90,max=90"`
	DropoffLon  float64 `json:"dropoff_lon" validate:"required,min=-180,max=180"`
}

type CreateRideResponse struct {
	RideID string `json:"ride_id"`
	Status string `json:"status"`
}

type AssignDriverRequest struct {
	DriverID string `json:"driver_id" validate:"required"`
}

type RideResponse struct {
	ID                 string  `json:"id"`
	PassengerID        string  `json:"passenger_id"`
	DriverID           string  `json:"driver_id,omitempty"`
	PickupLat          float64 `json:"pickup_lat"`
	PickupLon          float64 `json:"pickup_lon"`
	DropoffLat         float64 `json:"dropoff_lat"`
	DropoffLon         float64 `json:"dropoff_lon"`
	Status             string  `json:"status"`
	EstimatedFare      float32 `json:"estimated_fare"`
	ActualFare         float32 `json:"actual_fare,omitempty"`
	CancellationReason string  `json:"cancellation_reason,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// ===== HANDLERS =====

// CreateRide - POST /rides
func (s *HTTPServer) CreateRide(w http.ResponseWriter, r *http.Request) {
	var req CreateRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warn("failed to decode request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	cmd := application.CreateRideCommand{
		PassengerID: req.PassengerID,
		PickupLat:   req.PickupLat,
		PickupLon:   req.PickupLon,
		DropoffLat:  req.DropoffLat,
		DropoffLon:  req.DropoffLon,
	}

	rideID, err := s.createRideHandler.Handle(r.Context(), cmd)
	if err != nil {
		s.logger.Error("failed to create ride", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to create ride"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateRideResponse{
		RideID: rideID,
		Status: string(domain.RideStatusRequested),
	})
}

// GetRide - GET /rides/{rideID}
func (s *HTTPServer) GetRide(w http.ResponseWriter, r *http.Request) {
	rideID := chi.URLParam(r, "rideID")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ride_id required"})
		return
	}

	ride, err := s.getRideHandler.Handle(r.Context(), rideID)
	if err != nil {
		s.logger.Error("failed to get ride", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ride not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rideToResponse(ride))
}

// GetPassengerRides - GET /passengers/{passengerID}/rides
func (s *HTTPServer) GetPassengerRides(w http.ResponseWriter, r *http.Request) {
	passengerID := chi.URLParam(r, "passengerID")
	if passengerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "passenger_id required"})
		return
	}

	limit, offset := getPaginationParams(r)
	rides, err := s.getPassengerRidesHandler.Handle(r.Context(), passengerID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get passenger rides", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to get rides"})
		return
	}

	var responses []RideResponse
	for _, ride := range rides {
		responses = append(responses, *rideToResponse(&ride))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

// GetDriverRides - GET /drivers/{driverID}/rides
func (s *HTTPServer) GetDriverRides(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "driverID")
	if driverID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "driver_id required"})
		return
	}

	limit, offset := getPaginationParams(r)
	rides, err := s.getDriverRidesHandler.Handle(r.Context(), driverID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get driver rides", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to get rides"})
		return
	}

	var responses []RideResponse
	for _, ride := range rides {
		responses = append(responses, *rideToResponse(&ride))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

// AssignDriver - POST /rides/{rideID}/assign
func (s *HTTPServer) AssignDriver(w http.ResponseWriter, r *http.Request) {
	rideID := chi.URLParam(r, "rideID")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ride_id required"})
		return
	}

	var req AssignDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	cmd := application.AssignDriverCommand{
		RideID:   rideID,
		DriverID: req.DriverID,
	}

	if err := s.assignDriverHandler.Handle(r.Context(), cmd); err != nil {
		s.logger.Error("failed to assign driver", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to assign driver"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartRide - POST /rides/{rideID}/start
func (s *HTTPServer) StartRide(w http.ResponseWriter, r *http.Request) {
	rideID := chi.URLParam(r, "rideID")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ride_id required"})
		return
	}

	cmd := application.StartRideCommand{RideID: rideID}
	if err := s.startRideHandler.Handle(r.Context(), cmd); err != nil {
		s.logger.Error("failed to start ride", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to start ride"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CompleteRide - POST /rides/{rideID}/complete
func (s *HTTPServer) CompleteRide(w http.ResponseWriter, r *http.Request) {
	rideID := chi.URLParam(r, "rideID")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ride_id required"})
		return
	}

	var req struct {
		ActualFare float32 `json:"actual_fare" validate:"required,min=0"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	cmd := application.CompleteRideCommand{
		RideID:     rideID,
		ActualFare: req.ActualFare,
	}

	if err := s.completeRideHandler.Handle(r.Context(), cmd); err != nil {
		s.logger.Error("failed to complete ride", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to complete ride"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CancelRide - POST /rides/{rideID}/cancel
func (s *HTTPServer) CancelRide(w http.ResponseWriter, r *http.Request) {
	rideID := chi.URLParam(r, "rideID")
	if rideID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ride_id required"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	cmd := application.CancelRideCommand{
		RideID: rideID,
		Reason: req.Reason,
	}

	if err := s.cancelRideHandler.Handle(r.Context(), cmd); err != nil {
		s.logger.Error("failed to cancel ride", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to cancel ride"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ===== HELPERS =====

func rideToResponse(ride *domain.Ride) *RideResponse {
	return &RideResponse{
		ID:                 ride.ID,
		PassengerID:        ride.PassengerID,
		DriverID:           ride.DriverID,
		PickupLat:          ride.PickupLat,
		PickupLon:          ride.PickupLon,
		DropoffLat:         ride.DropoffLat,
		DropoffLon:         ride.DropoffLon,
		Status:             string(ride.Status),
		EstimatedFare:      ride.EstimatedFare,
		ActualFare:         ride.ActualFare,
		CancellationReason: ride.CancellationReason,
		CreatedAt:          ride.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          ride.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func getPaginationParams(r *http.Request) (limit, offset int) {
	limit = 10
	offset = 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if _, err := json.Unmarshal([]byte(l), &limit); err == nil && limit > 0 && limit <= 100 {
			// limit is valid
		} else {
			limit = 10
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if _, err := json.Unmarshal([]byte(o), &offset); err == nil && offset >= 0 {
			// offset is valid
		} else {
			offset = 0
		}
	}

	return
}
