ALTER TABLE inventories
  ADD COLUMN IF NOT EXISTS reorder_threshold BIGINT NOT NULL DEFAULT 10;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'inventories_reorder_threshold_check'
      AND conrelid = 'inventories'::regclass
  ) THEN
    ALTER TABLE inventories
      ADD CONSTRAINT inventories_reorder_threshold_check CHECK (reorder_threshold >= 0);
  END IF;
END $$;

UPDATE inventories
SET reorder_threshold = 10
WHERE reorder_threshold IS NULL;
