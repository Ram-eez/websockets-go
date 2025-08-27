package handlers

import (
	"fmt"
	"net/http"
	"websockets/middleware"
	"websockets/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var newUsr models.User
	newUsr.Username = c.PostForm("username")
	newUsr.Password = c.PostForm("password")

	user, err := FindUser(newUsr)
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

func FindUser(user models.User) (*models.User, error) {
	for _, founduser := range models.RegisteredUsers {
		if founduser.Username == user.Username {
			return &founduser, nil
		}
	}
	return nil, fmt.Errorf("could not find the user")
}
