package models

import (
	"time"

	"gorm.io/gorm"
)

type League struct {
	gorm.Model
	Name        string    `gorm:"not null"`
	Location    string    `gorm:"not null"`
	DateCreated time.Time `gorm:"not null"`
	OwnerID     uint      `gorm:"not null"`
	Owner       User      `gorm:"foreignKey:OwnerID"`
} 