package manager

import (
	"log"
	"net/http"
	"sync"
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
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
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

func (m *Manager) GetorCreateRoom(roomID string) *Room {
	m.Lock()
	defer m.Unlock()
	if r, ok := m.rooms[roomID]; ok {
		return r
	}
	r := NewRoom(m, roomID)
	m.rooms[roomID] = r
	go r.Run()
	return r

}

func (m *Manager) JoinRoom(client *Client, roomID string) {
	r := m.GetorCreateRoom(roomID)
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
	room := NewRoom(m, "room-"+uuid.NewString())

	m.Lock()
	m.rooms[room.id] = room
	m.Unlock()

	go room.Run()

	c.Header("HX-Redirect", "/room/"+room.id)
	c.Status(http.StatusSeeOther)

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

	c.HTML(http.StatusOK, "newroom.html", gin.H{
		"RoomID": roomID,
	})
}
