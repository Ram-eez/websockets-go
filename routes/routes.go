package routes

import (
	"websockets/handlers"
	manager "websockets/manager"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, repository *handlers.Handler) {
	manager := manager.NewManager()

	// (frontend)
	router.Static("/static", "./static")
	router.GET("/chat", ServeIndex)
	router.GET("/register", ServeRegister)
	router.GET("/login", ServeLogin)
	// websocket endpoint
	router.GET("/ws", manager.ServeWS)

	// auth routes
	router.POST("/register", repository.RegisterUsers)
	router.POST("/login", repository.Login)
	router.POST("/create-room", manager.CreateRoomHandler)
	router.GET("/rooms", manager.ListRooms)
	router.GET("/room/:id", manager.RoompageHandler)
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
