package dashboard

import (
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/analytics"
)

// SummaryResponse represents the dashboard summary payload.
type SummaryResponse struct {
	RestaurantID   uuid.UUID              `json:"restaurant_id"`
	RestaurantName string                 `json:"restaurant_name"`
	From           time.Time              `json:"from"`
	To             time.Time              `json:"to"`
	Sales          analytics.SalesSummary `json:"sales"`
}

// SalesByEmployeeResponse wraps the analytics result with date range.
type SalesByEmployeeResponse struct {
	From    time.Time                    `json:"from"`
	To      time.Time                    `json:"to"`
	Employees []analytics.SalesByEmployee  `json:"employees"`
}

// TopProductsResponse wraps the analytics result.
type TopProductsResponse struct {
	From     time.Time              `json:"from"`
	To       time.Time              `json:"to"`
	Products []analytics.TopProduct `json:"products"`
}

// CategoryCoverageResponse wraps the analytics result.
type CategoryCoverageResponse struct {
	From     time.Time                   `json:"from"`
	To       time.Time                   `json:"to"`
	Employees []analytics.CategoryCoverage `json:"employees"`
}

// TicketIdealResponse wraps the ticket ideal analysis.
type TicketIdealResponse struct {
	From   time.Time                  `json:"from"`
	To     time.Time                  `json:"to"`
	Result analytics.TicketIdealResult `json:"result"`
	Rules  map[string]int             `json:"rules,omitempty"`
}
