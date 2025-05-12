package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
	"backend/services"
)

type LeagueHandler struct {
	db          *gorm.DB
	ifpaService *services.IFPAService
}

func NewLeagueHandler(db *gorm.DB, ifpaService *services.IFPAService) *LeagueHandler {
	return &LeagueHandler{
		db:          db,
		ifpaService: ifpaService,
	}
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

// ListPlayers handles listing all players in a league
// @Summary List players in a league
// @Description Get a list of all players in a specific league
// @Tags leagues
// @Produce json
// @Param leagueID path string true "League ID"
// @Success 200 {object} ListResponse{data=[]PlayerResponse} "List of players"
// @Failure 400 {object} ErrorResponse "Invalid league ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/players [get]
func (h *LeagueHandler) ListPlayers(c *gin.Context) {
	leagueID := c.Param("leagueID")
	if leagueID == "" {
		log.Printf("ListPlayers error - No league ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "League ID is required"})
		return
	}

	var players []models.Player
	if err := h.db.Where("league_id = ?", leagueID).Find(&players).Error; err != nil {
		log.Printf("ListPlayers error - Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	log.Printf("ListPlayers success - Retrieved %d players for league %s", len(players), leagueID)
	c.JSON(http.StatusOK, gin.H{
		"data": players,
	})
}

// AddPlayersByIFPA handles adding players to a league by their IFPA numbers
// @Summary Add players to league by IFPA numbers
// @Description Add multiple players to a league using their IFPA numbers
// @Tags leagues
// @Accept json
// @Produce json
// @Security Bearer
// @Param leagueID path string true "League ID"
// @Param request body object true "IFPA numbers" {"ifpaNumbers": [1234, 5678]}
// @Success 200 {object} ListResponse{data=[]PlayerResponse} "Players added successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or league ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/players/ifpa [post]
func (h *LeagueHandler) AddPlayersByIFPA(c *gin.Context) {
	leagueID := c.Param("leagueID")
	if leagueID == "" {
		log.Printf("AddPlayersByIFPA error - No league ID provided")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "League ID is required"})
		return
	}

	// Convert leagueID to uint
	leagueIDUint, err := strconv.ParseUint(leagueID, 10, 32)
	if err != nil {
		log.Printf("AddPlayersByIFPA error - Invalid league ID: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid league ID"})
		return
	}

	var req struct {
		IFPANumbers []int `json:"ifpaNumbers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("AddPlayersByIFPA error - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	var addedPlayers []models.Player
	for _, ifpaNumber := range req.IFPANumbers {
		// Check if player already exists in the league
		var existingPlayer models.Player
		result := h.db.Where("league_id = ? AND ifpa_number = ?", leagueIDUint, strconv.Itoa(ifpaNumber)).First(&existingPlayer)
		if result.Error == nil {
			// Player already exists in the league
			addedPlayers = append(addedPlayers, existingPlayer)
			continue
		}

		// Get player details from IFPA API
		ifpaPlayer, err := h.ifpaService.GetPlayerByIFPANumber(ifpaNumber)
		if err != nil {
			log.Printf("AddPlayersByIFPA error - Failed to get IFPA player %d: %v", ifpaNumber, err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to get IFPA player %d", ifpaNumber)})
			return
		}

		// Create new player
		player := models.Player{
			LeagueID:   uint(leagueIDUint),
			IFPANumber: strconv.Itoa(ifpaNumber),
			Name:       fmt.Sprintf("%s %s", ifpaPlayer.FirstName, ifpaPlayer.LastName),
		}

		if err := h.db.Create(&player).Error; err != nil {
			log.Printf("AddPlayersByIFPA error - Failed to create player: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create player"})
			return
		}

		addedPlayers = append(addedPlayers, player)
	}

	log.Printf("AddPlayersByIFPA success - Added %d players to league %d", len(addedPlayers), leagueIDUint)
	c.JSON(http.StatusOK, gin.H{
		"data": addedPlayers,
	})
}
