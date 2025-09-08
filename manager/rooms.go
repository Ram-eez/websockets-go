package manager

import (
	"websockets/models"

	"github.com/google/uuid"
)

type Room struct {
	id      string
	clients map[string]*Client
	manager *Manager

	broadcast chan models.Message
}

func NewRoom(manager *Manager) *Room {
	NewUUID := uuid.New()
	return &Room{
		id:        NewUUID.String(),
		clients:   make(map[string]*Client),
		manager:   manager,
		broadcast: make(chan models.Message),
	}
}
