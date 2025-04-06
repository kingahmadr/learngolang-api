package test

import (
	"net/http"
	"testing"
)

func TestSendEmail(t *testing.T) {
	// Simulate sending an email
	resp, err := http.Get("http://localhost:8080/send-email?to=test@example.com")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, but got %v", resp.StatusCode)
	}
}
