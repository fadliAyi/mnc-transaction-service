package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName   string    `gorm:"not null"`
	LastName    string    `gorm:"not null"`
	PhoneNumber string    `gorm:"unique;not null"`
	Address     string    `gorm:"not null"`
	Pin         string    `gorm:"not null"`
	Balance     int       `gorm:"default:0"`
	gorm.Model
}
