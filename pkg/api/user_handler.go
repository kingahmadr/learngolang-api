package api

import (
	"encoding/json"
	"fmt"
	"io"
	"learngolang-api/internal/database"
	"learngolang-api/internal/utils"
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
	// var requestBody schema.UserRequest

	var createUser schema.UserRequest

	err := utils.DecodeAndValidateJSON(r, []string{"name", "email", "password"}, &createUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if password is empty
	if createUser.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Check if user already exists in the database
	var existingUser models.User
	if err := database.DB.Where("email = ?", createUser.Email).First(&existingUser).Error; err == nil {
		// User already exists
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash the password (consider using bcrypt for security)
	hashedPassword := hashPassword(createUser.Password)

	// Create user object
	user := models.User{
		Name:     createUser.Name,
		Email:    createUser.Email,
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
	// Split the URL and get the last part
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]

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
	// Split the URL and get the last part
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]

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

	// Split the URL and get the last part
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	// First, check if the user exists (excluding soft-deleted ones)
	if err := database.DB.First(&user, id).Error; err != nil {
		// Could be record not found, or already soft-deleted
		http.Error(w, "User not found or already deleted", http.StatusNotFound)
		return
	}

	// Check if already soft-deleted (safety net — optional)
	if user.DeletedAt.Valid {
		http.Error(w, "User already deleted", http.StatusConflict)
		return
	}

	// Proceed to delete
	if err := database.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
