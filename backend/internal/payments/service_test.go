package payments

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockRepository is a test double for Repository.
type mockRepository struct {
	payments     map[uuid.UUID]*Payment
	byClover     map[string]*Payment
	byRestaurant map[uuid.UUID][]*Payment
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		payments:     make(map[uuid.UUID]*Payment),
		byClover:     make(map[string]*Payment),
		byRestaurant: make(map[uuid.UUID][]*Payment),
	}
}

func (m *mockRepository) GetByID(_ context.Context, id uuid.UUID) (*Payment, error) {
	return m.payments[id], nil
}

func (m *mockRepository) GetByCloverPaymentID(_ context.Context, restaurantID uuid.UUID, cloverPaymentID string) (*Payment, error) {
	return m.byClover[cloverPaymentID], nil
}

func (m *mockRepository) ListByRestaurantAndDate(_ context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Payment, error) {
	var out []Payment
	for _, p := range m.byRestaurant[restaurantID] {
		if p.CreatedTime.After(from) && p.CreatedTime.Before(to) {
			out = append(out, *p)
		}
	}
	return out, nil
}

func (m *mockRepository) Upsert(_ context.Context, p *Payment) error {
	m.payments[p.ID] = p
	m.byClover[p.CloverPaymentID] = p
	m.byRestaurant[p.RestaurantID] = append(m.byRestaurant[p.RestaurantID], p)
	return nil
}

func (m *mockRepository) UpsertBatch(_ context.Context, payments []Payment) error {
	for i := range payments {
		if err := m.Upsert(nil, &payments[i]); err != nil {
			return err
		}
	}
	return nil
}

func TestService_GetByID(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	p := &Payment{ID: uuid.New(), CloverPaymentID: "clv_pay_1"}
	repo.payments[p.ID] = p

	got, err := svc.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.CloverPaymentID != "clv_pay_1" {
		t.Errorf("expected clv_pay_1, got %v", got)
	}
}

func TestService_UpsertBatch(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()
	restID := uuid.New()

	payments := []Payment{
		{ID: uuid.New(), CloverPaymentID: "clv_1", RestaurantID: restID, CreatedTime: time.Now().UTC()},
		{ID: uuid.New(), CloverPaymentID: "clv_2", RestaurantID: restID, CreatedTime: time.Now().UTC()},
	}

	if err := svc.UpsertBatch(ctx, payments); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.payments) != 2 {
		t.Errorf("expected 2 payments, got %d", len(repo.payments))
	}
}

func TestService_ListByRestaurantAndDate(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()
	restID := uuid.New()
	now := time.Now().UTC()

	repo.byRestaurant[restID] = []*Payment{
		{ID: uuid.New(), CloverPaymentID: "clv_1", RestaurantID: restID, CreatedTime: now},
	}

	got, err := svc.ListByRestaurantAndDate(ctx, restID, now.AddDate(0, 0, -1), now.AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 payment, got %d", len(got))
	}
}
