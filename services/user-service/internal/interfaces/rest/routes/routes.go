package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/interfaces/rest/handlers"
)

// RegisterRoutes registers all user service routes
func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)

	// Health check
	router.HandleFunc("/v1/health", userHandler.Health).Methods("GET")

	// User profile endpoints
	router.HandleFunc("/v1/users/{id}/profile", userHandler.GetProfile).Methods("GET")
	router.HandleFunc("/v1/users/{id}/profile", userHandler.UpdateProfile).Methods("PUT")

	// User preferences endpoints
	router.HandleFunc("/v1/users/{id}/preferences", userHandler.GetPreferences).Methods("GET")
	router.HandleFunc("/v1/users/{id}/preferences", userHandler.UpdatePreferences).Methods("PUT")

	// User history endpoints
	router.HandleFunc("/v1/users/{id}/ride-history", userHandler.GetRideHistory).Methods("GET")

	// Notification endpoints
	router.HandleFunc("/v1/users/{id}/notifications", userHandler.GetNotifications).Methods("GET")

	// Rating endpoints
	router.HandleFunc("/v1/users/{id}/ratings", userHandler.SubmitRating).Methods("POST")
}
