package repository

import (
	"context"
	"spy_cats_agency/internal/domain"
)

// CatRepository defines the interface for cat data operations.
type CatRepository interface {
	CreateCat(ctx context.Context, cat *domain.Cat) error
	GetCatByID(ctx context.Context, id int) (*domain.Cat, error)
	ListCats(ctx context.Context) ([]domain.Cat, error)
	UpdateCat(ctx context.Context, cat *domain.Cat) error
	DeleteCat(ctx context.Context, id int) error
}

// MissionRepository defines the interface for mission data operations.
type MissionRepository interface {
	CreateMission(ctx context.Context, mission *domain.Mission) error
	GetMissionByID(ctx context.Context, id int) (*domain.Mission, error)
	ListMissions(ctx context.Context) ([]domain.Mission, error)
	UpdateMission(ctx context.Context, mission *domain.Mission) error
	DeleteMission(ctx context.Context, id int) error
	AssignCatToMission(ctx context.Context, missionID, catID int) error
}

// TargetRepository defines the interface for target data operations.
type TargetRepository interface {
	AddTargetToMission(ctx context.Context, target *domain.Target) error
	GetTargetByID(ctx context.Context, id int) (*domain.Target, error)
	UpdateTarget(ctx context.Context, target *domain.Target) error
	DeleteTarget(ctx context.Context, id int) error
	GetTargetsByMissionID(ctx context.Context, missionID int) ([]domain.Target, error)
}
