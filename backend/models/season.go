package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model        `swaggerignore:"true"`
	Name              string               `json:"name"`
	DateCreated       time.Time            `json:"dateCreated" gorm:"not null"`
	LeagueID          uint                 `json:"leagueID"`
	League            League               `json:"league"`
	CountingGames     int                  `json:"countingGames"`
	EventCount        int                  `json:"eventCount"`
	HasFinals         bool                 `json:"hasFinals"`
	PointDistribution PointDistributionMap `json:"pointDistribution" gorm:"type:json"`
	CreatedAt         time.Time            `json:"created_at"`
}

type CreateSeasonRequest struct {
	Name          string `json:"name" binding:"required"`
	CountingGames int    `json:"countingGames" binding:"required"`
	HasFinals     bool   `json:"hasFinals"`
}
