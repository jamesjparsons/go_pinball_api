package models

import (
	"time"

	"gorm.io/gorm"
)

type Machine struct {
	gorm.Model `swaggerignore:"true"`
	OPDBID     string `json:"opdb_id" gorm:"uniqueIndex;not null"`
	Name       string `json:"name"`
	// Manufacturer  string    `json:"manufacturer"`
	Year      int    `json:"year"`
	IPDBID    int    `json:"ipdb_id"`
	Type      string `json:"type"`
	IsPinball bool   `json:"is_pinball"`
	IsGroup   bool   `json:"is_group"`
	IsAlias   bool   `json:"is_alias"`
	// LastUpdatedAt time.Time `json:"last_updated_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
