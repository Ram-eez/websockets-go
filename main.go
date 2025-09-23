package main

import (
	"websockets/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.RegisterRoutes(router)
	router.LoadHTMLGlob("views/*.html")
	router.Run(":8080")
}
