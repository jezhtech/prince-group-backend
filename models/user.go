package models

import (
	"time"

	"github.com/jezhtech/prince-group-backend/config"
)

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"not null;unique" json:"userId"`
	FirebaseID string    `gorm:"not null;unique" json:"firebaseId"`
	Role       string    `gorm:"not null;default:'user';enum:user,admin" json:"role"`
	FullName   string    `gorm:"not null" json:"fullName"`
	Email      string    `gorm:"not null;unique" json:"email"`
	Mobile     string    `gorm:"not null" json:"mobile"`
	Address    string    `gorm:"not null" json:"address"`
	City       string    `gorm:"not null" json:"city"`
	State      string    `gorm:"not null" json:"state"`
	Zip        string    `gorm:"not null" json:"zip"`
	Pincode    string    `gorm:"not null" json:"pincode"`
	Aadhaar    string    `gorm:"not null" json:"aadhaar"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func GetAllUsers() ([]User, error) {
	var users []User

	err := config.DB.Find(&users).Error
	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func GetUserByFirebaseId(firebaseID string) (User, error) {
	var user User

	err := config.DB.Where("firebase_id = ?", firebaseID).First(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateUser(user User) (User, error) {
	err := config.DB.Create(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UpdateUser(user User) (User, error) {
	err := config.DB.Save(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func DeleteUser(userID string) error {
	err := config.DB.Delete(&User{}, userID).Error
	if err != nil {
		return err
	}

	return nil
}
