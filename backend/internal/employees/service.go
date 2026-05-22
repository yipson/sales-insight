package employees

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service holds business logic for employees.
type Service struct {
	repo Repository
}

// NewService creates a new employee service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByID returns an employee by UUID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Employee, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCloverEmployeeID returns an employee by Clover employee ID.
func (s *Service) GetByCloverEmployeeID(ctx context.Context, restaurantID uuid.UUID, cloverEmployeeID string) (*Employee, error) {
	return s.repo.GetByCloverEmployeeID(ctx, restaurantID, cloverEmployeeID)
}

// ListByRestaurant returns all employees for a restaurant.
func (s *Service) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Employee, error) {
	return s.repo.ListByRestaurant(ctx, restaurantID)
}

// Upsert creates or updates a single employee (used by sync engine).
func (s *Service) Upsert(ctx context.Context, employee *Employee) error {
	return s.repo.Upsert(ctx, employee)
}

// UpsertBatch creates or updates multiple employees atomically (used by sync engine).
func (s *Service) UpsertBatch(ctx context.Context, employees []Employee) error {
	return s.repo.UpsertBatch(ctx, employees)
}

// Deactivate marks an employee as inactive.
func (s *Service) Deactivate(ctx context.Context, id uuid.UUID) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("lookup employee: %w", err)
	}
	if emp == nil {
		return fmt.Errorf("employee not found")
	}
	return s.repo.Deactivate(ctx, id)
}
