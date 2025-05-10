package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name            string        `json:"name" gorm:"not null"`
	Date            time.Time     `json:"date" gorm:"not null"`
	Players         []Player      `json:"players" gorm:"many2many:event_players;"`
	Machines        []Machine     `json:"machines" gorm:"many2many:event_machines;"`
	SeasonID        uint          `json:"seasonID" gorm:"not null"`
	Season          Season        `json:"season" gorm:"foreignKey:SeasonID"`
	IsFinals        bool          `json:"isFinals" gorm:"not null"`
	IsComplete      bool          `json:"isComplete" gorm:"not null"`
	HasWinnersGroup bool          `json:"hasWinnersGroup" gorm:"not null"`
	CompletedAt     time.Time     `json:"completedAt" gorm:""`
	SeedingMethod   SeedingMethod `json:"seedingMethod" gorm:"type:string;default:'AVERAGE'"`
	GroupOrdering   GroupOrdering `json:"groupOrdering" gorm:"type:string;default:'SEEDED'"`
}
