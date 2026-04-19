-- 000003_add_ai_runtime_tables.down.sql
-- Drop AI Runtime related tables

DROP INDEX IF EXISTS idx_call_logs_deleted_at;
DROP INDEX IF EXISTS idx_call_logs_created_at;
DROP INDEX IF EXISTS idx_call_logs_status;
DROP INDEX IF EXISTS idx_call_logs_call_type;
DROP INDEX IF EXISTS idx_call_logs_model;
DROP INDEX IF EXISTS idx_call_logs_agent;
DROP TABLE IF EXISTS call_logs;

DROP INDEX IF EXISTS idx_model_chains_deleted_at;
DROP INDEX IF EXISTS idx_model_chains_model;
DROP TABLE IF EXISTS model_chains;

DROP INDEX IF EXISTS idx_models_deleted_at;
DROP INDEX IF EXISTS idx_models_type;
DROP INDEX IF EXISTS idx_models_provider;
DROP INDEX IF EXISTS idx_models_name;
DROP TABLE IF EXISTS models;

DROP INDEX IF EXISTS idx_providers_deleted_at;
DROP TABLE IF EXISTS providers;

DROP INDEX IF EXISTS idx_skills_deleted_at;
DROP INDEX IF EXISTS idx_skills_agent_type;
DROP TABLE IF EXISTS skills;

DROP INDEX IF EXISTS idx_agent_executions_deleted_at;
DROP INDEX IF EXISTS idx_agent_executions_status;
DROP INDEX IF EXISTS idx_agent_executions_agent;
DROP INDEX IF EXISTS idx_agent_executions_task;
DROP TABLE IF EXISTS agent_executions;

DROP INDEX IF EXISTS idx_tasks_deleted_at;
DROP INDEX IF EXISTS idx_tasks_pipeline;
DROP INDEX IF EXISTS idx_tasks_project;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_name;
DROP TABLE IF EXISTS tasks;
