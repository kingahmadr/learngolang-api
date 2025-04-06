package mailer

import (
	"fmt"
	"net/http"
	"net/smtp"
	"learngolang-api/internal/utils"
)

// SendEmailHandler is a handler for sending emails
func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch email details from the request
	recipient := r.URL.Query().Get("to")
	if recipient == "" {
		http.Error(w, "Missing 'to' parameter", http.StatusBadRequest)
		return
	}

	// Send email
	err := SendEmail(recipient, "Test Subject", "This is a test email body.")
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Email sent successfully to %s!", recipient)
}

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	// SMTP Server settings
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	username := "your-email@example.com"
	password := "your-email-password"
	from := "your-email@example.com"

	// Message
	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	// Authentication
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Send email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}
