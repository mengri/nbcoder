-- 000002_add_indexes.up.sql
-- Add performance indexes for common queries

-- Cards table and related
CREATE TABLE IF NOT EXISTS cards (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    original TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',
    structured_output TEXT,
    pipeline_id VARCHAR(36),
    project_id VARCHAR(36) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    superseded_by VARCHAR(36),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (pipeline_id) REFERENCES pipelines(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_cards_title ON cards(title);
CREATE INDEX IF NOT EXISTS idx_cards_status ON cards(status);
CREATE INDEX IF NOT EXISTS idx_cards_priority ON cards(priority);
CREATE INDEX IF NOT EXISTS idx_cards_project ON cards(project_id);
CREATE INDEX IF NOT EXISTS idx_cards_pipeline ON cards(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_cards_superseded_by ON cards(superseded_by);
CREATE INDEX IF NOT EXISTS idx_cards_deleted_at ON cards(deleted_at);

CREATE TABLE IF NOT EXISTS card_dependencies (
    id VARCHAR(36) PRIMARY KEY,
    card_id VARCHAR(36) NOT NULL,
    depends_on_card_id VARCHAR(36) NOT NULL,
    dependency_type VARCHAR(50) NOT NULL DEFAULT 'BLOCKING',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on_card_id) REFERENCES cards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_card_dependencies_card ON card_dependencies(card_id);
CREATE INDEX IF NOT EXISTS idx_card_dependencies_depends_on ON card_dependencies(depends_on_card_id);
CREATE INDEX IF NOT EXISTS idx_card_dependencies_deleted_at ON card_dependencies(deleted_at);

-- Pipelines table and related
CREATE TABLE IF NOT EXISTS pipelines (
    id VARCHAR(36) PRIMARY KEY,
    card_id VARCHAR(36) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pipelines_deleted_at ON pipelines(deleted_at);

CREATE TABLE IF NOT EXISTS stage_records (
    id VARCHAR(36) PRIMARY KEY,
    pipeline_id VARCHAR(36) NOT NULL,
    stage_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'NOT_STARTED',
    started_at DATETIME,
    completed_at DATETIME,
    output TEXT,
    error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (pipeline_id) REFERENCES pipelines(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_stage_records_pipeline ON stage_records(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_stage_records_status ON stage_records(status);
CREATE INDEX IF NOT EXISTS idx_stage_records_deleted_at ON stage_records(deleted_at);
