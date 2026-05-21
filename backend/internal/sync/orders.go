package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/orders"
)

// SyncOrders extracts orders from Clover API incrementally.
// It returns the new cursor (highest modifiedTime seen) and count of orders processed.
func (e *Engine) SyncOrders(ctx context.Context, merchantID uuid.UUID, cursor time.Time) (time.Time, int, error) {
	token, err := e.getToken(merchantID)
	if err != nil {
		return cursor, 0, fmt.Errorf("get token: %w", err)
	}

	cloverMerchantID, err := e.getCloverMerchantID(ctx, merchantID)
	if err != nil {
		return cursor, 0, fmt.Errorf("get clover merchant id: %w", err)
	}

	path := fmt.Sprintf("/merchants/%s/orders?filter=modifiedTime>=%d&expand=lineItems", 
		cloverMerchantID, cursor.UnixMilli())

	body, err := e.cloverClient.Get(ctx, path, token)
	if err != nil {
		return cursor, 0, fmt.Errorf("fetch orders: %w", err)
	}

	var wrapper struct {
		Elements []clover.OrderResponse `json:"elements"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return cursor, 0, fmt.Errorf("unmarshal orders: %w", err)
	}

	if len(wrapper.Elements) == 0 {
		return cursor, 0, nil
	}

	newCursor := cursor
	var domainOrders []orders.Order
	var domainItems []orders.OrderItem

	for _, o := range wrapper.Elements {
		order, items := e.transformOrder(merchantID, o)
		domainOrders = append(domainOrders, order)
		domainItems = append(domainItems, items...)

		orderModified := time.UnixMilli(o.ModifiedTime)
		if orderModified.After(newCursor) {
			newCursor = orderModified
		}
	}

	// Upsert orders in batches
	if err := e.upsertOrders(ctx, domainOrders); err != nil {
		return cursor, 0, fmt.Errorf("upsert orders: %w", err)
	}

	// Upsert order items in batches
	if err := e.upsertOrderItems(ctx, domainItems); err != nil {
		return cursor, 0, fmt.Errorf("upsert order items: %w", err)
	}

	return newCursor, len(domainOrders), nil
}

func (e *Engine) transformOrder(merchantID uuid.UUID, o clover.OrderResponse) (orders.Order, []orders.OrderItem) {
	order := orders.Order{
		ID:            uuid.New(),
		RestaurantID:  merchantID,
		CloverOrderID: o.ID,
		TotalAmount:   o.Total,
		CreatedTime:   time.UnixMilli(o.ModifiedTime).UTC(),
		ModifiedTime:  time.UnixMilli(o.ModifiedTime).UTC(),
		SyncedAt:      time.Now().UTC(),
	}

	// Map employee if present
	if o.Employee != nil && o.Employee.ID != "" {
		// Employee mapping will be resolved later; store nil for now
		// In production, lookup employeeRepo.GetByCloverEmployeeID
		_ = o.Employee.ID
	}

	// Process line items and compute aggregates
	var domainItems []orders.OrderItem
	categorySet := make(map[string]struct{})
	var totalQuantity int32

	for _, li := range o.LineItems {
		item := orders.OrderItem{
			ID:               uuid.New(),
			RestaurantID:     merchantID,
			OrderID:          order.ID,
			CloverLineItemID: li.ID,
			Name:             li.Name,
			Quantity:         int32(li.Quantity),
			UnitPrice:        li.Price,
			TotalPrice:       li.Price * int64(li.Quantity),
			CreatedAt:        order.CreatedTime,
			SyncedAt:         time.Now().UTC(),
		}

		if li.Item != nil {
			item.CloverLineItemID = li.Item.ID
			for _, cat := range li.Item.Categories {
				categorySet[cat.ID] = struct{}{}
			}
		}

		totalQuantity += item.Quantity
		domainItems = append(domainItems, item)
	}

	order.ItemCount = int32(len(o.LineItems))
	order.TotalQuantity = totalQuantity
	order.UniqueCategoryCount = int32(len(categorySet))

	return order, domainItems
}

func (e *Engine) upsertOrders(ctx context.Context, domainOrders []orders.Order) error {
	for i := 0; i < len(domainOrders); i += e.batchSize {
		end := i + e.batchSize
		if end > len(domainOrders) {
			end = len(domainOrders)
		}
		batch := domainOrders[i:end]
		if err := e.orderRepo.UpsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (e *Engine) upsertOrderItems(ctx context.Context, domainItems []orders.OrderItem) error {
	for i := 0; i < len(domainItems); i += e.batchSize {
		end := i + e.batchSize
		if end > len(domainItems) {
			end = len(domainItems)
		}
		batch := domainItems[i:end]
		if err := e.orderItemRepo.UpsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}
