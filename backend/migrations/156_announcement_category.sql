ALTER TABLE announcements
    ADD COLUMN IF NOT EXISTS category VARCHAR(32) NOT NULL DEFAULT 'announcement';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'announcements_category_check'
          AND table_name = 'announcements'
    ) THEN
        ALTER TABLE announcements
            ADD CONSTRAINT announcements_category_check
            CHECK (category IN ('announcement', 'model_update', 'changelog'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_announcements_category ON announcements (category);
