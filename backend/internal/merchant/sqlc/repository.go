package sqlc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/merchant"
	"github.com/sales-insight/backend/internal/platform/security"
	"github.com/sqlc-dev/pqtype"
)

// SQLCRepository implements merchant.Repository using sqlc-generated code.
type SQLCRepository struct {
	queries   *Queries
	encrypter *security.Encrypter
}

// NewSQLCRepository creates a new sqlc-backed merchant repository.
func NewSQLCRepository(db DBTX, encrypter *security.Encrypter) *SQLCRepository {
	return &SQLCRepository{
		queries:   New(db),
		encrypter: encrypter,
	}
}

// toDomain converts a sqlc Restaurant to a domain Merchant.
func toDomain(r Restaurant) *merchant.Merchant {
	m := &merchant.Merchant{
		ID:          r.ID,
		CloverEnv:   r.CloverEnv,
		IsConnected: r.IsConnected,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	if r.Name.Valid {
		m.Name = r.Name.String
	}
	if r.CloverMerchantID.Valid {
		m.CloverMerchantID = r.CloverMerchantID.String
	}
	if len(r.AccessToken) > 0 {
		m.AccessToken = string(r.AccessToken)
	}
	if len(r.RefreshToken) > 0 {
		m.RefreshToken = string(r.RefreshToken)
	}
	if r.AccessTokenExpiresAt.Valid {
		t := r.AccessTokenExpiresAt.Time
		m.AccessTokenExpiresAt = &t
	}
	if r.RefreshTokenExpiresAt.Valid {
		t := r.RefreshTokenExpiresAt.Time
		m.RefreshTokenExpiresAt = &t
	}
	if r.TicketCompletoRules.Valid {
		_ = json.Unmarshal(r.TicketCompletoRules.RawMessage, &m.TicketCompletoRules)
	}

	return m
}

// fromDomain converts a domain Merchant to sqlc CreateMerchantParams.
func fromDomain(m *merchant.Merchant) CreateMerchantParams {
	params := CreateMerchantParams{
		ID:          m.ID,
		CloverEnv:   m.CloverEnv,
		IsConnected: m.IsConnected,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	if m.Name != "" {
		params.Name = sql.NullString{String: m.Name, Valid: true}
	}
	if m.CloverMerchantID != "" {
		params.CloverMerchantID = sql.NullString{String: m.CloverMerchantID, Valid: true}
	}
	if m.AccessToken != "" {
		params.AccessToken = []byte(m.AccessToken)
	}
	if m.RefreshToken != "" {
		params.RefreshToken = []byte(m.RefreshToken)
	}
	if m.AccessTokenExpiresAt != nil {
		params.AccessTokenExpiresAt = sql.NullTime{Time: *m.AccessTokenExpiresAt, Valid: true}
	}
	if m.RefreshTokenExpiresAt != nil {
		params.RefreshTokenExpiresAt = sql.NullTime{Time: *m.RefreshTokenExpiresAt, Valid: true}
	}
	if m.TicketCompletoRules != nil {
		raw, _ := json.Marshal(m.TicketCompletoRules)
		params.TicketCompletoRules = pqtype.NullRawMessage{RawMessage: raw, Valid: true}
	}

	return params
}

func (r *SQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*merchant.Merchant, error) {
	rest, err := r.queries.GetMerchantByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(rest), nil
}

func (r *SQLCRepository) GetByCloverMerchantID(ctx context.Context, cloverMerchantID string) (*merchant.Merchant, error) {
	rest, err := r.queries.GetMerchantByCloverMerchantID(ctx, sql.NullString{String: cloverMerchantID, Valid: true})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomain(rest), nil
}

func (r *SQLCRepository) List(ctx context.Context) ([]merchant.Merchant, error) {
	rows, err := r.queries.ListMerchants(ctx)
	if err != nil {
		return nil, err
	}

	var merchants []merchant.Merchant
	for _, row := range rows {
		merchants = append(merchants, *toDomain(row))
	}
	return merchants, nil
}

func (r *SQLCRepository) Create(ctx context.Context, m *merchant.Merchant) error {
	return r.queries.CreateMerchant(ctx, fromDomain(m))
}

func (r *SQLCRepository) Update(ctx context.Context, m *merchant.Merchant) error {
	params := UpdateMerchantParams{
		ID:        m.ID,
		CloverEnv: m.CloverEnv,
		IsConnected: m.IsConnected,
		IsActive:    m.IsActive,
		UpdatedAt:   time.Now().UTC(),
	}

	if m.Name != "" {
		params.Name = sql.NullString{String: m.Name, Valid: true}
	}
	if m.CloverMerchantID != "" {
		params.CloverMerchantID = sql.NullString{String: m.CloverMerchantID, Valid: true}
	}
	if m.TicketCompletoRules != nil {
		raw, _ := json.Marshal(m.TicketCompletoRules)
		params.TicketCompletoRules = pqtype.NullRawMessage{RawMessage: raw, Valid: true}
	}

	return r.queries.UpdateMerchant(ctx, params)
}

func (r *SQLCRepository) UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, accessExp, refreshExp *int64) error {
	params := UpdateMerchantTokensParams{
		ID:          id,
		AccessToken: []byte(accessToken),
		UpdatedAt:   time.Now().UTC(),
	}

	if refreshToken != "" {
		params.RefreshToken = []byte(refreshToken)
	}
	if accessExp != nil {
		params.AccessTokenExpiresAt = sql.NullTime{Time: time.Unix(*accessExp, 0), Valid: true}
	}
	if refreshExp != nil {
		params.RefreshTokenExpiresAt = sql.NullTime{Time: time.Unix(*refreshExp, 0), Valid: true}
	}

	return r.queries.UpdateMerchantTokens(ctx, params)
}

func (r *SQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteMerchant(ctx, DeleteMerchantParams{
		ID:        id,
		UpdatedAt: time.Now().UTC(),
	})
}
