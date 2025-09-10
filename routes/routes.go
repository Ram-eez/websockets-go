package routes

import (
	"websockets/handlers"
	manager "websockets/manager"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	manager := manager.NewManager()

	// (frontend)
	router.Static("/static", "./static")
	router.GET("/chat", ServeIndex)
	router.GET("/register", ServeRegister)
	router.GET("/login", ServeLogin)
	// websocket endpoint
	router.GET("/ws", manager.ServeWS)

	// auth routes
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)
	router.POST("/create-room", manager.CreateRoomHandler)
	router.GET("/room", handlers.GetAvalibleRooms)
}

func ServeIndex(c *gin.Context) {
	c.File("views/index.html")
}

func ServeRegister(c *gin.Context) {
	c.File("views/register.html")
}

func ServeLogin(c *gin.Context) {
	c.File("views/login.html")
}
