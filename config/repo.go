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
