package service

import (
	"context"
	"fmt"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/repository"
	"spy_cats_agency/pkg/catapi"
)

// catService is the implementation of the CatService interface.
type catService struct {
	catRepo    repository.CatRepository
	catAPIClient *catapi.Client
}

// NewCatService creates a new CatService.
func NewCatService(catRepo repository.CatRepository, catAPIClient *catapi.Client) CatService {
	return &catService{
		catRepo:    catRepo,
		catAPIClient: catAPIClient,
	}
}

// CreateCat validates the breed and creates a new cat.
func (s *catService) CreateCat(ctx context.Context, cat *domain.Cat) error {
	valid, err := s.catAPIClient.IsValidBreed(cat.Breed)
	if err != nil {
		return fmt.Errorf("failed to validate breed: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid cat breed: %s", cat.Breed)
	}

	return s.catRepo.CreateCat(ctx, cat)
}

// GetCat retrieves a cat by its ID.
func (s *catService) GetCat(ctx context.Context, id int) (*domain.Cat, error) {
	return s.catRepo.GetCatByID(ctx, id)
}

// ListCats retrieves all cats.
func (s *catService) ListCats(ctx context.Context) ([]domain.Cat, error) {
	return s.catRepo.ListCats(ctx)
}

// UpdateCatSalary updates a cat's salary.
func (s *catService) UpdateCatSalary(ctx context.Context, id int, salary float64) (*domain.Cat, error) {
	cat, err := s.catRepo.GetCatByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cat.Salary = salary
	if err := s.catRepo.UpdateCat(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

// DeleteCat deletes a cat.
func (s *catService) DeleteCat(ctx context.Context, id int) error {
	return s.catRepo.DeleteCat(ctx, id)
}
