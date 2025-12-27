-- Drop trigger if exist
DROP TRIGGER IF EXISTS set_user_updated_at on users;

-- Drop tigger function
DROP FUNCTION IF EXISTS update_user_updated_at_column;

-- Drop index (No need, because deleting the table will also delete the index, but writing it in makes it clearer.)
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_level;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_user_deleted_at;
DROP INDEX IF EXISTS idx_users_email_status;

-- Drop table
DROP TABLE IF EXISTS users;