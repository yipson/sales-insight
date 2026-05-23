package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	summary          *SalesSummary
	salesByEmployee  []SalesByEmployee
	topProducts      []TopProduct
	categoryCoverage []CategoryCoverageRow
	totalCategories  int32
	breakdown        []OrderCategoryBreakdown
}

func newMockRepository() *mockRepository {
	return &mockRepository{}
}

func (m *mockRepository) GetSalesSummary(_ context.Context, _ uuid.UUID, _, _ time.Time) (*SalesSummary, error) {
	return m.summary, nil
}

func (m *mockRepository) GetSalesByEmployee(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]SalesByEmployee, error) {
	return m.salesByEmployee, nil
}

func (m *mockRepository) GetTopProducts(_ context.Context, _ uuid.UUID, _, _ time.Time, _ int32) ([]TopProduct, error) {
	return m.topProducts, nil
}

func (m *mockRepository) GetCategoryCoverage(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]CategoryCoverageRow, error) {
	return m.categoryCoverage, nil
}

func (m *mockRepository) CountActiveAnalyticCategories(_ context.Context, _ uuid.UUID) (int32, error) {
	return m.totalCategories, nil
}

func (m *mockRepository) GetOrderCategoryBreakdown(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]OrderCategoryBreakdown, error) {
	return m.breakdown, nil
}

func TestService_GetSalesSummary(t *testing.T) {
	repo := newMockRepository()
	repo.summary = &SalesSummary{
		OrderCount: 10,
		TotalSales: 50000,
		AvgTicket:  5000,
	}
	svc := NewService(repo)
	ctx := context.Background()

	summary, err := svc.GetSalesSummary(ctx, uuid.New(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.OrderCount != 10 {
		t.Errorf("order_count = %d, want 10", summary.OrderCount)
	}
	if summary.TotalSales != 50000 {
		t.Errorf("total_sales = %d, want 50000", summary.TotalSales)
	}
}

func TestService_GetCategoryCoverage(t *testing.T) {
	repo := newMockRepository()
	repo.totalCategories = 5
	repo.categoryCoverage = []CategoryCoverageRow{
		{EmployeeID: uuid.New(), EmployeeName: "Alice", CategoriesCovered: 5},
		{EmployeeID: uuid.New(), EmployeeName: "Bob", CategoriesCovered: 3},
	}
	svc := NewService(repo)
	ctx := context.Background()

	coverage, err := svc.GetCategoryCoverage(ctx, uuid.New(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(coverage) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(coverage))
	}
	if coverage[0].CoveragePercent != 100 {
		t.Errorf("Alice coverage = %f, want 100", coverage[0].CoveragePercent)
	}
	if coverage[1].CoveragePercent != 60 {
		t.Errorf("Bob coverage = %f, want 60", coverage[1].CoveragePercent)
	}
}

func TestService_GetCategoryCoverage_ZeroCategories(t *testing.T) {
	repo := newMockRepository()
	repo.totalCategories = 0
	svc := NewService(repo)
	ctx := context.Background()

	coverage, err := svc.GetCategoryCoverage(ctx, uuid.New(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(coverage) != 0 {
		t.Errorf("expected empty coverage, got %d rows", len(coverage))
	}
}

func TestService_AnalyzeTicketIdeal(t *testing.T) {
	repo := newMockRepository()
	order1 := uuid.New()
	order2 := uuid.New()
	order3 := uuid.New()
	repo.breakdown = []OrderCategoryBreakdown{
		{OrderID: order1, CategorySlug: "tacos", TotalQuantity: 3},
		{OrderID: order1, CategorySlug: "bebidas", TotalQuantity: 2},
		{OrderID: order2, CategorySlug: "tacos", TotalQuantity: 1},
		{OrderID: order2, CategorySlug: "bebidas", TotalQuantity: 1},
		{OrderID: order3, CategorySlug: "tacos", TotalQuantity: 2},
		{OrderID: order3, CategorySlug: "postres", TotalQuantity: 1},
	}
	svc := NewService(repo)
	ctx := context.Background()

	rules := map[string]int{"tacos": 2, "bebidas": 1}
	result, err := svc.AnalyzeTicketIdeal(ctx, uuid.New(), time.Now(), time.Now(), rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalOrders != 3 {
		t.Errorf("total_orders = %d, want 3", result.TotalOrders)
	}
	if result.CompleteOrders != 1 {
		t.Errorf("complete_orders = %d, want 1", result.CompleteOrders)
	}
	if result.IncompleteOrders != 2 {
		t.Errorf("incomplete_orders = %d, want 2", result.IncompleteOrders)
	}
	wantRate := (1.0 / 3.0) * 100
	if result.CompletionRate < wantRate-0.001 || result.CompletionRate > wantRate+0.001 {
		t.Errorf("completion_rate = %f, want %f", result.CompletionRate, wantRate)
	}
}

func TestService_AnalyzeTicketIdeal_NoRules(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	result, err := svc.AnalyzeTicketIdeal(ctx, uuid.New(), time.Now(), time.Now(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalOrders != 0 {
		t.Errorf("expected zero totals with no rules, got %+v", result)
	}
}
