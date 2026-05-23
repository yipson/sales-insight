package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/analytics"
	"github.com/sales-insight/backend/internal/merchant"
)

// Service holds business logic for dashboard composition.
type Service struct {
	analyticsSvc *analytics.Service
	merchantSvc  *merchant.Service
}

// NewService creates a new dashboard service.
func NewService(analyticsSvc *analytics.Service, merchantSvc *merchant.Service) *Service {
	return &Service{
		analyticsSvc: analyticsSvc,
		merchantSvc:    merchantSvc,
	}
}

// GetSummary returns the dashboard summary for a restaurant and date range.
func (s *Service) GetSummary(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*SummaryResponse, error) {
	m, err := s.merchantSvc.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("lookup merchant: %w", err)
	}
	if m == nil {
		return nil, fmt.Errorf("merchant not found")
	}

	summary, err := s.analyticsSvc.GetSalesSummary(ctx, restaurantID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get sales summary: %w", err)
	}

	return &SummaryResponse{
		RestaurantID:   m.ID,
		RestaurantName: m.Name,
		From:           from,
		To:             to,
		Sales:          *summary,
	}, nil
}

// GetSalesByEmployee returns sales metrics per employee.
func (s *Service) GetSalesByEmployee(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*SalesByEmployeeResponse, error) {
	emps, err := s.analyticsSvc.GetSalesByEmployee(ctx, restaurantID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get sales by employee: %w", err)
	}
	return &SalesByEmployeeResponse{
		From:      from,
		To:        to,
		Employees: emps,
	}, nil
}

// GetTopProducts returns the best-selling products.
func (s *Service) GetTopProducts(ctx context.Context, restaurantID uuid.UUID, from, to time.Time, limit int32) (*TopProductsResponse, error) {
	products, err := s.analyticsSvc.GetTopProducts(ctx, restaurantID, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("get top products: %w", err)
	}
	return &TopProductsResponse{
		From:     from,
		To:       to,
		Products: products,
	}, nil
}

// GetCategoryCoverage returns category coverage per employee.
func (s *Service) GetCategoryCoverage(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*CategoryCoverageResponse, error) {
	coverage, err := s.analyticsSvc.GetCategoryCoverage(ctx, restaurantID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get category coverage: %w", err)
	}
	return &CategoryCoverageResponse{
		From:      from,
		To:        to,
		Employees: coverage,
	}, nil
}

// GetTicketIdeal returns the ticket ideal analysis.
func (s *Service) GetTicketIdeal(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*TicketIdealResponse, error) {
	m, err := s.merchantSvc.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("lookup merchant: %w", err)
	}
	if m == nil {
		return nil, fmt.Errorf("merchant not found")
	}

	result, err := s.analyticsSvc.AnalyzeTicketIdeal(ctx, restaurantID, from, to, m.TicketCompletoRules)
	if err != nil {
		return nil, fmt.Errorf("analyze ticket ideal: %w", err)
	}

	return &TicketIdealResponse{
		From:   from,
		To:     to,
		Result: *result,
		Rules:  m.TicketCompletoRules,
	}, nil
}
