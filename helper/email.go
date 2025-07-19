package helper

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/skip2/go-qrcode"
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
	subject := "Booking Confirmed - Rhythym of Kumari"

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
            <p>Rhythym of Kumari</p>
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

// generateQRCodeAsBase64 generates a QR code for the given text and returns it as a base64 encoded PNG
func generateQRCodeAsBase64(text string) (string, error) {
	// Generate QR code as PNG bytes
	qrBytes, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Convert to base64
	base64String := base64.StdEncoding.EncodeToString(qrBytes)
	return base64String, nil
}

// SendPaymentConfirmationEmail sends a payment confirmation email with ticket details
func SendPaymentConfirmationEmail(to, customerName, bookingNumber, ticketName, ticketCount, totalAmount, eventDate, eventLocation string, freeTickets int) error {
	subject := "Payment Successful - Rhythym of Kumari"

	// Calculate paid tickets
	paidTickets := ticketCount
	if freeTickets > 0 {
		// Extract the number from ticketCount string (e.g., "5" from "5 tickets")
		if count, err := strconv.Atoi(strings.Fields(ticketCount)[0]); err == nil {
			paidTickets = fmt.Sprintf("%d", count-freeTickets)
		}
	}

	// Generate QR code for the booking number
	qrCodeBase64, err := generateQRCodeAsBase64(bookingNumber)
	if err != nil {
		// If QR generation fails, use a placeholder
		qrCodeBase64 = ""
		fmt.Printf("Warning: Failed to generate QR code for booking %s: %v\n", bookingNumber, err)
	}

	// Create QR code image HTML
	qrCodeHTML := ""
	if qrCodeBase64 != "" {
		qrCodeHTML = fmt.Sprintf(`
			<div style="background: #fff; border: 1px solid #ddd; border-radius: 8px; padding: 20px; display: inline-block;">
				<img src="data:image/png;base64,%s" alt="QR Code for Booking %s" style="width: 120px; height: 120px; display: block; margin: 0 auto;" />
			</div>`, qrCodeBase64, bookingNumber)
	} else {
		qrCodeHTML = `
			<div style="background: #fff; border: 1px solid #ddd; border-radius: 8px; padding: 20px; display: inline-block;">
				<div style="width: 120px; height: 120px; background: #f8f9fa; border: 2px dashed #adb5bd; border-radius: 8px; display: flex; align-items: center; justify-content: center; color: #666; font-size: 12px;">
					QR Code<br>Unavailable
				</div>
			</div>`
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Payment Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #4eb4a7 0%%, #60afb4 100%%); color: white; padding: 40px; text-align: center; }
        .content { padding: 40px; }
        .success-icon { font-size: 48px; margin-bottom: 20px; }
        .ticket-card { background: linear-gradient(135deg, #f8f9fa 0%%, #e9ecef 100%%); border-radius: 15px; padding: 30px; margin: 30px 0; border-left: 5px solid #4eb4a7; }
        .booking-details { background: #f8f9fa; padding: 25px; border-radius: 10px; margin: 25px 0; }
        .free-ticket-badge { background: linear-gradient(135deg, #28a745 0%%, #20c997 100%%); color: white; padding: 8px 16px; border-radius: 20px; font-size: 14px; font-weight: bold; display: inline-block; margin: 10px 0; }
        .event-details { background: #fff3cd; border: 1px solid #ffeaa7; padding: 20px; border-radius: 10px; margin: 25px 0; }
        .footer { background: #f8f9fa; padding: 30px; text-align: center; color: #666; }
        .qr-section { background: #e9ecef; border: 2px dashed #adb5bd; border-radius: 10px; padding: 40px; text-align: center; margin: 20px 0; }
        .price-highlight { font-size: 24px; font-weight: bold; color: #28a745; }
        .ticket-count { font-size: 20px; font-weight: bold; color: #4eb4a7; }
        .important-note { background: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="success-icon">🎉</div>
            <h1>Payment Successful!</h1>
            <p>Your concert tickets have been confirmed</p>
        </div>
        
        <div class="content">
            <h2>Hello %s!</h2>
            <p>Thank you for your payment. Your booking has been successfully confirmed!</p>
            
            <div class="ticket-card">
                <h3 style="color: #4eb4a7; margin-top: 0;">🎫 Ticket Details</h3>
                <div style="display: flex; justify-content: space-between; align-items: center; margin: 15px 0;">
                    <span style="font-size: 18px; font-weight: bold;">%s</span>
                    <span class="ticket-count">%s</span>
                </div>
                %s
                <div style="margin-top: 20px;">
                    <span class="price-highlight">₹%s</span>
                    <span style="color: #666; font-size: 14px;"> (Total Amount Paid)</span>
                </div>
            </div>
            
            <div class="booking-details">
                <h3 style="color: #495057; margin-top: 0;">📋 Booking Information</h3>
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 15px; margin: 20px 0;">
                    <div>
                        <strong>Booking Number:</strong><br>
                        <span style="color: #4eb4a7; font-family: monospace; font-size: 16px;">%s</span>
                    </div>
                    <div>
                        <strong>Payment Status:</strong><br>
                        <span style="color: #28a745; font-weight: bold;">✅ Confirmed</span>
                    </div>
                    <div>
                        <strong>Tickets Paid For:</strong><br>
                        <span style="font-weight: bold;">%s</span>
                    </div>
                    <div>
                        <strong>Free Tickets:</strong><br>
                        <span style="color: #28a745; font-weight: bold;">%d</span>
                    </div>
                </div>
            </div>
            
            <div class="event-details">
                <h3 style="color: #856404; margin-top: 0;">🎵 Event Details</h3>
                <div style="margin: 15px 0;">
                    <strong>Date & Time:</strong> %s<br>
                    <strong>Venue:</strong> %s<br>
                    <strong>Artists:</strong> Aditya Rkay, Sri Nisha, Aparnaa Pratheep
                </div>
            </div>
            
            <div class="qr-section">
                <h4 style="margin-top: 0; color: #666;">🎫 Entry Pass</h4>
                <p style="color: #666; margin-bottom: 20px;">Scan this QR code at the venue for entry</p>
                %s
                <p style="color: #666; font-size: 12px; margin-top: 15px;">
                    <strong>Booking Number:</strong> %s
                </p>
            </div>
            
            <div class="important-note">
                <h4 style="margin-top: 0; color: #856404;">⚠️ Important Information</h4>
                <ul style="margin: 10px 0; padding-left: 20px;">
                    <li>Please arrive 30 minutes before the event starts</li>
                    <li>Bring a valid ID for verification</li>
                    <li>Entry will be denied without proper identification</li>
                    <li>No outside food or beverages allowed</li>
                    <li>Parking is available at the venue</li>
                </ul>
            </div>
        </div>
        
        <div class="footer">
            <p><strong>Prince Group Vista</strong></p>
            <p>Thank you for choosing us for your entertainment!</p>
            <p style="font-size: 12px; color: #999;">This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>`, customerName, ticketName, ticketCount,
		func() string {
			if freeTickets > 0 {
				return fmt.Sprintf(`<div class="free-ticket-badge">🎁 You got %d FREE ticket%s!</div>`, freeTickets, func() string {
					if freeTickets == 1 {
						return ""
					} else {
						return "s"
					}
				}())
			}
			return ""
		}(), totalAmount, bookingNumber, paidTickets, freeTickets, eventDate, eventLocation, qrCodeHTML, bookingNumber)

	return SendEmail(to, subject, htmlBody)
}
