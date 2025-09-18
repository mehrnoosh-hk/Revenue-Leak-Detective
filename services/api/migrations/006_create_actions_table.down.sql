-- Drop trigger first
DROP TRIGGER IF EXISTS update_actions_updated_at ON actions;

-- Drop indexes
DROP INDEX IF EXISTS idx_actions_leak_id;
DROP INDEX IF EXISTS idx_actions_action_type;
DROP INDEX IF EXISTS idx_actions_action_status;
DROP INDEX IF EXISTS idx_actions_action_result;

-- Drop table
DROP TABLE IF EXISTS actions;