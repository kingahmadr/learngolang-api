package main

import (
	"fmt"
	_ "learngolang-api/docs"
	"learngolang-api/internal/database"
	"learngolang-api/internal/middleware"
	"learngolang-api/pkg/api"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// @title My API
// @version 1.0
// @description This is a sample API server for demonstrating Swagger in Go
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// Load environment variables
	_ = godotenv.Load()
	env := os.Getenv("APP_ENV")
	// Load the appropriate .env file based on APP_ENV
	fmt.Print("APP_ENV: ", env)
	var envFile string
	fmt.Print("envFile: ", envFile)
	if env == "production" {
		envFile = ".env.production"
	} else if env == "development" {
		envFile = ".env.development"
	} else {
		envFile = ".env.development"
	}

	errEnv := godotenv.Load(envFile)
	if errEnv != nil {
		log.Fatalf("Error loading .env file: %v", errEnv)
	}
	// Initialize database
	database.InitDB()
	database.Migrate()

	// Create a new ServeMux
	mux := http.NewServeMux()

	apiURL := "/api/v1"

	// User routes
	mux.HandleFunc(apiURL+"/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.ListUsersHandler(w, r)
		case http.MethodPost:
			api.CreateUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(apiURL+"/users/", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc(apiURL+"/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.LoginHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve Swagger JSON (generated by swag init)
	mux.Handle("/swagger.json", http.FileServer(http.Dir("./docs")))

	// Serve Swagger UI (static files copied from swagger-ui /dist)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger-ui"))))

	// Wrap with CORS and Logger
	allowedOrigins := []string{
		"http://localhost:8080",
		"http://192.168.100.19:8080",
	}

	handler := middleware.CORS(allowedOrigins)(middleware.Logger(mux))

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080...")
	fmt.Println("Swagger UI available at: http://localhost:8080/swagger/")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
