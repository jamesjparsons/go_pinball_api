package handlers

import "backend/models"

// UserResponse represents the user data in API responses
type UserResponse struct {
	models.SwaggerUser
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// LeagueResponse represents the league data in API responses
type LeagueResponse struct {
	models.SwaggerLeague
}

// SeasonResponse represents the season data in API responses
type SeasonResponse struct {
	models.SwaggerSeason
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request body"`
}

// ListResponse represents a list response with data
type ListResponse struct {
	Data interface{} `json:"data"`
}
