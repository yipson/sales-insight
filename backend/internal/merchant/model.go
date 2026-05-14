package merchant

import (
	"time"

	"github.com/google/uuid"
)

// Merchant represents a restaurant connected to Clover.
type Merchant struct {
	ID                     uuid.UUID      `json:"id"`
	Name                   string         `json:"name"`
	CloverMerchantID       string         `json:"clover_merchant_id"`
	AccessToken            string         `json:"-"` // never serialize
	RefreshToken           string         `json:"-"` // never serialize
	AccessTokenExpiresAt   *time.Time     `json:"-"`
	RefreshTokenExpiresAt  *time.Time     `json:"-"`
	CloverEnv              string         `json:"clover_env"`
	IsConnected            bool           `json:"is_connected"`
	IsActive               bool           `json:"is_active"`
	TicketCompletoRules    map[string]int `json:"ticket_completo_rules,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

// TokenData holds decrypted tokens for in-memory cache.
type TokenData struct {
	MerchantID            uuid.UUID
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  *time.Time
	RefreshTokenExpiresAt *time.Time
}
