package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	registerRoutes(router)
}

func registerRoutes(router *gin.Engine) {
	router.GET("/", ServeIndex)
}

func ServeIndex(c *gin.Context) {
	c.File("index.html")
}
