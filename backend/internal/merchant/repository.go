package merchant

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the contract for merchant persistence.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	GetByCloverMerchantID(ctx context.Context, cloverMerchantID string) (*Merchant, error)
	List(ctx context.Context) ([]Merchant, error)
	Create(ctx context.Context, m *Merchant) error
	Update(ctx context.Context, m *Merchant) error
	UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, accessExp, refreshExp *int64) error
	Delete(ctx context.Context, id uuid.UUID) error
}
