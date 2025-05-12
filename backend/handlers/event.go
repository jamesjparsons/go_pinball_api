package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

type EventHandler struct {
	db *gorm.DB
}

func NewEventHandler(db *gorm.DB) *EventHandler {
	return &EventHandler{db: db}
}

// CreateEvent handles event creation
// @Summary Create a new event
// @Description Create a new event for a specific season
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param leagueID path string true "League ID"
// @Param seasonID path string true "Season ID"
// @Param request body object true "Event details" {"name": "string", "date": "string", "isFinals": "boolean", "hasWinnersGroup": "boolean", "seedingMethod": "string", "groupOrdering": "string"}
// @Success 201 {object} ListResponse{data=EventResponse} "Event created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/seasons/{seasonID}/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	leagueID := c.Param("leagueID")
	seasonID := c.Param("seasonID")

	// Convert IDs to uint
	_, err := strconv.ParseUint(leagueID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid league ID"})
		return
	}

	seasonIDUint, err := strconv.ParseUint(seasonID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid season ID"})
		return
	}

	var req struct {
		Name            string `json:"name" binding:"required"`
		Date            string `json:"date" binding:"required"`
		IsFinals        bool   `json:"isFinals"`
		HasWinnersGroup bool   `json:"hasWinnersGroup"`
		SeedingMethod   string `json:"seedingMethod"`
		GroupOrdering   string `json:"groupOrdering"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid date format. Use RFC3339 format"})
		return
	}

	event := models.Event{
		Name:            req.Name,
		Date:            date,
		SeasonID:        uint(seasonIDUint),
		IsFinals:        req.IsFinals,
		IsComplete:      false,
		HasWinnersGroup: req.HasWinnersGroup,
		SeedingMethod:   models.SeedingMethod(req.SeedingMethod),
		GroupOrdering:   models.GroupOrdering(req.GroupOrdering),
	}

	if err := h.db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": event,
	})
}

// ListEvents handles listing all events for a season
// @Summary List events for a season
// @Description Get a list of all events for a specific season
// @Tags events
// @Produce json
// @Param leagueID path string true "League ID"
// @Param seasonID path string true "Season ID"
// @Success 200 {object} ListResponse{data=[]EventResponse} "List of events"
// @Failure 400 {object} ErrorResponse "Invalid league ID or season ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/seasons/{seasonID}/events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	leagueID := c.Param("leagueID")
	seasonID := c.Param("seasonID")

	// Convert IDs to uint
	_, err := strconv.ParseUint(leagueID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid league ID"})
		return
	}

	seasonIDUint, err := strconv.ParseUint(seasonID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid season ID"})
		return
	}

	var events []models.Event
	if err := h.db.Where("season_id = ?", seasonIDUint).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": events,
	})
}

// GetEvent handles getting a single event by ID
// @Summary Get event by ID
// @Description Get detailed information about a specific event
// @Tags events
// @Produce json
// @Param leagueID path string true "League ID"
// @Param seasonID path string true "Season ID"
// @Param eventID path string true "Event ID"
// @Success 200 {object} ListResponse{data=EventResponse} "Event details"
// @Failure 400 {object} ErrorResponse "Invalid league ID, season ID, or event ID"
// @Failure 404 {object} ErrorResponse "Event not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /leagues/{leagueID}/seasons/{seasonID}/events/{eventID} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	leagueID := c.Param("leagueID")
	seasonID := c.Param("seasonID")
	eventID := c.Param("eventID")

	// Convert IDs to uint
	_, err := strconv.ParseUint(leagueID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid league ID"})
		return
	}

	seasonIDUint, err := strconv.ParseUint(seasonID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid season ID"})
		return
	}

	eventIDUint, err := strconv.ParseUint(eventID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid event ID"})
		return
	}

	var event models.Event
	if err := h.db.Where("id = ? AND season_id = ?", eventIDUint, seasonIDUint).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch event"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": event,
	})
}
