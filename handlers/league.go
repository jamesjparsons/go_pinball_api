package handlers

import (
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

// @Summary     List all leagues
// @Description Get a list of all leagues
// @Tags        leagues
// @Produce     json
// @Security    Bearer
// @Success     200 {array} models.League
// @Failure     401 {object} ErrorResponse
// @Router      /leagues [get]
func (h *LeagueHandler) ListLeagues(c *gin.Context) {
	var leagues []models.League
	if err := h.db.Preload("Owner").Find(&leagues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leagues"})
		return
	}

	c.JSON(http.StatusOK, leagues)
}

// @Summary     Get league details
// @Description Get details of a specific league
// @Tags        leagues
// @Produce     json
// @Security    Bearer
// @Param       leagueID path int true "League ID"
// @Success     200 {object} models.League
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /leagues/{leagueID} [get]
func (h *LeagueHandler) GetLeague(c *gin.Context) {
	leagueID := c.Param("leagueID")

	var league models.League
	if err := h.db.Preload("Owner").First(&league, leagueID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found"})
		return
	}

	c.JSON(http.StatusOK, league)
}

// @Summary     Create new league
// @Description Create a new league
// @Tags        leagues
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       league body CreateLeagueRequest true "League data"
// @Success     200 {object} models.League
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /leagues/create [post]
func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	var req CreateLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	league := models.League{
		Name:     req.Name,
		Location: req.Location,
		OwnerID:  userID.(uint),
	}

	if err := h.db.Create(&league).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	// Load the owner information
	if err := h.db.Preload("Owner").First(&league, league.Model.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load league details"})
		return
	}

	c.JSON(http.StatusOK, league)
}

type CreateLeagueRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}
