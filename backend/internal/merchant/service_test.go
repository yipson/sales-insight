package merchant

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	merchants map[uuid.UUID]*Merchant
	byClover  map[string]*Merchant
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		merchants: make(map[uuid.UUID]*Merchant),
		byClover:  make(map[string]*Merchant),
	}
}

func (m *mockRepository) GetByID(_ context.Context, id uuid.UUID) (*Merchant, error) {
	return m.merchants[id], nil
}

func (m *mockRepository) GetByCloverMerchantID(_ context.Context, cloverID string) (*Merchant, error) {
	return m.byClover[cloverID], nil
}

func (m *mockRepository) List(_ context.Context) ([]Merchant, error) {
	var out []Merchant
	for _, v := range m.merchants {
		out = append(out, *v)
	}
	return out, nil
}

func (m *mockRepository) Create(_ context.Context, merchant *Merchant) error {
	m.merchants[merchant.ID] = merchant
	if merchant.CloverMerchantID != "" {
		m.byClover[merchant.CloverMerchantID] = merchant
	}
	return nil
}

func (m *mockRepository) Update(_ context.Context, merchant *Merchant) error {
	m.merchants[merchant.ID] = merchant
	if merchant.CloverMerchantID != "" {
		m.byClover[merchant.CloverMerchantID] = merchant
	}
	return nil
}

func (m *mockRepository) UpdateTokens(_ context.Context, _ uuid.UUID, _, _ string, _, _ *int64) error {
	return nil
}

func (m *mockRepository) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func TestService_Bootstrap_CreatesNew(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	m, err := svc.Bootstrap(ctx, "Test Merchant", "clv_123", "atoken", "rtoken", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "Test Merchant" {
		t.Errorf("name = %q, want %q", m.Name, "Test Merchant")
	}
	if !m.IsConnected {
		t.Error("expected merchant to be connected")
	}
}

func TestService_Bootstrap_UpdatesExisting(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	existing := &Merchant{
		ID:               uuid.New(),
		Name:             "Old Name",
		CloverMerchantID: "clv_123",
		IsConnected:      false,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	_ = repo.Create(ctx, existing)

	m, err := svc.Bootstrap(ctx, "New Name", "clv_123", "new_token", "new_refresh", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "New Name" {
		t.Errorf("name = %q, want %q", m.Name, "New Name")
	}
	if !m.IsConnected {
		t.Error("expected merchant to be connected after bootstrap")
	}
}

func TestService_Revoke(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	existing := &Merchant{
		ID:               uuid.New(),
		Name:             "Merchant",
		CloverMerchantID: "clv_456",
		IsConnected:      true,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	_ = repo.Create(ctx, existing)

	if err := svc.Revoke(ctx, existing.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := repo.GetByID(ctx, existing.ID)
	if updated.IsConnected {
		t.Error("expected merchant to be disconnected")
	}
}
