package persistence

import (
	"sync"

	"github.com/mengri/nbcoder/domain/knowledge"
)

type InMemoryDocumentRepo struct {
	docs map[string]*knowledge.Document
	mu   sync.RWMutex
}

func NewInMemoryDocumentRepo() *InMemoryDocumentRepo {
	return &InMemoryDocumentRepo{
		docs: make(map[string]*knowledge.Document),
	}
}

func (r *InMemoryDocumentRepo) Save(doc *knowledge.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.docs[doc.ID] = doc
	return nil
}

func (r *InMemoryDocumentRepo) FindByID(id string) (*knowledge.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	doc, ok := r.docs[id]
	if !ok {
		return nil, nil
	}
	return doc, nil
}

func (r *InMemoryDocumentRepo) FindByProjectID(projectID string) ([]*knowledge.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*knowledge.Document
	for _, doc := range r.docs {
		if doc.ProjectID == projectID {
			result = append(result, doc)
		}
	}
	return result, nil
}

func (r *InMemoryDocumentRepo) Update(doc *knowledge.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.docs[doc.ID] = doc
	return nil
}

func (r *InMemoryDocumentRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.docs, id)
	return nil
}

type InMemoryChunkRepo struct {
	chunks map[string]*knowledge.Chunk
	mu     sync.RWMutex
}

func NewInMemoryChunkRepo() *InMemoryChunkRepo {
	return &InMemoryChunkRepo{
		chunks: make(map[string]*knowledge.Chunk),
	}
}

func (r *InMemoryChunkRepo) Save(chunk *knowledge.Chunk) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.chunks[chunk.ID] = chunk
	return nil
}

func (r *InMemoryChunkRepo) FindByDocumentID(documentID string) ([]*knowledge.Chunk, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*knowledge.Chunk
	for _, c := range r.chunks {
		if c.DocumentID == documentID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryChunkRepo) FindByID(id string) (*knowledge.Chunk, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.chunks[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryChunkRepo) DeleteByDocumentID(documentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, c := range r.chunks {
		if c.DocumentID == documentID {
			delete(r.chunks, id)
		}
	}
	return nil
}

type InMemoryDocumentIndexRepo struct {
	indices map[string]*knowledge.DocumentIndex
	mu      sync.RWMutex
}

func NewInMemoryDocumentIndexRepo() *InMemoryDocumentIndexRepo {
	return &InMemoryDocumentIndexRepo{
		indices: make(map[string]*knowledge.DocumentIndex),
	}
}

func (r *InMemoryDocumentIndexRepo) Save(index *knowledge.DocumentIndex) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.indices[index.ID] = index
	return nil
}

func (r *InMemoryDocumentIndexRepo) Search(query *knowledge.SearchQuery) ([]*knowledge.SearchResult, error) {
	return nil, nil
}

func (r *InMemoryDocumentIndexRepo) DeleteByDocumentID(documentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, idx := range r.indices {
		if idx.DocumentID == documentID {
			delete(r.indices, id)
		}
	}
	return nil
}

type InMemoryLineageRepo struct {
	lineages map[string]*knowledge.DocumentLineage
	mu       sync.RWMutex
}

func NewInMemoryLineageRepo() *InMemoryLineageRepo {
	return &InMemoryLineageRepo{
		lineages: make(map[string]*knowledge.DocumentLineage),
	}
}

func (r *InMemoryLineageRepo) Save(lineage *knowledge.DocumentLineage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lineages[lineage.ID] = lineage
	return nil
}

func (r *InMemoryLineageRepo) FindByDocumentID(documentID string) ([]*knowledge.DocumentLineage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*knowledge.DocumentLineage
	for _, l := range r.lineages {
		if l.DocumentID == documentID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (r *InMemoryLineageRepo) FindByParentDocumentID(parentDocumentID string) ([]*knowledge.DocumentLineage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*knowledge.DocumentLineage
	for _, l := range r.lineages {
		if l.ParentDocumentID == parentDocumentID {
			result = append(result, l)
		}
	}
	return result, nil
}

func (r *InMemoryLineageRepo) FindAllAncestors(documentID string) ([]*knowledge.DocumentLineage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	visited := make(map[string]bool)
	var result []*knowledge.DocumentLineage
	r.findAncestors(documentID, visited, &result)
	return result, nil
}

func (r *InMemoryLineageRepo) findAncestors(documentID string, visited map[string]bool, result *[]*knowledge.DocumentLineage) {
	for _, l := range r.lineages {
		if l.DocumentID == documentID && l.ParentDocumentID != "" && !visited[l.ParentDocumentID] {
			visited[l.ParentDocumentID] = true
			parentLineages, _ := r.FindByDocumentID(l.ParentDocumentID)
			for _, pl := range parentLineages {
				*result = append(*result, pl)
			}
			r.findAncestors(l.ParentDocumentID, visited, result)
		}
	}
}

func (r *InMemoryLineageRepo) FindAllDescendants(documentID string) ([]*knowledge.DocumentLineage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	visited := make(map[string]bool)
	var result []*knowledge.DocumentLineage
	r.findDescendants(documentID, visited, &result)
	return result, nil
}

func (r *InMemoryLineageRepo) findDescendants(documentID string, visited map[string]bool, result *[]*knowledge.DocumentLineage) {
	for _, l := range r.lineages {
		if l.ParentDocumentID == documentID && !visited[l.DocumentID] {
			visited[l.DocumentID] = true
			*result = append(*result, l)
			r.findDescendants(l.DocumentID, visited, result)
		}
	}
}
