package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// ProfileRequest represents profile update request
type ProfileRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Country         string `json:"country"`
	PhoneNumber     string `json:"phone_number"`
}

// PreferenceRequest represents preference update request
type PreferenceRequest struct {
	Language              string `json:"language"`
	Currency              string `json:"currency"`
	NotificationsEnabled  bool   `json:"notifications_enabled"`
	EmailNotifications    bool   `json:"email_notifications"`
	SMSNotifications      bool   `json:"sms_notifications"`
	PushNotifications     bool   `json:"push_notifications"`
	FemaleDriverPreference bool   `json:"female_driver_preference"`
	PoolingPreference     string `json:"pooling_preference"`
}

// RatingRequest represents rating submission request
type RatingRequest struct {
	RideID    string                 `json:"ride_id"`
	Rating    int                    `json:"rating"`
	Review    string                 `json:"review"`
	Categories map[string]int        `json:"categories"`
}

// UserHandler handles user-related endpoints
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetProfile handles GET /v1/users/:id/profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// TODO: Fetch user profile from database
	// TODO: Return profile JSON

	writeSuccessResponse(w, http.StatusOK, "User profile retrieved", map[string]interface{}{
		"id":        userID,
		"first_name": "John",
		"last_name": "Doe",
		"email":    "john@example.com",
		"phone":    "+251911234567",
		"bio":      "I love riding",
		"rating":   4.8,
		"total_rides": 45,
	})
}

// UpdateProfile handles PUT /v1/users/:id/profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req ProfileRequest
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

	// TODO: Update profile in database
	// TODO: Emit user.profile.updated event

	writeSuccessResponse(w, http.StatusOK, "Profile updated successfully", map[string]string{
		"user_id": userID,
	})
}

// GetPreferences handles GET /v1/users/:id/preferences
func (h *UserHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// TODO: Fetch user preferences from database

	writeSuccessResponse(w, http.StatusOK, "Preferences retrieved", map[string]interface{}{
		"language":               "en",
		"currency":              "ETB",
		"notifications_enabled": true,
		"female_driver_preference": false,
		"pooling_preference":    "standard",
	})
}

// UpdatePreferences handles PUT /v1/users/:id/preferences
func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req PreferenceRequest
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

	// TODO: Update preferences in database

	writeSuccessResponse(w, http.StatusOK, "Preferences updated successfully", map[string]string{
		"user_id": userID,
	})
}

// GetRideHistory handles GET /v1/users/:id/ride-history
func (h *UserHandler) GetRideHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// TODO: Fetch ride history from database

	writeSuccessResponse(w, http.StatusOK, "Ride history retrieved", map[string]interface{}{
		"user_id": userID,
		"total_rides": 45,
		"rides": []map[string]interface{}{
			{
				"ride_id": uuid.New().String(),
				"date": "2024-01-15",
				"from": "Addis Mall",
				"to": "Bole Airport",
				"fare": 150.00,
			},
		},
	})
}

// GetNotifications handles GET /v1/users/:id/notifications
func (h *UserHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// TODO: Fetch notifications from database

	writeSuccessResponse(w, http.StatusOK, "Notifications retrieved", map[string]interface{}{
		"user_id": userID,
		"unread_count": 3,
		"notifications": []map[string]interface{}{},
	})
}

// SubmitRating handles POST /v1/users/:id/ratings
func (h *UserHandler) SubmitRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req RatingRequest
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

	if req.Rating < 1 || req.Rating > 5 {
		writeErrorResponse(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	// TODO: Save rating to database
	// TODO: Emit feedback.submitted event
	// TODO: Update user average rating

	writeSuccessResponse(w, http.StatusCreated, "Rating submitted successfully", map[string]string{
		"user_id": userID,
		"ride_id": req.RideID,
	})
}

// Health handles health check
func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeSuccessResponse(w, http.StatusOK, "User service is healthy", map[string]string{
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
