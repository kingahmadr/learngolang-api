package test

import (
	"testing"
	"your-project/internal/auth"
)

func TestAuthenticate(t *testing.T) {
	valid := auth.Authenticate("admin", "password")
	if !valid {
		t.Errorf("Expected true, but got false")
	}

	invalid := auth.Authenticate("user", "wrongpassword")
	if invalid {
		t.Errorf("Expected false, but got true")
	}
}
