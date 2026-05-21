package sync

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/employees"
	"github.com/sales-insight/backend/internal/merchant"
	"github.com/sales-insight/backend/internal/orders"
	"github.com/sales-insight/backend/internal/payments"
	"github.com/sales-insight/backend/internal/products"
	"github.com/sales-insight/backend/internal/token_cache"
)

// Engine coordinates data extraction from Clover API and persistence.
// It is designed as a separable component for future worker extraction.
type Engine struct {
	cloverClient  *clover.Client
	tokenCache    *tokencache.Cache
	merchantRepo  merchant.Repository
	batchSize     int

	// Domain repositories (interfaces — will be implemented in Phase 4)
	orderRepo       orders.Repository
	orderItemRepo   orders.OrderItemRepository
	categorySummaryRepo orders.CategorySummaryRepository
	productRepo     products.ProductRepository
	categoryRepo    products.CategoryRepository
	employeeRepo    employees.Repository
	paymentRepo     payments.Repository
}

// NewEngine creates a new sync engine.
func NewEngine(
	cloverClient *clover.Client,
	tokenCache *tokencache.Cache,
	merchantRepo merchant.Repository,
	orderRepo orders.Repository,
	orderItemRepo orders.OrderItemRepository,
	categorySummaryRepo orders.CategorySummaryRepository,
	productRepo products.ProductRepository,
	categoryRepo products.CategoryRepository,
	employeeRepo employees.Repository,
	paymentRepo payments.Repository,
) *Engine {
	return &Engine{
		cloverClient:        cloverClient,
		tokenCache:          tokenCache,
		merchantRepo:        merchantRepo,
		orderRepo:           orderRepo,
		orderItemRepo:       orderItemRepo,
		categorySummaryRepo: categorySummaryRepo,
		productRepo:         productRepo,
		categoryRepo:        categoryRepo,
		employeeRepo:        employeeRepo,
		paymentRepo:         paymentRepo,
		batchSize:           100,
	}
}

// SetBatchSize configures the batch size for transactional writes (default: 100).
func (e *Engine) SetBatchSize(size int) {
	e.batchSize = size
}

// getToken returns the access token for a merchant from cache.
func (e *Engine) getToken(merchantID uuid.UUID) (string, error) {
	td, ok := e.tokenCache.Get(merchantID)
	if !ok {
		return "", fmt.Errorf("no token found for merchant %s", merchantID)
	}
	return td.AccessToken, nil
}

// getCloverMerchantID resolves the Clover merchant ID from our internal merchant record.
func (e *Engine) getCloverMerchantID(ctx context.Context, merchantID uuid.UUID) (string, error) {
	m, err := e.merchantRepo.GetByID(ctx, merchantID)
	if err != nil {
		return "", fmt.Errorf("lookup merchant: %w", err)
	}
	if m == nil {
		return "", fmt.Errorf("merchant not found")
	}
	if m.CloverMerchantID == "" {
		return "", fmt.Errorf("merchant has no clover_merchant_id")
	}
	return m.CloverMerchantID, nil
}
