package config

import (
	"database/sql"
	"websockets/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(User *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, id, password) VALUES ($1, $2, $3)", User.Username, User.ID, User.Password)
	return err
}

func (r *UserRepository) SearchUser(User *models.User) (*models.User, error) {
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
