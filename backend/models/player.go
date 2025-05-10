package models

import (
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Name       string `json:"name" gorm:"not null"`
	LeagueID   uint   `json:"leagueID" gorm:"not null"`
	League     League `json:"league" gorm:"foreignKey:LeagueID"`
	IFPANumber string `json:"ifpaNumber" gorm:""`
}
