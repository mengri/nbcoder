package sqlite

import (
	"encoding/json"

	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type DocumentChunkRepo struct {
	db *gorm.DB
}

func NewDocumentChunkRepo(db *gorm.DB) *DocumentChunkRepo {
	return &DocumentChunkRepo{db: db}
}

func (r *DocumentChunkRepo) Save(chunk *knowledge.Chunk) error {
	embeddingJSON, _ := json.Marshal(chunk.Embedding)
	model := &models.DocumentChunk{
		ID:         chunk.ID,
		DocumentID: chunk.DocumentID,
		Content:    chunk.Content,
		ChunkIndex: chunk.ChunkIndex,
		Embedding:  chunk.Embedding,
		CreatedAt:  chunk.CreatedAt,
		UpdatedAt:  chunk.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DocumentChunkRepo) FindByID(id string) (*knowledge.Chunk, error) {
	var model models.DocumentChunk
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.modelToDomain(&model), nil
}

func (r *DocumentChunkRepo) FindByDocumentID(documentID string) ([]*knowledge.Chunk, error) {
	var models []models.DocumentChunk
	result := r.db.Where("document_id = ?", documentID).Order("chunk_index ASC").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentChunkRepo) Update(chunk *knowledge.Chunk) error {
	embeddingJSON, _ := json.Marshal(chunk.Embedding)
	model := &models.DocumentChunk{
		ID:         chunk.ID,
		DocumentID: chunk.DocumentID,
		Content:    chunk.Content,
		ChunkIndex: chunk.ChunkIndex,
		Embedding:  chunk.Embedding,
		CreatedAt:  chunk.CreatedAt,
		UpdatedAt:  chunk.UpdatedAt,
	}

	result := r.db.Model(&models.DocumentChunk{}).Where("id = ?", chunk.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DocumentChunkRepo) Delete(id string) error {
	result := r.db.Delete(&models.DocumentChunk{}, "id = ?", id)
	return result.Error
}

func (r *DocumentChunkRepo) modelToDomain(m *models.DocumentChunk) *knowledge.Chunk {
	return &knowledge.Chunk{
		ID:         m.ID,
		DocumentID: m.DocumentID,
		Content:    m.Content,
		ChunkIndex: m.ChunkIndex,
		Embedding:  m.Embedding,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
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
	db *gorm.DB
}

func NewDocumentIndexRepo(db *gorm.DB) *DocumentIndexRepo {
	return &DocumentIndexRepo{db: db}
}

func (r *DocumentIndexRepo) Save(index *knowledge.DocumentIndex) error {
	model := &models.DocumentIndex{
		ID:         index.ID,
		DocumentID: index.DocumentID,
		IndexName:  index.IndexName,
		IndexData:  models.JSONMap(index.IndexData),
		CreatedAt:  index.CreatedAt,
		UpdatedAt:  index.UpdatedAt,
	}

	result := r.db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DocumentIndexRepo) FindByID(id string) (*knowledge.DocumentIndex, error) {
	var model models.DocumentIndex
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return r.modelToDomain(&model), nil
}

func (r *DocumentIndexRepo) FindByDocumentID(documentID string) ([]*knowledge.DocumentIndex, error) {
	var models []models.DocumentIndex
	result := r.db.Where("document_id = ?", documentID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.modelsToDomain(models), nil
}

func (r *DocumentIndexRepo) Update(index *knowledge.DocumentIndex) error {
	model := &models.DocumentIndex{
		ID:         index.ID,
		DocumentID: index.DocumentID,
		IndexName:  index.IndexName,
		IndexData:  models.JSONMap(index.IndexData),
		CreatedAt:  index.CreatedAt,
		UpdatedAt:  index.UpdatedAt,
	}

	result := r.db.Model(&models.DocumentIndex{}).Where("id = ?", index.ID).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DocumentIndexRepo) Delete(id string) error {
	result := r.db.Delete(&models.DocumentIndex{}, "id = ?", id)
	return result.Error
}

func (r *DocumentIndexRepo) modelToDomain(m *models.DocumentIndex) *knowledge.DocumentIndex {
	return &knowledge.DocumentIndex{
		ID:         m.ID,
		DocumentID: m.DocumentID,
		IndexName:  m.IndexName,
		IndexData:  map[string]interface{}(m.IndexData),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (r *DocumentIndexRepo) modelsToDomain(models []models.DocumentIndex) []*knowledge.DocumentIndex {
	result := make([]*knowledge.DocumentIndex, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
