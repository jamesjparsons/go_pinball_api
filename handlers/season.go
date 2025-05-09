package handlers

import (
	"net/http"

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

// @Summary     List seasons for a league
// @Description Get all seasons for a specific league
// @Tags        seasons
// @Produce     json
// @Security    Bearer
// @Param       leagueID path int true "League ID"
// @Success     200 {array} models.Season
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /leagues/{leagueID}/seasons [get]
func (h *SeasonHandler) ListSeasons(c *gin.Context) {
	leagueID := c.Param("leagueID")

	var seasons []models.Season
	if err := h.db.Where("league_id = ?", leagueID).Find(&seasons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seasons"})
		return
	}

	c.JSON(http.StatusOK, seasons)
}

// @Summary     Create new season
// @Description Create a new season for a league
// @Tags        seasons
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       leagueID path int true "League ID"
// @Param       season body CreateSeasonRequest true "Season data"
// @Success     200 {object} models.Season
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /leagues/{leagueID}/seasons [post]
func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	leagueID := c.Param("leagueID")

	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	season := models.Season{
		Name:     req.Name,
		LeagueID: uint(leagueID),
	}

	if err := h.db.Create(&season).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create season"})
		return
	}

	c.JSON(http.StatusOK, season)
}

type CreateSeasonRequest struct {
	Name string `json:"name" binding:"required"`
}
