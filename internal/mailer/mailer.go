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

	// Log the email sending attempt
    utils.LogInfo("Sending email to: " + recipient)

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
    smtpHost := "192.168.100.19"
    smtpPort := "25"
    from := "no-reply@ahmadcloud.my.id"
    message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

    // No auth
    return smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, message)
}

