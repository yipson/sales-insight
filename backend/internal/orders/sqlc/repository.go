package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/orders"
)

// OrderSQLCRepository implements orders.Repository.
type OrderSQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewOrderSQLCRepository creates a new sqlc-backed order repository.
func NewOrderSQLCRepository(db *sql.DB) *OrderSQLCRepository {
	return &OrderSQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func orderToDomain(o Order) *orders.Order {
	ord := &orders.Order{
		ID:                  o.ID,
		RestaurantID:        o.RestaurantID,
		CloverOrderID:       o.CloverOrderID,
		State:               o.State,
		TotalAmount:         o.TotalAmount,
		TaxAmount:           o.TaxAmount,
		DiscountAmount:      o.DiscountAmount,
		TipAmount:           o.TipAmount,
		CreatedTime:         o.CreatedTime,
		ModifiedTime:        o.ModifiedTime,
		ItemCount:           o.ItemCount,
		UniqueCategoryCount: o.UniqueCategoryCount,
		TotalQuantity:       o.TotalQuantity,
		SyncedAt:            o.SyncedAt,
	}
	if o.EmployeeID.Valid {
		id := o.EmployeeID.UUID
		ord.EmployeeID = &id
	}
	if o.OrderNumber.Valid {
		ord.OrderNumber = o.OrderNumber.String
	}
	if o.OrderType.Valid {
		ord.OrderType = o.OrderType.String
	}
	if o.PayType.Valid {
		ord.PayType = o.PayType.String
	}
	return ord
}

func orderFromDomain(o *orders.Order) UpsertOrderParams {
	params := UpsertOrderParams{
		ID:                  o.ID,
		RestaurantID:        o.RestaurantID,
		CloverOrderID:       o.CloverOrderID,
		State:               o.State,
		TotalAmount:         o.TotalAmount,
		TaxAmount:           o.TaxAmount,
		DiscountAmount:      o.DiscountAmount,
		TipAmount:           o.TipAmount,
		CreatedTime:         o.CreatedTime,
		ModifiedTime:        o.ModifiedTime,
		ItemCount:           o.ItemCount,
		UniqueCategoryCount: o.UniqueCategoryCount,
		TotalQuantity:       o.TotalQuantity,
		SyncedAt:            o.SyncedAt,
	}
	if o.EmployeeID != nil {
		params.EmployeeID = uuid.NullUUID{UUID: *o.EmployeeID, Valid: true}
	}
	params.OrderNumber = sql.NullString{String: o.OrderNumber, Valid: o.OrderNumber != ""}
	params.OrderType = sql.NullString{String: o.OrderType, Valid: o.OrderType != ""}
	params.PayType = sql.NullString{String: o.PayType, Valid: o.PayType != ""}
	return params
}

func (r *OrderSQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return orderToDomain(row), nil
}

func (r *OrderSQLCRepository) GetByCloverOrderID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) (*orders.Order, error) {
	row, err := r.queries.GetOrderByCloverOrderID(ctx, GetOrderByCloverOrderIDParams{
		RestaurantID:  restaurantID,
		CloverOrderID: cloverOrderID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return orderToDomain(row), nil
}

func (r *OrderSQLCRepository) ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]orders.Order, error) {
	rows, err := r.queries.ListOrdersByRestaurantAndDate(ctx, ListOrdersByRestaurantAndDateParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	var out []orders.Order
	for _, row := range rows {
		out = append(out, *orderToDomain(row))
	}
	return out, nil
}

func (r *OrderSQLCRepository) ListByEmployee(ctx context.Context, employeeID uuid.UUID, from, to time.Time) ([]orders.Order, error) {
	rows, err := r.queries.ListOrdersByEmployeeAndDate(ctx, ListOrdersByEmployeeAndDateParams{
		EmployeeID:    uuid.NullUUID{UUID: employeeID, Valid: true},
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	var out []orders.Order
	for _, row := range rows {
		out = append(out, *orderToDomain(row))
	}
	return out, nil
}

func (r *OrderSQLCRepository) Upsert(ctx context.Context, order *orders.Order) error {
	_, err := r.queries.UpsertOrder(ctx, orderFromDomain(order))
	return err
}

func (r *OrderSQLCRepository) UpsertBatch(ctx context.Context, orders []orders.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range orders {
		if _, err := qtx.UpsertOrder(ctx, orderFromDomain(&orders[i])); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *OrderSQLCRepository) DeleteByCloverID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) error {
	return r.queries.DeleteOrderByCloverID(ctx, DeleteOrderByCloverIDParams{
		RestaurantID:  restaurantID,
		CloverOrderID: cloverOrderID,
	})
}

// ============================================================
// OrderItemSQLCRepository implements orders.OrderItemRepository.
// ============================================================

type OrderItemSQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewOrderItemSQLCRepository creates a new sqlc-backed order item repository.
func NewOrderItemSQLCRepository(db *sql.DB) *OrderItemSQLCRepository {
	return &OrderItemSQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func orderItemToDomain(oi OrderItem) *orders.OrderItem {
	item := &orders.OrderItem{
		ID:           oi.ID,
		RestaurantID: oi.RestaurantID,
		OrderID:      oi.OrderID,
		Name:         oi.Name,
		Quantity:     oi.Quantity,
		UnitPrice:    oi.UnitPrice,
		TotalPrice:   oi.TotalPrice,
		CreatedAt:    oi.CreatedAt,
		SyncedAt:     oi.SyncedAt,
	}
	if oi.ProductID.Valid {
		id := oi.ProductID.UUID
		item.ProductID = &id
	}
	if oi.CloverLineItemID.Valid {
		item.CloverLineItemID = oi.CloverLineItemID.String
	}
	if oi.AnalyticCategoryID.Valid {
		id := oi.AnalyticCategoryID.UUID
		item.AnalyticCategoryID = &id
	}
	return item
}

func orderItemFromDomain(oi *orders.OrderItem) UpsertOrderItemParams {
	params := UpsertOrderItemParams{
		ID:           oi.ID,
		RestaurantID: oi.RestaurantID,
		OrderID:      oi.OrderID,
		Name:         oi.Name,
		Quantity:     oi.Quantity,
		UnitPrice:    oi.UnitPrice,
		TotalPrice:   oi.TotalPrice,
		CreatedAt:    oi.CreatedAt,
		SyncedAt:     oi.SyncedAt,
	}
	if oi.ProductID != nil {
		params.ProductID = uuid.NullUUID{UUID: *oi.ProductID, Valid: true}
	}
	params.CloverLineItemID = sql.NullString{String: oi.CloverLineItemID, Valid: oi.CloverLineItemID != ""}
	if oi.AnalyticCategoryID != nil {
		params.AnalyticCategoryID = uuid.NullUUID{UUID: *oi.AnalyticCategoryID, Valid: true}
	}
	return params
}

func (r *OrderItemSQLCRepository) UpsertBatch(ctx context.Context, items []orders.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range items {
		if _, err := qtx.UpsertOrderItem(ctx, orderItemFromDomain(&items[i])); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *OrderItemSQLCRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]orders.OrderItem, error) {
	rows, err := r.queries.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	var out []orders.OrderItem
	for _, row := range rows {
		out = append(out, *orderItemToDomain(row))
	}
	return out, nil
}

// ============================================================
// CategorySummarySQLCRepository implements orders.CategorySummaryRepository.
// ============================================================

type CategorySummarySQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewCategorySummarySQLCRepository creates a new sqlc-backed category summary repository.
func NewCategorySummarySQLCRepository(db *sql.DB) *CategorySummarySQLCRepository {
	return &CategorySummarySQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func summaryToDomain(s OrderCategorySummary) *orders.OrderCategorySummary {
	return &orders.OrderCategorySummary{
		ID:                 s.ID,
		RestaurantID:       s.RestaurantID,
		OrderID:            s.OrderID,
		AnalyticCategoryID: s.AnalyticCategoryID,
		ItemCount:          s.ItemCount,
		TotalQuantity:      s.TotalQuantity,
		TotalAmount:        s.TotalAmount,
	}
}

func summaryFromDomain(s *orders.OrderCategorySummary) UpsertOrderCategorySummaryParams {
	return UpsertOrderCategorySummaryParams{
		ID:                 s.ID,
		RestaurantID:       s.RestaurantID,
		OrderID:            s.OrderID,
		AnalyticCategoryID: s.AnalyticCategoryID,
		ItemCount:          s.ItemCount,
		TotalQuantity:      s.TotalQuantity,
		TotalAmount:        s.TotalAmount,
	}
}

func (r *CategorySummarySQLCRepository) UpsertBatch(ctx context.Context, summaries []orders.OrderCategorySummary) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range summaries {
		if _, err := qtx.UpsertOrderCategorySummary(ctx, summaryFromDomain(&summaries[i])); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *CategorySummarySQLCRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]orders.OrderCategorySummary, error) {
	rows, err := r.queries.GetOrderCategorySummaryByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	var out []orders.OrderCategorySummary
	for _, row := range rows {
		out = append(out, *summaryToDomain(row))
	}
	return out, nil
}
