package handler

// CreateCatRequest defines the request body for creating a cat.
type CreateCatRequest struct {
	Name              string  `json:"name" binding:"required"`
	YearsOfExperience int     `json:"years_of_experience" binding:"required,gte=0"`
	Breed             string  `json:"breed" binding:"required"`
	Salary            float64 `json:"salary" binding:"required,gt=0"`
}

// UpdateCatSalaryRequest defines the request body for updating a cat's salary.
type UpdateCatSalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,gt=0"`
}

// CreateMissionRequest represents the request to create a new mission.
type CreateMissionRequest struct {
	CatID   *int                  `json:"cat_id,omitempty"`
	Targets []CreateTargetRequest `json:"targets" binding:"required,min=1,max=3,dive"`
}

// CreateTargetRequest defines the structure for a target within a mission creation request.
type CreateTargetRequest struct {
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
	Notes   string `json:"notes"`
}

// AssignCatRequest defines the request body for assigning a cat to a mission.
type AssignCatRequest struct {
	CatID int `json:"cat_id" binding:"required"`
}

// UpdateMissionRequest defines the request body for updating a mission.
type UpdateMissionRequest struct {
	CatID *int `json:"cat_id,omitempty"`
}

// CompleteMissionRequest defines the request body for marking a mission as completed.
type CompleteMissionRequest struct {
	Completed bool `json:"completed" binding:"required"`
}

// UpdateTargetNotesRequest defines the request body for updating a target's notes.
type UpdateTargetNotesRequest struct {
	Notes string `json:"notes" binding:"required"`
}
