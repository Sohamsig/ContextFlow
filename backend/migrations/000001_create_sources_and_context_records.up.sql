CREATE TABLE sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',

    config JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT sources_status_check
        CHECK (status IN ('active', 'disabled', 'error'))
);

CREATE TABLE context_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    source_id UUID NOT NULL REFERENCES sources(id)
        ON DELETE CASCADE,

    external_id VARCHAR(512) NOT NULL,

    content TEXT NOT NULL,
    content_type VARCHAR(100) NOT NULL DEFAULT 'text',

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    content_hash VARCHAR(64) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT context_records_source_external_id_unique
        UNIQUE (source_id, external_id)
);

CREATE INDEX idx_sources_status
    ON sources(status);

CREATE INDEX idx_context_records_source_id
    ON context_records(source_id);

CREATE INDEX idx_context_records_ingested_at
    ON context_records(ingested_at DESC);
