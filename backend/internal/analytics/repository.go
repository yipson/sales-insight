package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the contract for analytics queries.
type Repository interface {
	GetSalesSummary(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) (*SalesSummary, error)
	GetSalesByEmployee(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]SalesByEmployee, error)
	GetTopProducts(ctx context.Context, restaurantID uuid.UUID, from, to time.Time, limit int32) ([]TopProduct, error)
	GetCategoryCoverage(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]CategoryCoverageRow, error)
	CountActiveAnalyticCategories(ctx context.Context, restaurantID uuid.UUID) (int32, error)
	GetOrderCategoryBreakdown(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]OrderCategoryBreakdown, error)
}

// CategoryCoverageRow is the raw row from the category-coverage query.
type CategoryCoverageRow struct {
	EmployeeID        uuid.UUID `json:"employee_id"`
	EmployeeName      string    `json:"employee_name"`
	CategoriesCovered int32     `json:"categories_covered"`
}
