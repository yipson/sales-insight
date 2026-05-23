package analytics

import "github.com/google/uuid"

// SalesSummary represents aggregated sales KPIs for a date range.
type SalesSummary struct {
	OrderCount      int32   `json:"order_count"`
	TotalSales      int64   `json:"total_sales"`
	AvgTicket       float64 `json:"avg_ticket"`
	TotalTax        int64   `json:"total_tax"`
	TotalTips       int64   `json:"total_tips"`
	TotalDiscounts  int64   `json:"total_discounts"`
}

// SalesByEmployee represents sales metrics per employee.
type SalesByEmployee struct {
	EmployeeID   uuid.UUID `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	OrderCount   int32     `json:"order_count"`
	TotalSales   int64     `json:"total_sales"`
	AvgTicket    float64   `json:"avg_ticket"`
}

// TopProduct represents a best-selling product.
type TopProduct struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	TotalQuantity int64     `json:"total_quantity"`
	TotalRevenue  int64     `json:"total_revenue"`
}

// CategoryCoverage represents how many analytic categories an employee sold.
type CategoryCoverage struct {
	EmployeeID         uuid.UUID `json:"employee_id"`
	EmployeeName       string    `json:"employee_name"`
	CategoriesCovered  int32     `json:"categories_covered"`
	TotalCategories    int32     `json:"total_categories"`
	CoveragePercent    float64   `json:"coverage_percent"`
}

// TicketIdealResult represents the analysis of orders against ideal-ticket rules.
type TicketIdealResult struct {
	TotalOrders      int32            `json:"total_orders"`
	CompleteOrders   int32            `json:"complete_orders"`
	IncompleteOrders int32            `json:"incomplete_orders"`
	CompletionRate   float64          `json:"completion_rate"`
}

// OrderCategoryBreakdown is used internally for ticket-ideal analysis.
type OrderCategoryBreakdown struct {
	OrderID        uuid.UUID `json:"order_id"`
	CategorySlug   string    `json:"category_slug"`
	TotalQuantity  int32     `json:"total_quantity"`
}
