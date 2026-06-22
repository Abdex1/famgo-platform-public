// services/ride-service/internal/transport/websocket.go
// WebSocket Real-Time Handlers

package transport

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/application"
	"github.com/Abdex1/FamGo-platform/services/ride-service/internal/domain"
)

// WebSocketHandler handles real-time ride updates
type WebSocketHandler struct {
	getRideHandler *application.GetRideHandler
	logger         *zap.Logger
	upgrader       websocket.Upgrader
	clients        map[string][]*WebSocketClient // rideID -> clients
	clientsMutex   sync.RWMutex
}

// WebSocketClient represents a connected WebSocket client
type WebSocketClient struct {
	conn     *websocket.Conn
	rideID   string
	userID   string
	send     chan interface{}
	done     chan bool
}

func NewWebSocketHandler(getRideHandler *application.GetRideHandler, logger *zap.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		getRideHandler: getRideHandler,
		logger:         logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // SECURITY: In production, implement proper CORS
			},
		},
		clients: make(map[string][]*WebSocketClient),
	}
}

// Subscribe handles WebSocket subscription for ride updates
func (h *WebSocketHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	rideID := r.URL.Query().Get("ride_id")
	userID := r.URL.Query().Get("user_id")

	if rideID == "" || userID == "" {
		http.Error(w, "ride_id and user_id required", http.StatusBadRequest)
		return
	}

	// Verify user has access to this ride
	ctx := r.Context()
	ride, err := h.getRideHandler.Handle(ctx, rideID)
	if err != nil {
		http.Error(w, "ride not found", http.StatusNotFound)
		return
	}

	if ride.PassengerID != userID && ride.DriverID != userID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	// Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", zap.Error(err))
		return
	}

	client := &WebSocketClient{
		conn:   conn,
		rideID: rideID,
		userID: userID,
		send:   make(chan interface{}, 10),
		done:   make(chan bool),
	}

	// Register client
	h.clientsMutex.Lock()
	h.clients[rideID] = append(h.clients[rideID], client)
	h.clientsMutex.Unlock()

	h.logger.Info("websocket client connected",
		zap.String("rideID", rideID),
		zap.String("userID", userID))

	// Send initial ride state
	client.send <- map[string]interface{}{
		"type": "RIDE_STATE",
		"data": ride,
	}

	// Handle client
	go h.handleClient(client)
}

// handleClient manages a single WebSocket client
func (h *WebSocketHandler) handleClient(client *WebSocketClient) {
	defer func() {
		// Unregister client
		h.clientsMutex.Lock()
		clients := h.clients[client.rideID]
		for i, c := range clients {
			if c == client {
				h.clients[client.rideID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(h.clients[client.rideID]) == 0 {
			delete(h.clients, client.rideID)
		}
		h.clientsMutex.Unlock()

		client.conn.Close()
		close(client.done)
		h.logger.Info("websocket client disconnected",
			zap.String("rideID", client.rideID),
			zap.String("userID", client.userID))
	}()

	// Send loop
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg := <-client.send:
				client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := client.conn.WriteJSON(msg); err != nil {
					h.logger.Error("websocket write failed", zap.Error(err))
					client.conn.Close()
					return
				}

			case <-ticker.C:
				// Heartbeat
				client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					h.logger.Error("websocket heartbeat failed", zap.Error(err))
					client.conn.Close()
					return
				}

			case <-client.done:
				return
			}
		}
	}()

	// Read loop (for client commands)
	for {
		var msg map[string]interface{}
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("websocket read error", zap.Error(err))
			}
			return
		}

		// Handle client commands (e.g., location updates for drivers)
		h.handleClientMessage(client, msg)
	}
}

// handleClientMessage processes incoming WebSocket messages
func (h *WebSocketHandler) handleClientMessage(client *WebSocketClient, msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "LOCATION_UPDATE":
		// Driver sending location updates
		// This would be relayed to GPS service and other interested clients
		h.broadcastToRide(client.rideID, map[string]interface{}{
			"type": "LOCATION_UPDATE",
			"user": client.userID,
			"data": msg,
		})

	case "PING":
		// Client heartbeat
		client.send <- map[string]interface{}{
			"type": "PONG",
		}
	}
}

// broadcastToRide sends a message to all clients watching a ride
func (h *WebSocketHandler) broadcastToRide(rideID string, msg interface{}) {
	h.clientsMutex.RLock()
	clients := h.clients[rideID]
	h.clientsMutex.RUnlock()

	for _, client := range clients {
		select {
		case client.send <- msg:
		default:
			// Client's send channel is full, drop message
			h.logger.Warn("websocket send channel full, dropping message",
				zap.String("rideID", rideID),
				zap.String("userID", client.userID))
		}
	}
}

// BroadcastRideUpdate sends a ride state update to all watchers
func (h *WebSocketHandler) BroadcastRideUpdate(ride *domain.Ride) {
	h.broadcastToRide(ride.ID, map[string]interface{}{
		"type": "RIDE_UPDATE",
		"data": ride,
	})
}

// Message types for WebSocket communication
type WebSocketMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
