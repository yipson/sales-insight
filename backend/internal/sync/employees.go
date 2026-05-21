package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sales-insight/backend/internal/employees"
)

// SyncEmployees extracts employees from Clover API.
// It returns the count of employees processed.
func (e *Engine) SyncEmployees(ctx context.Context, merchantID uuid.UUID) (int, error) {
	token, err := e.getToken(merchantID)
	if err != nil {
		return 0, fmt.Errorf("get token: %w", err)
	}

	cloverMerchantID, err := e.getCloverMerchantID(ctx, merchantID)
	if err != nil {
		return 0, fmt.Errorf("get clover merchant id: %w", err)
	}

	body, err := e.cloverClient.Get(ctx, fmt.Sprintf("/merchants/%s/employees", cloverMerchantID), token)
	if err != nil {
		return 0, fmt.Errorf("fetch employees: %w", err)
	}

	var wrapper struct {
		Elements []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Role   string `json:"role"`
			Active bool   `json:"active"`
		} `json:"elements"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return 0, fmt.Errorf("unmarshal employees: %w", err)
	}

	if len(wrapper.Elements) == 0 {
		return 0, nil
	}

	var domainEmployees []employees.Employee
	for _, emp := range wrapper.Elements {
		domainEmployees = append(domainEmployees, employees.Employee{
			ID:               uuid.New(),
			RestaurantID:     merchantID,
			CloverEmployeeID: emp.ID,
			Name:             emp.Name,
			Role:             emp.Role,
			IsActive:         emp.Active,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
			SyncedAt:         time.Now().UTC(),
		})
	}

	if err := e.upsertEmployees(ctx, domainEmployees); err != nil {
		return 0, fmt.Errorf("upsert employees: %w", err)
	}

	return len(domainEmployees), nil
}

func (e *Engine) upsertEmployees(ctx context.Context, domainEmployees []employees.Employee) error {
	for i := 0; i < len(domainEmployees); i += e.batchSize {
		end := i + e.batchSize
		if end > len(domainEmployees) {
			end = len(domainEmployees)
		}
		batch := domainEmployees[i:end]
		if err := e.employeeRepo.UpsertBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}
