package config

import (
	"os"
	"log"
)

// LoadConfig loads the environment configurations
func LoadConfig() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	// Other config settings
}
