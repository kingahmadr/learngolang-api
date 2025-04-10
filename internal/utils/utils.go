// package utils

// import "fmt"

// // Example utility function
// func PrintMessage(msg string) {
// 	fmt.Println("Message:", msg)
// }

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LogInfo logs an informational message
func LogInfo(message string) {
	fmt.Println("[INFO] " + message)
}

// DecodeAndValidateJSON decodes request body into the given struct and validates allowed fields
func DecodeAndValidateJSON(r *http.Request, allowedFields []string, target interface{}) error {
	// Decode into map first to check for extra fields
	var requestMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestMap); err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}

	// Check for extra/unexpected fields
	allowed := make(map[string]bool)
	for _, field := range allowedFields {
		allowed[field] = true
	}
	for key := range requestMap {
		if !allowed[key] {
			return fmt.Errorf("invalid field: %s", key)
		}
	}

	// Marshal and decode again into target struct
	bodyBytes, _ := json.Marshal(requestMap)
	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}

	return nil
}
