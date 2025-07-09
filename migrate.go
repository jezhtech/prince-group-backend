package main

import (
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/models"
)

func InitAutoMigrate() {
	// Migrate models in order to handle foreign key dependencies
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.Referral{})
	config.DB.AutoMigrate(&models.Ticket{})
	config.DB.AutoMigrate(&models.Booking{})
}
