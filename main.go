package main

import (
    "log"
    "os"

    "github.com/chidi150c/database/gorm"
    "github.com/chidi150c/database/server"
)

func main() {
    // Get the port and host site from environment variables
    port := os.Getenv("PORT4")
    hostSite := os.Getenv("HOSTSITE")

    // Initialize your DBServices
    dbName := "myapp.db" // Replace with your desired database name
    ts, err := gorm.NewDBServices(dbName)
    if err != nil {
        log.Fatalf("Failed to initialize DBServices: %v", err)
    }

    // Setup and Start Web Server
    th := server.NewTradeHandler(ts, hostSite)
    server := server.NewServer(port, th)

    // Start the web server
    if err := server.Open(); err != nil {
        log.Fatalf("Unable to open server for listen and serve: %v", err)
    }
}
