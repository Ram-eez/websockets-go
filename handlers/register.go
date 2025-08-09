package handlers

import (
	"fmt"
	"websockets/manager"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var newUsr manager.User

	newUsr.Username = c.PostForm("username")
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 10)
	if err != nil {
		fmt.Println("could not hash the password: ", err)
		return
	}

	newUsr.Password = string(hashedpass)
	newUsr.ID = uuid.New().String()

	manager.RegisteredUsers = append(manager.RegisteredUsers, newUsr)
}
