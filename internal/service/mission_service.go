package service

import (
	"context"
	"fmt"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/repository"
)

// missionService is the implementation of the MissionService interface.
type missionService struct {
	missionRepo repository.MissionRepository
	catRepo     repository.CatRepository
}

// NewMissionService creates a new MissionService.
func NewMissionService(missionRepo repository.MissionRepository, catRepo repository.CatRepository) MissionService {
	return &missionService{
		missionRepo: missionRepo,
		catRepo:     catRepo,
	}
}

// CreateMission creates a new mission, ensuring it has between 1 and 3 targets.
func (s *missionService) CreateMission(ctx context.Context, mission *domain.Mission) error {
	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		return fmt.Errorf("a mission must have between 1 and 3 targets")
	}

	// If a cat ID is provided, validate that the cat exists and is available
	if mission.CatID != nil {
		cat, err := s.catRepo.GetCatByID(ctx, *mission.CatID)
		if err != nil {
			return fmt.Errorf("cat not found: %w", err)
		}
		if cat.Status != "available" {
			return fmt.Errorf("cat is not available for a mission")
		}
	}

	return s.missionRepo.CreateMission(ctx, mission)
}

// GetMission retrieves a mission by its ID.
func (s *missionService) GetMission(ctx context.Context, id int) (*domain.Mission, error) {
	return s.missionRepo.GetMissionByID(ctx, id)
}

// ListMissions retrieves all missions.
func (s *missionService) ListMissions(ctx context.Context) ([]domain.Mission, error) {
	return s.missionRepo.ListMissions(ctx)
}

// UpdateMission is a placeholder for more complex mission update logic if needed.
func (s *missionService) UpdateMission(ctx context.Context, mission *domain.Mission) error {
	return s.missionRepo.UpdateMission(ctx, mission)
}

// DeleteMission deletes a mission if it's not assigned to a cat.
func (s *missionService) DeleteMission(ctx context.Context, id int) error {
	mission, err := s.missionRepo.GetMissionByID(ctx, id)
	if err != nil {
		return err
	}
	if mission.CatID != nil {
		return fmt.Errorf("cannot delete a mission that is assigned to a cat")
	}
	return s.missionRepo.DeleteMission(ctx, id)
}

// AssignCatToMission assigns an available cat to a mission.
func (s *missionService) AssignCatToMission(ctx context.Context, missionID, catID int) error {
	cat, err := s.catRepo.GetCatByID(ctx, catID)
	if err != nil {
		return err
	}
	if cat.Status != "available" {
		return fmt.Errorf("cat is not available for a mission")
	}
	return s.missionRepo.AssignCatToMission(ctx, missionID, catID)
}





// CompleteMission manually marks a mission as completed or uncompleted.
func (s *missionService) CompleteMission(ctx context.Context, missionID int, completed bool) (*domain.Mission, error) {
	// Get the current mission
	mission, err := s.missionRepo.GetMissionByID(ctx, missionID)
	if err != nil {
		return nil, err
	}

	// Update the completion status
	mission.Completed = completed
	if err := s.missionRepo.UpdateMission(ctx, mission); err != nil {
		return nil, err
	}

	return mission, nil
}
