package knowledge

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/pkg/uid"
)

type KnowledgeService struct {
	docRepo knowledge.DocumentRepo
	dirRepo knowledge.DirectoryRepo
	chunkRepo knowledge.ChunkRepo
	indexRepo knowledge.DocumentIndexRepo
}

func NewKnowledgeService(
	docRepo knowledge.DocumentRepo,
	dirRepo knowledge.DirectoryRepo,
	chunkRepo knowledge.ChunkRepo,
	indexRepo knowledge.DocumentIndexRepo,
) *KnowledgeService {
	return &KnowledgeService{
		docRepo:   docRepo,
		dirRepo:   dirRepo,
		chunkRepo: chunkRepo,
		indexRepo: indexRepo,
	}
}

func (s *KnowledgeService) UploadDocument(name, path, projectID, directoryID, content string) (*knowledge.Document, error) {
	doc := knowledge.NewDocument(uid.NewID(), name, path, projectID)
	if directoryID != "" {
		dir, err := s.dirRepo.FindByID(directoryID)
		if err != nil {
			return nil, err
		}
		if dir == nil {
			return nil, fmt.Errorf("directory not found: %s", directoryID)
		}
		doc.SetDirectory(directoryID)
	}
	if content != "" {
		doc.SetContent(content)
	}
	if err := doc.Validate(); err != nil {
		return nil, err
	}
	if err := s.docRepo.Save(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *KnowledgeService) BatchUploadDocuments(docs []struct {
	Name, Path, ProjectID, DirectoryID, Content string
}) ([]*knowledge.Document, error) {
	var results []*knowledge.Document
	for _, d := range docs {
		doc, err := s.UploadDocument(d.Name, d.Path, d.ProjectID, d.DirectoryID, d.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to upload document %s: %w", d.Name, err)
		}
		results = append(results, doc)
	}
	return results, nil
}

func (s *KnowledgeService) GetDocument(id string) (*knowledge.Document, error) {
	return s.docRepo.FindByID(id)
}

func (s *KnowledgeService) ListDocuments(projectID string) ([]*knowledge.Document, error) {
	return s.docRepo.FindByProjectID(projectID)
}

func (s *KnowledgeService) ListDocumentsByDirectory(directoryID string) ([]*knowledge.Document, error) {
	return s.docRepo.FindByDirectoryID(directoryID)
}

func (s *KnowledgeService) DeleteDocument(id string) error {
	doc, err := s.docRepo.FindByID(id)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("document not found: %s", id)
	}
	_ = s.chunkRepo.DeleteByDocumentID(id)
	_ = s.indexRepo.DeleteByDocumentID(id)
	return s.docRepo.Delete(id)
}

func (s *KnowledgeService) CreateDirectory(name, parentID, projectID string) (*knowledge.Directory, error) {
	if parentID != "" {
		parent, err := s.dirRepo.FindByID(parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("parent directory not found: %s", parentID)
		}
	}
	dir := knowledge.NewDirectory(uid.NewID(), name, parentID, projectID)
	if err := dir.Validate(); err != nil {
		return nil, err
	}
	if err := s.dirRepo.Save(dir); err != nil {
		return nil, err
	}
	return dir, nil
}

func (s *KnowledgeService) GetDirectory(id string) (*knowledge.Directory, error) {
	return s.dirRepo.FindByID(id)
}

func (s *KnowledgeService) ListDirectories(projectID string) ([]*knowledge.Directory, error) {
	return s.dirRepo.FindByProjectID(projectID)
}

func (s *KnowledgeService) ListSubDirectories(parentID string) ([]*knowledge.Directory, error) {
	return s.dirRepo.FindByParentID(parentID)
}

func (s *KnowledgeService) DeleteDirectory(id string) error {
	docs, _ := s.docRepo.FindByDirectoryID(id)
	for _, doc := range docs {
		_ = s.DeleteDocument(doc.ID)
	}
	subs, _ := s.dirRepo.FindByParentID(id)
	for _, sub := range subs {
		_ = s.DeleteDirectory(sub.ID)
	}
	return s.dirRepo.Delete(id)
}

func (s *KnowledgeService) Search(query *knowledge.SearchQuery) ([]*knowledge.SearchResult, error) {
	return s.indexRepo.Search(query)
}
