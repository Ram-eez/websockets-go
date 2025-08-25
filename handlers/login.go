package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"websockets/manager"
	"websockets/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var newUsr manager.User
	newUsr.Username = c.PostForm("username")
	newUsr.Password = c.PostForm("password")

	user, err := FindUser(newUsr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("retrieved user: ", user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUsr.Password)); err != nil {
		fmt.Println("invalid password")
		return
	}

	tokenString, err := middleware.CreateToken(newUsr.Username)
	if err != nil {
		fmt.Println("err : ", err)
		return
	}

	U, err := json.Marshal(user)
	if err != nil {
		fmt.Println("could not marshal user struct")
		return
	}

	session := sessions.Default(c)
	session.Set("user", U)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/chat")
}

func FindUser(user manager.User) (*manager.User, error) {
	for _, founduser := range manager.RegisteredUsers {
		if founduser.Username == user.Username {
			return &founduser, nil
		}
	}
	return nil, fmt.Errorf("could not find the user")
}
