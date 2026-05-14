package merchant

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service holds business logic for merchants.
type Service struct {
	repo Repository
}

// NewService creates a new merchant service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByID returns a merchant by UUID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCloverID returns a merchant by Clover merchant ID.
func (s *Service) GetByCloverID(ctx context.Context, cloverMerchantID string) (*Merchant, error) {
	return s.repo.GetByCloverMerchantID(ctx, cloverMerchantID)
}

// List returns all active merchants.
func (s *Service) List(ctx context.Context) ([]Merchant, error) {
	return s.repo.List(ctx)
}

// Bootstrap creates or updates a merchant with tokens (for local dev).
func (s *Service) Bootstrap(ctx context.Context, name, cloverMerchantID, accessToken, refreshToken string, accessExp, refreshExp *time.Time) (*Merchant, error) {
	existing, err := s.repo.GetByCloverMerchantID(ctx, cloverMerchantID)
	if err != nil {
		return nil, fmt.Errorf("lookup merchant: %w", err)
	}

	now := time.Now().UTC()
	if existing != nil {
		existing.Name = name
		existing.AccessToken = accessToken
		existing.RefreshToken = refreshToken
		existing.AccessTokenExpiresAt = accessExp
		existing.RefreshTokenExpiresAt = refreshExp
		existing.IsConnected = true
		existing.UpdatedAt = now
		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("update merchant: %w", err)
		}
		return existing, nil
	}

	m := &Merchant{
		ID:                    uuid.New(),
		Name:                  name,
		CloverMerchantID:      cloverMerchantID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
		CloverEnv:             "sandbox",
		IsConnected:           true,
		IsActive:              true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("create merchant: %w", err)
	}
	return m, nil
}

// Revoke marks a merchant as disconnected.
func (s *Service) Revoke(ctx context.Context, id uuid.UUID) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("lookup merchant: %w", err)
	}
	if m == nil {
		return fmt.Errorf("merchant not found")
	}
	m.IsConnected = false
	m.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, m)
}
