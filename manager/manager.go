package manager

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"websockets/config"
	"websockets/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Manager struct {
	rooms map[string]*Room
	repo  *config.Repository
	sync.RWMutex
}

func NewManager(repo *config.Repository) *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
		repo:  repo,
	}
}

func (m *Manager) ServeWS(c *gin.Context) {
	log.Println("starting websocket new conn")

	// Upgrade the connection
	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("could not upgrade connection: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Authentication and validation phase
	// If anything fails here, we need to close the connection
	token, err := c.Cookie("Authorization")
	if err != nil {
		log.Printf("could not get jwt token from cookies: %v", err)
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Authentication required"))
		conn.Close()
		return
	}

	user, err := middleware.GetUserFromToken(token)
	if err != nil {
		log.Printf("user not authenticated for chat: %v", err)
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Invalid authentication token"))
		conn.Close()
		return
	}

	// Get room ID from query parameter
	roomID := c.Query("room")
	if roomID == "" {
		roomID = "lobby"
	}

	// Create client - from this point on, client owns the connection lifecycle
	client := NewClient(conn, m, user)

	// Start client read/write goroutines
	go client.readMessages()
	go client.writeMessages()

	// Join the specified room
	m.JoinRoom(client, roomID)
}

func (m *Manager) GetOrCreateRoom(roomID string) *Room {
	// First quick check with read lock (fast path for existing rooms)
	m.RLock()
	if r, ok := m.rooms[roomID]; ok {
		m.RUnlock()
		return r
	}
	m.RUnlock()

	// Acquire write lock to potentially create room
	m.Lock()
	defer m.Unlock()

	// Double-check: another goroutine might have created it while we waited for the lock
	if r, ok := m.rooms[roomID]; ok {
		return r
	}

	// Try to create room in database
	// Use a query that checks if room exists first
	existingRoom, err := m.repo.GetRoom(roomID)
	if err != nil {
		log.Printf("Error checking room existence: %v", err)
	}

	// If room doesn't exist in DB, create it
	if existingRoom == "" {
		err = m.repo.CreateRoom(roomID)
		if err != nil {
			log.Printf("Error creating room in DB: %v", err)
			// Continue anyway - room will exist in memory
		}
	}

	// Create room in memory
	r := NewRoom(m, roomID)
	m.rooms[roomID] = r
	go r.Run()

	return r
}

func (m *Manager) JoinRoom(client *Client, roomID string) {
	r := m.GetOrCreateRoom(roomID)
	r.register <- client
}

func (m *Manager) LeaveRoom(client *Client, roomID string) {
	m.RLock()
	r, ok := m.rooms[roomID]
	m.RUnlock()
	if ok {
		r.unregister <- client
	}
}

func (m *Manager) UnregisterEverywhere(client *Client) {
	client.closeOnce.Do(func() {
		m.RLock()
		for _, r := range m.rooms {
			r.unregister <- client
		}
		m.RUnlock()

		close(client.egress)
		_ = client.connection.Close()
	})
}

func (m *Manager) CreateRoomHandler(c *gin.Context) {
	roomID := "room-" + uuid.NewString()

	// Just create in DB
	err := m.repo.CreateRoom(roomID)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		c.String(http.StatusInternalServerError, "Failed to create room")
		return
	}

	// Return the new room content and trigger room list refresh
	c.Header("HX-Trigger", "refreshRooms")
	c.HTML(http.StatusOK, "room-content.html", gin.H{
		"RoomID": roomID,
	})
}

func (m *Manager) RoompageHandler(c *gin.Context) {
	roomID := c.Param("id")

	// Check if room exists in DB
	_, err := m.repo.GetRoom(roomID)
	if err != nil {
		c.String(http.StatusNotFound, "room not found")
		return
	}

	// Return the room content for embedding
	c.HTML(http.StatusOK, "room-content.html", gin.H{
		"RoomID": roomID,
	})
}

func (m *Manager) ListRooms(c *gin.Context) {
	rooms, err := m.repo.GetAllRooms()
	if err != nil {
		fmt.Println("could not fetch rooms: ", err)
		return
	}

	for _, id := range rooms {
		fmt.Fprintf(c.Writer,
			`<li>
				<button hx-get="/room/%s" hx-target="#room-output" hx-swap="innerHTML">%s</button>
			</li>`,
			id, id,
		)
	}
}

func (m *Manager) RemoveRoom(RoomID string) {
	m.Lock()
	defer m.Unlock()
	delete(m.rooms, RoomID)
}

// func isDuplicateError(err error) bool {
// 	errStr := err.Error()
// 	return strings.Contains(errStr, "duplicate") ||
// 		strings.Contains(errStr, "already exists") ||
// 		strings.Contains(errStr, "UNIQUE constraint") ||
// 		strings.Contains(errStr, "Duplicate entry")
// }
