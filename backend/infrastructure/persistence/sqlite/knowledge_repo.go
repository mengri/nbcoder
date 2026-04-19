package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type DocumentRepo struct {
	db *gorm.DB
}

func NewDocumentRepo(db *gorm.DB) knowledge.DocumentRepo {
	return &DocumentRepo{db: db}
}

func (r *DocumentRepo) Save(doc *knowledge.Document) error {
	model := &models.Document{
		ID:          doc.ID,
		Name:        doc.Name,
		Path:        doc.Path,
		ProjectID:   doc.ProjectID,
		DirectoryID: doc.DirectoryID,
		Content:     doc.Content,
		Version:     doc.Version,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save document: %w", result.Error)
	}
	return nil
}

func (r *DocumentRepo) FindByID(id string) (*knowledge.Document, error) {
	var model models.Document
	result := r.db.Preload("Chunks").Preload("Indices").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find document by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DocumentRepo) FindByProjectID(projectID string) ([]*knowledge.Document, error) {
	var models []models.Document
	result := r.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find documents by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentRepo) FindByDirectoryID(directoryID string) ([]*knowledge.Document, error) {
	var models []models.Document
	result := r.db.Where("directory_id = ?", directoryID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find documents by directory id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentRepo) Update(doc *knowledge.Document) error {
	model := &models.Document{
		ID:          doc.ID,
		Name:        doc.Name,
		Path:        doc.Path,
		ProjectID:   doc.ProjectID,
		DirectoryID: doc.DirectoryID,
		Content:     doc.Content,
		Version:     doc.Version,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	result := r.db.Model(&models.Document{}).Where("id = ?", doc.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update document: %w", result.Error)
	}
	return nil
}

func (r *DocumentRepo) Delete(id string) error {
	result := r.db.Delete(&models.Document{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete document: %w", result.Error)
	}
	return nil
}

func (r *DocumentRepo) modelToDomain(m *models.Document) *knowledge.Document {
	return &knowledge.Document{
		ID:          m.ID,
		Name:        m.Name,
		Path:        m.Path,
		ProjectID:   m.ProjectID,
		DirectoryID: m.DirectoryID,
		Content:     m.Content,
		Version:     m.Version,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *DocumentRepo) modelsToDomain(models []models.Document) []*knowledge.Document {
	result := make([]*knowledge.Document, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type DirectoryRepo struct {
	db *gorm.DB
}

func NewDirectoryRepo(db *gorm.DB) knowledge.DirectoryRepo {
	return &DirectoryRepo{db: db}
}

func (r *DirectoryRepo) Save(dir *knowledge.Directory) error {
	model := &models.Directory{
		ID:        dir.ID,
		Name:      dir.Name,
		ParentID:  dir.ParentID,
		ProjectID: dir.ProjectID,
		Path:      dir.Path,
		CreatedAt: dir.CreatedAt,
		UpdatedAt: dir.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save directory: %w", result.Error)
	}
	return nil
}

func (r *DirectoryRepo) FindByID(id string) (*knowledge.Directory, error) {
	var model models.Directory
	result := r.db.Preload("Children").Preload("Documents").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find directory by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DirectoryRepo) FindByProjectID(projectID string) ([]*knowledge.Directory, error) {
	var models []models.Directory
	result := r.db.Where("project_id = ?", projectID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find directories by project id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DirectoryRepo) FindByParentID(parentID string) ([]*knowledge.Directory, error) {
	var models []models.Directory
	result := r.db.Where("parent_id = ?", parentID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find directories by parent id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DirectoryRepo) FindRootByProjectID(projectID string) (*knowledge.Directory, error) {
	var model models.Directory
	result := r.db.Where("project_id = ? AND parent_id = ''", projectID).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find root directory: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DirectoryRepo) Update(dir *knowledge.Directory) error {
	model := &models.Directory{
		ID:        dir.ID,
		Name:      dir.Name,
		ParentID:  dir.ParentID,
		ProjectID: dir.ProjectID,
		Path:      dir.Path,
		CreatedAt: dir.CreatedAt,
		UpdatedAt: dir.UpdatedAt,
	}

	result := r.db.Model(&models.Directory{}).Where("id = ?", dir.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update directory: %w", result.Error)
	}
	return nil
}

func (r *DirectoryRepo) Delete(id string) error {
	result := r.db.Delete(&models.Directory{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete directory: %w", result.Error)
	}
	return nil
}

func (r *DirectoryRepo) modelToDomain(m *models.Directory) *knowledge.Directory {
	return &knowledge.Directory{
		ID:        m.ID,
		Name:      m.Name,
		ParentID:  m.ParentID,
		ProjectID: m.ProjectID,
		Path:      m.Path,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *DirectoryRepo) modelsToDomain(models []models.Directory) []*knowledge.Directory {
	result := make([]*knowledge.Directory, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
