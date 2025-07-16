package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jezhtech/prince-group-backend/models"
)

// Cashfree Payment Links API Request/Response structures based on official documentation
type CashfreePaymentLinkRequest struct {
	LinkID          string `json:"link_id,omitempty"`
	CustomerDetails struct {
		CustomerEmail string `json:"customer_email"`
		CustomerName  string `json:"customer_name"`
		CustomerPhone string `json:"customer_phone"`
	} `json:"customer_details"`
	LinkAmount   float64 `json:"link_amount"`
	LinkCurrency string  `json:"link_currency"`
	LinkPurpose  string  `json:"link_purpose"`
	LinkMeta     struct {
		NotifyURL string `json:"notify_url"`
		ReturnURL string `json:"return_url"`
	} `json:"link_meta"`
	LinkNotify struct {
		SendEmail bool `json:"send_email"`
		SendSMS   bool `json:"send_sms"`
	} `json:"link_notify"`
	LinkAutoReminders bool `json:"link_auto_reminders"`
}

type CashfreePaymentLinkResponse struct {
	CfLinkID        string  `json:"cf_link_id"`
	LinkID          string  `json:"link_id"`
	LinkStatus      string  `json:"link_status"`
	LinkCurrency    string  `json:"link_currency"`
	LinkAmount      float64 `json:"link_amount"`
	LinkAmountPaid  float64 `json:"link_amount_paid"`
	LinkPurpose     string  `json:"link_purpose"`
	LinkCreatedAt   string  `json:"link_created_at"`
	CustomerDetails struct {
		CustomerName  string `json:"customer_name"`
		CustomerPhone string `json:"customer_phone"`
		CustomerEmail string `json:"customer_email"`
	} `json:"customer_details"`
	LinkMeta struct {
		NotifyURL string `json:"notify_url"`
		ReturnURL string `json:"return_url"`
	} `json:"link_meta"`
	LinkURL        string `json:"link_url"`
	LinkExpiryTime string `json:"link_expiry_time"`
	LinkQRCode     string `json:"link_qrcode"`
	LinkNotify     struct {
		SendSMS   bool `json:"send_sms"`
		SendEmail bool `json:"send_email"`
	} `json:"link_notify"`
}

type CashfreePaymentDetails struct {
	CfPaymentID     string  `json:"cf_payment_id"`
	OrderID         string  `json:"order_id"`
	Entity          string  `json:"entity"`
	OrderAmount     float64 `json:"order_amount"`
	PaymentAmount   float64 `json:"payment_amount"`
	PaymentCurrency string  `json:"payment_currency"`
	PaymentStatus   string  `json:"payment_status"`
	PaymentMessage  string  `json:"payment_message"`
	BankReference   string  `json:"bank_reference"`
	AuthID          string  `json:"auth_id"`
	PaymentMethod   struct {
		PaymentMethod string `json:"payment_method"`
		CardNetwork   string `json:"card_network,omitempty"`
		CardIssuer    string `json:"card_issuer,omitempty"`
		CardType      string `json:"card_type,omitempty"`
	} `json:"payment_method"`
	PaymentTime     string `json:"payment_time"`
	PaymentGroup    string `json:"payment_group"`
	CustomerDetails struct {
		CustomerID    string `json:"customer_id"`
		CustomerName  string `json:"customer_name"`
		CustomerEmail string `json:"customer_email"`
		CustomerPhone string `json:"customer_phone"`
	} `json:"customer_details"`
}

type PaymentStatus struct {
	OrderID       string    `json:"orderId"`
	PaymentID     string    `json:"paymentId"`
	TransactionID string    `json:"transactionId"`
	Status        string    `json:"status"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"paymentMethod"`
	PaymentDate   time.Time `json:"paymentDate"`
	Message       string    `json:"message,omitempty"`
}

type CreatePaymentRequest struct {
	BookingID     string  `json:"bookingId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	CustomerName  string  `json:"customerName"`
	CustomerEmail string  `json:"customerEmail"`
	CustomerPhone string  `json:"customerPhone"`
	OrderNote     string  `json:"orderNote,omitempty"`
}

// CreatePaymentLink creates a new payment link with Cashfree
func CreatePaymentLink(c *gin.Context) {
	// Get Firebase ID from context (set by middleware)
	_, exists := c.Get("firebaseId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Generate a unique LinkID
	linkID := uuid.New().String()

	booking, err := models.GetBookingByBookingNumber(req.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking"})
		return
	}

	booking.PaymentLinkID = linkID
	models.UpdateBooking(booking)

	// Prepare Cashfree Payment Link request according to their API documentation
	cashfreeReq := CashfreePaymentLinkRequest{
		LinkID:       linkID,
		LinkAmount:   req.Amount,
		LinkCurrency: req.Currency,
		LinkPurpose:  req.OrderNote,
		CustomerDetails: struct {
			CustomerEmail string `json:"customer_email"`
			CustomerName  string `json:"customer_name"`
			CustomerPhone string `json:"customer_phone"`
		}{
			CustomerEmail: req.CustomerEmail,
			CustomerName:  req.CustomerName,
			CustomerPhone: req.CustomerPhone,
		},
		LinkMeta: struct {
			NotifyURL string `json:"notify_url"`
			ReturnURL string `json:"return_url"`
		}{
			NotifyURL: os.Getenv("CASHFREE_NOTIFY_URL"),
			ReturnURL: os.Getenv("CASHFREE_RETURN_URL") + "?orderId=" + linkID + "&bookingId=" + req.BookingID,
		},
		LinkNotify: struct {
			SendEmail bool `json:"send_email"`
			SendSMS   bool `json:"send_sms"`
		}{
			SendEmail: false,
			SendSMS:   false,
		},
		LinkAutoReminders: true,
	}

	// Call Cashfree Payment Links API
	response, err := createCashfreePaymentLink(cashfreeReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment link: " + err.Error()})
		return
	}

	// Return response in our expected format with the custom linkID
	c.JSON(http.StatusOK, gin.H{
		"paymentSessionId": response.CfLinkID,
		"orderId":          linkID, // Use our custom linkID instead of Cashfree's
		"orderAmount":      response.LinkAmount,
		"orderCurrency":    response.LinkCurrency,
		"status":           response.LinkStatus,
		"paymentLink":      response.LinkURL,
		"linkId":           linkID, // Include the linkID for frontend reference
	})
}

// CheckPaymentStatus checks the status of a payment
func CheckPaymentStatus(c *gin.Context) {
	// Get Firebase ID from context (set by middleware)
	_, exists := c.Get("firebaseId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	linkID := c.Param("linkId")
	if linkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link ID is required"})
		return
	}

	// Call Cashfree API to check status
	status, err := getCashfreePaymentLinkStatus(linkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check payment status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetPaymentHistory gets payment history for the authenticated user
func GetPaymentHistory(c *gin.Context) {
	// Get Firebase ID from context (set by middleware)
	_, exists := c.Get("firebaseId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Implement payment history retrieval from database
	// For now, return empty array
	payments := []PaymentStatus{}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// PaymentWebhook handles Cashfree webhook notifications
func PaymentWebhook(c *gin.Context) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Verify webhook signature (implement signature verification)
	if !verifyWebhookSignature(c.Request.Header.Get("X-Webhook-Signature"), body) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
		return
	}

	// Parse webhook data
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
		return
	}

	// Process webhook data
	// Update booking status based on payment status
	// Send confirmation emails, etc.

	fmt.Printf("Webhook received: %+v\n", webhookData)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// PaymentCallback handles the redirect from Cashfree after payment
func PaymentCallback(c *gin.Context) {
	// Get query parameters from Cashfree redirect
	orderID := c.Query("order_id")
	paymentStatus := c.Query("payment_status")
	cfLinkID := c.Query("cf_link_id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Get the booking by our custom orderID (link_id)
	booking, err := models.GetBookingByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking"})
		return
	}

	// Update booking status based on payment status
	var paymentStatusStr string

	switch paymentStatus {
	case "SUCCESS":
		paymentStatusStr = "success"
	case "FAILED":
		paymentStatusStr = "failed"
	default:
		paymentStatusStr = "pending"
	}

	// Update booking
	booking.PaymentStatus = paymentStatusStr
	booking.PaymentLinkID = cfLinkID

	updatedBooking, err := models.UpdateBooking(booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	// Redirect to frontend with status and orderID
	redirectURL := os.Getenv("FRONTEND_URL") + "/payment/result"
	if paymentStatus == "SUCCESS" {
		redirectURL += "?status=success&bookingId=" + updatedBooking.BookingNumber + "&orderId=" + orderID
	} else {
		redirectURL += "?status=failed&bookingId=" + updatedBooking.BookingNumber + "&orderId=" + orderID
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// Helper function to create Cashfree payment link using their API
func createCashfreePaymentLink(req CashfreePaymentLinkRequest) (*CashfreePaymentLinkResponse, error) {
	clientID := os.Getenv("CASHFREE_CLIENT_ID")
	clientSecret := os.Getenv("CASHFREE_CLIENT_SECRET")
	apiURL := os.Getenv("CASHFREE_API_URL")

	if clientID == "" || clientSecret == "" || apiURL == "" {
		return nil, fmt.Errorf("cashfree configuration missing")
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Debug: Print the JSON being sent
	fmt.Printf("Sending JSON to Cashfree: %s\n", string(jsonData))

	// Create HTTP request to Cashfree's Create Payment Link API
	httpReq, err := http.NewRequest("POST", apiURL+"/links", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers according to Cashfree documentation
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-client-secret", clientSecret)
	httpReq.Header.Set("x-api-version", "2025-01-01")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cashfree API error: %s", string(respBody))
	}

	// Parse response
	var cashfreeResp CashfreePaymentLinkResponse
	if err := json.Unmarshal(respBody, &cashfreeResp); err != nil {
		return nil, err
	}

	return &cashfreeResp, nil
}

// Helper function to get payment link status from Cashfree
func getCashfreePaymentLinkStatus(linkID string) (*PaymentStatus, error) {
	clientID := os.Getenv("CASHFREE_CLIENT_ID")
	clientSecret := os.Getenv("CASHFREE_CLIENT_SECRET")
	apiURL := os.Getenv("CASHFREE_API_URL")

	if clientID == "" || clientSecret == "" || apiURL == "" {
		return nil, fmt.Errorf("cashfree configuration missing")
	}

	// Create HTTP request to Cashfree's Get Orders for a Payment Link API
	// Using the correct endpoint from the documentation: /links/{link_id}/orders
	httpReq, err := http.NewRequest("GET", apiURL+"/links/"+linkID+"/orders", nil)
	if err != nil {
		return nil, err
	}

	// Set headers according to Cashfree documentation
	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-client-secret", clientSecret)
	httpReq.Header.Set("x-api-version", "2023-08-01") // Using the correct API version

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cashfree API error: %s", string(respBody))
	}

	// Parse response - Cashfree returns an array of orders
	var ordersResponse []struct {
		CfOrderID        string  `json:"cf_order_id"`
		OrderID          string  `json:"order_id"`
		OrderAmount      float64 `json:"order_amount"`
		OrderCurrency    string  `json:"order_currency"`
		OrderStatus      string  `json:"order_status"`
		OrderNote        string  `json:"order_note"`
		CreatedAt        string  `json:"created_at"`
		PaymentSessionID string  `json:"payment_session_id"`
		CustomerDetails  struct {
			CustomerName  string `json:"customer_name"`
			CustomerEmail string `json:"customer_email"`
			CustomerPhone string `json:"customer_phone"`
		} `json:"customer_details"`
	}

	if err := json.Unmarshal(respBody, &ordersResponse); err != nil {
		return nil, err
	}

	// If no orders found, return a pending status
	if len(ordersResponse) == 0 {
		return &PaymentStatus{
			OrderID:       linkID,
			PaymentID:     linkID,
			TransactionID: linkID,
			Status:        "pending",
			Amount:        0,
			Currency:      "INR",
			PaymentMethod: "cashfree",
			PaymentDate:   time.Now(),
			Message:       "No orders found for this payment link",
		}, nil
	}

	// Get the first order (most recent)
	order := ordersResponse[0]

	// Convert to our PaymentStatus format
	status := &PaymentStatus{
		OrderID:       order.OrderID,
		PaymentID:     order.CfOrderID,
		TransactionID: order.PaymentSessionID,
		Status:        mapOrderStatus(order.OrderStatus),
		Amount:        order.OrderAmount,
		Currency:      order.OrderCurrency,
		PaymentMethod: "cashfree",
		PaymentDate:   time.Now(),
		Message:       order.OrderNote,
	}

	return status, nil
}

// Helper function to get payment details from Cashfree
func getCashfreePaymentDetails(orderID string) (*CashfreePaymentDetails, error) {
	clientID := os.Getenv("CASHFREE_CLIENT_ID")
	clientSecret := os.Getenv("CASHFREE_CLIENT_SECRET")
	apiURL := os.Getenv("CASHFREE_API_URL")

	if clientID == "" || clientSecret == "" || apiURL == "" {
		return nil, fmt.Errorf("cashfree configuration missing")
	}

	// Create HTTP request to Cashfree's Get Payments for Order API
	httpReq, err := http.NewRequest("GET", apiURL+"/orders/"+orderID+"/payments", nil)
	if err != nil {
		return nil, err
	}

	// Set headers according to Cashfree documentation
	httpReq.Header.Set("x-client-id", clientID)
	httpReq.Header.Set("x-client-secret", clientSecret)
	httpReq.Header.Set("x-api-version", "2023-08-01")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cashfree API error: %s", string(respBody))
	}

	// Parse response - Cashfree returns an array of payments
	var paymentsResponse struct {
		Payments []CashfreePaymentDetails `json:"payments"`
	}
	if err := json.Unmarshal(respBody, &paymentsResponse); err != nil {
		return nil, err
	}

	// Return the first successful payment
	if len(paymentsResponse.Payments) > 0 {
		return &paymentsResponse.Payments[0], nil
	}

	return nil, fmt.Errorf("no payment details found")
}

// Helper function to map Cashfree payment link status to our status
func mapPaymentLinkStatus(cashfreeStatus string) string {
	switch cashfreeStatus {
	case "ACTIVE":
		return "pending"
	case "PAID":
		return "success"
	case "EXPIRED":
		return "failed"
	case "CANCELLED":
		return "failed"
	default:
		return "pending"
	}
}

// Helper function to map Cashfree order status to our status
func mapOrderStatus(cashfreeOrderStatus string) string {
	switch cashfreeOrderStatus {
	case "ACTIVE":
		return "pending"
	case "PAID":
		return "success"
	case "EXPIRED":
		return "failed"
	case "CANCELLED":
		return "failed"
	case "PENDING":
		return "pending"
	case "FAILED":
		return "failed"
	default:
		return "pending"
	}
}

// Helper function to verify webhook signature
func verifyWebhookSignature(signature string, body []byte) bool {
	// Implement signature verification based on Cashfree documentation
	// This is a placeholder - implement according to Cashfree's webhook signature verification
	webhookSecret := os.Getenv("CASHFREE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return false
	}

	// Create signature
	h := sha256.New()
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}
