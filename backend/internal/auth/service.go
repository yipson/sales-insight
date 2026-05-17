package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/merchant"
	"github.com/sales-insight/backend/internal/platform/security"
	tokencache "github.com/sales-insight/backend/internal/token_cache"
)

// Service handles Clover OAuth authentication and token management.
type Service struct {
	oauthClient  *clover.OAuthClient
	merchantRepo merchant.Repository
	encrypter    *security.Encrypter
	tokenCache   *tokencache.Cache
	clientID     string
	clientSecret string
}

// NewService creates an auth service.
func NewService(
	oauthClient *clover.OAuthClient,
	merchantRepo merchant.Repository,
	encrypter *security.Encrypter,
	tokenCache *tokencache.Cache,
	clientID, clientSecret string,
) *Service {
	return &Service{
		oauthClient:  oauthClient,
		merchantRepo: merchantRepo,
		encrypter:    encrypter,
		tokenCache:   tokenCache,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// HandleCallback processes the OAuth callback from Clover.
// It exchanges the authorization code for tokens and stores them.
func (s *Service) HandleCallback(ctx context.Context, code, cloverMerchantID string) (*merchant.Merchant, error) {
	// Exchange code for tokens
	tokenResp, err := s.oauthClient.ExchangeCode(ctx, s.clientID, s.clientSecret, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	// Encrypt tokens
	encAccess, err := s.encrypter.Encrypt(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("encrypt access token: %w", err)
	}
	encRefresh, err := s.encrypter.Encrypt(tokenResp.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("encrypt refresh token: %w", err)
	}

	// Find or create merchant
	m, err := s.merchantRepo.GetByCloverMerchantID(ctx, cloverMerchantID)
	if err != nil {
		return nil, fmt.Errorf("lookup merchant: %w", err)
	}

	now := time.Now().UTC()
	accessExp := time.Unix(tokenResp.AccessTokenExpiration, 0)
	refreshExp := time.Unix(tokenResp.RefreshTokenExpiration, 0)

	if m != nil {
		// Update existing merchant
		m.AccessToken = encAccess
		m.RefreshToken = encRefresh
		m.AccessTokenExpiresAt = &accessExp
		m.RefreshTokenExpiresAt = &refreshExp
		m.IsConnected = true
		m.UpdatedAt = now
		if err := s.merchantRepo.Update(ctx, m); err != nil {
			return nil, fmt.Errorf("update merchant: %w", err)
		}
	} else {
		// Create new merchant
		m = &merchant.Merchant{
			ID:                    uuid.New(),
			CloverMerchantID:      cloverMerchantID,
			AccessToken:           encAccess,
			RefreshToken:          encRefresh,
			AccessTokenExpiresAt:  &accessExp,
			RefreshTokenExpiresAt: &refreshExp,
			CloverEnv:             "sandbox",
			IsConnected:           true,
			IsActive:              true,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		if err := s.merchantRepo.Create(ctx, m); err != nil {
			return nil, fmt.Errorf("create merchant: %w", err)
		}
	}

	// Update in-memory cache
	s.tokenCache.Set(m.ID, tokencache.TokenData{
		MerchantID:            m.ID,
		AccessToken:           tokenResp.AccessToken,
		RefreshToken:          tokenResp.RefreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
	})

	return m, nil
}

// RefreshAccessToken refreshes the access token for a merchant if needed.
// It updates both the database and the in-memory cache.
func (s *Service) RefreshAccessToken(ctx context.Context, merchantID uuid.UUID) error {
	// Check cache first
	if td, ok := s.tokenCache.Get(merchantID); ok {
		if !s.tokenCache.IsAccessTokenExpired(merchantID) {
			return nil // token still valid
		}
		// Use cached refresh token
		return s.doRefresh(ctx, merchantID, td.RefreshToken)
	}

	// Fallback to DB
	m, err := s.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return fmt.Errorf("lookup merchant: %w", err)
	}
	if m == nil {
		return fmt.Errorf("merchant not found")
	}

	refreshToken, err := s.encrypter.Decrypt(m.RefreshToken)
	if err != nil {
		return fmt.Errorf("decrypt refresh token: %w", err)
	}

	return s.doRefresh(ctx, merchantID, refreshToken)
}

func (s *Service) doRefresh(ctx context.Context, merchantID uuid.UUID, refreshToken string) error {
	tokenResp, err := s.oauthClient.RefreshTokens(ctx, s.clientID, refreshToken)
	if err != nil {
		return fmt.Errorf("refresh tokens: %w", err)
	}

	// Encrypt new tokens
	encAccess, err := s.encrypter.Encrypt(tokenResp.AccessToken)
	if err != nil {
		return fmt.Errorf("encrypt access token: %w", err)
	}
	encRefresh, err := s.encrypter.Encrypt(tokenResp.RefreshToken)
	if err != nil {
		return fmt.Errorf("encrypt refresh token: %w", err)
	}

	// Update database
	now := time.Now().UTC()
	accessExp := time.Unix(tokenResp.AccessTokenExpiration, 0)
	refreshExp := time.Unix(tokenResp.RefreshTokenExpiration, 0)

	if err := s.merchantRepo.UpdateTokens(ctx, merchantID, encAccess, encRefresh, &tokenResp.AccessTokenExpiration, &tokenResp.RefreshTokenExpiration); err != nil {
		return fmt.Errorf("update tokens: %w", err)
	}

	// Update cache
	s.tokenCache.Set(merchantID, tokencache.TokenData{
		MerchantID:            merchantID,
		AccessToken:           tokenResp.AccessToken,
		RefreshToken:          tokenResp.RefreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
	})

	// Touch updated_at
	m, _ := s.merchantRepo.GetByID(ctx, merchantID)
	if m != nil {
		m.UpdatedAt = now
		_ = s.merchantRepo.Update(ctx, m)
	}

	return nil
}

// RebuildCache loads all connected merchants from DB into memory cache.
// Call this at startup.
func (s *Service) RebuildCache(ctx context.Context) error {
	merchants, err := s.merchantRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("list merchants: %w", err)
	}

	for _, m := range merchants {
		if !m.IsConnected || m.AccessToken == "" {
			continue
		}

		accessToken, err := s.encrypter.Decrypt(m.AccessToken)
		if err != nil {
			fmt.Printf("Error: decrypt access token for merchant %s: %v\n", m.ID, err)
			continue // skip merchants with undecryptable tokens
		}
		refreshToken, _ := s.encrypter.Decrypt(m.RefreshToken)

		s.tokenCache.Set(m.ID, tokencache.TokenData{
			MerchantID:            m.ID,
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  time.Time{},
			RefreshTokenExpiresAt: time.Time{},
		})
	}

	return nil
}

// GetAccessToken returns the current access token for a merchant (from cache).
func (s *Service) GetAccessToken(merchantID uuid.UUID) (string, bool) {
	td, ok := s.tokenCache.Get(merchantID)
	if !ok {
		return "", false
	}
	return td.AccessToken, true
}
