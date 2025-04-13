package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"learngolang-api/config"
	_ "learngolang-api/docs"
	"learngolang-api/internal/database"
	"learngolang-api/internal/middleware"
	"learngolang-api/pkg/api"

	"github.com/joho/godotenv"
)

// @title My API
// @version 1.0
// @description This is a sample API server for demonstrating Swagger in Go
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load initial .env
	_ = godotenv.Load()

	// Load environment-specific .env
	env := os.Getenv("APP_ENV")
	var envFile string
	if env == "production" {
		envFile = ".env.production"
	} else if env == "development" {
		envFile = ".env.development"
	} else {
		envFile = ".env.development"
	}
	fmt.Println("APP_ENV:", env)
	fmt.Println("Using env file:", envFile)

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("❌ Error loading %s: %v", envFile, err)
	}

	// Validate critical config
	config.LoadConfig()

	// Get the port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ PORT environment variable not set")
	}

	// Init DB
	database.InitDB()
	database.Migrate()

	// Routing
	mux := http.NewServeMux()
	apiURL := "/api/v1"

	// Routes
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
			middleware.RequireJWT(http.HandlerFunc(api.DeleteUserHandler)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(apiURL+"/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.LoginHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Swagger JSON + UI
	mux.Handle("/swagger.json", http.FileServer(http.Dir("./docs")))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger-ui"))))

	// Apply middleware
	allowedOrigins := []string{
		"http://localhost:" + port,
		"http://192.168.100.19:" + port,
		"http://127.0.0.1:" + port,
		"https://apigo.ahmadcloud.my.id",
		"http://apigo.ahmadcloud.my.id",
		"http://10.10.64.182:" + port,
	}

	handler := middleware.CORS(allowedOrigins)(middleware.Logger(mux))

	// Start server
	fmt.Printf("✅ Server is running on http://localhost:%s\n", port)
	fmt.Printf("📘 Swagger UI available at: http://localhost:%s/swagger/\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
