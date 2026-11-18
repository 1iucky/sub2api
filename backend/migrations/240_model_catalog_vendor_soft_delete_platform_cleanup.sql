ALTER TABLE model_vendors
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_model_vendors_deleted_at ON model_vendors(deleted_at);

WITH normalized AS (
    SELECT
        id,
        CASE
            WHEN LOWER(TRIM(platform)) IN ('anthropic', 'claude') THEN 'anthropic'
            WHEN LOWER(TRIM(platform)) IN ('google', 'gemini', 'vertex_ai', 'vertex-ai', 'vertex') THEN 'gemini'
            WHEN LOWER(TRIM(platform)) = 'antigravity' THEN 'antigravity'
            ELSE 'openai'
        END AS normalized_platform,
        ROW_NUMBER() OVER (
            PARTITION BY
                CASE
                    WHEN LOWER(TRIM(platform)) IN ('anthropic', 'claude') THEN 'anthropic'
                    WHEN LOWER(TRIM(platform)) IN ('google', 'gemini', 'vertex_ai', 'vertex-ai', 'vertex') THEN 'gemini'
                    WHEN LOWER(TRIM(platform)) = 'antigravity' THEN 'antigravity'
                    ELSE 'openai'
                END,
                normalized_model_id
            ORDER BY
                CASE source WHEN 'manual' THEN 0 ELSE 1 END,
                updated_at DESC,
                id DESC
        ) AS rn
    FROM model_catalogs
)
DELETE FROM model_catalogs mc
USING normalized n
WHERE mc.id = n.id AND n.rn > 1;

UPDATE model_catalogs
SET platform = CASE
    WHEN LOWER(TRIM(platform)) IN ('anthropic', 'claude') THEN 'anthropic'
    WHEN LOWER(TRIM(platform)) IN ('google', 'gemini', 'vertex_ai', 'vertex-ai', 'vertex') THEN 'gemini'
    WHEN LOWER(TRIM(platform)) = 'antigravity' THEN 'antigravity'
    ELSE 'openai'
END
WHERE LOWER(TRIM(platform)) NOT IN ('openai', 'anthropic', 'gemini', 'antigravity');
