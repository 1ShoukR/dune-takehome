package services

import (
	"log"
	"sync"

	"dune-takehome-server/models"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebSocketService struct {
	clients map[string]*websocket.Conn
	rooms   map[string]map[string]*websocket.Conn // roomID -> clientID -> connection
	mutex   sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[string]*websocket.Conn),
		rooms:   make(map[string]map[string]*websocket.Conn),
	}
}

// HandleConnection handles new WebSocket connections
func (ws *WebSocketService) HandleConnection(c *websocket.Conn) {
	clientID := generateClientID()
	
	ws.mutex.Lock()
	ws.clients[clientID] = c
	ws.mutex.Unlock()
	
	log.Printf("🔌 Client connected: %s", clientID)

	ws.sendMessage(c, map[string]interface{}{
		"type": "connected",
		"client_id": clientID,
	})

	for {
		var msg map[string]interface{}
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ WebSocket read error: %v", err)
			break
		}

		ws.handleMessage(clientID, c, msg)
	}

	ws.mutex.Lock()
	delete(ws.clients, clientID)
	// Remove from all rooms
	for roomID, room := range ws.rooms {
		delete(room, clientID)
		if len(room) == 0 {
			delete(ws.rooms, roomID)
		}
	}
	ws.mutex.Unlock()
	
	log.Printf("🔌❌ Client disconnected: %s", clientID)
}

func (ws *WebSocketService) handleMessage(clientID string, c *websocket.Conn, msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "join-analytics":
		if formID, ok := msg["form_id"].(string); ok {
			ws.joinRoom(clientID, c, formID)
			log.Printf("📊 Client %s joined analytics room for form: %s", clientID, formID)
		}
	case "leave-analytics":
		if formID, ok := msg["form_id"].(string); ok {
			ws.leaveRoom(clientID, formID)
			log.Printf("📤 Client %s left analytics room for form: %s", clientID, formID)
		}
	}
}

func (ws *WebSocketService) joinRoom(clientID string, c *websocket.Conn, roomID string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	
	if ws.rooms[roomID] == nil {
		ws.rooms[roomID] = make(map[string]*websocket.Conn)
	}
	ws.rooms[roomID][clientID] = c
}

func (ws *WebSocketService) leaveRoom(clientID, roomID string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	
	if room, exists := ws.rooms[roomID]; exists {
		delete(room, clientID)
		if len(room) == 0 {
			delete(ws.rooms, roomID)
		}
	}
}

func (ws *WebSocketService) sendMessage(c *websocket.Conn, message map[string]interface{}) {
	if err := c.WriteJSON(message); err != nil {
		log.Printf("❌ WebSocket write error: %v", err)
	}
}

// BroadcastNewResponse sends real-time analytics update when a new response is submitted
func (ws *WebSocketService) BroadcastNewResponse(formID primitive.ObjectID, analytics *models.FormAnalytics) {
	roomID := formID.Hex()
	log.Printf("📡 Broadcasting analytics update to room: %s", roomID)
	
	message := map[string]interface{}{
		"type":      "analytics-update",
		"form_id":   formID.Hex(),
		"analytics": analytics,
		"timestamp": analytics.CreatedAt,
	}

	ws.broadcastToRoom(roomID, message)
}

// BroadcastFormUpdate sends updates when form structure changes
func (ws *WebSocketService) BroadcastFormUpdate(formID primitive.ObjectID, form *models.Form) {
	roomID := formID.Hex()
	log.Printf("📡 Broadcasting form update to room: %s", roomID)
	
	message := map[string]interface{}{
		"type":    "form-update",
		"form_id": formID.Hex(),
		"form":    form.ToResponse(),
	}

	ws.broadcastToRoom(roomID, message)
}

func (ws *WebSocketService) broadcastToRoom(roomID string, message map[string]interface{}) {
	ws.mutex.RLock()
	room, exists := ws.rooms[roomID]
	if !exists {
		ws.mutex.RUnlock()
		return
	}
	
	// Create a copy of connections to avoid holding the lock while sending
	connections := make(map[string]*websocket.Conn)
	for clientID, conn := range room {
		connections[clientID] = conn
	}
	ws.mutex.RUnlock()

	// Send to all connections in the room
	for clientID, conn := range connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("❌ Failed to send message to client %s: %v", clientID, err)
			// Remove failed connection
			ws.mutex.Lock()
			delete(ws.rooms[roomID], clientID)
			delete(ws.clients, clientID)
			ws.mutex.Unlock()
		}
	}
}

func generateClientID() string {
	// Simple client ID generation
	return primitive.NewObjectID().Hex()
}