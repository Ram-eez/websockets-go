package routes

import (
	"websockets/config"
	"websockets/handlers"
	manager "websockets/manager"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, handlers *handlers.Handler, repository *config.Repository) {
	manager := manager.NewManager(repository)

	// (frontend)
	router.Static("/static", "./static")
	router.GET("/chat", ServeIndex)
	router.GET("/register", ServeRegister)
	router.GET("/login", ServeLogin)
	// websocket endpoint
	router.GET("/ws", manager.ServeWS)

	// auth routes
	router.POST("/register", handlers.RegisterUsers)
	router.POST("/login", handlers.Login)
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
