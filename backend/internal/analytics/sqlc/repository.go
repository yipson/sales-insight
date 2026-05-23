package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/analytics"
)

// SQLCRepository implements analytics.Repository using sqlc-generated code.
type SQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewSQLCRepository creates a new sqlc-backed analytics repository.
func NewSQLCRepository(db *sql.DB) *SQLCRepository {
	return &SQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func (r *SQLCRepository) GetSalesSummary(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*analytics.SalesSummary, error) {
	row, err := r.queries.GetSalesSummary(ctx, GetSalesSummaryParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	return &analytics.SalesSummary{
		OrderCount:     row.OrderCount,
		TotalSales:     row.TotalSales,
		AvgTicket:      row.AvgTicket,
		TotalTax:       row.TotalTax,
		TotalTips:      row.TotalTips,
		TotalDiscounts: row.TotalDiscounts,
	}, nil
}

func (r *SQLCRepository) GetSalesByEmployee(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]analytics.SalesByEmployee, error) {
	rows, err := r.queries.GetSalesByEmployee(ctx, GetSalesByEmployeeParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	var out []analytics.SalesByEmployee
	for _, row := range rows {
		out = append(out, analytics.SalesByEmployee{
			EmployeeID:   row.EmployeeID,
			EmployeeName: row.EmployeeName,
			OrderCount:   row.OrderCount,
			TotalSales:   row.TotalSales,
			AvgTicket:    row.AvgTicket,
		})
	}
	return out, nil
}

func (r *SQLCRepository) GetTopProducts(ctx context.Context, restaurantID uuid.UUID, from, to time.Time, limit int32) ([]analytics.TopProduct, error) {
	rows, err := r.queries.GetTopProducts(ctx, GetTopProductsParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}
	var out []analytics.TopProduct
	for _, row := range rows {
		out = append(out, analytics.TopProduct{
			ProductID:     row.ProductID,
			ProductName:   row.ProductName,
			TotalQuantity: row.TotalQuantity,
			TotalRevenue:  row.TotalRevenue,
		})
	}
	return out, nil
}

func (r *SQLCRepository) GetCategoryCoverage(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]analytics.CategoryCoverageRow, error) {
	rows, err := r.queries.GetCategoryCoverage(ctx, GetCategoryCoverageParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	var out []analytics.CategoryCoverageRow
	for _, row := range rows {
		out = append(out, analytics.CategoryCoverageRow{
			EmployeeID:        row.EmployeeID,
			EmployeeName:      row.EmployeeName,
			CategoriesCovered: row.CategoriesCovered,
		})
	}
	return out, nil
}

func (r *SQLCRepository) CountActiveAnalyticCategories(ctx context.Context, restaurantID uuid.UUID) (int32, error) {
	return r.queries.CountActiveAnalyticCategories(ctx, restaurantID)
}

func (r *SQLCRepository) GetOrderCategoryBreakdown(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]analytics.OrderCategoryBreakdown, error) {
	rows, err := r.queries.GetOrderCategoryBreakdown(ctx, GetOrderCategoryBreakdownParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	var out []analytics.OrderCategoryBreakdown
	for _, row := range rows {
		out = append(out, analytics.OrderCategoryBreakdown{
			OrderID:       row.OrderID,
			CategorySlug:  row.CategorySlug,
			TotalQuantity: row.TotalQuantity,
		})
	}
	return out, nil
}
