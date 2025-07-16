package models

import (
	"time"

	"github.com/jezhtech/prince-group-backend/config"
)

type Referral struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReferralID  string    `gorm:"not null;unique" json:"referralId"`
	Name        string    `gorm:"not null" json:"name"`
	SocialMedia string    `gorm:"not null" json:"socialMedia"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func GetReferralByID(id string) (Referral, error) {
	var referral Referral

	err := config.DB.Where("id = ?", id).First(&referral).Error
	if err != nil {
		return Referral{}, err
	}

	return referral, nil
}

func GetReferralByCode(code string) (Referral, error) {
	var referral Referral

	err := config.DB.Where("referral_id = ?", code).First(&referral).Error
	if err != nil {
		return Referral{}, err
	}

	return referral, nil
}

func GetAllReferrals() ([]Referral, error) {
	var referrals []Referral

	err := config.DB.Find(&referrals).Error
	if err != nil {
		return []Referral{}, err
	}
	return referrals, nil
}

func CreateReferral(referral Referral) (Referral, error) {
	err := config.DB.Create(&referral).Error
	if err != nil {
		return Referral{}, err
	}

	return referral, nil
}

func UpdateReferral(referral Referral) (Referral, error) {
	err := config.DB.Save(&referral).Error
	if err != nil {
		return Referral{}, err
	}

	return referral, nil
}

func DeleteReferral(id uint) error {
	err := config.DB.Delete(&Referral{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
