-- 000001_create_initial_schema.down.sql
-- Drop initial database schema tables

DROP INDEX IF EXISTS idx_change_logs_deleted_at;
DROP INDEX IF EXISTS idx_change_logs_changed_at;
DROP INDEX IF EXISTS idx_change_logs_key;
DROP INDEX IF EXISTS idx_change_logs_project;
DROP TABLE IF EXISTS config_change_logs;

DROP INDEX IF EXISTS idx_lifecycles_deleted_at;
DROP INDEX IF EXISTS idx_lifecycles_status;
DROP TABLE IF EXISTS project_lifecycles;

DROP INDEX IF EXISTS idx_branch_policy_deleted_at;
DROP INDEX IF EXISTS idx_branch_policy_project;
DROP TABLE IF EXISTS branch_policy_configs;

DROP INDEX IF EXISTS idx_dev_standards_deleted_at;
DROP INDEX IF EXISTS idx_dev_standards_type;
DROP INDEX IF EXISTS idx_dev_standards_project;
DROP TABLE IF EXISTS dev_standards;

DROP INDEX IF EXISTS idx_standards_deleted_at;
DROP TABLE IF EXISTS standards;

DROP INDEX IF EXISTS idx_project_configs_deleted_at;
DROP INDEX IF EXISTS idx_project_configs_key;
DROP INDEX IF EXISTS idx_project_configs_project;
DROP TABLE IF EXISTS project_configs;

DROP INDEX IF EXISTS idx_projects_deleted_at;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_repo_url;
DROP INDEX IF EXISTS idx_projects_name;
DROP TABLE IF EXISTS projects;
