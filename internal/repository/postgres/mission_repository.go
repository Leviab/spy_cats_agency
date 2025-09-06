package postgres

import (
	"context"
	"fmt"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/repository"
)

// MissionRepository implements the repository.MissionRepository interface.
type MissionRepository struct {
	db *DB
}

// NewMissionRepository creates a new mission repository.
func NewMissionRepository(db *DB) repository.MissionRepository {
	return &MissionRepository{db: db}
}

// CreateMission creates a new mission and its associated targets within a transaction.
func (r *MissionRepository) CreateMission(ctx context.Context, mission *domain.Mission) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create the mission
	missionQuery := `INSERT INTO missions (cat_id, completed) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err = tx.QueryRowxContext(ctx, missionQuery, mission.CatID, mission.Completed).Scan(&mission.ID, &mission.CreatedAt, &mission.UpdatedAt)
	if err != nil {
		return err
	}

	// Create the targets
	if len(mission.Targets) > 0 {
		for i := range mission.Targets {
			mission.Targets[i].MissionID = mission.ID
			targetQuery := `INSERT INTO targets (mission_id, name, country, notes, completed)
							VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
			err = tx.QueryRowxContext(ctx, targetQuery, mission.ID, mission.Targets[i].Name, mission.Targets[i].Country, mission.Targets[i].Notes, mission.Targets[i].Completed).
				Scan(&mission.Targets[i].ID, &mission.Targets[i].CreatedAt, &mission.Targets[i].UpdatedAt)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// GetMissionByID retrieves a mission and its targets.
func (r *MissionRepository) GetMissionByID(ctx context.Context, id int) (*domain.Mission, error) {
	var mission domain.Mission
	query := `SELECT id, cat_id, completed, created_at, updated_at FROM missions WHERE id = $1`
	if err := r.db.GetContext(ctx, &mission, query, id); err != nil {
		return nil, err
	}

	// Get targets for this mission
	var targets []domain.Target
	targetQuery := `SELECT id, mission_id, name, country, notes, completed, created_at, updated_at 
			 FROM targets WHERE mission_id = $1 ORDER BY created_at`
	err := r.db.SelectContext(ctx, &targets, targetQuery, id)
	if err != nil {
		return nil, err
	}
	mission.Targets = targets

	return &mission, nil
}

// ListMissions retrieves all missions.
func (r *MissionRepository) ListMissions(ctx context.Context) ([]domain.Mission, error) {
	var missions []domain.Mission
	query := `SELECT id, cat_id, completed, created_at, updated_at FROM missions ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &missions, query); err != nil {
		return nil, err
	}
	return missions, nil
}

// UpdateMission updates a mission's state.
func (r *MissionRepository) UpdateMission(ctx context.Context, mission *domain.Mission) error {
	query := `UPDATE missions SET cat_id = $1, completed = $2, updated_at = now() WHERE id = $3 RETURNING updated_at`
	return r.db.QueryRowxContext(ctx, query, mission.CatID, mission.Completed, mission.ID).Scan(&mission.UpdatedAt)
}

// DeleteMission deletes a mission.
func (r *MissionRepository) DeleteMission(ctx context.Context, id int) error {
	query := `DELETE FROM missions WHERE id = $1 AND cat_id IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return fmt.Errorf("mission is assigned to a cat or does not exist")
	}
	return err
}

// AssignCatToMission assigns a cat to a mission and updates the cat's status.
func (r *MissionRepository) AssignCatToMission(ctx context.Context, missionID, catID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Assign cat to mission
	missionQuery := `UPDATE missions SET cat_id = $1, updated_at = now() WHERE id = $2`
	_, err = tx.ExecContext(ctx, missionQuery, catID, missionID)
	if err != nil {
		return err
	}

	// Update cat status
	catQuery := `UPDATE cats SET status = 'on_mission', updated_at = now() WHERE id = $1 AND status = 'available'`
	result, err := tx.ExecContext(ctx, catQuery, catID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return fmt.Errorf("cat is not available")
	}

	return tx.Commit()
}
