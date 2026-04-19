package knowledge

import (
	"github.com/mengri/nbcoder/domain/knowledge"
)

type KnowledgeService struct {
	docRepo    knowledge.DocumentRepo
	chunkRepo  knowledge.ChunkRepo
	indexRepo  knowledge.DocumentIndexRepo
}

func NewKnowledgeService(
	docRepo knowledge.DocumentRepo,
	chunkRepo knowledge.ChunkRepo,
	indexRepo knowledge.DocumentIndexRepo,
) *KnowledgeService {
	return &KnowledgeService{
		docRepo:   docRepo,
		chunkRepo: chunkRepo,
		indexRepo: indexRepo,
	}
}

func (s *KnowledgeService) UploadDocument(doc *knowledge.Document) error {
	return s.docRepo.Save(doc)
}

func (s *KnowledgeService) Search(query *knowledge.SearchQuery) ([]*knowledge.SearchResult, error) {
	return s.indexRepo.Search(query)
}

func (s *KnowledgeService) GetDocument(id string) (*knowledge.Document, error) {
	return s.docRepo.FindByID(id)
}
