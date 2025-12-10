package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spidey52/service-discovery/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.RWMutex
)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Add client to the list
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Printf("WebSocket client connected. Total clients: %d", len(clients))

	// Remove client when function returns
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		log.Printf("WebSocket client disconnected. Total clients: %d", len(clients))
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

type ServiceUpdateAction string

const (
	ActionRegister   ServiceUpdateAction = "register"
	ActionDeregister ServiceUpdateAction = "deregister"
	ActionHeartbeat  ServiceUpdateAction = "heartbeat"
)

type ServiceUpdate struct {
	Action  ServiceUpdateAction `json:"action"`
	Service models.Instance     `json:"service"`
}

func BroadcastMessage(msg ServiceUpdate) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("WebSocket send error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
