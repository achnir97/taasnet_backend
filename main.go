package main

import (
	"fmt"
	"log"
	"net/http"

	"taas-api/config"
	"taas-api/routes"
)

func main() {
	// Step 1: Connect to the database
	config.ConnectDB()

	// Step 2: Setup routes
	router := routes.SetupRoutes()

	// Step 3: Start the server
	port := ":8086" // Change the port as needed
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Printf("Server is running on http://localhost%s\n", port)
}
