package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service holds business logic for analytics calculations.
type Service struct {
	repo Repository
}

// NewService creates a new analytics service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetSalesSummary returns aggregated sales KPIs.
func (s *Service) GetSalesSummary(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*SalesSummary, error) {
	return s.repo.GetSalesSummary(ctx, restaurantID, from, to)
}

// GetSalesByEmployee returns sales metrics grouped by employee.
func (s *Service) GetSalesByEmployee(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]SalesByEmployee, error) {
	return s.repo.GetSalesByEmployee(ctx, restaurantID, from, to)
}

// GetTopProducts returns the best-selling products.
func (s *Service) GetTopProducts(ctx context.Context, restaurantID uuid.UUID, from, to time.Time, limit int32) ([]TopProduct, error) {
	return s.repo.GetTopProducts(ctx, restaurantID, from, to, limit)
}

// GetCategoryCoverage returns per-employee category coverage with percentages.
func (s *Service) GetCategoryCoverage(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]CategoryCoverage, error) {
	totalCategories, err := s.repo.CountActiveAnalyticCategories(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("count active categories: %w", err)
	}
	if totalCategories == 0 {
		return []CategoryCoverage{}, nil
	}

	rows, err := s.repo.GetCategoryCoverage(ctx, restaurantID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get category coverage: %w", err)
	}

	var out []CategoryCoverage
	for _, row := range rows {
		percent := float64(row.CategoriesCovered) / float64(totalCategories) * 100
		out = append(out, CategoryCoverage{
			EmployeeID:        row.EmployeeID,
			EmployeeName:      row.EmployeeName,
			CategoriesCovered: row.CategoriesCovered,
			TotalCategories:   totalCategories,
			CoveragePercent:   percent,
		})
	}
	return out, nil
}

// AnalyzeTicketIdeal calculates how many orders meet the restaurant's ticket-ideal rules.
// rules is a map[categorySlug]minimumQuantity, e.g. {"tacos": 2, "bebidas": 1}.
func (s *Service) AnalyzeTicketIdeal(ctx context.Context, restaurantID uuid.UUID, from, to time.Time, rules map[string]int) (*TicketIdealResult, error) {
	if len(rules) == 0 {
		return &TicketIdealResult{}, nil
	}

	rows, err := s.repo.GetOrderCategoryBreakdown(ctx, restaurantID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get order category breakdown: %w", err)
	}

	// Group by order_id
	orderCategories := make(map[uuid.UUID]map[string]int32)
	for _, row := range rows {
		if orderCategories[row.OrderID] == nil {
			orderCategories[row.OrderID] = make(map[string]int32)
		}
		orderCategories[row.OrderID][row.CategorySlug] = row.TotalQuantity
	}

	var complete, incomplete int32
	for _, cats := range orderCategories {
		isComplete := true
		for slug, minQty := range rules {
			if cats[slug] < int32(minQty) {
				isComplete = false
				break
			}
		}
		if isComplete {
			complete++
		} else {
			incomplete++
		}
	}

	total := complete + incomplete
	var rate float64
	if total > 0 {
		rate = float64(complete) / float64(total) * 100
	}

	return &TicketIdealResult{
		TotalOrders:      total,
		CompleteOrders:   complete,
		IncompleteOrders: incomplete,
		CompletionRate:   rate,
	}, nil
}
