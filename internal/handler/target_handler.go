package handler

import (
	"net/http"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TargetHandler handles the HTTP requests for targets.
type TargetHandler struct {
	targetService service.TargetService
}

// NewTargetHandler creates a new TargetHandler.
func NewTargetHandler(targetService service.TargetService) *TargetHandler {
	return &TargetHandler{targetService: targetService}
}

// AddTargetToMission handles adding a target to a mission.
// @Summary Add a target to a mission
// @Description Adds a new target to an existing, non-completed mission.
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param target body CreateTargetRequest true "Target to add"
// @Success 201 {object} domain.Target
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /missions/{id}/targets [post]
func (h *TargetHandler) AddTargetToMission(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Invalid mission ID format", err))
		return
	}

	var req CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	target := &domain.Target{Name: req.Name, Country: req.Country, Notes: req.Notes}
	if err := h.targetService.AddTargetToMission(c.Request.Context(), missionID, target); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, err.Error(), err))
		return
	}

	c.JSON(http.StatusCreated, target)
}

// UpdateTargetNotes handles updating a target's notes.
// @Summary Update target notes
// @Description Updates the notes for a specific target if it is not completed.
// @Tags targets
// @Accept json
// @Produce json
// @Param id path int true "Target ID"
// @Param notes body UpdateTargetNotesRequest true "New notes"
// @Success 200 {object} domain.Target
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /targets/{id}/notes [patch]
func (h *TargetHandler) UpdateTargetNotes(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "invalid id format", err))
		return
	}

	var req UpdateTargetNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	target, err := h.targetService.UpdateTargetNotes(c.Request.Context(), targetID, req.Notes)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, target)
}

// CompleteTarget handles marking a target as complete.
// @Summary Complete a target
// @Description Marks a target as complete. If all targets in the mission are complete, the mission is also marked as complete.
// @Tags targets
// @Param id path int true "Target ID"
// @Success 200 {object} domain.Target
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /targets/{id}/complete [patch]
func (h *TargetHandler) CompleteTarget(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Invalid target ID format", err))
		return
	}

	target, err := h.targetService.CompleteTarget(c.Request.Context(), targetID)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, target)
}

// DeleteTarget handles deleting a target from a mission.
// @Summary Delete a target
// @Description Deletes a target from a mission if it is not yet completed.
// @Tags targets
// @Param id path int true "Target ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /targets/{id} [delete]
func (h *TargetHandler) DeleteTarget(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Invalid target ID format", err))
		return
	}

	if err := h.targetService.DeleteTarget(c.Request.Context(), targetID); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, err.Error(), err))
		return
	}

	c.Status(http.StatusNoContent)
}
