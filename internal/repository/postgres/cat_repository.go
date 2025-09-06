package postgres

import (
	"context"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/repository"
)

// CatRepository implements the repository.CatRepository interface.
type CatRepository struct {
	db *DB
}

// NewCatRepository creates a new cat repository.
func NewCatRepository(db *DB) repository.CatRepository {
	return &CatRepository{db: db}
}

// CreateCat creates a new cat in the database.
func (r *CatRepository) CreateCat(ctx context.Context, cat *domain.Cat) error {
	query := `INSERT INTO cats (name, years_of_experience, breed, salary)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at, status`
	return r.db.QueryRowxContext(ctx, query, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary).
		Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt, &cat.Status)
}

// GetCatByID retrieves a cat by its ID.
func (r *CatRepository) GetCatByID(ctx context.Context, id int) (*domain.Cat, error) {
	var cat domain.Cat
	query := `SELECT id, name, years_of_experience, breed, salary, status, created_at, updated_at
			  FROM cats WHERE id = $1`
	err := r.db.GetContext(ctx, &cat, query, id)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// ListCats retrieves all cats from the database.
func (r *CatRepository) ListCats(ctx context.Context) ([]domain.Cat, error) {
	var cats []domain.Cat
	query := `SELECT id, name, years_of_experience, breed, salary, status, created_at, updated_at
			  FROM cats ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &cats, query)
	return cats, err
}

// UpdateCat updates a cat's information.
func (r *CatRepository) UpdateCat(ctx context.Context, cat *domain.Cat) error {
	query := `UPDATE cats SET salary = $1, updated_at = now() WHERE id = $2 RETURNING updated_at`
	return r.db.QueryRowxContext(ctx, query, cat.Salary, cat.ID).Scan(&cat.UpdatedAt)
}

// DeleteCat removes a cat from the database.
func (r *CatRepository) DeleteCat(ctx context.Context, id int) error {
	query := `DELETE FROM cats WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
