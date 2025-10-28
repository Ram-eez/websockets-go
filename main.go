package main

import (
	"fmt"
	"log"
	"websockets/config"
	"websockets/handlers"
	"websockets/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}
	defer conn.Close()
	err = conn.Ping()
	if err != nil {
		log.Fatal("failed to ping the database: ", err)
	}
	fmt.Println("Connected to the database successfully")
	userRepository := config.NewUserRepository(conn)
	h := handlers.NewHandler(userRepository)

	router := gin.Default()
	routes.RegisterRoutes(router, h, userRepository)
	router.LoadHTMLGlob("views/*.html")
	router.Run(":8080")
}
