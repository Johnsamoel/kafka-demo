package utils

import (
	"log"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"fmt"
)

// LoadEnv loads environment variables from the .env file
func LoadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// SendEmail sends an email using Gmail's SMTP server
func SendEmail(recipientEmail, subject, body string) error {
	// Load environment variables
	LoadEnv()

	// Read sender email and password from .env
	senderEmail := os.Getenv("GMAIL_USER")
	senderPassword := os.Getenv("GMAIL_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		log.Fatal("SENDER_EMAIL or SENDER_PASSWORD is not set in the .env file")
		return nil
	}

	// Gmail SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	// Create a new email message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderEmail)
	mailer.SetHeader("To", recipientEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	// Create a dialer for Gmail's SMTP server
	dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// Send the email
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully!")
	return nil
}


// GenerateSystemTestEmail generates a system test email template for a new user
func GenerateSystemTestEmail(userName, systemName string) (subject, body string) {
	subject = fmt.Sprintf("Welcome to %s, %s! (System Test Email)", systemName, userName)

	bodyTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				background-color: #f9f9f9;
				color: #333;
				padding: 20px;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				background: #ffffff;
				border-radius: 8px;
				box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
				padding: 20px;
			}
			.header {
				background-color: #2196F3;
				color: white;
				padding: 10px;
				text-align: center;
				border-radius: 8px 8px 0 0;
			}
			.content {
				padding: 20px;
				text-align: left;
			}
			.footer {
				text-align: center;
				padding: 10px;
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Welcome to {{.SystemName}}!</h1>
			</div>
			<div class="content">
				<p>Dear {{.UserName}},</p>
				<p>Thank you for joining <strong>{{.SystemName}}</strong>!</p>
				<p>This email is being sent to confirm that the system is operating correctly and to ensure that our Kafka messaging service is functioning as expected.</p>
				<p><strong>No further action is required from your side.</strong></p>
				<p>If you receive this email, it means that everything is set up and working as intended.</p>
				<p>Thank you for being a part of our test!</p>
				<p>Best Regards,<br>The {{.SystemName}} Team</p>
			</div>
			<div class="footer">
				<p>If you have any questions, feel free to contact us at <a href="mailto:support@example.com">support@example.com</a>.</p>
				<p>&copy; {{.Year}} {{.SystemName}}. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`

	// Replace placeholders with actual data
	body = strings.ReplaceAll(bodyTemplate, "{{.SystemName}}", systemName)
	body = strings.ReplaceAll(body, "{{.UserName}}", userName)
	body = strings.ReplaceAll(body, "{{.Year}}", fmt.Sprintf("%d", 2025)) // Replace with the current year

	return subject, body
}
