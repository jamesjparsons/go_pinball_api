package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

type LeagueHandler struct {
	db *gorm.DB
}

func NewLeagueHandler(db *gorm.DB) *LeagueHandler {
	return &LeagueHandler{db: db}
}

// CreateLeague handles league creation
func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("CreateLeague error - No user ID in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		Location string `json:"location" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateLeague error - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	league := models.League{
		Name:     req.Name,
		Location: req.Location,
		OwnerID:  userID.(uint),
	}

	if err := h.db.Create(&league).Error; err != nil {
		log.Printf("CreateLeague error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	// Load the owner information
	if err := h.db.Preload("Owner").First(&league, league.ID).Error; err != nil {
		log.Printf("CreateLeague error - Failed to load owner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	log.Printf("CreateLeague success - League created: %s by user %d", league.Name, userID)
	c.JSON(http.StatusCreated, gin.H{
		"data": league,
	})
}

// ListLeagues handles listing all leagues
func (h *LeagueHandler) ListLeagues(c *gin.Context) {
	var leagues []models.League
	if err := h.db.Preload("Owner").Find(&leagues).Error; err != nil {
		log.Printf("ListLeagues error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leagues"})
		return
	}

	log.Printf("ListLeagues success - Retrieved %d leagues", len(leagues))
	c.JSON(http.StatusOK, gin.H{
		"data": leagues,
	})
}

// GetLeague handles getting a single league by ID
func (h *LeagueHandler) GetLeague(c *gin.Context) {
	leagueID := c.Param("leagueID")
	if leagueID == "" {
		log.Printf("GetLeague error - No league ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No league ID provided"})
		return
	}

	var league models.League
	if err := h.db.Preload("Owner").First(&league, leagueID).Error; err != nil {
		log.Printf("GetLeague error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch league"})
		return
	}

	log.Printf("GetLeague success - Retrieved league: %s", league.Name)
	c.JSON(http.StatusOK, gin.H{
		"data": league,
	})
}
