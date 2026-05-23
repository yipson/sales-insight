package orders

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	orders     map[uuid.UUID]*Order
	byClover   map[string]*Order
	byRestaurant map[uuid.UUID][]Order
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		orders:       make(map[uuid.UUID]*Order),
		byClover:     make(map[string]*Order),
		byRestaurant: make(map[uuid.UUID][]Order),
	}
}

func (m *mockRepository) GetByID(_ context.Context, id uuid.UUID) (*Order, error) {
	return m.orders[id], nil
}

func (m *mockRepository) GetByCloverOrderID(_ context.Context, restaurantID uuid.UUID, cloverOrderID string) (*Order, error) {
	return m.byClover[cloverOrderID], nil
}

func (m *mockRepository) ListByRestaurantAndDate(_ context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Order, error) {
	var out []Order
	for _, o := range m.byRestaurant[restaurantID] {
		if o.CreatedTime.After(from) && o.CreatedTime.Before(to) {
			out = append(out, o)
		}
	}
	return out, nil
}

func (m *mockRepository) ListByEmployee(_ context.Context, employeeID uuid.UUID, from, to time.Time) ([]Order, error) {
	var out []Order
	for _, o := range m.orders {
		if o.EmployeeID != nil && *o.EmployeeID == employeeID && o.CreatedTime.After(from) && o.CreatedTime.Before(to) {
			out = append(out, *o)
		}
	}
	return out, nil
}

func (m *mockRepository) Upsert(_ context.Context, o *Order) error {
	m.orders[o.ID] = o
	m.byClover[o.CloverOrderID] = o
	m.byRestaurant[o.RestaurantID] = append(m.byRestaurant[o.RestaurantID], *o)
	return nil
}

func (m *mockRepository) UpsertBatch(_ context.Context, orders []Order) error {
	for i := range orders {
		if err := m.Upsert(nil, &orders[i]); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockRepository) DeleteByCloverID(_ context.Context, restaurantID uuid.UUID, cloverOrderID string) error {
	delete(m.byClover, cloverOrderID)
	return nil
}

// mockOrderItemRepository is a test double for OrderItemRepository.
type mockOrderItemRepository struct {
	items map[uuid.UUID][]OrderItem
}

func newMockOrderItemRepository() *mockOrderItemRepository {
	return &mockOrderItemRepository{
		items: make(map[uuid.UUID][]OrderItem),
	}
}

func (m *mockOrderItemRepository) UpsertBatch(_ context.Context, items []OrderItem) error {
	for _, item := range items {
		m.items[item.OrderID] = append(m.items[item.OrderID], item)
	}
	return nil
}

func (m *mockOrderItemRepository) GetByOrderID(_ context.Context, orderID uuid.UUID) ([]OrderItem, error) {
	return m.items[orderID], nil
}

// mockCategorySummaryRepository is a test double for CategorySummaryRepository.
type mockCategorySummaryRepository struct {
	summaries map[uuid.UUID][]OrderCategorySummary
}

func newMockCategorySummaryRepository() *mockCategorySummaryRepository {
	return &mockCategorySummaryRepository{
		summaries: make(map[uuid.UUID][]OrderCategorySummary),
	}
}

func (m *mockCategorySummaryRepository) UpsertBatch(_ context.Context, summaries []OrderCategorySummary) error {
	for _, s := range summaries {
		m.summaries[s.OrderID] = append(m.summaries[s.OrderID], s)
	}
	return nil
}

func (m *mockCategorySummaryRepository) GetByOrderID(_ context.Context, orderID uuid.UUID) ([]OrderCategorySummary, error) {
	return m.summaries[orderID], nil
}

func newTestService() (*Service, *mockRepository, *mockOrderItemRepository, *mockCategorySummaryRepository) {
	repo := newMockRepository()
	itemRepo := newMockOrderItemRepository()
	summaryRepo := newMockCategorySummaryRepository()
	svc := NewService(repo, itemRepo, summaryRepo)
	return svc, repo, itemRepo, summaryRepo
}

func TestService_GetByID(t *testing.T) {
	svc, repo, _, _ := newTestService()
	ctx := context.Background()

	o := &Order{ID: uuid.New(), CloverOrderID: "clv_1"}
	repo.orders[o.ID] = o

	got, err := svc.GetByID(ctx, o.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.CloverOrderID != "clv_1" {
		t.Errorf("expected clv_1, got %v", got)
	}
}

func TestService_UpsertBatch(t *testing.T) {
	svc, repo, _, _ := newTestService()
	ctx := context.Background()
	restID := uuid.New()

	orders := []Order{
		{ID: uuid.New(), CloverOrderID: "clv_1", RestaurantID: restID, CreatedTime: time.Now().UTC()},
		{ID: uuid.New(), CloverOrderID: "clv_2", RestaurantID: restID, CreatedTime: time.Now().UTC()},
	}

	if err := svc.UpsertBatch(ctx, orders); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(repo.orders))
	}
}

func TestService_GetOrderWithDetails(t *testing.T) {
	svc, repo, itemRepo, summaryRepo := newTestService()
	ctx := context.Background()

	orderID := uuid.New()
	order := &Order{ID: orderID, CloverOrderID: "clv_1", CreatedTime: time.Now().UTC()}
	repo.orders[orderID] = order

	items := []OrderItem{
		{ID: uuid.New(), OrderID: orderID, Name: "Taco", Quantity: 2},
	}
	itemRepo.items[orderID] = items

	summaries := []OrderCategorySummary{
		{ID: uuid.New(), OrderID: orderID, AnalyticCategoryID: uuid.New(), ItemCount: 1, TotalQuantity: 2, TotalAmount: 1000},
	}
	summaryRepo.summaries[orderID] = summaries

	detail, err := svc.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.Order.CloverOrderID != "clv_1" {
		t.Errorf("order id mismatch")
	}
	if len(detail.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(detail.Items))
	}
	if len(detail.CategorySummary) != 1 {
		t.Errorf("expected 1 summary, got %d", len(detail.CategorySummary))
	}
}

func TestService_GetOrderWithDetails_NotFound(t *testing.T) {
	svc, _, _, _ := newTestService()
	ctx := context.Background()

	_, err := svc.GetOrderWithDetails(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for non-existent order")
	}
}

func TestService_DeleteByCloverID(t *testing.T) {
	svc, repo, _, _ := newTestService()
	ctx := context.Background()
	restID := uuid.New()

	order := &Order{ID: uuid.New(), CloverOrderID: "clv_del", RestaurantID: restID}
	repo.byClover["clv_del"] = order

	if err := svc.DeleteByCloverID(ctx, restID, "clv_del"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := repo.byClover["clv_del"]; ok {
		t.Error("expected order to be deleted")
	}
}
