package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/payments"
)

// SyncPayments extracts payments from Clover API incrementally.
// It returns the new cursor (highest modifiedTime seen) and count of payments processed.
func (e *Engine) SyncPayments(ctx context.Context, merchantID uuid.UUID, cursor time.Time) (time.Time, int, error) {
	token, err := e.getToken(merchantID)
	if err != nil {
		return cursor, 0, fmt.Errorf("get token: %w", err)
	}

	cloverMerchantID, err := e.getCloverMerchantID(ctx, merchantID)
	if err != nil {
		return cursor, 0, fmt.Errorf("get clover merchant id: %w", err)
	}

	path := fmt.Sprintf("/merchants/%s/payments?filter=modifiedTime>=%d", cloverMerchantID, cursor.UnixMilli())

	body, err := e.cloverClient.Get(ctx, path, token)
	if err != nil {
		return cursor, 0, fmt.Errorf("fetch payments: %w", err)
	}

	var wrapper struct {
		Elements []struct {
			ID        string `json:"id"`
			Amount    int64  `json:"amount"`
			TipAmount int64  `json:"tipAmount"`
			TaxAmount int64  `json:"taxAmount"`
			Order struct {
				ID string `json:"id"`
			} `json:"order"`
			Employee struct {
				ID string `json:"id"`
			} `json:"employee"`
			PaymentType       string `json:"tender"`
			CardType          string `json:"cardType"`
			Result            string `json:"result"`
			ExternalPaymentID string `json:"externalPaymentId"`
			CreatedTime       int64  `json:"createdTime"`
			ModifiedTime      int64  `json:"modifiedTime"`
		} `json:"elements"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return cursor, 0, fmt.Errorf("unmarshal payments: %w", err)
	}

	if len(wrapper.Elements) == 0 {
		return cursor, 0, nil
	}

	newCursor := cursor
	var domainPayments []payments.Payment

	for _, p := range wrapper.Elements {
		payment := payments.Payment{
			ID:              uuid.New(),
			RestaurantID:    merchantID,
			CloverPaymentID: p.ID,
			Amount:          p.Amount,
			TipAmount:       p.TipAmount,
			TaxAmount:       p.TaxAmount,
			PaymentType:     p.PaymentType,
			CardType:        p.CardType,
			Result:          p.Result,
			ExternalPaymentID: p.ExternalPaymentID,
			CreatedTime:     time.UnixMilli(p.CreatedTime).UTC(),
			ModifiedTime:    time.UnixMilli(p.ModifiedTime).UTC(),
			SyncedAt:        time.Now().UTC(),
		}

		// Resolve internal UUIDs from Clover external IDs
		if p.Order.ID != "" {
			payment.OrderID = e.resolveOrderID(ctx, merchantID, p.Order.ID)
		}
		if p.Employee.ID != "" {
			payment.EmployeeID = e.resolveEmployeeID(ctx, merchantID, p.Employee.ID)
		}

		domainPayments = append(domainPayments, payment)

		paymentModified := time.UnixMilli(p.ModifiedTime)
		if paymentModified.After(newCursor) {
			newCursor = paymentModified
		}
	}

	if err := e.upsertPayments(ctx, domainPayments); err != nil {
		return cursor, 0, fmt.Errorf("upsert payments: %w", err)
	}

	return newCursor, len(domainPayments), nil
}

func (e *Engine) upsertPayments(ctx context.Context, domainPayments []payments.Payment) error {
	for i := 0; i < len(domainPayments); i += e.batchSize {
		end := i + e.batchSize
		if end > len(domainPayments) {
			end = len(domainPayments)
		}
		batch := domainPayments[i:end]
		if err := e.paymentRepo.UpsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}
