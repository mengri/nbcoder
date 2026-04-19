-- 000002_add_indexes.down.sql
-- Drop performance indexes and related tables

DROP INDEX IF EXISTS idx_stage_records_deleted_at;
DROP INDEX IF EXISTS idx_stage_records_status;
DROP INDEX IF EXISTS idx_stage_records_pipeline;
DROP TABLE IF EXISTS stage_records;

DROP INDEX IF EXISTS idx_pipelines_deleted_at;
DROP TABLE IF EXISTS pipelines;

DROP INDEX IF EXISTS idx_card_dependencies_deleted_at;
DROP INDEX IF EXISTS idx_card_dependencies_depends_on;
DROP INDEX IF EXISTS idx_card_dependencies_card;
DROP TABLE IF EXISTS card_dependencies;

DROP INDEX IF EXISTS idx_cards_deleted_at;
DROP INDEX IF EXISTS idx_cards_superseded_by;
DROP INDEX IF EXISTS idx_cards_pipeline;
DROP INDEX IF EXISTS idx_cards_project;
DROP INDEX IF EXISTS idx_cards_priority;
DROP INDEX IF EXISTS idx_cards_status;
DROP INDEX IF EXISTS idx_cards_title;
DROP TABLE IF EXISTS cards;
