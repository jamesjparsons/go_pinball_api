package models

import (
	"time"
)

type Season struct {
	ID                uint             `json:"id" gorm:"primaryKey"`
	Name              string           `json:"name"`
	DateCreated       time.Time        `json:"dateCreated" gorm:"not null"`
	LeagueID          uint             `json:"leagueID"`
	League            League           `json:"league"`
	CountingGames     int              `json:"countingGames"`
	EventCount        int              `json:"eventCount"`
	HasFinals         bool             `json:"hasFinals"`
	PointDistribution map[string][]int `json:"pointDistribution" gorm:"type:json"`
	CreatedAt         time.Time        `json:"created_at"`
}

type CreateSeasonRequest struct {
	Name          string `json:"name" binding:"required"`
	CountingGames int    `json:"countingGames" binding:"required"`
	HasFinals     bool   `json:"hasFinals"`
}
