package orders

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service holds business logic for orders.
type Service struct {
	repo               Repository
	itemRepo           OrderItemRepository
	categorySummaryRepo CategorySummaryRepository
}

// NewService creates a new order service.
func NewService(
	repo Repository,
	itemRepo OrderItemRepository,
	categorySummaryRepo CategorySummaryRepository,
) *Service {
	return &Service{
		repo:               repo,
		itemRepo:           itemRepo,
		categorySummaryRepo: categorySummaryRepo,
	}
}

// ============================================================
// Orders
// ============================================================

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByCloverOrderID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) (*Order, error) {
	return s.repo.GetByCloverOrderID(ctx, restaurantID, cloverOrderID)
}

func (s *Service) ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Order, error) {
	return s.repo.ListByRestaurantAndDate(ctx, restaurantID, from, to)
}

func (s *Service) ListByEmployee(ctx context.Context, employeeID uuid.UUID, from, to time.Time) ([]Order, error) {
	return s.repo.ListByEmployee(ctx, employeeID, from, to)
}

func (s *Service) Upsert(ctx context.Context, order *Order) error {
	return s.repo.Upsert(ctx, order)
}

func (s *Service) UpsertBatch(ctx context.Context, orders []Order) error {
	return s.repo.UpsertBatch(ctx, orders)
}

func (s *Service) DeleteByCloverID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) error {
	return s.repo.DeleteByCloverID(ctx, restaurantID, cloverOrderID)
}

// OrderDetail represents an order with its items and category summaries.
type OrderDetail struct {
	Order            Order                  `json:"order"`
	Items            []OrderItem            `json:"items"`
	CategorySummary  []OrderCategorySummary `json:"category_summary"`
}

// GetOrderWithDetails returns an order with its items and category summaries.
func (s *Service) GetOrderWithDetails(ctx context.Context, id uuid.UUID) (*OrderDetail, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lookup order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	items, err := s.itemRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lookup order items: %w", err)
	}

	summary, err := s.categorySummaryRepo.GetByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lookup category summary: %w", err)
	}

	return &OrderDetail{
		Order:           *order,
		Items:           items,
		CategorySummary: summary,
	}, nil
}

// ============================================================
// Order Items
// ============================================================

func (s *Service) UpsertItemBatch(ctx context.Context, items []OrderItem) error {
	return s.itemRepo.UpsertBatch(ctx, items)
}

func (s *Service) GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderItem, error) {
	return s.itemRepo.GetByOrderID(ctx, orderID)
}

// ============================================================
// Category Summary
// ============================================================

func (s *Service) UpsertCategorySummaryBatch(ctx context.Context, summaries []OrderCategorySummary) error {
	return s.categorySummaryRepo.UpsertBatch(ctx, summaries)
}
