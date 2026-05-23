package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/analytics"
	"github.com/sales-insight/backend/internal/merchant"
)

// mockMerchantRepository is a test double for merchant.Repository.
type mockMerchantRepository struct {
	merchants map[uuid.UUID]*merchant.Merchant
}

func newMockMerchantRepository() *mockMerchantRepository {
	return &mockMerchantRepository{merchants: make(map[uuid.UUID]*merchant.Merchant)}
}

func (m *mockMerchantRepository) GetByID(_ context.Context, id uuid.UUID) (*merchant.Merchant, error) {
	return m.merchants[id], nil
}
func (m *mockMerchantRepository) GetByCloverMerchantID(_ context.Context, _ string) (*merchant.Merchant, error) {
	return nil, nil
}
func (m *mockMerchantRepository) List(_ context.Context) ([]merchant.Merchant, error) {
	return nil, nil
}
func (m *mockMerchantRepository) Create(_ context.Context, _ *merchant.Merchant) error { return nil }
func (m *mockMerchantRepository) Update(_ context.Context, _ *merchant.Merchant) error { return nil }
func (m *mockMerchantRepository) UpdateTokens(_ context.Context, _ uuid.UUID, _, _ string, _, _ *int64) error {
	return nil
}
func (m *mockMerchantRepository) Delete(_ context.Context, _ uuid.UUID) error { return nil }

// mockAnalyticsRepository is a test double for analytics.Repository.
type mockAnalyticsRepository struct {
	summary          *analytics.SalesSummary
	salesByEmployee  []analytics.SalesByEmployee
	topProducts      []analytics.TopProduct
	categoryCoverage []analytics.CategoryCoverageRow
	totalCategories  int32
	breakdown        []analytics.OrderCategoryBreakdown
}

func newMockAnalyticsRepository() *mockAnalyticsRepository {
	return &mockAnalyticsRepository{}
}

func (m *mockAnalyticsRepository) GetSalesSummary(_ context.Context, _ uuid.UUID, _, _ time.Time) (*analytics.SalesSummary, error) {
	return m.summary, nil
}
func (m *mockAnalyticsRepository) GetSalesByEmployee(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]analytics.SalesByEmployee, error) {
	return m.salesByEmployee, nil
}
func (m *mockAnalyticsRepository) GetTopProducts(_ context.Context, _ uuid.UUID, _, _ time.Time, _ int32) ([]analytics.TopProduct, error) {
	return m.topProducts, nil
}
func (m *mockAnalyticsRepository) GetCategoryCoverage(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]analytics.CategoryCoverageRow, error) {
	return m.categoryCoverage, nil
}
func (m *mockAnalyticsRepository) CountActiveAnalyticCategories(_ context.Context, _ uuid.UUID) (int32, error) {
	return m.totalCategories, nil
}
func (m *mockAnalyticsRepository) GetOrderCategoryBreakdown(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]analytics.OrderCategoryBreakdown, error) {
	return m.breakdown, nil
}

func newTestServices() (*Service, *mockMerchantRepository, *mockAnalyticsRepository) {
	merchantRepo := newMockMerchantRepository()
	merchantSvc := merchant.NewService(merchantRepo)
	analyticsRepo := newMockAnalyticsRepository()
	analyticsSvc := analytics.NewService(analyticsRepo)
	return NewService(analyticsSvc, merchantSvc), merchantRepo, analyticsRepo
}

func TestService_GetSummary(t *testing.T) {
	svc, merchantRepo, analyticsRepo := newTestServices()
	ctx := context.Background()
	restID := uuid.New()

	merchantRepo.merchants[restID] = &merchant.Merchant{
		ID:   restID,
		Name: "Test Restaurant",
	}
	analyticsRepo.summary = &analytics.SalesSummary{
		OrderCount: 5,
		TotalSales: 25000,
		AvgTicket:  5000,
	}

	result, err := svc.GetSummary(ctx, restID, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RestaurantName != "Test Restaurant" {
		t.Errorf("name = %q, want %q", result.RestaurantName, "Test Restaurant")
	}
	if result.Sales.OrderCount != 5 {
		t.Errorf("order_count = %d, want 5", result.Sales.OrderCount)
	}
}

func TestService_GetSalesByEmployee(t *testing.T) {
	svc, _, analyticsRepo := newTestServices()
	ctx := context.Background()
	restID := uuid.New()

	analyticsRepo.salesByEmployee = []analytics.SalesByEmployee{
		{EmployeeName: "Alice", TotalSales: 10000},
	}

	result, err := svc.GetSalesByEmployee(ctx, restID, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Employees) != 1 {
		t.Errorf("expected 1 employee, got %d", len(result.Employees))
	}
}

func TestService_GetTopProducts(t *testing.T) {
	svc, _, analyticsRepo := newTestServices()
	ctx := context.Background()
	restID := uuid.New()

	analyticsRepo.topProducts = []analytics.TopProduct{
		{ProductName: "Taco", TotalQuantity: 50},
	}

	result, err := svc.GetTopProducts(ctx, restID, time.Now(), time.Now(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Products) != 1 {
		t.Errorf("expected 1 product, got %d", len(result.Products))
	}
}

func TestService_GetTicketIdeal(t *testing.T) {
	svc, merchantRepo, analyticsRepo := newTestServices()
	ctx := context.Background()
	restID := uuid.New()

	merchantRepo.merchants[restID] = &merchant.Merchant{
		ID:                    restID,
		Name:                  "Test",
		TicketCompletoRules:   map[string]int{"tacos": 2, "bebidas": 1},
	}
	order1 := uuid.New()
	analyticsRepo.breakdown = []analytics.OrderCategoryBreakdown{
		{OrderID: order1, CategorySlug: "tacos", TotalQuantity: 3},
		{OrderID: order1, CategorySlug: "bebidas", TotalQuantity: 2},
	}

	result, err := svc.GetTicketIdeal(ctx, restID, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result.TotalOrders != 1 {
		t.Errorf("total_orders = %d, want 1", result.Result.TotalOrders)
	}
	if result.Rules["tacos"] != 2 {
		t.Errorf("rules[tacos] = %d, want 2", result.Rules["tacos"])
	}
}

func TestService_GetSummary_MerchantNotFound(t *testing.T) {
	svc, _, _ := newTestServices()
	ctx := context.Background()

	_, err := svc.GetSummary(ctx, uuid.New(), time.Now(), time.Now())
	if err == nil {
		t.Error("expected error for non-existent merchant")
	}
}
