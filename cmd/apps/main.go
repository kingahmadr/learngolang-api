package main

import (
	"fmt"
	"log"
	"net/http"
	"learngolang-api/internal/mailer"
	"learngolang-api/internal/middleware"
	"learngolang-api/pkg/api"
)

func main() {
	// Load configuration
	// config.LoadConfig() // Example if you use external config

	// Initialize HTTP router and middleware
	http.HandleFunc("/", api.HomeHandler)
	http.HandleFunc("/send-email", mailer.SendEmailHandler)

	// Middleware for logging requests
	http.Handle("/", middleware.Logger(http.DefaultServeMux))

	// Start the server
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
