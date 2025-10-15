package main

import (
	"log"
	"websockets/config"
	"websockets/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	conn, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	routes.RegisterRoutes(router)
	router.LoadHTMLGlob("views/*.html")
	router.Run(":8080")
}
