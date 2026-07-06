ALTER TABLE channel_model_pricing
    ADD COLUMN IF NOT EXISTS display_currency VARCHAR(3) NOT NULL DEFAULT 'USD';

ALTER TABLE usage_logs
    ADD COLUMN IF NOT EXISTS display_currency VARCHAR(3) NOT NULL DEFAULT 'USD';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'channel_model_pricing_display_currency_check'
          AND table_name = 'channel_model_pricing'
    ) THEN
        ALTER TABLE channel_model_pricing
            ADD CONSTRAINT channel_model_pricing_display_currency_check
            CHECK (display_currency IN ('USD', 'CNY'));
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'usage_logs_display_currency_check'
          AND table_name = 'usage_logs'
    ) THEN
        ALTER TABLE usage_logs
            ADD CONSTRAINT usage_logs_display_currency_check
            CHECK (display_currency IN ('USD', 'CNY'));
    END IF;
END $$;
