-- Add configurable retry attempts for channel monitor checks.
-- retry_count means extra attempts after the first run; 0 disables retry.

ALTER TABLE channel_monitors
    ADD COLUMN IF NOT EXISTS retry_count INT NOT NULL DEFAULT 0;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'channel_monitors_retry_count_check'
          AND table_name = 'channel_monitors'
    ) THEN
        ALTER TABLE channel_monitors
            ADD CONSTRAINT channel_monitors_retry_count_check
            CHECK (retry_count BETWEEN 0 AND 5);
    END IF;
END $$;
