-- 145_subscription_request_limits.sql
-- Add request count limits to subscription_plans and request usage tracking to user_subscriptions.

-- Add request count limits to subscription_plans (product catalog template).
ALTER TABLE subscription_plans
    ADD COLUMN IF NOT EXISTS daily_request_limit BIGINT,
    ADD COLUMN IF NOT EXISTS weekly_request_limit BIGINT,
    ADD COLUMN IF NOT EXISTS monthly_request_limit BIGINT;

-- Add plan live reference and request usage counters to user_subscriptions.
ALTER TABLE user_subscriptions
    ADD COLUMN IF NOT EXISTS plan_id BIGINT,
    ADD COLUMN IF NOT EXISTS daily_usage_requests BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS weekly_usage_requests BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS monthly_usage_requests BIGINT NOT NULL DEFAULT 0;

-- Non-negative checks for plan request limits.
DO $$ BEGIN
    ALTER TABLE subscription_plans
        ADD CONSTRAINT subscription_plans_daily_request_limit_nonnegative
            CHECK (daily_request_limit IS NULL OR daily_request_limit >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE subscription_plans
        ADD CONSTRAINT subscription_plans_weekly_request_limit_nonnegative
            CHECK (weekly_request_limit IS NULL OR weekly_request_limit >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE subscription_plans
        ADD CONSTRAINT subscription_plans_monthly_request_limit_nonnegative
            CHECK (monthly_request_limit IS NULL OR monthly_request_limit >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Non-negative checks for subscription request usage.
DO $$ BEGIN
    ALTER TABLE user_subscriptions
        ADD CONSTRAINT user_subscriptions_daily_usage_requests_nonnegative
            CHECK (daily_usage_requests >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE user_subscriptions
        ADD CONSTRAINT user_subscriptions_weekly_usage_requests_nonnegative
            CHECK (weekly_usage_requests >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE user_subscriptions
        ADD CONSTRAINT user_subscriptions_monthly_usage_requests_nonnegative
            CHECK (monthly_usage_requests >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Partial index for plan_id lookups (skip soft-deleted rows).
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_plan_id
    ON user_subscriptions(plan_id)
    WHERE deleted_at IS NULL;
