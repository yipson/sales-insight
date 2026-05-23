package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/merchant"
	"github.com/sales-insight/backend/internal/platform/security"
	tokencache "github.com/sales-insight/backend/internal/token_cache"
)

// mockOAuthClient is a test double for clover.TokenExchanger.
type mockOAuthClient struct {
	exchangeFunc  func(ctx context.Context, clientID, clientSecret, code string) (*clover.TokenResponse, error)
	refreshFunc   func(ctx context.Context, clientID, refreshToken string) (*clover.TokenResponse, error)
}

func (m *mockOAuthClient) ExchangeCode(ctx context.Context, clientID, clientSecret, code string) (*clover.TokenResponse, error) {
	if m.exchangeFunc != nil {
		return m.exchangeFunc(ctx, clientID, clientSecret, code)
	}
	return nil, errors.New("exchange not mocked")
}

func (m *mockOAuthClient) RefreshTokens(ctx context.Context, clientID, refreshToken string) (*clover.TokenResponse, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, clientID, refreshToken)
	}
	return nil, errors.New("refresh not mocked")
}

// mockMerchantRepository reuses the same pattern as merchant/service_test.go.
type mockMerchantRepository struct {
	merchants map[uuid.UUID]*merchant.Merchant
	byClover  map[string]*merchant.Merchant
}

func newMockMerchantRepository() *mockMerchantRepository {
	return &mockMerchantRepository{
		merchants: make(map[uuid.UUID]*merchant.Merchant),
		byClover:  make(map[string]*merchant.Merchant),
	}
}

func (m *mockMerchantRepository) GetByID(_ context.Context, id uuid.UUID) (*merchant.Merchant, error) {
	return m.merchants[id], nil
}

func (m *mockMerchantRepository) GetByCloverMerchantID(_ context.Context, cloverID string) (*merchant.Merchant, error) {
	return m.byClover[cloverID], nil
}

func (m *mockMerchantRepository) List(_ context.Context) ([]merchant.Merchant, error) {
	var out []merchant.Merchant
	for _, v := range m.merchants {
		out = append(out, *v)
	}
	return out, nil
}

func (m *mockMerchantRepository) Create(_ context.Context, merch *merchant.Merchant) error {
	m.merchants[merch.ID] = merch
	if merch.CloverMerchantID != "" {
		m.byClover[merch.CloverMerchantID] = merch
	}
	return nil
}

func (m *mockMerchantRepository) Update(_ context.Context, merch *merchant.Merchant) error {
	m.merchants[merch.ID] = merch
	if merch.CloverMerchantID != "" {
		m.byClover[merch.CloverMerchantID] = merch
	}
	return nil
}

func (m *mockMerchantRepository) UpdateTokens(_ context.Context, _ uuid.UUID, _, _ string, _, _ *int64) error {
	return nil
}

func (m *mockMerchantRepository) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func setupTestService() (*Service, *mockOAuthClient, *mockMerchantRepository, *tokencache.Cache, *security.Encrypter) {
	enc, _ := security.NewEncrypter("01234567890123456789012345678901")
	cache := tokencache.NewCache()
	oauth := &mockOAuthClient{}
	repo := newMockMerchantRepository()
	svc := NewService(oauth, repo, enc, cache, "client_id", "client_secret")
	return svc, oauth, repo, cache, enc
}

func TestService_HandleCallback_CreatesNewMerchant(t *testing.T) {
	svc, oauth, repo, cache, _ := setupTestService()
	ctx := context.Background()

	oauth.exchangeFunc = func(_ context.Context, _, _, _ string) (*clover.TokenResponse, error) {
		return &clover.TokenResponse{
			AccessToken:            "access123",
			AccessTokenExpiration:  time.Now().Add(1 * time.Hour).Unix(),
			RefreshToken:           "refresh456",
			RefreshTokenExpiration: time.Now().Add(24 * time.Hour).Unix(),
		}, nil
	}

	m, err := svc.HandleCallback(ctx, "authcode", "clv_merchant_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected merchant, got nil")
	}
	if m.CloverMerchantID != "clv_merchant_123" {
		t.Errorf("clover_merchant_id = %q, want %q", m.CloverMerchantID, "clv_merchant_123")
	}
	if !m.IsConnected {
		t.Error("expected merchant to be connected")
	}

	// Verify DB
	stored, _ := repo.GetByCloverMerchantID(ctx, "clv_merchant_123")
	if stored == nil {
		t.Fatal("expected merchant in repository")
	}

	// Verify cache
	td, ok := cache.Get(m.ID)
	if !ok {
		t.Fatal("expected token in cache")
	}
	if td.AccessToken != "access123" {
		t.Errorf("cached access token = %q, want %q", td.AccessToken, "access123")
	}
}

func TestService_HandleCallback_UpdatesExistingMerchant(t *testing.T) {
	svc, oauth, repo, cache, _ := setupTestService()
	ctx := context.Background()

	existingID := uuid.New()
	existing := &merchant.Merchant{
		ID:               existingID,
		Name:             "Old Name",
		CloverMerchantID: "clv_merchant_456",
		IsConnected:      false,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	_ = repo.Create(ctx, existing)

	oauth.exchangeFunc = func(_ context.Context, _, _, _ string) (*clover.TokenResponse, error) {
		return &clover.TokenResponse{
			AccessToken:            "new_access",
			AccessTokenExpiration:  time.Now().Add(1 * time.Hour).Unix(),
			RefreshToken:           "new_refresh",
			RefreshTokenExpiration: time.Now().Add(24 * time.Hour).Unix(),
		}, nil
	}

	m, err := svc.HandleCallback(ctx, "authcode", "clv_merchant_456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ID != existingID {
		t.Errorf("expected existing merchant ID to be reused, got %s", m.ID)
	}
	if !m.IsConnected {
		t.Error("expected merchant to be connected after update")
	}

	td, ok := cache.Get(existingID)
	if !ok {
		t.Fatal("expected token in cache")
	}
	if td.AccessToken != "new_access" {
		t.Errorf("cached access token = %q, want %q", td.AccessToken, "new_access")
	}
}

func TestService_HandleCallback_ExchangeError(t *testing.T) {
	svc, oauth, _, _, _ := setupTestService()
	ctx := context.Background()

	oauth.exchangeFunc = func(_ context.Context, _, _, _ string) (*clover.TokenResponse, error) {
		return nil, errors.New("invalid code")
	}

	_, err := svc.HandleCallback(ctx, "badcode", "clv_merchant_789")
	if err == nil {
		t.Fatal("expected error from exchange")
	}
}

func TestService_RefreshAccessToken_FromCache_NotExpired(t *testing.T) {
	svc, _, _, cache, _ := setupTestService()
	ctx := context.Background()
	id := uuid.New()

	cache.Set(id, tokencache.TokenData{
		MerchantID:           id,
		AccessToken:          "cached_token",
		RefreshToken:         "refresh_token",
		AccessTokenExpiresAt: time.Now().Add(1 * time.Hour),
	})

	err := svc.RefreshAccessToken(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Token should still be in cache unchanged
	td, ok := cache.Get(id)
	if !ok {
		t.Fatal("expected token in cache")
	}
	if td.AccessToken != "cached_token" {
		t.Errorf("token was refreshed when it should not have been: %q", td.AccessToken)
	}
}

func TestService_RefreshAccessToken_FromCache_Expired(t *testing.T) {
	svc, oauth, _, cache, _ := setupTestService()
	ctx := context.Background()
	id := uuid.New()

	cache.Set(id, tokencache.TokenData{
		MerchantID:           id,
		AccessToken:          "old_token",
		RefreshToken:         "refresh_xyz",
		AccessTokenExpiresAt: time.Now().Add(-1 * time.Hour),
	})

	oauth.refreshFunc = func(_ context.Context, _, _ string) (*clover.TokenResponse, error) {
		return &clover.TokenResponse{
			AccessToken:            "refreshed_token",
			AccessTokenExpiration:  time.Now().Add(1 * time.Hour).Unix(),
			RefreshToken:           "new_refresh",
			RefreshTokenExpiration: time.Now().Add(24 * time.Hour).Unix(),
		}, nil
	}

	err := svc.RefreshAccessToken(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	td, ok := cache.Get(id)
	if !ok {
		t.Fatal("expected token in cache after refresh")
	}
	if td.AccessToken != "refreshed_token" {
		t.Errorf("access token = %q, want %q", td.AccessToken, "refreshed_token")
	}
}

func TestService_RefreshAccessToken_MerchantNotFound(t *testing.T) {
	svc, _, _, _, _ := setupTestService()
	ctx := context.Background()
	id := uuid.New()

	err := svc.RefreshAccessToken(ctx, id)
	if err == nil {
		t.Fatal("expected error when merchant not found")
	}
}

func TestService_RebuildCache(t *testing.T) {
	svc, _, repo, cache, enc := setupTestService()
	ctx := context.Background()

	// Create connected merchant with encrypted tokens
	m := &merchant.Merchant{
		ID:               uuid.New(),
		CloverMerchantID: "clv_abc",
		Name:             "Test Merchant",
		IsConnected:      true,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	encAccess, _ := enc.Encrypt("access_abc")
	encRefresh, _ := enc.Encrypt("refresh_abc")
	m.AccessToken = encAccess
	m.RefreshToken = encRefresh
	_ = repo.Create(ctx, m)

	// Also create a disconnected merchant that should be skipped
	m2 := &merchant.Merchant{
		ID:               uuid.New(),
		CloverMerchantID: "clv_def",
		Name:             "Disconnected",
		IsConnected:      false,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	_ = repo.Create(ctx, m2)

	err := svc.RebuildCache(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Connected merchant should be in cache
	td, ok := cache.Get(m.ID)
	if !ok {
		t.Fatal("expected connected merchant in cache")
	}
	if td.AccessToken != "access_abc" {
		t.Errorf("access token = %q, want %q", td.AccessToken, "access_abc")
	}
	if td.RefreshToken != "refresh_abc" {
		t.Errorf("refresh token = %q, want %q", td.RefreshToken, "refresh_abc")
	}

	// Disconnected merchant should NOT be in cache
	_, ok = cache.Get(m2.ID)
	if ok {
		t.Error("expected disconnected merchant to NOT be in cache")
	}
}

func TestService_GetAccessToken(t *testing.T) {
	svc, _, _, cache, _ := setupTestService()
	id := uuid.New()

	// Not found
	_, ok := svc.GetAccessToken(id)
	if ok {
		t.Error("expected not found")
	}

	// Set and get
	cache.Set(id, tokencache.TokenData{
		MerchantID:  id,
		AccessToken: "my_token",
	})

	token, ok := svc.GetAccessToken(id)
	if !ok {
		t.Fatal("expected token to be found")
	}
	if token != "my_token" {
		t.Errorf("token = %q, want %q", token, "my_token")
	}
}
