CREATE TABLE IF NOT EXISTS model_vendors (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    provider_key VARCHAR(80) NOT NULL DEFAULT '',
    icon_key VARCHAR(80) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_model_vendors_provider_key ON model_vendors(provider_key);

CREATE TABLE IF NOT EXISTS model_catalogs (
    id BIGSERIAL PRIMARY KEY,
    model_id VARCHAR(200) NOT NULL,
    normalized_model_id VARCHAR(200) NOT NULL,
    display_name VARCHAR(200) NOT NULL DEFAULT '',
    platform VARCHAR(50) NOT NULL DEFAULT '',
    provider VARCHAR(100) NOT NULL DEFAULT '',
    vendor_id BIGINT NULL REFERENCES model_vendors(id) ON DELETE SET NULL,
    mode VARCHAR(50) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    capabilities JSONB NOT NULL DEFAULT '{}'::jsonb,
    endpoints JSONB NOT NULL DEFAULT '[]'::jsonb,
    pricing JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',
    source VARCHAR(50) NOT NULL DEFAULT 'manual',
    icon_key VARCHAR(80) NOT NULL DEFAULT '',
    last_synced_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT model_catalogs_status_check CHECK (status IN ('active', 'disabled')),
    CONSTRAINT model_catalogs_visibility_check CHECK (visibility IN ('public', 'admin')),
    CONSTRAINT model_catalogs_platform_model_unique UNIQUE (platform, normalized_model_id)
);

CREATE INDEX IF NOT EXISTS idx_model_catalogs_platform ON model_catalogs(platform);
CREATE INDEX IF NOT EXISTS idx_model_catalogs_provider ON model_catalogs(provider);
CREATE INDEX IF NOT EXISTS idx_model_catalogs_vendor_id ON model_catalogs(vendor_id);
CREATE INDEX IF NOT EXISTS idx_model_catalogs_status_visibility ON model_catalogs(status, visibility);
CREATE INDEX IF NOT EXISTS idx_model_catalogs_source ON model_catalogs(source);
