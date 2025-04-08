package main

import (
	"fmt"
	"log"
	"net/http"
	"learngolang-api/internal/mailer"
	"learngolang-api/internal/middleware"
	"learngolang-api/pkg/api"
	"github.com/joho/godotenv"
	"learngolang-api/internal/database"
    // "learngolang-api/pkg/models"
)

func main() {

	// Load environment variables from .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file:", errEnv)
	}

	// Initialize database
	database.InitDB()
	database.Migrate()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", api.HomeHandler)
	mux.HandleFunc("/send-email", mailer.SendEmailHandler)

	// Register user routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.ListUsersHandler(w, r)
		} else if r.Method == http.MethodPost {
			api.CreateUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetUserHandler(w, r)
		case http.MethodPut:
			api.UpdateUserHandler(w, r)
		case http.MethodDelete:
			api.DeleteUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Wrap mux with middleware
	handler := middleware.Logger(mux)

	// Start the server
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
