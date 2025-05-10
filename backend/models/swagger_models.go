package models

// BaseModel represents the common fields from gorm.Model for Swagger documentation
type BaseModel struct {
	ID        uint   `json:"id" example:"1"`
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}

// SwaggerUser represents the User model for Swagger documentation
type SwaggerUser struct {
	BaseModel
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
}

// SwaggerLeague represents the League model for Swagger documentation
type SwaggerLeague struct {
	BaseModel
	Name        string      `json:"name" example:"Downtown Pinball League"`
	Location    string      `json:"location" example:"123 Main St"`
	DateCreated string      `json:"dateCreated" example:"2024-01-01T00:00:00Z"`
	OwnerID     uint        `json:"ownerID" example:"1"`
	Owner       SwaggerUser `json:"owner"`
}

// SwaggerSeason represents the Season model for Swagger documentation
type SwaggerSeason struct {
	BaseModel
	Name              string               `json:"name" example:"Spring 2024"`
	DateCreated       string               `json:"dateCreated" example:"2024-01-01T00:00:00Z"`
	LeagueID          uint                 `json:"leagueID" example:"1"`
	League            SwaggerLeague        `json:"league"`
	CountingGames     uint                 `json:"countingGames" example:"5"`
	EventCount        uint                 `json:"eventCount" example:"0"`
	HasFinals         bool                 `json:"hasFinals" example:"false"`
	PointDistribution PointDistributionMap `json:"pointDistribution"`
}
