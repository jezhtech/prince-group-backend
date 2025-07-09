package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/models"
)

func GetReferral(c *gin.Context) {
	referralID := c.Param("id")

	referral, err := models.GetReferralByID(referralID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral"})
		return
	}

	c.JSON(200, gin.H{
		"referral": referral,
	})
}

func GetAllReferrals(c *gin.Context) {
	referrals, err := models.GetAllReferrals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referrals"})
		return
	}

	c.JSON(200, gin.H{
		"referrals": referrals,
	})
}

func CreateReferral(c *gin.Context) {
	var referral models.Referral

	if err := c.ShouldBindJSON(&referral); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	referral, err := models.CreateReferral(referral)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create referral"})
		return
	}

	c.JSON(200, gin.H{
		"referral": referral,
	})
}

func UpdateReferral(c *gin.Context) {
	referralID := c.Param("id")

	referral, err := models.GetReferralByID(referralID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral"})
		return
	}

	if err := c.ShouldBindJSON(&referral); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	referral, err = models.UpdateReferral(referral)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update referral"})
		return
	}

	c.JSON(200, gin.H{
		"referral": referral,
	})
}

func DeleteReferral(c *gin.Context) {
	referralID := c.Param("id")

	referral, err := models.GetReferralByID(referralID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral"})
		return
	}

	err = models.DeleteReferral(referral.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete referral"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Referral deleted successfully",
	})
}
