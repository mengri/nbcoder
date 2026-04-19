-- 000004_add_knowledge_git_notify_tables.down.sql
-- Drop Knowledge, Git, Notify, and Clone Pool tables

DROP INDEX IF EXISTS idx_clone_instances_deleted_at;
DROP INDEX IF EXISTS idx_clone_instances_status;
DROP INDEX IF EXISTS idx_clone_instances_project;
DROP TABLE IF EXISTS clone_instances;

DROP INDEX IF EXISTS idx_notification_histories_deleted_at;
DROP INDEX IF EXISTS idx_notification_histories_recipient;
DROP INDEX IF EXISTS idx_notification_histories_notification;
DROP TABLE IF EXISTS notification_histories;

DROP INDEX IF EXISTS idx_notification_templates_deleted_at;
DROP INDEX IF EXISTS idx_notification_templates_event_type;
DROP TABLE IF EXISTS notification_templates;

DROP INDEX IF EXISTS idx_subscription_preferences_deleted_at;
DROP INDEX IF EXISTS idx_subscription_preferences_event_type;
DROP INDEX IF EXISTS idx_subscription_preferences_recipient;
DROP TABLE IF EXISTS subscription_preferences;

DROP INDEX IF EXISTS idx_subscriptions_deleted_at;
DROP INDEX IF EXISTS idx_subscriptions_event_type;
DROP INDEX IF EXISTS idx_subscriptions_recipient;
DROP TABLE IF EXISTS subscriptions;

DROP INDEX IF EXISTS idx_notifications_deleted_at;
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_status;
DROP INDEX IF EXISTS idx_notifications_recipient;
DROP INDEX IF EXISTS idx_notifications_event_type;
DROP TABLE IF EXISTS notifications;

DROP INDEX IF EXISTS idx_commits_deleted_at;
DROP INDEX IF EXISTS idx_commits_commit_time;
DROP INDEX IF EXISTS idx_commits_hash;
DROP INDEX IF EXISTS idx_commits_project;
DROP TABLE IF EXISTS commits;

DROP INDEX IF EXISTS idx_pull_requests_deleted_at;
DROP INDEX IF EXISTS idx_pull_requests_source_branch;
DROP INDEX IF EXISTS idx_pull_requests_status;
DROP INDEX IF EXISTS idx_pull_requests_project;
DROP TABLE IF EXISTS pull_requests;

DROP INDEX IF EXISTS idx_repositories_deleted_at;
DROP TABLE IF EXISTS repositories;

DROP INDEX IF EXISTS idx_directories_deleted_at;
DROP INDEX IF EXISTS idx_directories_parent;
DROP INDEX IF EXISTS idx_directories_path;
DROP INDEX IF EXISTS idx_directories_project;
DROP TABLE IF EXISTS directories;

DROP INDEX IF EXISTS idx_document_indices_deleted_at;
DROP INDEX IF EXISTS idx_document_indices_document;
DROP TABLE IF EXISTS document_indices;

DROP INDEX IF EXISTS idx_document_chunks_deleted_at;
DROP INDEX IF EXISTS idx_document_chunks_document;
DROP TABLE IF EXISTS document_chunks;

DROP INDEX IF EXISTS idx_documents_deleted_at;
DROP INDEX IF EXISTS idx_documents_path;
DROP INDEX IF EXISTS idx_documents_project;
DROP TABLE IF EXISTS documents;
