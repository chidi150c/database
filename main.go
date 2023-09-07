package main

import (
	"log"
	"os"

	"github.com/chidi150c/database/gorm"
	"github.com/chidi150c/database/server"
    "github.com/chidi150c/database/webclient"
)

func main() {
	// Get the port and host site from environment variables
	port := os.Getenv("PORT4")
	hostSite := os.Getenv("HOSTSITE")

	// Initialize your DBServices
	dbName := "myapp.db" // Replace with your desired database name
	dbs, err := gorm.NewDBServices(dbName)
	if err != nil {
		log.Fatalf("Failed to initialize DBServices: %v", err)
	}

	// Check if the required tables (TradingSystem and AppData) exist in the database.
	// If they don't exist, create them.
	if err := dbs.CheckAndCreateTables(); err != nil {
		log.Fatalf("Error creating database tables: %v", err)
	}

	// Initialize your WebSocket service
	webSocketService := server.NewWebSocketService(hostSite) 

	// Initialize your TradeHandler
	th := server.NewTradeHandler(dbs, webSocketService, hostSite)

	// Setup and Start Web Server
	server := server.NewServer(port, th)

	// Start the web server
	if err := server.Open(); err != nil {
		log.Fatalf("Unable to open server for listen and serve: %v", err)
	}
}
