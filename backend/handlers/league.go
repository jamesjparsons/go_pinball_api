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
// @Summary Create a new league
// @Description Create a new pinball league
// @Tags leagues
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body object true "League details" {"name": "string", "location": "string"}
// @Success 201 {object} ListResponse{data=LeagueResponse} "League created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/create [post]
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
// @Summary List all leagues
// @Description Get a list of all pinball leagues
// @Tags leagues
// @Produce json
// @Success 200 {object} ListResponse{data=[]LeagueResponse} "List of leagues"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues [get]
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
// @Summary Get league by ID
// @Description Get detailed information about a specific league
// @Tags leagues
// @Produce json
// @Param leagueID path string true "League ID"
// @Success 200 {object} ListResponse{data=LeagueResponse} "League details"
// @Failure 400 {object} ErrorResponse "Invalid league ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID} [get]
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
