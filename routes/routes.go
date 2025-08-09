package routes

import (
	manager "websockets/manager"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	Manager := manager.NewManager()
	router.GET("/", ServeIndex)
	router.GET("/ws", Manager.ServeWS)

	authRoutes := router.Group("/")
	{
		authRoutes.GET("/register")
		authRoutes.GET("/login")
	}
}

func ServeIndex(c *gin.Context) {
	c.File("views/index.html")
}
