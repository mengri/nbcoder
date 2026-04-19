-- 000001_create_initial_schema.up.sql
-- Create initial database schema for NBCoder

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    repo_url VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_repo_url ON projects(repo_url);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);

-- Project configs table
CREATE TABLE IF NOT EXISTS project_configs (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_project_configs_project ON project_configs(project_id);
CREATE INDEX IF NOT EXISTS idx_project_configs_key ON project_configs(key);
CREATE INDEX IF NOT EXISTS idx_project_configs_deleted_at ON project_configs(deleted_at);

-- Standards table
CREATE TABLE IF NOT EXISTS standards (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL UNIQUE,
    branch_strategy TEXT,
    tech_stack TEXT,
    coding_conventions TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_standards_deleted_at ON standards(deleted_at);

-- Dev standards table
CREATE TABLE IF NOT EXISTS dev_standards (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    standard_type VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_dev_standards_project ON dev_standards(project_id);
CREATE INDEX IF NOT EXISTS idx_dev_standards_type ON dev_standards(standard_type);
CREATE INDEX IF NOT EXISTS idx_dev_standards_deleted_at ON dev_standards(deleted_at);

-- Branch policy configs table
CREATE TABLE IF NOT EXISTS branch_policy_configs (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    policy_name VARCHAR(255) NOT NULL,
    policy_config TEXT,
    require_reviews BOOLEAN DEFAULT 0,
    min_reviewers INTEGER DEFAULT 1,
    require_tests BOOLEAN DEFAULT 0,
    auto_merge_enabled BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_branch_policy_project ON branch_policy_configs(project_id);
CREATE INDEX IF NOT EXISTS idx_branch_policy_deleted_at ON branch_policy_configs(deleted_at);

-- Project lifecycles table
CREATE TABLE IF NOT EXISTS project_lifecycles (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATING',
    activated_at DATETIME,
    suspended_at DATETIME,
    archived_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_lifecycles_status ON project_lifecycles(status);
CREATE INDEX IF NOT EXISTS idx_lifecycles_deleted_at ON project_lifecycles(deleted_at);

-- Config change logs table
CREATE TABLE IF NOT EXISTS config_change_logs (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    config_key VARCHAR(255) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    changed_by VARCHAR(255),
    deleted_at DATETIME,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_change_logs_project ON config_change_logs(project_id);
CREATE INDEX IF NOT EXISTS idx_change_logs_key ON config_change_logs(config_key);
CREATE INDEX IF NOT EXISTS idx_change_logs_changed_at ON config_change_logs(changed_at);
CREATE INDEX IF NOT EXISTS idx_change_logs_deleted_at ON config_change_logs(deleted_at);
