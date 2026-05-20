package products

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a synchronized product/item from Clover.
type Product struct {
	ID           uuid.UUID  `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	CloverItemID string     `json:"clover_item_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Price        int64      `json:"price"`
	Cost         *int64     `json:"cost,omitempty"`
	SKU          string     `json:"sku"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	IsAvailable  bool       `json:"is_available"`
	IsDeleted    bool       `json:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	SyncedAt     time.Time  `json:"synced_at"`
}

// Category represents a native Clover category.
type Category struct {
	ID               uuid.UUID `json:"id"`
	RestaurantID     uuid.UUID `json:"restaurant_id"`
	CloverCategoryID string    `json:"clover_category_id"`
	Name             string    `json:"name"`
	SortOrder        *int32    `json:"sort_order,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	SyncedAt         time.Time `json:"synced_at"`
}

// AnalyticCategory represents a customizable analytic category per restaurant.
type AnalyticCategory struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	DisplayOrder int32     `json:"display_order"`
	Color        string    `json:"color"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CategoryMapping maps a native Clover category to an analytic category.
type CategoryMapping struct {
	ID                 uuid.UUID `json:"id"`
	RestaurantID       uuid.UUID `json:"restaurant_id"`
	CategoryID         uuid.UUID `json:"category_id"`
	AnalyticCategoryID uuid.UUID `json:"analytic_category_id"`
	MappedBy           string    `json:"mapped_by"`
	MappedAt           time.Time `json:"mapped_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
