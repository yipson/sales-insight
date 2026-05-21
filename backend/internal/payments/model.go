package payments

import (
	"time"

	"github.com/google/uuid"
)

// Payment represents a synchronized payment from Clover.
type Payment struct {
	ID                uuid.UUID  `json:"id"`
	RestaurantID      uuid.UUID  `json:"restaurant_id"`
	CloverPaymentID   string     `json:"clover_payment_id"`
	OrderID           *uuid.UUID `json:"order_id,omitempty"`
	EmployeeID        *uuid.UUID `json:"employee_id,omitempty"`
	Amount            int64      `json:"amount"`
	TipAmount         int64      `json:"tip_amount"`
	TaxAmount         int64      `json:"tax_amount"`
	PaymentType       string     `json:"payment_type"`
	CardType          string     `json:"card_type"`
	Result            string     `json:"result"`
	ExternalPaymentID string     `json:"external_payment_id"`
	CreatedTime       time.Time  `json:"created_time"`
	ModifiedTime      time.Time  `json:"modified_time"`
	SyncedAt          time.Time  `json:"synced_at"`
}
