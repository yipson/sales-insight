package employees

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/orders"
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

// mockOrderLister is a test double for orderLister.
type mockOrderLister struct {
	orders map[uuid.UUID][]orders.Order
}

func newMockOrderLister() *mockOrderLister {
	return &mockOrderLister{orders: make(map[uuid.UUID][]orders.Order)}
}

func (m *mockOrderLister) ListByEmployee(_ context.Context, employeeID uuid.UUID, from, to time.Time) ([]orders.Order, error) {
	return m.orders[employeeID], nil
}

func newTestService() (*Service, *mockRepository, *mockOrderLister) {
	repo := newMockRepository()
	orderLister := newMockOrderLister()
	return NewService(repo, orderLister), repo, orderLister
}

func TestService_GetByID(t *testing.T) {
	svc, repo, _ := newTestService()
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
	svc, repo, _ := newTestService()
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
	svc, repo, _ := newTestService()
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
	svc, repo, _ := newTestService()
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
	svc, _, _ := newTestService()
	ctx := context.Background()

	err := svc.Deactivate(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for non-existent employee")
	}
}

func TestService_GetEmployeeOrders(t *testing.T) {
	svc, repo, orderLister := newTestService()
	ctx := context.Background()
	now := time.Now().UTC()

	empID := uuid.New()
	repo.employees[empID] = &Employee{ID: empID, Name: "Alice"}

	orderLister.orders[empID] = []orders.Order{
		{ID: uuid.New(), CloverOrderID: "clv_1", CreatedTime: now},
		{ID: uuid.New(), CloverOrderID: "clv_2", CreatedTime: now},
	}

	got, err := svc.GetEmployeeOrders(ctx, empID, now.AddDate(0, 0, -1), now.AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 orders, got %d", len(got))
	}
}

func TestService_GetEmployeeOrders_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	ctx := context.Background()
	now := time.Now().UTC()

	_, err := svc.GetEmployeeOrders(ctx, uuid.New(), now.AddDate(0, 0, -1), now.AddDate(0, 0, 1))
	if err == nil {
		t.Error("expected error for non-existent employee")
	}
}
