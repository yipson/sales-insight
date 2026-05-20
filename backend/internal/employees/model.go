package employees

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents a synchronized employee from Clover.
type Employee struct {
	ID               uuid.UUID `json:"id"`
	RestaurantID     uuid.UUID `json:"restaurant_id"`
	CloverEmployeeID string    `json:"clover_employee_id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Role             string    `json:"role"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	SyncedAt         time.Time `json:"synced_at"`
}
