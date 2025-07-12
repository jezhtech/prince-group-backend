package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/models"
)

func GetTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := models.GetTicketByID(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"ticket":  ticket,
	})
}

func GetAllTickets(c *gin.Context) {
	tickets, err := models.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tickets": tickets,
	})
}

func CreateTicket(c *gin.Context) {
	var ticket models.Ticket

	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Debug: Log the received data
	fmt.Printf("Received ticket data: %+v\n", ticket)
	fmt.Printf("Benefits type: %T, value: %+v\n", ticket.Benefits, ticket.Benefits)

	// Validate required fields
	if ticket.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket name is required"})
		return
	}
	if ticket.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket price must be greater than 0"})
		return
	}
	if ticket.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket type is required"})
		return
	}
	if ticket.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket description is required"})
		return
	}
	if ticket.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket status is required"})
		return
	}
	if ticket.TotalTickets <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total tickets must be greater than 0"})
		return
	}
	if ticket.AvailableTickets <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available tickets must be greater than 0"})
		return
	}
	if ticket.AvailableTickets > ticket.TotalTickets {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available tickets cannot be greater than total tickets"})
		return
	}

	// Ensure Benefits is not nil and is a proper slice
	if ticket.Benefits == nil {
		ticket.Benefits = make([]string, 0)
	}

	createdTicket, err := models.CreateTicket(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Ticket created successfully",
		"ticket":  createdTicket,
	})
}

func UpdateTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Get existing ticket
	existingTicket, err := models.GetTicketByID(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Bind the update data
	var updateData models.Ticket
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Update only the provided fields
	if updateData.Name != "" {
		existingTicket.Name = updateData.Name
	}
	if updateData.Price > 0 {
		existingTicket.Price = updateData.Price
	}
	if updateData.Type != "" {
		existingTicket.Type = updateData.Type
	}
	if updateData.Description != "" {
		existingTicket.Description = updateData.Description
	}
	if updateData.Status != "" {
		existingTicket.Status = updateData.Status
	}
	if updateData.TotalTickets > 0 {
		existingTicket.TotalTickets = updateData.TotalTickets
	}
	if updateData.AvailableTickets >= 0 {
		existingTicket.AvailableTickets = updateData.AvailableTickets
	}
	if updateData.OfferPriceWithReferral >= 0 {
		existingTicket.OfferPriceWithReferral = updateData.OfferPriceWithReferral
	}
	if updateData.OfferPriceWithReferralAndYoutube >= 0 {
		existingTicket.OfferPriceWithReferralAndYoutube = updateData.OfferPriceWithReferralAndYoutube
	}

	// Update Benefits if provided
	if updateData.Benefits != nil {
		existingTicket.Benefits = updateData.Benefits
	}

	// Validate business rules
	if existingTicket.AvailableTickets > existingTicket.TotalTickets {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Available tickets cannot be greater than total tickets"})
		return
	}

	updatedTicket, err := models.UpdateTicket(existingTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket updated successfully",
		"ticket":  updatedTicket,
	})
}

func DeleteTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Check if ticket exists before deleting
	_, err = models.GetTicketByID(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	err = models.DeleteTicket(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket deleted successfully",
	})
}
