package handlers

import (
	"fmt"
	"net/http"
	"websockets/config"
	"websockets/middleware"
	"websockets/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	users *config.Repository
}

func NewHandler(userRepo *config.Repository) *Handler {
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

func (h *Handler) Login(c *gin.Context) {
	var newUsr models.User
	newUsr.Username = c.PostForm("username")
	newUsr.Password = c.PostForm("password")

	user, err := h.users.SearchUser(&newUsr)
	if err != nil {
		fmt.Println("could not find valid user: ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	fmt.Println("retrieved user: ", user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUsr.Password)); err != nil {
		fmt.Println("invalid password")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	tokenString, err := middleware.CreateToken(user)
	if err != nil {
		fmt.Println("err : ", err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/chat")
}
