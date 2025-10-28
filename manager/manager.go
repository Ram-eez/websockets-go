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

	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("could not upgrade conn: ", conn)
	}
	token, err := c.Cookie("Authorization")
	if err != nil {
		log.Fatal("could not get jwt token from cookies", err)
	}

	user, err := middleware.GetUserFromToken(token)
	if err != nil {
		log.Fatal("user not authenticated for chat: ", err)
	}

	client := NewClient(conn, m, user)

	go client.readMessages()
	go client.writeMessages()

	roomID := c.Query("room")
	if roomID == "" {
		roomID = "lobby"
	}
	m.JoinRoom(client, roomID)

}

func (m *Manager) GetOrCreateRoom(roomID string) *Room {
	m.Lock()
	defer m.Unlock()
	// checking in memory first
	if r, ok := m.rooms[roomID]; ok {
		return r
	}
	room, err := m.repo.GetRoom(roomID)
	if err != nil {
		log.Printf("Error checking room in DB: %v", err)
	}

	// Create in DB if doesnt exist
	if room == "" {
		if err := m.repo.CreateRoom(roomID); err != nil {
			log.Printf("Error creating room in DB: %v", err)
		}
	}
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

	m.GetOrCreateRoom(roomID)

	// Return the new room content and trigger room list refresh
	c.Header("HX-Trigger", "refreshRooms")
	c.HTML(http.StatusOK, "room-content.html", gin.H{
		"RoomID": roomID,
	})
}

func (m *Manager) RoompageHandler(c *gin.Context) {
	roomID := c.Param("id")

	m.RLock()
	_, ok := m.rooms[roomID]
	m.RUnlock()

	if !ok {
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

	for _, roomID := range rooms {
		_ = m.GetOrCreateRoom(roomID)
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
