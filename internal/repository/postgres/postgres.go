package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"spy_cats_agency/internal/config"
	"spy_cats_agency/internal/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver
)

// DB is a wrapper for the sqlx.DB that will implement our repository interfaces.
type DB struct {
	*sqlx.DB
}

// New creates a new database connection.
func New(cfg config.Config) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{db}, nil
}

// --- Cat Repository Implementation ---

// CreateCat creates a new cat in the database.
func (db *DB) CreateCat(ctx context.Context, cat *domain.Cat) error {
	query := `INSERT INTO cats (name, years_of_experience, breed, salary)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at, status`
	return db.QueryRowxContext(ctx, query, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary).
		Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt, &cat.Status)
}

// GetCatByID retrieves a cat by its ID.
func (db *DB) GetCatByID(ctx context.Context, id int) (*domain.Cat, error) {
	var cat domain.Cat
	query := `SELECT id, name, years_of_experience, breed, salary, status, created_at, updated_at
			  FROM cats WHERE id = $1`
	err := db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// ListCats retrieves all cats from the database.
func (db *DB) ListCats(ctx context.Context) ([]domain.Cat, error) {
	var cats []domain.Cat
	query := `SELECT id, name, years_of_experience, breed, salary, status, created_at, updated_at
			  FROM cats ORDER BY created_at DESC`
	err := db.SelectContext(ctx, &cats, query)
	return cats, err
}

// UpdateCat updates a cat's information.
func (db *DB) UpdateCat(ctx context.Context, cat *domain.Cat) error {
	query := `UPDATE cats SET salary = $1, updated_at = now() WHERE id = $2 RETURNING updated_at`
	return db.QueryRowxContext(ctx, query, cat.Salary, cat.ID).Scan(&cat.UpdatedAt)
}

// DeleteCat removes a cat from the database.
func (db *DB) DeleteCat(ctx context.Context, id int) error {
	query := `DELETE FROM cats WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// --- Mission Repository Implementation ---

// CreateMission creates a new mission and its associated targets within a transaction.
func (db *DB) CreateMission(ctx context.Context, mission *domain.Mission) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback is a no-op if the transaction is committed.

	// Create the mission
	missionQuery := `INSERT INTO missions (completed) VALUES ($1) RETURNING id, created_at, updated_at`
	err = tx.QueryRowxContext(ctx, missionQuery, mission.Completed).Scan(&mission.ID, &mission.CreatedAt, &mission.UpdatedAt)
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
func (db *DB) GetMissionByID(ctx context.Context, id int) (*domain.Mission, error) {
	var mission domain.Mission
	query := `SELECT id, cat_id, completed, created_at, updated_at FROM missions WHERE id = $1`
	if err := db.GetContext(ctx, &mission, query, id); err != nil {
		return nil, err
	}

	targets, err := db.GetTargetsByMissionID(ctx, id)
	if err != nil {
		return nil, err
	}
	mission.Targets = targets

	return &mission, nil
}

// ListMissions retrieves all missions.
func (db *DB) ListMissions(ctx context.Context) ([]domain.Mission, error) {
	var missions []domain.Mission
	query := `SELECT id, cat_id, completed, created_at, updated_at FROM missions ORDER BY created_at DESC`
	if err := db.SelectContext(ctx, &missions, query); err != nil {
		return nil, err
	}
	return missions, nil
}

// UpdateMission updates a mission's state.
func (db *DB) UpdateMission(ctx context.Context, mission *domain.Mission) error {
	query := `UPDATE missions SET completed = $1, updated_at = now() WHERE id = $2 RETURNING updated_at`
	return db.QueryRowxContext(ctx, query, mission.Completed, mission.ID).Scan(&mission.UpdatedAt)
}

// DeleteMission deletes a mission.
func (db *DB) DeleteMission(ctx context.Context, id int) error {
	query := `DELETE FROM missions WHERE id = $1 AND cat_id IS NULL`
	result, err := db.ExecContext(ctx, query, id)
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
func (db *DB) AssignCatToMission(ctx context.Context, missionID, catID int) error {
	tx, err := db.BeginTxx(ctx, nil)
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

// --- Target Repository Implementation ---

// AddTargetToMission adds a new target to an existing mission.
func (db *DB) AddTargetToMission(ctx context.Context, target *domain.Target) error {
	query := `INSERT INTO targets (mission_id, name, country, notes, completed)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return db.QueryRowxContext(ctx, query, target.MissionID, target.Name, target.Country, target.Notes, target.Completed).
		Scan(&target.ID, &target.CreatedAt, &target.UpdatedAt)
}

// GetTargetByID retrieves a single target by its ID.
func (db *DB) GetTargetByID(ctx context.Context, id int) (*domain.Target, error) {
	var target domain.Target
	query := `SELECT id, mission_id, name, country, notes, completed, created_at, updated_at FROM targets WHERE id = $1`
	err := db.GetContext(ctx, &target, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target not found")
		}
		return nil, err
	}
	return &target, nil
}

// UpdateTarget updates a target's notes or completion status.
func (db *DB) UpdateTarget(ctx context.Context, target *domain.Target) error {
	query := `UPDATE targets SET notes = $1, completed = $2, updated_at = now() WHERE id = $3 RETURNING updated_at`
	return db.QueryRowxContext(ctx, query, target.Notes, target.Completed, target.ID).Scan(&target.UpdatedAt)
}

// DeleteTarget removes a target from a mission.
func (db *DB) DeleteTarget(ctx context.Context, id int) error {
	query := `DELETE FROM targets WHERE id = $1 AND completed = false`
	result, err := db.ExecContext(ctx, query, id)
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
func (db *DB) GetTargetsByMissionID(ctx context.Context, missionID int) ([]domain.Target, error) {
	var targets []domain.Target
	query := `SELECT id, mission_id, name, country, notes, completed, created_at, updated_at 
			 FROM targets WHERE mission_id = $1 ORDER BY created_at`
	err := db.SelectContext(ctx, &targets, query, missionID)
	return targets, err
}
