package main

import (
	manager "websockets/Manager"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	registerRoutes(router)
	router.Run(":8080")
}

func registerRoutes(router *gin.Engine) {
	Manager := manager.NewManager()
	router.GET("/", ServeIndex)
	router.GET("/ws", Manager.ServeWS)
}

func ServeIndex(c *gin.Context) {
	c.File("index.html")
}
