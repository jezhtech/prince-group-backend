package controllers

import (
	"net/http"

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
