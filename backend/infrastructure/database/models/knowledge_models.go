package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Path        string         `gorm:"type:varchar(500);not null;index" json:"path"`
	ProjectID   string         `gorm:"type:varchar(36);not null;index:idx_document_project" json:"project_id"`
	DirectoryID string         `gorm:"type:varchar(36);index" json:"directory_id"`
	Content     string         `gorm:"type:text" json:"content"`
	Version     int            `gorm:"not null;default:1" json:"version"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project   *Project        `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Directory *Directory      `gorm:"foreignKey:DirectoryID" json:"directory,omitempty"`
	Chunks    []DocumentChunk `gorm:"foreignKey:DocumentID" json:"chunks,omitempty"`
	Indices   []DocumentIndex `gorm:"foreignKey:DocumentID" json:"indices,omitempty"`
}

func (Document) TableName() string {
	return "documents"
}

type DocumentChunk struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	DocumentID string         `gorm:"type:varchar(36);not null;index:idx_chunk_document" json:"document_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	ChunkIndex int            `gorm:"not null" json:"chunk_index"`
	Embedding  []float64      `gorm:"type:json" json:"embedding,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Document *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
}

func (DocumentChunk) TableName() string {
	return "document_chunks"
}

type DocumentIndex struct {
	ID         string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	DocumentID string         `gorm:"type:varchar(36);not null;index:idx_index_document" json:"document_id"`
	IndexName  string         `gorm:"type:varchar(255);not null" json:"index_name"`
	IndexData  JSONMap        `gorm:"type:json" json:"index_data,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Document *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
}

func (DocumentIndex) TableName() string {
	return "document_indices"
}

type Directory struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	ParentID  string         `gorm:"type:varchar(36);index" json:"parent_id"`
	ProjectID string         `gorm:"type:varchar(36);not null;index:idx_directory_project" json:"project_id"`
	Path      string         `gorm:"type:varchar(500);not null;index" json:"path"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Project   *Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Parent    *Directory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Directory `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Documents []Document  `gorm:"foreignKey:DirectoryID" json:"documents,omitempty"`
}

func (Directory) TableName() string {
	return "directories"
}
