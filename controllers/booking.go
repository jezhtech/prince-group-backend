package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/helper"
	"github.com/jezhtech/prince-group-backend/models"
)

func GetBooking(c *gin.Context) {
	bookingNumber := c.Param("bookingNumber")

	booking, err := models.GetBookingByBookingNumber(bookingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking"})
		return
	}

	c.JSON(200, gin.H{
		"booking": booking,
	})
}

func GetAllBookings(c *gin.Context) {
	bookings, err := models.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookings"})
		return
	}

	c.JSON(200, gin.H{
		"bookings": bookings,
	})
}

// GetAllBookingsPaginated returns paginated bookings with user and ticket data
func GetAllBookingsPaginated(c *gin.Context) {
	// Get pagination parameters from query string
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get paginated bookings
	paginatedBookings, err := models.GetAllBookingsPaginated(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookings"})
		return
	}

	c.JSON(200, gin.H{
		"bookings": paginatedBookings.Bookings,
		"pagination": gin.H{
			"total":       paginatedBookings.Total,
			"page":        paginatedBookings.Page,
			"pageSize":    paginatedBookings.PageSize,
			"totalPages":  paginatedBookings.TotalPages,
			"hasNext":     paginatedBookings.HasNext,
			"hasPrevious": paginatedBookings.HasPrevious,
		},
	})
}

// GetClientBookingsPaginated returns paginated bookings for client access
func GetClientBookingsPaginated(c *gin.Context) {
	// Get pagination parameters from query string
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get paginated bookings
	paginatedBookings, err := models.GetAllBookingsPaginated(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookings"})
		return
	}

	c.JSON(200, gin.H{
		"bookings": paginatedBookings.Bookings,
		"pagination": gin.H{
			"total":       paginatedBookings.Total,
			"page":        paginatedBookings.Page,
			"pageSize":    paginatedBookings.PageSize,
			"totalPages":  paginatedBookings.TotalPages,
			"hasNext":     paginatedBookings.HasNext,
			"hasPrevious": paginatedBookings.HasPrevious,
		},
	})
}

// GetClientBookingsStats returns overall booking statistics for client
func GetClientBookingsStats(c *gin.Context) {
	// Get all bookings to calculate stats
	allBookings, err := models.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking stats"})
		return
	}

	// Calculate stats
	total := len(allBookings)
	paid := 0
	pending := 0
	failed := 0

	for _, booking := range allBookings {
		switch booking.PaymentStatus {
		case "success":
			paid++
		case "pending":
			pending++
		case "failed":
			failed++
		}
	}

	c.JSON(200, gin.H{
		"stats": gin.H{
			"total":   total,
			"paid":    paid,
			"pending": pending,
			"failed":  failed,
		},
	})
}

func CreateBooking(c *gin.Context) {
	var booking models.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	booking.BookingNumber = helper.GenerateBookingNumber()

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		booking.BookingNumber = helper.GenerateBookingNumber()

		// Check if UserID already exists
		var existingBooking models.Booking
		err := config.DB.Where("booking_number = ?", booking.BookingNumber).First(&existingBooking).Error

		if err != nil {
			// UserID doesn't exist, we can use it
			break
		}

		// If we've tried maxRetries times, return an error
		if i == maxRetries-1 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique user ID after multiple attempts"})
			return
		}
	}

	booking, err := models.CreateBooking(booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(200, gin.H{
		"booking": booking,
	})
}

func UpdateBooking(c *gin.Context) {
	bookingNumber := c.Param("bookingNumber")

	booking, err := models.GetBookingByBookingNumber(bookingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking"})
		return
	}

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	booking, err = models.UpdateBooking(booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	c.JSON(200, gin.H{
		"booking": booking,
	})
}

func DeleteBooking(c *gin.Context) {
	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	err = models.DeleteBooking(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(200, gin.H{"message": "Booking deleted successfully"})
}

func GetBookingsByUserId(c *gin.Context) {
	firebaseId := c.GetString("firebaseId")

	if firebaseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase ID not found in context"})
		return
	}

	bookings, err := models.GetBookingsByUserId(firebaseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookings"})
		return
	}

	c.JSON(200, gin.H{
		"bookings": bookings,
	})
}

func CheckPayment(c *gin.Context) {
	bookingNumber := c.Param("bookingNumber")
	paymentLinkID := c.Query("paymentLinkId")

	if paymentLinkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment"})
		return
	}

	booking, err := models.GetBookingByBookingNumber(bookingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking"})
		return
	}

	if booking.PaymentLinkID != paymentLinkID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment"})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
	})
}
