package sqlite

import (
	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type DocumentChunkRepo struct {
	dbProvider DBProvider
}

func NewDocumentChunkRepo(dbProvider DBProvider) *DocumentChunkRepo {
	return &DocumentChunkRepo{dbProvider: dbProvider}
}

func (r *DocumentChunkRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *DocumentChunkRepo) DeleteByDocumentID(documentID string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}
	result := db.Where("document_id = ?", documentID).Delete(&models.DocumentChunk{})
	return result.Error
}

func (r *DocumentChunkRepo) Save(chunk *knowledge.Chunk) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.DocumentChunk{
		ID:         chunk.ID,
		DocumentID: chunk.DocumentID,
		Content:    chunk.Content,
		ChunkIndex: chunk.Index,
		Embedding:  []float64{},
	}

	result := db.Save(model)
	return result.Error
}

func (r *DocumentChunkRepo) FindByID(id string) (*knowledge.Chunk, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.DocumentChunk
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.modelToDomain(&model), nil
}

func (r *DocumentChunkRepo) FindByDocumentID(documentID string) ([]*knowledge.Chunk, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.DocumentChunk
	result := db.Where("document_id = ?", documentID).Order("chunk_index ASC").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentChunkRepo) Update(chunk *knowledge.Chunk) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.DocumentChunk{
		DocumentID: chunk.DocumentID,
		Content:    chunk.Content,
		ChunkIndex: chunk.Index,
		Embedding:  []float64{},
	}

	result := db.Model(&models.DocumentChunk{}).Where("id = ?", chunk.ID).Updates(model)
	return result.Error
}

func (r *DocumentChunkRepo) Delete(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.Delete(&models.DocumentChunk{}, "id = ?", id)
	return result.Error
}

func (r *DocumentChunkRepo) modelToDomain(m *models.DocumentChunk) *knowledge.Chunk {
	return &knowledge.Chunk{
		ID:         m.ID,
		DocumentID: m.DocumentID,
		Content:    m.Content,
		Index:      m.ChunkIndex,
	}
}

func (r *DocumentChunkRepo) modelsToDomain(models []models.DocumentChunk) []*knowledge.Chunk {
	result := make([]*knowledge.Chunk, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}

type DocumentIndexRepo struct {
	dbProvider DBProvider
}

func NewDocumentIndexRepo(dbProvider DBProvider) *DocumentIndexRepo {
	return &DocumentIndexRepo{dbProvider: dbProvider}
}

func (r *DocumentIndexRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *DocumentIndexRepo) DeleteByDocumentID(documentID string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}
	result := db.Where("document_id = ?", documentID).Delete(&models.DocumentIndex{})
	return result.Error
}

func (r *DocumentIndexRepo) Search(query *knowledge.SearchQuery) ([]*knowledge.SearchResult, error) {
	return []*knowledge.SearchResult{}, nil
}

func (r *DocumentIndexRepo) Save(index *knowledge.DocumentIndex) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.DocumentIndex{
		ID:         index.ID,
		DocumentID: index.DocumentID,
		IndexName:  "",
		IndexData:  models.JSONMap{},
	}

	result := db.Save(model)
	return result.Error
}

func (r *DocumentIndexRepo) FindByID(id string) (*knowledge.DocumentIndex, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.DocumentIndex
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.modelToDomain(&model), nil
}

func (r *DocumentIndexRepo) FindByDocumentID(documentID string) ([]*knowledge.DocumentIndex, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.DocumentIndex
	result := db.Where("document_id = ?", documentID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentIndexRepo) Update(index *knowledge.DocumentIndex) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.DocumentIndex{
		DocumentID: index.DocumentID,
	}

	result := db.Model(&models.DocumentIndex{}).Where("id = ?", index.ID).Updates(model)
	return result.Error
}

func (r *DocumentIndexRepo) Delete(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.Delete(&models.DocumentIndex{}, "id = ?", id)
	return result.Error
}

func (r *DocumentIndexRepo) modelToDomain(m *models.DocumentIndex) *knowledge.DocumentIndex {
	return &knowledge.DocumentIndex{
		ID:        m.ID,
		DocumentID: m.DocumentID,
		ChunkID:   "",
		Embedding: []float64{},
	}
}

func (r *DocumentIndexRepo) modelsToDomain(models []models.DocumentIndex) []*knowledge.DocumentIndex {
	result := make([]*knowledge.DocumentIndex, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
