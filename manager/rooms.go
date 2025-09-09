package manager

import (
	"websockets/models"

	"github.com/google/uuid"
)

type Room struct {
	id      string
	clients map[string]*Client
	manager *Manager

	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
}

func NewRoom(manager *Manager, name string) *Room {
	NewUUID := uuid.New()
	return &Room{
		id:         NewUUID.String(),
		clients:    make(map[string]*Client),
		manager:    manager,
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.register:
			r.clients[client.id] = client

		case client := <-r.unregister:
			delete(r.clients, client.id)

		case msg := <-r.broadcast:
			for _, client := range r.clients {
				client.egress <- msg
			}
		}
	}
}
