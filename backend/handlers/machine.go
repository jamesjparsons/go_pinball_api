package handlers

import (
	"net/http"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type MachineHandler struct {
	opdbService *services.OPDBService
}

func NewMachineHandler(opdbService *services.OPDBService) *MachineHandler {
	return &MachineHandler{
		opdbService: opdbService,
	}
}

// GetMachine godoc
// @Summary Get machine details from OPDB
// @Description Get machine details from OPDB API and cache in database
// @Tags machines
// @Accept json
// @Produce json
// @Param opdb_id path string true "OPDB ID"
// @Success 200 {object} models.Machine
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /machines/{opdb_id} [get]
func (h *MachineHandler) GetMachine(c *gin.Context) {
	opdbID := c.Param("opdb_id")
	if opdbID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "OPDB ID is required"})
		return
	}

	machine, err := h.opdbService.GetMachine(opdbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, machine)
}
