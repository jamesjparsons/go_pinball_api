package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

type SeasonHandler struct {
	db *gorm.DB
}

func NewSeasonHandler(db *gorm.DB) *SeasonHandler {
	return &SeasonHandler{db: db}
}

// CreateSeason handles season creation
func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	// Check if user is authenticated
	if _, exists := c.Get("userID"); !exists {
		log.Printf("CreateSeason error - No user ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get league ID from URL parameter
	leagueIDStr := c.Param("leagueID")
	if leagueIDStr == "" {
		log.Printf("CreateSeason error - No league ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "League ID is required"})
		return
	}

	// Convert leagueID to uint
	leagueID, err := strconv.ParseUint(leagueIDStr, 10, 32)
	if err != nil {
		log.Printf("CreateSeason error - Invalid league ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateSeason error - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create the season
	season := models.Season{
		Name:        req.Name,
		LeagueID:    uint(leagueID),
		DateCreated: time.Now(),
		// Set default values
		CountingGames: 5,
		EventCount:    0,
		HasFinals:     false,
		PointDistribution: models.PointDistributionMap{
			"4": []float64{4, 3, 2, 1},
			"3": []float64{4, 2.5, 1},
		},
	}

	if err := h.db.Create(&season).Error; err != nil {
		log.Printf("CreateSeason error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create season"})
		return
	}

	// Load the league information
	if err := h.db.Preload("League").First(&season, season.ID).Error; err != nil {
		log.Printf("CreateSeason error - Failed to load league: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create season"})
		return
	}

	log.Printf("CreateSeason success - Season created: %s in league %d", season.Name, season.LeagueID)
	c.JSON(http.StatusCreated, gin.H{
		"data": season,
	})
}

// ListSeasons handles listing all seasons for a league
func (h *SeasonHandler) ListSeasons(c *gin.Context) {
	leagueIDStr := c.Param("leagueID")
	if leagueIDStr == "" {
		log.Printf("ListSeasons error - No league ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "League ID is required"})
		return
	}

	// Convert leagueID to uint
	leagueID, err := strconv.ParseUint(leagueIDStr, 10, 32)
	if err != nil {
		log.Printf("ListSeasons error - Invalid league ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	var seasons []models.Season
	if err := h.db.Where("league_id = ?", leagueID).Preload("League").Find(&seasons).Error; err != nil {
		log.Printf("ListSeasons error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seasons"})
		return
	}

	log.Printf("ListSeasons success - Retrieved %d seasons for league %d", len(seasons), leagueID)
	c.JSON(http.StatusOK, gin.H{
		"data": seasons,
	})
}
