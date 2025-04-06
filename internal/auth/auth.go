package auth

import "fmt"

// Simple auth function for example
func Authenticate(username, password string) bool {
	// Imagine this connects to a DB or external service to authenticate
	if username == "admin" && password == "password" {
		return true
	}
	return false
}

func ExampleAuth() {
	if Authenticate("admin", "password") {
		fmt.Println("Authenticated!")
	} else {
		fmt.Println("Authentication Failed!")
	}
}
