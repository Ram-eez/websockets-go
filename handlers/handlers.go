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

func NewHandler(userRepo *config.UserRepository) *Handler {
	return &Handler{users: userRepo}
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

	if _, err := h.users.SearchUser(&newUsr); err == nil {
		fmt.Println("Username already in use pick another: ", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"Error": "Username already taken",
		})
		return
	}

	err = h.users.CreateUser(&newUsr)
	if err != nil {
		fmt.Println("could not create a new user : ", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"Error": "Database Error",
		})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
