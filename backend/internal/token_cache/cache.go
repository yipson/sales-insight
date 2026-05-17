package tokencache

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// TokenData holds decrypted tokens for a merchant.
type TokenData struct {
	MerchantID            uuid.UUID
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

// Cache stores decrypted tokens in memory with safe concurrent access.
type Cache struct {
	mu     sync.RWMutex
	tokens map[uuid.UUID]TokenData
}

// NewCache creates an empty token cache.
func NewCache() *Cache {
	return &Cache{
		tokens: make(map[uuid.UUID]TokenData),
	}
}

// Get returns the token data for a merchant.
func (c *Cache) Get(merchantID uuid.UUID) (TokenData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	td, ok := c.tokens[merchantID]
	return td, ok
}

// Set stores token data for a merchant.
func (c *Cache) Set(merchantID uuid.UUID, data TokenData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tokens[merchantID] = data
}

// Delete removes a merchant's tokens from the cache.
func (c *Cache) Delete(merchantID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.tokens, merchantID)
}

// IsAccessTokenExpired checks if the access token is expired (with 5 min buffer).
func (c *Cache) IsAccessTokenExpired(merchantID uuid.UUID) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	td, ok := c.tokens[merchantID]
	if !ok {
		return true
	}
	return time.Now().Add(5 * time.Minute).After(td.AccessTokenExpiresAt)
}
