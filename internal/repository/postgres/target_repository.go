package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/repository"
)

// TargetRepository implements the repository.TargetRepository interface.
type TargetRepository struct {
	db *DB
}

// NewTargetRepository creates a new target repository.
func NewTargetRepository(db *DB) repository.TargetRepository {
	return &TargetRepository{db: db}
}

// AddTargetToMission adds a new target to an existing mission.
func (r *TargetRepository) AddTargetToMission(ctx context.Context, target *domain.Target) error {
	query := `INSERT INTO targets (mission_id, name, country, notes, completed)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, target.MissionID, target.Name, target.Country, target.Notes, target.Completed).
		Scan(&target.ID, &target.CreatedAt, &target.UpdatedAt)
}

// GetTargetByID retrieves a single target by its ID.
func (r *TargetRepository) GetTargetByID(ctx context.Context, id int) (*domain.Target, error) {
	var target domain.Target
	query := `SELECT id, mission_id, name, country, notes, completed, created_at, updated_at FROM targets WHERE id = $1`
	err := r.db.GetContext(ctx, &target, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target not found")
		}
		return nil, err
	}
	return &target, nil
}

// UpdateTarget updates a target's notes or completion status.
func (r *TargetRepository) UpdateTarget(ctx context.Context, target *domain.Target) error {
	query := `UPDATE targets SET notes = $1, completed = $2, updated_at = now() WHERE id = $3 RETURNING updated_at`
	return r.db.QueryRowxContext(ctx, query, target.Notes, target.Completed, target.ID).Scan(&target.UpdatedAt)
}

// DeleteTarget removes a target from a mission.
func (r *TargetRepository) DeleteTarget(ctx context.Context, id int) error {
	query := `DELETE FROM targets WHERE id = $1 AND completed = false`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return fmt.Errorf("target is already completed or does not exist")
	}
	return err
}

// GetTargetsByMissionID retrieves all targets for a given mission.
func (r *TargetRepository) GetTargetsByMissionID(ctx context.Context, missionID int) ([]domain.Target, error) {
	var targets []domain.Target
	query := `SELECT id, mission_id, name, country, notes, completed, created_at, updated_at 
			 FROM targets WHERE mission_id = $1 ORDER BY created_at`
	err := r.db.SelectContext(ctx, &targets, query, missionID)
	return targets, err
}
