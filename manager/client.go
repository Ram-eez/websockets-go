package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"websockets/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id         string
	connection *websocket.Conn
	user       *models.User
	manager    *Manager
	// egress channel is an unbuffered channel which is used to avoid concurrent writes on the websocket conn
	egress    chan models.Message
	closeOnce sync.Once
}

// type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager, user *models.User) *Client {
	NewUUID := uuid.New()
	return &Client{
		id:         NewUUID.String(),
		connection: conn,
		user:       user,
		manager:    manager,
		egress:     make(chan models.Message),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.UnregisterEverywhere(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}

		var msg models.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			fmt.Println("json unmarshalling err: ", err)
			continue
		}
		msg.Username = c.user.Username

		roomID := msg.RoomID
		if roomID == "" {
			roomID = msg.Username + "'s lobby"
		}

		r := c.manager.GetorCreateRoom(roomID)
		r.broadcast <- msg

		fmt.Println(messageType)
		fmt.Println(string(payload))

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.UnregisterEverywhere(c)
	}()

	for msg := range c.egress {
		if err := c.connection.WriteMessage(
			websocket.TextMessage,
			[]byte(msg.GetMessageHTML()),
		); err != nil {
			fmt.Println("failed to send the message:", err)
			return
		}
	}
}
