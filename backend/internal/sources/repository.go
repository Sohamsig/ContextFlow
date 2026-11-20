package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("source not found")

type Source struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(
	ctx context.Context,
	name string,
	sourceType string,
	config json.RawMessage,
) (*Source, error) {
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}

	var source Source
	var configBytes []byte

	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO sources (name, type, config)
		 VALUES ($1, $2, $3::jsonb)
		 RETURNING id, name, type, status, config, created_at, updated_at`,
		name,
		sourceType,
		string(config),
	).Scan(
		&source.ID,
		&source.Name,
		&source.Type,
		&source.Status,
		&configBytes,
		&source.CreatedAt,
		&source.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create source: %w", err)
	}

	source.Config = json.RawMessage(configBytes)
	return &source, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*Source, error) {
	var source Source
	var configBytes []byte

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, name, type, status, config, created_at, updated_at
		 FROM sources
		 WHERE id = $1`,
		id,
	).Scan(
		&source.ID,
		&source.Name,
		&source.Type,
		&source.Status,
		&configBytes,
		&source.CreatedAt,
		&source.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get source: %w", err)
	}

	source.Config = json.RawMessage(configBytes)
	return &source, nil
}

func (r *Repository) List(
	ctx context.Context,
) ([]Source, error) {
	rows, err := r.pool.Query(
		ctx,
		`SELECT id, name, type, status, config, created_at, updated_at
		 FROM sources
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list sources: %w", err)
	}
	defer rows.Close()

	sources := make([]Source, 0)

	for rows.Next() {
		var source Source
		var configBytes []byte

		if err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.Type,
			&source.Status,
			&configBytes,
			&source.CreatedAt,
			&source.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan source: %w", err)
		}

		source.Config = json.RawMessage(configBytes)
		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sources: %w", err)
	}

	return sources, nil
}
