package models

import (
	"time"

	"gorm.io/gorm"
)

type League struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Location    string    `json:"location" gorm:"not null"`
	DateCreated time.Time `json:"dateCreated" gorm:"not null"`
	OwnerID     uint      `json:"ownerID" gorm:"not null"`
	Owner       User      `json:"owner" gorm:"foreignKey:OwnerID"`
}
