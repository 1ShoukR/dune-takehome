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
	rooms   map[string]map[string]bool // formID -> clientID -> bool
	mu      sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[string]*websocket.Conn),
		rooms:   make(map[string]map[string]bool),
	}
}

func (ws *WebSocketService) HandleConnection(c *websocket.Conn) {
	clientID := c.RemoteAddr().String()

	ws.mu.Lock()
	ws.clients[clientID] = c
	ws.mu.Unlock()

	log.Printf("🔌 Client connected: %s", clientID)

	defer func() {
		ws.mu.Lock()
		delete(ws.clients, clientID)
		// Remove from all rooms
		for formID, clients := range ws.rooms {
			if clients[clientID] {
				delete(clients, clientID)
				if len(clients) == 0 {
					delete(ws.rooms, formID)
				}
			}
		}
		ws.mu.Unlock()
		c.Close()
		log.Printf("🔌❌ Client disconnected: %s", clientID)
	}()

	for {
		var msg map[string]interface{}
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("❌ WebSocket read error: %v", err)
			break
		}

		ws.handleMessage(clientID, msg)
	}
}

func (ws *WebSocketService) handleMessage(clientID string, msg map[string]interface{}) {
	eventType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch eventType {
	case "join-analytics":
		if formID, ok := msg["formId"].(string); ok {
			ws.joinRoom(clientID, formID)
			log.Printf("📊 Client %s joined analytics room for form: %s", clientID, formID)
		}
	case "leave-analytics":
		if formID, ok := msg["formId"].(string); ok {
			ws.leaveRoom(clientID, formID)
			log.Printf("📤 Client %s left analytics room for form: %s", clientID, formID)
		}
	}
}

func (ws *WebSocketService) joinRoom(clientID, formID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.rooms[formID] == nil {
		ws.rooms[formID] = make(map[string]bool)
	}
	ws.rooms[formID][clientID] = true
}

func (ws *WebSocketService) leaveRoom(clientID, formID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.rooms[formID] != nil {
		delete(ws.rooms[formID], clientID)
		if len(ws.rooms[formID]) == 0 {
			delete(ws.rooms, formID)
		}
	}
}

// BroadcastNewResponse sends real-time analytics update when a new response is submitted
func (ws *WebSocketService) BroadcastNewResponse(formID primitive.ObjectID, analytics *models.FormAnalytics) {
	roomName := formID.Hex()
	log.Printf("📡 Broadcasting analytics update to room: %s", roomName)

	message := map[string]interface{}{
		"type":      "analytics-update",
		"form_id":   formID.Hex(),
		"analytics": analytics,
		"timestamp": analytics.CreatedAt,
	}

	ws.broadcastToRoom(roomName, message)
}

// BroadcastFormUpdate sends updates when form structure changes
func (ws *WebSocketService) BroadcastFormUpdate(formID primitive.ObjectID, form *models.Form) {
	roomName := formID.Hex()
	log.Printf("📡 Broadcasting form update to room: %s", roomName)

	message := map[string]interface{}{
		"type":    "form-update",
		"form_id": formID.Hex(),
		"form":    form.ToResponse(),
	}

	ws.broadcastToRoom(roomName, message)
}

func (ws *WebSocketService) broadcastToRoom(formID string, message map[string]interface{}) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	if clients, exists := ws.rooms[formID]; exists {
		for clientID := range clients {
			if conn, exists := ws.clients[clientID]; exists {
				if err := conn.WriteJSON(message); err != nil {
					log.Printf("❌ Error broadcasting to client %s: %v", clientID, err)
				}
			}
		}
	}
}
