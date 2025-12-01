package main

import (
    
	"UserService/internal/database"
	"UserService/internal/routes"
	"UserService/internal/config"
)


func main() {

	// Initialize database connection
    database.InitDB()
	

	serverPort := config.Load().Port

    r := routes.Setup()

    r.Run(":" + serverPort) // Use the port from the config
}