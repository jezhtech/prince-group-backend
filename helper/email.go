package helper

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SendGridAPIKey string
	FromEmail      string
	FromName       string
}

// GetEmailConfig returns email configuration from environment variables
func GetEmailConfig() EmailConfig {
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	// Provide fallback values if not set
	if fromEmail == "" {
		fromEmail = "noreply@princegroup.com"
	}
	if fromName == "" {
		fromName = "Prince Group Vista"
	}
	return EmailConfig{
		SendGridAPIKey: os.Getenv("SENDGRID_API_KEY"),
		FromEmail:      fromEmail,
		FromName:       fromName,
	}
}

// SendEmail sends an email using SendGrid API
func SendEmail(to, subject, body string) error {
	config := GetEmailConfig()

	// If no SendGrid config, just log the email (for development)
	if config.SendGridAPIKey == "" {
		fmt.Printf("=== EMAIL SENT ===\nTo: %s\nSubject: %s\nBody: %s\n================\n", to, subject, body)
		return nil
	}

	// Validate required fields
	if config.FromEmail == "" {
		return fmt.Errorf("FROM_EMAIL environment variable is not set")
	}
	if to == "" {
		return fmt.Errorf("recipient email is required")
	}

	// Create email using SendGrid helpers
	from := mail.NewEmail(config.FromName, config.FromEmail)
	toEmail := mail.NewEmail("", to)                                 // Empty name for recipient
	message := mail.NewSingleEmail(from, subject, toEmail, "", body) // Empty plain text, HTML body

	// Create SendGrid client
	client := sendgrid.NewSendClient(config.SendGridAPIKey)

	// Send email
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	// Check response status
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid API error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

// SendOTPEmail sends an OTP email with a nice HTML template
func SendOTPEmail(to, otp string) error {
	subject := "Your OTP Code - Prince Group Vista"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Code</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #4eb4a7 0%%, #60afb4 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .otp-code { background: #4eb4a7; color: white; font-size: 32px; font-weight: bold; padding: 20px; text-align: center; border-radius: 10px; margin: 20px 0; letter-spacing: 5px; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Your OTP Code</h1>
            <p>Prince Group Vista - Secure Login</p>
        </div>
        <div class="content">
            <h2>Hello!</h2>
            <p>You requested a one-time password (OTP) to access your account. Here's your verification code:</p>
            
            <div class="otp-code">%s</div>
            
            <p><strong>This code will expire in 5 minutes.</strong></p>
            
            <div class="warning">
                <strong>⚠️ Security Notice:</strong><br>
                • Never share this code with anyone<br>
                • Prince Group Vista will never ask for this code via phone or email<br>
                • If you didn't request this code, please ignore this email
            </div>
            
            <p>If you're having trouble, you can request a new code from the login page.</p>
        </div>
        <div class="footer">
            <p>© 2024 Prince Group Vista. All rights reserved.</p>
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>`, otp)

	return SendEmail(to, subject, htmlBody)
}

// SendBookingConfirmationEmail sends a booking confirmation email
func SendBookingConfirmationEmail(to, customerName, bookingID, eventName, eventDate, eventLocation string) error {
	subject := "Booking Confirmed - Prince Group Mega Music Festival"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Booking Confirmed</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #4eb4a7 0%%, #60afb4 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .booking-details { background: white; padding: 20px; border-radius: 10px; margin: 20px 0; border-left: 4px solid #4eb4a7; }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
        .success-icon { font-size: 48px; margin-bottom: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="success-icon">🎉</div>
            <h1>Booking Confirmed!</h1>
            <p>Prince Group Mega Music Festival</p>
        </div>
        <div class="content">
            <h2>Hello %s!</h2>
            <p>Your booking has been successfully confirmed. We're excited to see you at the event!</p>
            
            <div class="booking-details">
                <h3>Booking Details</h3>
                <p><strong>Booking ID:</strong> %s</p>
                <p><strong>Event:</strong> %s</p>
                <p><strong>Date:</strong> %s</p>
                <p><strong>Location:</strong> %s</p>
            </div>
            
            <p><strong>Important Information:</strong></p>
            <ul>
                <li>Please arrive 30 minutes before the event starts</li>
                <li>Bring a valid ID for verification</li>
                <li>Entry will be denied without proper identification</li>
                <li>No outside food or beverages allowed</li>
            </ul>
            
            <p>If you have any questions, please contact our support team.</p>
        </div>
        <div class="footer">
            <p>© 2024 Prince Group Vista. All rights reserved.</p>
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>`, customerName, bookingID, eventName, eventDate, eventLocation)

	return SendEmail(to, subject, htmlBody)
}
