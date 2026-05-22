package employees

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	employees    map[uuid.UUID]*Employee
	byCloverID   map[string]*Employee
	byRestaurant map[uuid.UUID][]Employee
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		employees:    make(map[uuid.UUID]*Employee),
		byCloverID:   make(map[string]*Employee),
		byRestaurant: make(map[uuid.UUID][]Employee),
	}
}

func (m *mockRepository) GetByID(_ context.Context, id uuid.UUID) (*Employee, error) {
	return m.employees[id], nil
}

func (m *mockRepository) GetByCloverEmployeeID(_ context.Context, restaurantID uuid.UUID, cloverEmployeeID string) (*Employee, error) {
	return m.byCloverID[cloverEmployeeID], nil
}

func (m *mockRepository) ListByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]Employee, error) {
	return m.byRestaurant[restaurantID], nil
}

func (m *mockRepository) Upsert(_ context.Context, emp *Employee) error {
	m.employees[emp.ID] = emp
	if emp.CloverEmployeeID != "" {
		m.byCloverID[emp.CloverEmployeeID] = emp
	}
	m.byRestaurant[emp.RestaurantID] = append(m.byRestaurant[emp.RestaurantID], *emp)
	return nil
}

func (m *mockRepository) UpsertBatch(_ context.Context, emps []Employee) error {
	for i := range emps {
		if err := m.Upsert(nil, &emps[i]); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockRepository) Deactivate(_ context.Context, id uuid.UUID) error {
	if emp, ok := m.employees[id]; ok {
		emp.IsActive = false
	}
	return nil
}

func TestService_GetByID(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	emp := &Employee{
		ID:   uuid.New(),
		Name: "John Doe",
	}
	repo.employees[emp.ID] = emp

	got, err := svc.GetByID(ctx, emp.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "John Doe" {
		t.Errorf("expected John Doe, got %v", got)
	}
}

func TestService_ListByRestaurant(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()
	restID := uuid.New()

	repo.byRestaurant[restID] = []Employee{
		{ID: uuid.New(), Name: "Alice", RestaurantID: restID},
		{ID: uuid.New(), Name: "Bob", RestaurantID: restID},
	}

	got, err := svc.ListByRestaurant(ctx, restID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 employees, got %d", len(got))
	}
}

func TestService_UpsertBatch(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()
	restID := uuid.New()

	emps := []Employee{
		{ID: uuid.New(), Name: "Alice", RestaurantID: restID, CloverEmployeeID: "clv_1"},
		{ID: uuid.New(), Name: "Bob", RestaurantID: restID, CloverEmployeeID: "clv_2"},
	}

	if err := svc.UpsertBatch(ctx, emps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.employees) != 2 {
		t.Errorf("expected 2 employees in repo, got %d", len(repo.employees))
	}
}

func TestService_Deactivate(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	emp := &Employee{
		ID:       uuid.New(),
		Name:     "Alice",
		IsActive: true,
	}
	repo.employees[emp.ID] = emp

	if err := svc.Deactivate(ctx, emp.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := repo.GetByID(ctx, emp.ID)
	if updated.IsActive {
		t.Error("expected employee to be deactivated")
	}
}

func TestService_Deactivate_NotFound(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	err := svc.Deactivate(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for non-existent employee")
	}
}
