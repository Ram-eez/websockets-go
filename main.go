package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	registerRoutes(router)

	router.Run(":8080")
}

func registerRoutes(router *gin.Engine) {
	router.GET("/", ServeIndex)
	router.GET("/ws")
}

func ServeIndex(c *gin.Context) {
	c.File("index.html")
}
