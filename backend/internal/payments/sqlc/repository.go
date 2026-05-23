package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/payments"
)

// SQLCRepository implements payments.Repository using sqlc-generated code.
type SQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewSQLCRepository creates a new sqlc-backed payment repository.
func NewSQLCRepository(db *sql.DB) *SQLCRepository {
	return &SQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func toDomain(p Payment) *payments.Payment {
	pay := &payments.Payment{
		ID:              p.ID,
		RestaurantID:    p.RestaurantID,
		CloverPaymentID: p.CloverPaymentID,
		Amount:          p.Amount,
		TipAmount:       p.TipAmount,
		TaxAmount:       p.TaxAmount,
		Result:          p.Result,
		CreatedTime:     p.CreatedTime,
		ModifiedTime:    p.ModifiedTime,
		SyncedAt:        p.SyncedAt,
	}
	if p.OrderID.Valid {
		id := p.OrderID.UUID
		pay.OrderID = &id
	}
	if p.EmployeeID.Valid {
		id := p.EmployeeID.UUID
		pay.EmployeeID = &id
	}
	if p.PaymentType.Valid {
		pay.PaymentType = p.PaymentType.String
	}
	if p.CardType.Valid {
		pay.CardType = p.CardType.String
	}
	if p.ExternalPaymentID.Valid {
		pay.ExternalPaymentID = p.ExternalPaymentID.String
	}
	return pay
}

func fromDomain(p *payments.Payment) UpsertPaymentParams {
	params := UpsertPaymentParams{
		ID:              p.ID,
		RestaurantID:    p.RestaurantID,
		CloverPaymentID: p.CloverPaymentID,
		Amount:          p.Amount,
		TipAmount:       p.TipAmount,
		TaxAmount:       p.TaxAmount,
		Result:          p.Result,
		CreatedTime:     p.CreatedTime,
		ModifiedTime:    p.ModifiedTime,
		SyncedAt:        p.SyncedAt,
	}
	if p.OrderID != nil {
		params.OrderID = uuid.NullUUID{UUID: *p.OrderID, Valid: true}
	}
	if p.EmployeeID != nil {
		params.EmployeeID = uuid.NullUUID{UUID: *p.EmployeeID, Valid: true}
	}
	params.PaymentType = sql.NullString{String: p.PaymentType, Valid: p.PaymentType != ""}
	params.CardType = sql.NullString{String: p.CardType, Valid: p.CardType != ""}
	params.ExternalPaymentID = sql.NullString{String: p.ExternalPaymentID, Valid: p.ExternalPaymentID != ""}
	return params
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*payments.Payment, error) {
	row, err := r.queries.GetPaymentByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}

func (r *SQLCRepository) GetByCloverPaymentID(ctx context.Context, restaurantID uuid.UUID, cloverPaymentID string) (*payments.Payment, error) {
	row, err := r.queries.GetPaymentByCloverID(ctx, GetPaymentByCloverIDParams{
		RestaurantID:    restaurantID,
		CloverPaymentID: cloverPaymentID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}

func (r *SQLCRepository) ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]payments.Payment, error) {
	rows, err := r.queries.ListPaymentsByRestaurantAndDate(ctx, ListPaymentsByRestaurantAndDateParams{
		RestaurantID:  restaurantID,
		CreatedTime:   from,
		CreatedTime_2: to,
	})
	if err != nil {
		return nil, err
	}
	var out []payments.Payment
	for _, row := range rows {
		out = append(out, *toDomain(row))
	}
	return out, nil
}

func (r *SQLCRepository) Upsert(ctx context.Context, payment *payments.Payment) error {
	_, err := r.queries.UpsertPayment(ctx, fromDomain(payment))
	return err
}

func (r *SQLCRepository) UpsertBatch(ctx context.Context, payments []payments.Payment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range payments {
		if _, err := qtx.UpsertPayment(ctx, fromDomain(&payments[i])); err != nil {
			return err
		}
	}
	return tx.Commit()
}
