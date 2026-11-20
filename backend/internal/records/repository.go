package records

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("context record not found")

type Record struct {
	ID          string          `json:"id"`
	SourceID    string          `json:"source_id"`
	ExternalID  string          `json:"external_id"`
	Content     string          `json:"content"`
	ContentType string          `json:"content_type"`
	Metadata    json.RawMessage `json:"metadata"`
	ContentHash string          `json:"content_hash"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	IngestedAt  time.Time       `json:"ingested_at"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Upsert(
	ctx context.Context,
	sourceID string,
	externalID string,
	content string,
	contentType string,
	metadata json.RawMessage,
) (*Record, error) {
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	if !json.Valid(metadata) {
		return nil, errors.New("metadata must be valid JSON")
	}

	var metadataObject map[string]json.RawMessage
	if err := json.Unmarshal(metadata, &metadataObject); err != nil || metadataObject == nil {
		return nil, errors.New("metadata must be a JSON object")
	}

	hash := sha256.Sum256([]byte(content))
	contentHash := hex.EncodeToString(hash[:])

	var record Record
	var metadataBytes []byte

	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO context_records (
			source_id,
			external_id,
			content,
			content_type,
			metadata,
			content_hash
		)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6)
		ON CONFLICT (source_id, external_id)
		DO UPDATE SET
			content = EXCLUDED.content,
			content_type = EXCLUDED.content_type,
			metadata = EXCLUDED.metadata,
			content_hash = EXCLUDED.content_hash,
			updated_at = NOW(),
			ingested_at = NOW()
		WHERE context_records.content_hash IS DISTINCT FROM EXCLUDED.content_hash
		   OR context_records.content_type IS DISTINCT FROM EXCLUDED.content_type
		   OR context_records.metadata IS DISTINCT FROM EXCLUDED.metadata
		RETURNING
			id,
			source_id,
			external_id,
			content,
			content_type,
			metadata,
			content_hash,
			created_at,
			updated_at,
			ingested_at`,
		sourceID,
		externalID,
		content,
		contentType,
		string(metadata),
		contentHash,
	).Scan(
		&record.ID,
		&record.SourceID,
		&record.ExternalID,
		&record.Content,
		&record.ContentType,
		&metadataBytes,
		&record.ContentHash,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.IngestedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		// The record already existed with identical content and metadata.
		// Fetch it without modifying its timestamps.
		return r.GetByExternalID(ctx, sourceID, externalID)
	}

	if err != nil {
		return nil, fmt.Errorf("upsert context record: %w", err)
	}

	record.Metadata = json.RawMessage(metadataBytes)
	return &record, nil
}

func (r *Repository) GetByExternalID(
	ctx context.Context,
	sourceID string,
	externalID string,
) (*Record, error) {
	var record Record
	var metadataBytes []byte

	err := r.pool.QueryRow(
		ctx,
		`SELECT
			id,
			source_id,
			external_id,
			content,
			content_type,
			metadata,
			content_hash,
			created_at,
			updated_at,
			ingested_at
		FROM context_records
		WHERE source_id = $1 AND external_id = $2`,
		sourceID,
		externalID,
	).Scan(
		&record.ID,
		&record.SourceID,
		&record.ExternalID,
		&record.Content,
		&record.ContentType,
		&metadataBytes,
		&record.ContentHash,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.IngestedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get context record: %w", err)
	}

	record.Metadata = json.RawMessage(metadataBytes)
	return &record, nil
}
