package routes

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"famgo/driver-service/internal/interfaces/rest/handlers"
)

// RegisterRoutes registers all driver service routes
func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	driverHandler := handlers.NewDriverHandler(db)

	// Health check
	router.HandleFunc("/v1/health", driverHandler.Health).Methods("GET")

	// Driver profile endpoints
	router.HandleFunc("/v1/drivers/{id}", driverHandler.GetProfile).Methods("GET")
	router.HandleFunc("/v1/drivers/{id}", driverHandler.UpdateProfile).Methods("PUT")

	// Driver status endpoints
	router.HandleFunc("/v1/drivers/{id}/go-online", driverHandler.GoOnline).Methods("POST")
	router.HandleFunc("/v1/drivers/{id}/go-offline", driverHandler.GoOffline).Methods("POST")

	// Vehicle endpoints
	router.HandleFunc("/v1/drivers/{id}/vehicles", driverHandler.GetVehicles).Methods("GET")
	router.HandleFunc("/v1/drivers/{id}/vehicles", driverHandler.AddVehicle).Methods("POST")

	// Document endpoints
	router.HandleFunc("/v1/drivers/{id}/documents", driverHandler.GetDocuments).Methods("GET")

	// Earnings endpoints
	router.HandleFunc("/v1/drivers/{id}/earnings", driverHandler.GetEarnings).Methods("GET")

	// Location endpoint
	router.HandleFunc("/v1/drivers/{id}/location", driverHandler.UpdateLocation).Methods("POST")
}
