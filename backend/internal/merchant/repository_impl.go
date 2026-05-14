package merchant

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/platform/security"
)

// SQLRepository implements Repository using database/sql.
// This is a manual implementation; will be replaced by sqlc-generated code.
type SQLRepository struct {
	db        *sql.DB
	encrypter *security.Encrypter
}

// NewSQLRepository creates a new SQL-backed merchant repository.
func NewSQLRepository(db *sql.DB, encrypter *security.Encrypter) *SQLRepository {
	return &SQLRepository{db: db, encrypter: encrypter}
}

// scanMerchant scans a single merchant row.
func (r *SQLRepository) scanMerchant(row *sql.Row) (*Merchant, error) {
	var m Merchant
	var accessTokenEnc, refreshTokenEnc []byte
	var rulesRaw []byte
	var accessExp, refreshExp sql.NullInt64

	err := row.Scan(
		&m.ID, &m.Name, &m.CloverMerchantID,
		&accessTokenEnc, &refreshTokenEnc,
		&accessExp, &refreshExp,
		&m.CloverEnv, &m.IsConnected, &m.IsActive,
		&rulesRaw, &m.CreatedAt, &m.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if accessExp.Valid {
		t := time.Unix(accessExp.Int64, 0)
		m.AccessTokenExpiresAt = &t
	}
	if refreshExp.Valid {
		t := time.Unix(refreshExp.Int64, 0)
		m.RefreshTokenExpiresAt = &t
	}
	if len(rulesRaw) > 0 {
		_ = json.Unmarshal(rulesRaw, &m.TicketCompletoRules)
	}

	return &m, nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, clover_merchant_id,
			access_token, refresh_token,
			access_token_expires_at, refresh_token_expires_at,
			clover_env, is_connected, is_active,
			ticket_completo_rules, created_at, updated_at
		FROM restaurants WHERE id = $1`, id)
	return r.scanMerchant(row)
}

func (r *SQLRepository) GetByCloverMerchantID(ctx context.Context, cloverMerchantID string) (*Merchant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, clover_merchant_id,
			access_token, refresh_token,
			access_token_expires_at, refresh_token_expires_at,
			clover_env, is_connected, is_active,
			ticket_completo_rules, created_at, updated_at
		FROM restaurants WHERE clover_merchant_id = $1`, cloverMerchantID)
	return r.scanMerchant(row)
}

func (r *SQLRepository) List(ctx context.Context) ([]Merchant, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, clover_merchant_id,
			access_token, refresh_token,
			access_token_expires_at, refresh_token_expires_at,
			clover_env, is_connected, is_active,
			ticket_completo_rules, created_at, updated_at
		FROM restaurants WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []Merchant
	for rows.Next() {
		var m Merchant
		var accessTokenEnc, refreshTokenEnc []byte
		var rulesRaw []byte
		var accessExp, refreshExp sql.NullInt64

		if err := rows.Scan(
			&m.ID, &m.Name, &m.CloverMerchantID,
			&accessTokenEnc, &refreshTokenEnc,
			&accessExp, &refreshExp,
			&m.CloverEnv, &m.IsConnected, &m.IsActive,
			&rulesRaw, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if accessExp.Valid {
			t := time.Unix(accessExp.Int64, 0)
			m.AccessTokenExpiresAt = &t
		}
		if refreshExp.Valid {
			t := time.Unix(refreshExp.Int64, 0)
			m.RefreshTokenExpiresAt = &t
		}
		if len(rulesRaw) > 0 {
			_ = json.Unmarshal(rulesRaw, &m.TicketCompletoRules)
		}
		merchants = append(merchants, m)
	}
	return merchants, rows.Err()
}

func (r *SQLRepository) Create(ctx context.Context, m *Merchant) error {
	var rulesRaw []byte
	if m.TicketCompletoRules != nil {
		var err error
		rulesRaw, err = json.Marshal(m.TicketCompletoRules)
		if err != nil {
			return fmt.Errorf("marshal rules: %w", err)
		}
	}

	var accessExp, refreshExp *int64
	if m.AccessTokenExpiresAt != nil {
		v := m.AccessTokenExpiresAt.Unix()
		accessExp = &v
	}
	if m.RefreshTokenExpiresAt != nil {
		v := m.RefreshTokenExpiresAt.Unix()
		refreshExp = &v
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO restaurants (
			id, name, clover_merchant_id,
			access_token, refresh_token,
			access_token_expires_at, refresh_token_expires_at,
			clover_env, is_connected, is_active,
			ticket_completo_rules, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		m.ID, m.Name, m.CloverMerchantID,
		sql.NullString{String: m.AccessToken, Valid: m.AccessToken != ""},
		sql.NullString{String: m.RefreshToken, Valid: m.RefreshToken != ""},
		accessExp, refreshExp,
		m.CloverEnv, m.IsConnected, m.IsActive,
		rulesRaw, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *SQLRepository) Update(ctx context.Context, m *Merchant) error {
	var rulesRaw []byte
	if m.TicketCompletoRules != nil {
		var err error
		rulesRaw, err = json.Marshal(m.TicketCompletoRules)
		if err != nil {
			return fmt.Errorf("marshal rules: %w", err)
		}
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE restaurants SET
			name = $2,
			clover_merchant_id = $3,
			clover_env = $4,
			is_connected = $5,
			is_active = $6,
			ticket_completo_rules = $7,
			updated_at = $8
		WHERE id = $1`,
		m.ID, m.Name, m.CloverMerchantID,
		m.CloverEnv, m.IsConnected, m.IsActive,
		rulesRaw, time.Now().UTC(),
	)
	return err
}

func (r *SQLRepository) UpdateTokens(ctx context.Context, id uuid.UUID, accessToken, refreshToken string, accessExp, refreshExp *int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE restaurants SET
			access_token = $2,
			refresh_token = $3,
			access_token_expires_at = $4,
			refresh_token_expires_at = $5,
			updated_at = $6
		WHERE id = $1`,
		id,
		sql.NullString{String: accessToken, Valid: accessToken != ""},
		sql.NullString{String: refreshToken, Valid: refreshToken != ""},
		accessExp, refreshExp,
		time.Now().UTC(),
	)
	return err
}

func (r *SQLRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE restaurants SET is_active = false, updated_at = $2 WHERE id = $1`, id, time.Now().UTC())
	return err
}
