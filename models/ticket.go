package models

import (
	"time"

	"github.com/jezhtech/prince-group-backend/config"
)

type Ticket struct {
	ID                               uint      `gorm:"primaryKey" json:"id"`
	Name                             string    `gorm:"not null" json:"name"`
	Price                            int       `gorm:"not null" json:"price"`
	Type                             string    `gorm:"not null" json:"type"`
	Description                      string    `gorm:"not null" json:"description"`
	Benefits                         []string  `gorm:"serializer:json;not null" json:"benefits"`
	Status                           string    `gorm:"not null" json:"status"`
	TotalTickets                     int       `gorm:"not null" json:"totalTickets"`
	OfferPriceWithReferral           int       `gorm:"not null" json:"offerPriceWithReferral"`
	OfferPriceWithReferralAndYoutube int       `gorm:"not null" json:"offerPriceWithReferralAndYoutube"`
	AvailableTickets                 int       `gorm:"not null" json:"availableTickets"`
	CreatedAt                        time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt                        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func GetTicketByID(id uint) (Ticket, error) {
	var ticket Ticket

	err := config.DB.Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func GetAllTickets() ([]Ticket, error) {
	var tickets []Ticket
	err := config.DB.Find(&tickets).Error
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

func CreateTicket(ticket Ticket) (Ticket, error) {
	err := config.DB.Create(&ticket).Error
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func UpdateTicket(ticket Ticket) (Ticket, error) {
	err := config.DB.Save(&ticket).Error
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func DeleteTicket(id uint) error {
	err := config.DB.Delete(&Ticket{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
