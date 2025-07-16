package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jezhtech/prince-group-backend/config"
)

type Booking struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BookingNumber string    `gorm:"column:booking_number;not null;unique" json:"bookingNumber"`
	UserID        string    `gorm:"column:user_id;not null" json:"userId"`
	ReferralID    string    `gorm:"column:referral_id;not null" json:"referralId"`
	TicketID      uint      `gorm:"column:ticket_id;not null" json:"ticketId"`
	TicketCount   int       `gorm:"column:ticket_count;not null" json:"ticketCount"`
	PaymentMethod string    `gorm:"column:payment_method;not null;" json:"paymentMethod"`
	PaymentStatus string    `gorm:"column:payment_status;not null;enum:pending,success,failed" json:"paymentStatus"`
	PaymentLinkID string    `gorm:"column:payment_link_id;not null" json:"paymentLinkId"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	User     User     `gorm:"foreignKey:UserID;references:FirebaseID" json:"user"`
	Ticket   Ticket   `gorm:"foreignKey:TicketID;references:ID" json:"ticket"`
	Referral Referral `gorm:"foreignKey:ReferralID;references:ReferralID" json:"referral"`
}

func GetBookingByBookingNumber(bookingNumber string) (Booking, error) {
	var booking Booking

	err := config.DB.Where("booking_number = ?", bookingNumber).
		Preload("User").
		Preload("Ticket").
		Preload("Referral").
		First(&booking).Error
	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}

func GetAllBookings() ([]Booking, error) {
	var bookings []Booking
	err := config.DB.Preload("User").Preload("Ticket").Preload("Referral").Find(&bookings).Error
	if err != nil {
		return []Booking{}, err
	}

	return bookings, nil
}

func CreateBooking(booking Booking) (Booking, error) {
	// Generate UUID if not provided
	if booking.ID == uuid.Nil {
		booking.ID = uuid.New()
	}

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

func DeleteBooking(id uuid.UUID) error {
	err := config.DB.Delete(&Booking{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

func GetBookingsByUserId(userId string) ([]Booking, error) {
	var bookings []Booking
	err := config.DB.Where("bookings.user_id = ?", userId).Preload("User").Preload("Ticket").Preload("Referral").Find(&bookings).Error
	if err != nil {
		return []Booking{}, err
	}

	return bookings, nil
}

func GetBookingByPaymentID(paymentID string) (Booking, error) {
	var booking Booking
	err := config.DB.Where("payment_id = ?", paymentID).Preload("User").Preload("Ticket").Preload("Referral").First(&booking).Error
	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}

func GetBookingByOrderID(orderID string) (Booking, error) {
	var booking Booking
	err := config.DB.Where("payment_link_id = ?", orderID).
		Preload("User").
		Preload("Ticket").
		Preload("Referral").
		First(&booking).Error
	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}

// GetBookingWithEmailData gets booking with preloaded user and ticket data
func GetBookingWithEmailData(bookingNumber string) (Booking, error) {
	var booking Booking

	err := config.DB.Where("booking_number = ?", bookingNumber).
		Preload("User").
		Preload("Ticket").
		First(&booking).Error

	if err != nil {
		return Booking{}, err
	}

	return booking, nil
}
