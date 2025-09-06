package service

import (
	"context"
	"fmt"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/repository"
)

// targetService is the implementation of the TargetService interface.
type targetService struct {
	targetRepo  repository.TargetRepository
	missionRepo repository.MissionRepository
}

// NewTargetService creates a new TargetService.
func NewTargetService(targetRepo repository.TargetRepository, missionRepo repository.MissionRepository) TargetService {
	return &targetService{
		targetRepo:  targetRepo,
		missionRepo: missionRepo,
	}
}

// AddTargetToMission adds a target to an existing, non-completed mission.
func (s *targetService) AddTargetToMission(ctx context.Context, missionID int, target *domain.Target) error {
	mission, err := s.missionRepo.GetMissionByID(ctx, missionID)
	if err != nil {
		return err
	}
	if mission.Completed {
		return fmt.Errorf("cannot add a target to a completed mission")
	}
	if len(mission.Targets) >= 3 {
		return fmt.Errorf("a mission cannot have more than 3 targets")
	}
	target.MissionID = missionID
	return s.targetRepo.AddTargetToMission(ctx, target)
}

// UpdateTargetNotes updates the notes of a target if it and its mission are not complete.
func (s *targetService) UpdateTargetNotes(ctx context.Context, targetID int, notes string) (*domain.Target, error) {
	target, err := s.targetRepo.GetTargetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if target.Completed {
		return nil, fmt.Errorf("cannot update notes on a completed target")
	}

	mission, err := s.missionRepo.GetMissionByID(ctx, target.MissionID)
	if err != nil {
		return nil, err
	}
	if mission.Completed {
		return nil, fmt.Errorf("cannot update notes on a target in a completed mission")
	}

	target.Notes = notes
	if err := s.targetRepo.UpdateTarget(ctx, target); err != nil {
		return nil, err
	}
	return target, nil
}

// CompleteTarget marks a target as complete and checks if the entire mission is now complete.
func (s *targetService) CompleteTarget(ctx context.Context, targetID int) (*domain.Target, error) {
	target, err := s.targetRepo.GetTargetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	target.Completed = true
	if err := s.targetRepo.UpdateTarget(ctx, target); err != nil {
		return nil, err
	}

	// Check if all targets in the mission are now complete
	mission, err := s.missionRepo.GetMissionByID(ctx, target.MissionID)
	if err != nil {
		return nil, err // The target was updated, but we failed to check the mission status
	}

	allTargetsComplete := true
	for _, t := range mission.Targets {
		if !t.Completed {
			allTargetsComplete = false
			break
		}
	}

	if allTargetsComplete {
		mission.Completed = true
		if err := s.missionRepo.UpdateMission(ctx, mission); err != nil {
			return target, fmt.Errorf("failed to mark mission as complete: %w", err)
		}
	}

	return target, nil
}

// DeleteTarget deletes a target if it is not yet completed.
func (s *targetService) DeleteTarget(ctx context.Context, targetID int) error {
	target, err := s.targetRepo.GetTargetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target.Completed {
		return fmt.Errorf("cannot delete a completed target")
	}
	return s.targetRepo.DeleteTarget(ctx, targetID)
}
