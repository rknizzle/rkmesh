ALTER TABLE models DROP CONSTRAINT IF EXISTS models_user_id_fkey;
ALTER TABLE models DROP COLUMN IF EXISTS user_id;