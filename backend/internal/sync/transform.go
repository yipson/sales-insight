package sync

import (
	"fmt"
	"time"
)

// Transform provides utilities for normalizing Clover data.
// These helpers ensure consistent data format across all sync operations.

// PriceToCents validates and converts a price value.
// Clover returns prices in cents as integers; this ensures non-negative values.
func PriceToCents(price int64) (int64, error) {
	if price < 0 {
		return 0, fmt.Errorf("price cannot be negative: %d", price)
	}
	return price, nil
}

// TimestampToUTC converts a Unix millisecond timestamp to UTC time.
func TimestampToUTC(ms int64) time.Time {
	return time.UnixMilli(ms).UTC()
}

// RoundToMicrosecond rounds a time to microsecond precision.
// PostgreSQL TIMESTAMPTZ has microsecond precision; rounding prevents
// sub-microsecond differences causing cursor drift.
func RoundToMicrosecond(t time.Time) time.Time {
	return t.Round(time.Microsecond)
}

// CategorySummaryInput holds data needed to compute category aggregates.
type CategorySummaryInput struct {
	OrderID    string
	ProductID  string
	Quantity   int32
	TotalPrice int64
	CategoryID string
}

// CalculateCategorySummary computes aggregated data per analytic category for an order.
// Reserved for Phase 4 when category mappings are fully implemented.
func CalculateCategorySummary(inputs []CategorySummaryInput) error {
	// TODO: Implement in Phase 4 with full category mapping resolution
	return nil
}
