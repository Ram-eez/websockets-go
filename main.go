package main

import (
	"websockets/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret-key"))
	router.Use(sessions.Sessions("mysess", store))
	routes.RegisterRoutes(router)
	router.Run(":8080")
}
