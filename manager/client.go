package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	Message string `json:"message"`
}

type Client struct {
	id         string
	connection *websocket.Conn
	manager    *Manager
	// egress channel is an unbuffered channel which is used to avoid concurrent writes on the websocket conn
	egress chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	NewUUID := uuid.New()
	return &Client{
		id:         NewUUID.String(),
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		for wsclient := range c.manager.clients {
			wsclient.egress <- payload
		}

		fmt.Println(messageType)
		fmt.Println(string(payload))

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					fmt.Println("conn closed :", err)
				}
				return
			}
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println("json unmarshalling err: ", err)
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, msg.getMessageHTML()); err != nil {
				fmt.Println("failed to send the message : ", err)

			}

			fmt.Println("message sent")
		default:
		}
	}
}

func (msg *Message) getMessageHTML() []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		fmt.Println("templete parsing err: ", err)
		return nil
	}

	var renderedMessage bytes.Buffer

	if err := tmpl.Execute(&renderedMessage, msg); err != nil {
		fmt.Println("execution err could not replace : ", err)
		return nil
	}

	fmt.Println("generated HTML with replaced obj: ", renderedMessage.String())

	return renderedMessage.Bytes()
}
