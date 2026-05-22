package sqlc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/employees"
)

// SQLCRepository implements employees.Repository using sqlc-generated code.
type SQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewSQLCRepository creates a new sqlc-backed employee repository.
func NewSQLCRepository(db *sql.DB) *SQLCRepository {
	return &SQLCRepository{
		db:      db,
		queries: New(db),
	}
}

// toDomain converts a sqlc Employee to a domain Employee.
func toDomain(e Employee) *employees.Employee {
	emp := &employees.Employee{
		ID:               e.ID,
		RestaurantID:     e.RestaurantID,
		CloverEmployeeID: e.CloverEmployeeID,
		Name:             e.Name,
		IsActive:         e.IsActive,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
		SyncedAt:         e.SyncedAt,
	}
	if e.Email.Valid {
		emp.Email = e.Email.String
	}
	if e.Phone.Valid {
		emp.Phone = e.Phone.String
	}
	if e.Role.Valid {
		emp.Role = e.Role.String
	}
	return emp
}

// fromDomain converts a domain Employee to sqlc UpsertEmployeeParams.
func fromDomain(e *employees.Employee) UpsertEmployeeParams {
	return UpsertEmployeeParams{
		ID:               e.ID,
		RestaurantID:     e.RestaurantID,
		CloverEmployeeID: e.CloverEmployeeID,
		Name:             e.Name,
		Email:            sql.NullString{String: e.Email, Valid: e.Email != ""},
		Phone:            sql.NullString{String: e.Phone, Valid: e.Phone != ""},
		Role:             sql.NullString{String: e.Role, Valid: e.Role != ""},
		IsActive:         e.IsActive,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
		SyncedAt:         e.SyncedAt,
	}
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*employees.Employee, error) {
	row, err := r.queries.GetEmployeeByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}

func (r *SQLCRepository) GetByCloverEmployeeID(ctx context.Context, restaurantID uuid.UUID, cloverEmployeeID string) (*employees.Employee, error) {
	row, err := r.queries.GetEmployeeByCloverID(ctx, GetEmployeeByCloverIDParams{
		RestaurantID:     restaurantID,
		CloverEmployeeID: cloverEmployeeID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}

func (r *SQLCRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]employees.Employee, error) {
	rows, err := r.queries.ListEmployeesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	var out []employees.Employee
	for _, row := range rows {
		out = append(out, *toDomain(row))
	}
	return out, nil
}

func (r *SQLCRepository) Upsert(ctx context.Context, employee *employees.Employee) error {
	_, err := r.queries.UpsertEmployee(ctx, fromDomain(employee))
	return err
}

func (r *SQLCRepository) UpsertBatch(ctx context.Context, employees []employees.Employee) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range employees {
		if _, err := qtx.UpsertEmployee(ctx, fromDomain(&employees[i])); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLCRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeactivateEmployee(ctx, id)
}
