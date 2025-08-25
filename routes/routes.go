package routes

import (
	"websockets/handlers"
	manager "websockets/manager"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	Manager := manager.NewManager()
	router.GET("/", ServeIndex)
	router.GET("/ws", Manager.ServeWS)

	authRoutes := router.Group("/")
	{
		authRoutes.GET("/register", handlers.RegisterHandler)
		authRoutes.GET("/login", handlers.LoginHandler)
	}
}

func ServeIndex(c *gin.Context) {
	c.File("views/index.html")
}
