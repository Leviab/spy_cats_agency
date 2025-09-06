package handler

import (
	"net/http"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MissionHandler handles the HTTP requests for missions.
type MissionHandler struct {
	missionService service.MissionService
}

// NewMissionHandler creates a new MissionHandler.
func NewMissionHandler(missionService service.MissionService) *MissionHandler {
	return &MissionHandler{missionService: missionService}
}

// CreateMission handles the creation of a new mission.
// @Summary Create a new mission
// @Description Creates a new mission with 1 to 3 targets. Optionally assign a cat during creation.
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body CreateMissionRequest true "Mission to create"
// @Success 201 {object} domain.Mission
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /missions [post]
func (h *MissionHandler) CreateMission(c *gin.Context) {
	var req CreateMissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	mission := &domain.Mission{
		CatID: req.CatID,
	}
	for _, t := range req.Targets {
		mission.Targets = append(mission.Targets, domain.Target{Name: t.Name, Country: t.Country, Notes: t.Notes})
	}

	if err := h.missionService.CreateMission(c.Request.Context(), mission); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to create mission", err))
		return
	}

	c.JSON(http.StatusCreated, mission)
}

// GetMission handles retrieving a single mission by its ID.
// @Summary Get a mission by ID
// @Description Retrieves details of a specific mission, including its targets.
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} domain.Mission
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /missions/{id} [get]
func (h *MissionHandler) GetMission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "invalid id format", err))
		return
	}

	mission, err := h.missionService.GetMission(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to get mission", err))
		return
	}

	c.JSON(http.StatusOK, mission)
}

// ListMissions handles listing all missions.
// @Summary List all missions
// @Description Retrieves a list of all missions.
// @Tags missions
// @Produce json
// @Success 200 {array} domain.Mission
// @Failure 500 {object} ErrorResponse
// @Router /missions [get]
func (h *MissionHandler) ListMissions(c *gin.Context) {
	missions, err := h.missionService.ListMissions(c.Request.Context())
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to list missions", err))
		return
	}

	c.JSON(http.StatusOK, missions)
}

// DeleteMission handles deleting a mission.
// @Summary Delete a mission
// @Description Deletes a mission if it is not assigned to a cat.
// @Tags missions
// @Param id path int true "Mission ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /missions/{id} [delete]
func (h *MissionHandler) DeleteMission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "invalid id format", err))
		return
	}

	if err := h.missionService.DeleteMission(c.Request.Context(), id); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to delete mission", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignCatToMission handles assigning a cat to a mission.
// @Summary Assign a cat to a mission
// @Description Assigns an available spy cat to an existing mission.
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param cat body AssignCatRequest true "Cat to assign"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /missions/{id}/assign-cat [patch]
func (h *MissionHandler) AssignCatToMission(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "invalid id format", err))
		return
	}

	var req AssignCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	if err := h.missionService.AssignCatToMission(c.Request.Context(), missionID, req.CatID); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to assign cat to mission", err))
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Cat assigned successfully"})
}

// CompleteMission handles marking a mission as completed or uncompleted.
// @Summary Complete/uncomplete a mission
// @Description Manually marks a mission as completed or uncompleted.
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param completion body CompleteMissionRequest true "Completion status"
// @Success 200 {object} domain.Mission
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /missions/{id}/complete [patch]
func (h *MissionHandler) CompleteMission(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Invalid mission ID format", err))
		return
	}

	var req CompleteMissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, err.Error(), err))
		return
	}

	mission, err := h.missionService.CompleteMission(c.Request.Context(), missionID, req.Completed)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to complete mission", err))
		return
	}

	c.JSON(http.StatusOK, mission)
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
func (h *MissionHandler) AddTargetToMission(c *gin.Context) {
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
	if err := h.missionService.AddTargetToMission(c.Request.Context(), missionID, target); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to add target to mission", err))
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
func (h *MissionHandler) UpdateTargetNotes(c *gin.Context) {
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

	target, err := h.missionService.UpdateTargetNotes(c.Request.Context(), targetID, req.Notes)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to update target notes", err))
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
func (h *MissionHandler) CompleteTarget(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Invalid target ID format", err))
		return
	}

	target, err := h.missionService.CompleteTarget(c.Request.Context(), targetID)
	if err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to complete target", err))
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
func (h *MissionHandler) DeleteTarget(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(NewAppError(http.StatusBadRequest, "Invalid target ID format", err))
		return
	}

	if err := h.missionService.DeleteTarget(c.Request.Context(), targetID); err != nil {
		_ = c.Error(NewAppError(http.StatusInternalServerError, "Failed to delete target", err))
		return
	}

	c.Status(http.StatusNoContent)
}
