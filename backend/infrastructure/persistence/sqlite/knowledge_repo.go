package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type DocumentRepo struct {
	dbProvider DBProvider
}

func NewDocumentRepo(dbProvider DBProvider) knowledge.DocumentRepo {
	return &DocumentRepo{dbProvider: dbProvider}
}

func (r *DocumentRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *DocumentRepo) Save(doc *knowledge.Document) error {
	db, err := r.getDB(doc.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Document{
		ID:          doc.ID,
		Name:        doc.Name,
		Path:        doc.Path,
		ProjectName: doc.ProjectName,
		DirectoryID: doc.DirectoryID,
		Content:     doc.Content,
		Version:     doc.Version,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save document: %w", result.Error)
	}
	return nil
}

func (r *DocumentRepo) FindByID(id string, projectName string) (*knowledge.Document, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.Document
	result := db.Preload("Chunks").Preload("Indices").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find document by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DocumentRepo) FindByProjectName(projectName string) ([]*knowledge.Document, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.Document
	result := db.Where("project_name = ?", projectName).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find documents by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentRepo) FindByDirectoryID(directoryID string) ([]*knowledge.Document, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Document
	result := db.Where("directory_id = ?", directoryID).Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find documents by directory id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentRepo) Update(doc *knowledge.Document) error {
	db, err := r.getDB(doc.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Document{
		ID:          doc.ID,
		Name:        doc.Name,
		Path:        doc.Path,
		ProjectName: doc.ProjectName,
		DirectoryID: doc.DirectoryID,
		Content:     doc.Content,
		Version:     doc.Version,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	result := db.Model(&models.Document{}).Where("id = ?", doc.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update document: %w", result.Error)
	}
	return nil
}

func (r *DocumentRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.Document{}, "id = ?", id)
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
		ProjectName: m.ProjectName,
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
	dbProvider DBProvider
}

func NewDirectoryRepo(dbProvider DBProvider) knowledge.DirectoryRepo {
	return &DirectoryRepo{dbProvider: dbProvider}
}

func (r *DirectoryRepo) getDB(projectName string) (*gorm.DB, error) {
	if projectName == "" {
		return r.dbProvider.GetGlobalDB(), nil
	}
	return r.dbProvider.GetProjectDB(projectName)
}

func (r *DirectoryRepo) Save(dir *knowledge.Directory) error {
	db, err := r.getDB(dir.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Directory{
		ID:          dir.ID,
		Name:        dir.Name,
		ParentID:    dir.ParentID,
		ProjectName: dir.ProjectName,
		Path:        dir.Path,
		CreatedAt:   dir.CreatedAt,
		UpdatedAt:   dir.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save directory: %w", result.Error)
	}
	return nil
}

func (r *DirectoryRepo) FindByID(id string, projectName string) (*knowledge.Directory, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.Directory
	result := db.Preload("Children").Preload("Documents").First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find directory by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DirectoryRepo) FindByProjectName(projectName string) ([]*knowledge.Directory, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var models []models.Directory
	result := db.Where("project_name = ?", projectName).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find directories by project name: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DirectoryRepo) FindByParentID(parentID string) ([]*knowledge.Directory, error) {
	db, err := r.getDB("")
	if err != nil {
		return nil, err
	}

	var models []models.Directory
	result := db.Where("parent_id = ?", parentID).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find directories by parent id: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *DirectoryRepo) FindRootByProjectName(projectName string) (*knowledge.Directory, error) {
	db, err := r.getDB(projectName)
	if err != nil {
		return nil, err
	}

	var model models.Directory
	result := db.Where("project_name = ? AND parent_id = ''", projectName).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find root directory: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *DirectoryRepo) Update(dir *knowledge.Directory) error {
	db, err := r.getDB(dir.ProjectName)
	if err != nil {
		return err
	}

	model := &models.Directory{
		ID:          dir.ID,
		Name:        dir.Name,
		ParentID:    dir.ParentID,
		ProjectName: dir.ProjectName,
		Path:        dir.Path,
		CreatedAt:   dir.CreatedAt,
		UpdatedAt:   dir.UpdatedAt,
	}

	result := db.Model(&models.Directory{}).Where("id = ?", dir.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update directory: %w", result.Error)
	}
	return nil
}

func (r *DirectoryRepo) Delete(id string, projectName string) error {
	db, err := r.getDB(projectName)
	if err != nil {
		return err
	}

	result := db.Delete(&models.Directory{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete directory: %w", result.Error)
	}
	return nil
}

func (r *DirectoryRepo) modelToDomain(m *models.Directory) *knowledge.Directory {
	return &knowledge.Directory{
		ID:          m.ID,
		Name:        m.Name,
		ParentID:    m.ParentID,
		ProjectName: m.ProjectName,
		Path:        m.Path,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *DirectoryRepo) modelsToDomain(models []models.Directory) []*knowledge.Directory {
	result := make([]*knowledge.Directory, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
