package orders

import (
	"time"

	"github.com/google/uuid"
)

// Order represents a synchronized order from Clover.
type Order struct {
	ID                  uuid.UUID   `json:"id"`
	RestaurantID        uuid.UUID   `json:"restaurant_id"`
	CloverOrderID       string      `json:"clover_order_id"`
	EmployeeID          *uuid.UUID  `json:"employee_id,omitempty"`
	OrderNumber         string      `json:"order_number"`
	OrderType           string      `json:"order_type"`
	State               string      `json:"state"`
	TotalAmount         int64       `json:"total_amount"`
	TaxAmount           int64       `json:"tax_amount"`
	DiscountAmount      int64       `json:"discount_amount"`
	TipAmount           int64       `json:"tip_amount"`
	CreatedTime         time.Time   `json:"created_time"`
	ModifiedTime        time.Time   `json:"modified_time"`
	PayType             string      `json:"pay_type"`
	ItemCount           int32       `json:"item_count"`
	UniqueCategoryCount int32       `json:"unique_category_count"`
	TotalQuantity       int32       `json:"total_quantity"`
	SyncedAt            time.Time   `json:"synced_at"`
}

// OrderItem represents a line item within an order.
type OrderItem struct {
	ID                 uuid.UUID  `json:"id"`
	RestaurantID       uuid.UUID  `json:"restaurant_id"`
	OrderID            uuid.UUID  `json:"order_id"`
	ProductID          *uuid.UUID `json:"product_id,omitempty"`
	CloverLineItemID   string     `json:"clover_line_item_id"`
	Name               string     `json:"name"`
	Quantity           int32      `json:"quantity"`
	UnitPrice          int64      `json:"unit_price"`
	TotalPrice         int64      `json:"total_price"`
	AnalyticCategoryID *uuid.UUID `json:"analytic_category_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	SyncedAt           time.Time  `json:"synced_at"`
}

// OrderCategorySummary represents aggregated data per category for an order.
type OrderCategorySummary struct {
	ID                 uuid.UUID `json:"id"`
	RestaurantID       uuid.UUID `json:"restaurant_id"`
	OrderID            uuid.UUID `json:"order_id"`
	AnalyticCategoryID uuid.UUID `json:"analytic_category_id"`
	ItemCount          int32     `json:"item_count"`
	TotalQuantity      int32     `json:"total_quantity"`
	TotalAmount        int64     `json:"total_amount"`
}
