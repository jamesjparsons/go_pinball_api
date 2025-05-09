package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	Name              string               `json:"name" gorm:"not null"`
	DateCreated       time.Time            `json:"dateCreated" gorm:"not null"`
	LeagueID          uint                 `json:"leagueID" gorm:"not null"`
	League            League               `json:"league" gorm:"foreignKey:LeagueID"`
	CountingGames     uint                 `json:"countingGames" gorm:"not null"`
	EventCount        uint                 `json:"eventCount" gorm:"not null"`
	HasFinals         bool                 `json:"hasFinals" gorm:"not null"`
	PointDistribution PointDistributionMap `json:"pointDistribution" gorm:"type:json"`
}
