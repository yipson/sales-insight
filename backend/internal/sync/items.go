package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/products"
)

// SyncItems extracts products and native categories from Clover API.
// It returns the count of products and categories processed.
func (e *Engine) SyncItems(ctx context.Context, merchantID uuid.UUID) (int, int, error) {
	token, err := e.getToken(merchantID)
	if err != nil {
		return 0, 0, fmt.Errorf("get token: %w", err)
	}

	cloverMerchantID, err := e.getCloverMerchantID(ctx, merchantID)
	if err != nil {
		return 0, 0, fmt.Errorf("get clover merchant id: %w", err)
	}

	// Fetch categories
	catBody, err := e.cloverClient.Get(ctx, fmt.Sprintf("/merchants/%s/categories", cloverMerchantID), token)
	if err != nil {
		return 0, 0, fmt.Errorf("fetch categories: %w", err)
	}

	var catWrapper struct {
		Elements []clover.CategoryResponse `json:"elements"`
	}
	if err := json.Unmarshal(catBody, &catWrapper); err != nil {
		return 0, 0, fmt.Errorf("unmarshal categories: %w", err)
	}

	// Fetch items (products)
	itemBody, err := e.cloverClient.Get(ctx, fmt.Sprintf("/merchants/%s/items?expand=categories", cloverMerchantID), token)
	if err != nil {
		return 0, 0, fmt.Errorf("fetch items: %w", err)
	}

	var itemWrapper struct {
		Elements []clover.ItemResponse `json:"elements"`
	}
	if err := json.Unmarshal(itemBody, &itemWrapper); err != nil {
		return 0, 0, fmt.Errorf("unmarshal items: %w", err)
	}

	// Transform categories
	var domainCategories []products.Category
	for _, c := range catWrapper.Elements {
		domainCategories = append(domainCategories, products.Category{
			ID:               uuid.New(),
			RestaurantID:     merchantID,
			CloverCategoryID: c.ID,
			Name:             c.Name,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
			SyncedAt:         time.Now().UTC(),
		})
	}

	// Upsert categories in batches
	if len(domainCategories) > 0 {
		if err := e.upsertCategories(ctx, domainCategories); err != nil {
			return 0, 0, fmt.Errorf("upsert categories: %w", err)
		}
	}

	// Transform products and build category mappings
	var domainProducts []products.Product
	for _, item := range itemWrapper.Elements {
		product := products.Product{
			ID:           uuid.New(),
			RestaurantID: merchantID,
			CloverItemID: item.ID,
			Name:         item.Name,
			IsAvailable:  true,
			IsDeleted:    false,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			SyncedAt:     time.Now().UTC(),
		}

		// Map first category if available (Clover items can have multiple categories)
		if len(item.Categories) > 0 {
			for _, cat := range item.Categories {
				// Find matching domain category by Clover ID
				for i := range domainCategories {
					if domainCategories[i].CloverCategoryID == cat.ID {
						product.CategoryID = &domainCategories[i].ID
						break
					}
				}
				if product.CategoryID != nil {
					break
				}
			}
		}

		domainProducts = append(domainProducts, product)
	}

	// Upsert products in batches
	if len(domainProducts) > 0 {
		if err := e.upsertProducts(ctx, domainProducts); err != nil {
			return 0, 0, fmt.Errorf("upsert products: %w", err)
		}
	}

	// TODO: Soft delete products not present in this sync ( Phase 4 enhancement )

	return len(domainProducts), len(domainCategories), nil
}

func (e *Engine) upsertCategories(ctx context.Context, domainCategories []products.Category) error {
	for i := 0; i < len(domainCategories); i += e.batchSize {
		end := i + e.batchSize
		if end > len(domainCategories) {
			end = len(domainCategories)
		}
		batch := domainCategories[i:end]
		if err := e.categoryRepo.UpsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (e *Engine) upsertProducts(ctx context.Context, domainProducts []products.Product) error {
	for i := 0; i < len(domainProducts); i += e.batchSize {
		end := i + e.batchSize
		if end > len(domainProducts) {
			end = len(domainProducts)
		}
		batch := domainProducts[i:end]
		if err := e.productRepo.UpsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}
