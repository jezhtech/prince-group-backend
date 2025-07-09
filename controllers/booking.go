package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/models"
)

func GetBooking(c *gin.Context) {
	bookingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := models.GetBookingByID(uint(bookingID))
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
	bookingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := models.GetBookingByID(uint(bookingID))
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
}

func DeleteBooking(c *gin.Context) {
	bookingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	err = models.DeleteBooking(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(200, gin.H{"message": "Booking deleted successfully"})
}
