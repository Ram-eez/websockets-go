package config

import (
	"database/sql"
	"websockets/models"
)

type Repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(User *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, id, password) VALUES ($1, $2, $3)", User.Username, User.ID, User.Password)
	return err
}

func (r *Repository) SearchUser(User *models.User) (*models.User, error) {
	var userRecord models.User
	err := r.db.QueryRow("SELECT username, id, password FROM users where username = $1", User.Username).Scan(&userRecord.Username, &userRecord.ID, &userRecord.Password)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	} else {
		return &userRecord, nil
	}
}

func (r *Repository) CreateRoom(id string) error {
	_, err := r.db.Exec("INSERT INTO rooms (id) VALUES ($1)", id)
	return err
}

func (r *Repository) GetAllRooms() ([]string, error) {
	rows, err := r.db.Query("SELECT id FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []string
	for rows.Next() {
		var roomid string
		if err := rows.Scan(&roomid); err != nil {
			continue
		}
		rooms = append(rooms, roomid)
	}
	return rooms, nil
}

func (r *Repository) GetRoom(roomID string) (string, error) {
	var id string
	err := r.db.QueryRow("SELECT id FROM rooms WHERE id = $1", roomID).Scan(&id)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", nil
	}
	return id, nil
}

func (r *Repository) AddMessage(Message *models.Message) error {
	_, err := r.db.Exec("INSERT INTO messages (username, message, roomid) VALUES ($1, $2, $3)", Message.Username, Message.Message, Message.RoomID)
	return err
}

func (r *Repository) GetAllRoomMessages(roomID string, limit int) ([]*models.Message, error) {
	// Get recent messages in reverse chronological order, then reverse in code
	query := `
		SELECT username, message, roomid 
		FROM messages 
		WHERE roomid = $1 
		ORDER BY id DESC 
		LIMIT $2
	`

	rows, err := r.db.Query(query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.Username, &m.Message, &m.RoomID); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse the slice so oldest message is first (chronological order for chat)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
