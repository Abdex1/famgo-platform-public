package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/infrastructure/redis"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/interfaces/websocket"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func main() {
	servicePort := getEnv("SERVICE_PORT", "3012")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	locationStore := redis.NewLocationStore(redisClient)

	// Create WebSocket server
	wsServer := websocket.NewWebSocketServer(locationStore)
	go wsServer.Run()

	// Setup router
	router := mux.NewRouter()
	
	// WebSocket endpoint
	router.HandleFunc("/ws/location", wsServer.HandleConnection)

	// Health endpoint
	router.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"gps-service"}`))
	})

	// Start server
	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: router,
	}

	log.Printf("🚀 GPS Service starting on port %s\n", servicePort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
