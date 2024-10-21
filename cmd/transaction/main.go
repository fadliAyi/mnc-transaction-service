package main

import (
	"transaction-service/internal/background"
	"transaction-service/internal/config"
	"transaction-service/internal/database"
	"transaction-service/internal/routes"
)

func main() {
	config.Load()                  // Load environment variables
	database.ConnectDatabase()     // Connect to the database
	background.InitTransferQueue() // Initialize the transfer queue and start the worker
	router := routes.Setup()       // Initialize routes
	router.Run(":8080")            // Start the web server on port 8080
}
