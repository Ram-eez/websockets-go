package manager

import (
	"log"
	"net/http"
	"sync"
	"websockets/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) ServeWS(c *gin.Context) {

	log.Println("starting websocket new conn")

	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("could not upgrade conn: ", conn)
		return
	}
	token, err := c.Cookie("Authorization")
	if err != nil {
		log.Fatal("could not get jwt token from cookies", err)
		return
	}

	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(models.SecretKey), nil
	})

	claims := tokenString.Claims.(jwt.MapClaims)

	var user models.User
	username := claims["username"].(string)
	userID := claims["userID"].(string)
	user.Username = username
	user.ID = userID
	client := NewClient(conn, m, &user)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
