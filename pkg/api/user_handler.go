package api

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"learngolang-api/internal/database"
	"learngolang-api/pkg/models"
	"learngolang-api/schema"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Hash password with bcrypt
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
	}
	return string(hashedPassword)
}

// POST /users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody schema.UserRequest

	// Decode the request body into a map to catch extra fields
	var requestMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestMap); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check for extra fields
	for key := range requestMap {
		if key != "name" && key != "email" && key != "password" {
			http.Error(w, fmt.Sprintf("Invalid field: %s", key), http.StatusBadRequest)
			return
		}
	}

	// Re-read the request body and decode into the struct
	// After checking for extra fields, now it's safe to read the body again
	bodyBytes, _ := json.Marshal(requestMap)
	if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if password is empty
	if requestBody.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Check if user already exists in the database
	var existingUser models.User
	if err := database.DB.Where("email = ?", requestBody.Email).First(&existingUser).Error; err == nil {
		// User already exists
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash the password (consider using bcrypt for security)
	hashedPassword := hashPassword(requestBody.Password)

	// Create user object
	user := models.User{
		Name:     requestBody.Name,
		Email:    requestBody.Email,
		Password: hashedPassword,
	}

	// Save user in the database
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Do not return password in the response
	user.Password = ""

	// Send back the user object without password
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GET /users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	var response []schema.UserResponse
	for _, user := range users {
		response = append(response, schema.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /users/{id}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// PUT /users/{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var updates models.User
	json.Unmarshal(body, &updates)

	user.Name = updates.Name
	user.Email = updates.Email

	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DELETE /users/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
