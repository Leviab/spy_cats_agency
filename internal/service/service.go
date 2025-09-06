package service

import (
	"context"
	"spy_cats_agency/internal/domain"
)

// CatService defines the interface for cat-related business logic.
type CatService interface {
	CreateCat(ctx context.Context, cat *domain.Cat) error
	GetCat(ctx context.Context, id int) (*domain.Cat, error)
	ListCats(ctx context.Context) ([]domain.Cat, error)
	UpdateCatSalary(ctx context.Context, id int, salary float64) (*domain.Cat, error)
	DeleteCat(ctx context.Context, id int) error
}

// MissionService defines the interface for mission-related business logic.
type MissionService interface {
	CreateMission(ctx context.Context, mission *domain.Mission) error
	GetMission(ctx context.Context, id int) (*domain.Mission, error)
	ListMissions(ctx context.Context) ([]domain.Mission, error)
	UpdateMission(ctx context.Context, mission *domain.Mission) error
	DeleteMission(ctx context.Context, id int) error
	AssignCatToMission(ctx context.Context, missionID, catID int) error

	// CompleteMission manually marks a mission as completed or uncompleted.
	CompleteMission(ctx context.Context, missionID int, completed bool) (*domain.Mission, error)
	AddTargetToMission(ctx context.Context, missionID int, target *domain.Target) error
	UpdateTargetNotes(ctx context.Context, targetID int, notes string) (*domain.Target, error)
	CompleteTarget(ctx context.Context, targetID int) (*domain.Target, error)
	DeleteTarget(ctx context.Context, targetID int) error
}
