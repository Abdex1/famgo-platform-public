package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/infrastructure/redis"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// LocationUpdate message from client
type LocationUpdateMsg struct {
	DriverID  string  `json:"driver_id"`
	RideID    string  `json:"ride_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Speed     float64 `json:"speed"`
	Accuracy  float64 `json:"accuracy"`
}

// WebSocketServer manages WebSocket connections
type WebSocketServer struct {
	locationStore *redis.LocationStore
	clients       map[*Client]bool
	mu            sync.RWMutex
	broadcast     chan interface{}
	register      chan *Client
	unregister    chan *Client
	upgrader      websocket.Upgrader
}

// Client represents connected WebSocket client
type Client struct {
	conn     *websocket.Conn
	send     chan interface{}
	driverID string
	rideID   string
}

// NewWebSocketServer creates new WebSocket server
func NewWebSocketServer(locationStore *redis.LocationStore) *WebSocketServer {
	return &WebSocketServer{
		locationStore: locationStore,
		clients:       make(map[*Client]bool),
		broadcast:     make(chan interface{}, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// Run starts WebSocket server event loop
func (s *WebSocketServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			log.Printf("Client connected: %s", client.driverID)

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()
			log.Printf("Client disconnected: %s", client.driverID)

		case message := <-s.broadcast:
			s.mu.RLock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					// Client's send channel is full, skip
				}
			}
			s.mu.RUnlock()
		}
	}
}

// HandleConnection upgrades HTTP to WebSocket
func (s *WebSocketServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	driverID := r.URL.Query().Get("driver_id")
	rideID := r.URL.Query().Get("ride_id")

	if driverID == "" || rideID == "" {
		http.Error(w, "Missing driver_id or ride_id", http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan interface{}, 256),
		driverID: driverID,
		rideID:   rideID,
	}

	s.register <- client

	// Handle client messages
	go s.handleClientMessages(client)
	go s.handleClientSend(client)
}

// handleClientMessages reads messages from client
func (s *WebSocketServer) handleClientMessages(client *Client) {
	defer func() {
		s.unregister <- client
		client.conn.Close()
	}()

	for {
		var msg LocationUpdateMsg
		if err := client.conn.ReadJSON(&msg); err != nil {
			break
		}

		// Save to Redis
		driverLoc := &domain.DriverLocation{
			DriverID:  client.driverID,
			Latitude:  msg.Latitude,
			Longitude: msg.Longitude,
			Heading:   msg.Heading,
			Speed:     msg.Speed,
			UpdatedAt: time.Now(),
		}

		if err := s.locationStore.SaveDriverLocation(context.Background(), driverLoc); err != nil {
			log.Printf("Error saving location: %v", err)
		}

		// Save to history
		locUpdate := &domain.LocationUpdate{
			ID:        uuid.New().String(),
			DriverID:  client.driverID,
			RideID:    client.rideID,
			Latitude:  msg.Latitude,
			Longitude: msg.Longitude,
			Heading:   msg.Heading,
			Speed:     msg.Speed,
			Accuracy:  msg.Accuracy,
			Timestamp: time.Now(),
		}

		if err := s.locationStore.SaveLocationHistory(context.Background(), locUpdate); err != nil {
			log.Printf("Error saving location history: %v", err)
		}

		// Broadcast to other clients
		s.broadcast <- map[string]interface{}{
			"type":      "location_update",
			"driver_id": client.driverID,
			"ride_id":   client.rideID,
			"latitude":  msg.Latitude,
			"longitude": msg.Longitude,
			"heading":   msg.Heading,
			"speed":     msg.Speed,
			"timestamp": time.Now().Unix(),
		}
	}
}

// handleClientSend writes messages to client
func (s *WebSocketServer) handleClientSend(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Import context package
import "context"
