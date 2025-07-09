package models

import (
	"time"

	"github.com/jezhtech/prince-group-backend/config"
)

type Booking struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        string    `gorm:"not null" json:"userId"`
	ReferralID    string    `gorm:"not null" json:"referralId"`
	TicketID      uint      `gorm:"not null" json:"ticketId"`
	Status        string    `gorm:"not null;enum:pending,confirmed,cancelled" json:"status"`
	PaymentMethod string    `gorm:"not null;" json:"paymentMethod"`
	PaymentStatus string    `gorm:"not null;enum:pending,success,failed" json:"paymentStatus"`
	PaymentDate   time.Time `gorm:"not null" json:"paymentDate"`
	PaymentURL    string    `gorm:"not null" json:"paymentUrl"`
	PaymentID     string    `gorm:"not null" json:"paymentId"`
	TransactionID string    `gorm:"not null" json:"transactionId"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	User     User     `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Ticket   Ticket   `gorm:"foreignKey:TicketID" json:"ticket"`
	Referral Referral `gorm:"foreignKey:ReferralID;references:ReferralID" json:"referral"`
}

func GetBookingByID(id uint) (Booking, error) {
	var booking Booking

	err := config.DB.Where("id = ?", id).First(&booking).Error
	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}

func GetAllBookings() ([]Booking, error) {
	var bookings []Booking
	err := config.DB.Find(&bookings).Error
	if err != nil {
		return []Booking{}, err
	}

	return bookings, nil
}

func CreateBooking(booking Booking) (Booking, error) {
	err := config.DB.Create(&booking).Error
	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}

func UpdateBooking(booking Booking) (Booking, error) {
	err := config.DB.Save(&booking).Error
	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}

func DeleteBooking(id uint) error {
	err := config.DB.Delete(&Booking{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
