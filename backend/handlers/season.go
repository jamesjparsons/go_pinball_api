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
// @Summary Create a new season
// @Description Create a new season for a specific league
// @Tags seasons
// @Accept json
// @Produce json
// @Security Bearer
// @Param leagueID path string true "League ID"
// @Param request body object true "Season details" {"name": "string"}
// @Success 201 {object} ListResponse{data=SeasonResponse} "Season created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/seasons [post]
func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	leagueID := c.Param("leagueID")
	if leagueID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "League ID is required"})
		return
	}

	// Convert leagueID to uint
	leagueIDUint, err := strconv.ParseUint(leagueID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid league ID"})
		return
	}

	var req models.CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	season := models.Season{
		Name:              req.Name,
		DateCreated:       time.Now(),
		LeagueID:          uint(leagueIDUint),
		CountingGames:     req.CountingGames,
		EventCount:        0,
		HasFinals:         req.HasFinals,
		PointDistribution: make(map[string][]int),
	}

	if err := h.db.Create(&season).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create season"})
		return
	}

	c.JSON(http.StatusCreated, season)
}

// ListSeasons handles listing all seasons for a league
// @Summary List seasons for a league
// @Description Get a list of all seasons for a specific league
// @Tags seasons
// @Produce json
// @Param leagueID path string true "League ID"
// @Success 200 {object} ListResponse{data=[]SeasonResponse} "List of seasons"
// @Failure 400 {object} ErrorResponse "Invalid league ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/seasons [get]
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
