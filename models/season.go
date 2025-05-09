package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	gorm.Model
	Name      string    `json:"name"`
	LeagueID  uint      `json:"leagueId"`
	League    League    `json:"league" gorm:"foreignKey:LeagueID"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
