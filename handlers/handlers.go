package handlers

import (
	"fmt"
	"net/http"
	"websockets/config"
	"websockets/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	users *config.UserRepository
}

func (h *Handler) RegisterUsers(c *gin.Context) {
	var newUsr models.User

	newUsr.Username = c.PostForm("username")
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 10)
	if err != nil {
		fmt.Println("could not hash the password: ", err)
		return
	}

	newUsr.Password = string(hashedpass)
	newUsr.ID = uuid.New().String()

	h.users.CreateUser(&newUsr)
	if err != nil {
		fmt.Println("could not create a new user : ", err)
	}

	c.Redirect(http.StatusFound, "/login")
}
