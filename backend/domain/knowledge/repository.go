package knowledge

type DocumentIndex struct {
	ID         string   `json:"id"`
	DocumentID string   `json:"document_id"`
	ChunkID    string   `json:"chunk_id"`
	Embedding  []float64 `json:"embedding,omitempty"`
}

type SearchQuery struct {
	Query     string  `json:"query"`
	TopK      int     `json:"top_k"`
	Threshold float64 `json:"threshold,omitempty"`
}

type SearchResult struct {
	Chunk    *Chunk   `json:"chunk"`
	Score    float64  `json:"score"`
	Document *Document `json:"document,omitempty"`
}

type DocumentRepo interface {
	Save(doc *Document) error
	FindByID(id string) (*Document, error)
	FindByProjectID(projectID string) ([]*Document, error)
	FindByDirectoryID(directoryID string) ([]*Document, error)
	Update(doc *Document) error
	Delete(id string) error
}

type DirectoryRepo interface {
	Save(dir *Directory) error
	FindByID(id string) (*Directory, error)
	FindByProjectID(projectID string) ([]*Directory, error)
	FindByParentID(parentID string) ([]*Directory, error)
	Update(dir *Directory) error
	Delete(id string) error
}

type ChunkRepo interface {
	Save(chunk *Chunk) error
	FindByDocumentID(documentID string) ([]*Chunk, error)
	FindByID(id string) (*Chunk, error)
	DeleteByDocumentID(documentID string) error
}

type DocumentIndexRepo interface {
	Save(index *DocumentIndex) error
	Search(query *SearchQuery) ([]*SearchResult, error)
	DeleteByDocumentID(documentID string) error
}

type LineageRepo interface {
	Save(lineage *DocumentLineage) error
	FindByDocumentID(documentID string) ([]*DocumentLineage, error)
	FindByParentDocumentID(parentDocumentID string) ([]*DocumentLineage, error)
	FindAllAncestors(documentID string) ([]*DocumentLineage, error)
	FindAllDescendants(documentID string) ([]*DocumentLineage, error)
}
