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

func (s *KnowledgeService) UploadDocument(name, path, projectName, directoryID, content string) (*knowledge.Document, error) {
	doc := knowledge.NewDocument(uid.NewID(), name, path, projectName)
	if directoryID != "" {
		dir, err := s.dirRepo.FindByID(directoryID, projectName)
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
	Name, Path, ProjectName, DirectoryID, Content string
}) ([]*knowledge.Document, error) {
	var results []*knowledge.Document
	for _, d := range docs {
		doc, err := s.UploadDocument(d.Name, d.Path, d.ProjectName, d.DirectoryID, d.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to upload document %s: %w", d.Name, err)
		}
		results = append(results, doc)
	}
	return results, nil
}

func (s *KnowledgeService) GetDocument(id, projectName string) (*knowledge.Document, error) {
	return s.docRepo.FindByID(id, projectName)
}

func (s *KnowledgeService) ListDocuments(projectName string) ([]*knowledge.Document, error) {
	return s.docRepo.FindByProjectName(projectName)
}

func (s *KnowledgeService) ListDocumentsByDirectory(directoryID string) ([]*knowledge.Document, error) {
	return s.docRepo.FindByDirectoryID(directoryID)
}

func (s *KnowledgeService) DeleteDocument(id, projectName string) error {
	doc, err := s.docRepo.FindByID(id, projectName)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("document not found: %s", id)
	}
	_ = s.chunkRepo.DeleteByDocumentID(id)
	_ = s.indexRepo.DeleteByDocumentID(id)
	return s.docRepo.Delete(id, projectName)
}

func (s *KnowledgeService) CreateDirectory(name, parentID, projectName string) (*knowledge.Directory, error) {
	if parentID != "" {
		parent, err := s.dirRepo.FindByID(parentID, projectName)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("parent directory not found: %s", parentID)
		}
	}
	dir := knowledge.NewDirectory(uid.NewID(), name, parentID, projectName)
	if err := dir.Validate(); err != nil {
		return nil, err
	}
	if err := s.dirRepo.Save(dir); err != nil {
		return nil, err
	}
	return dir, nil
}

func (s *KnowledgeService) GetDirectory(id, projectName string) (*knowledge.Directory, error) {
	return s.dirRepo.FindByID(id, projectName)
}

func (s *KnowledgeService) ListDirectories(projectName string) ([]*knowledge.Directory, error) {
	return s.dirRepo.FindByProjectName(projectName)
}

func (s *KnowledgeService) ListSubDirectories(parentID string) ([]*knowledge.Directory, error) {
	return s.dirRepo.FindByParentID(parentID)
}

func (s *KnowledgeService) DeleteDirectory(id, projectName string) error {
	docs, _ := s.docRepo.FindByDirectoryID(id)
	for _, doc := range docs {
		_ = s.DeleteDocument(doc.ID, projectName)
	}
	subs, _ := s.dirRepo.FindByParentID(id)
	for _, sub := range subs {
		_ = s.DeleteDirectory(sub.ID, projectName)
	}
	return s.dirRepo.Delete(id, projectName)
}

func (s *KnowledgeService) Search(query *knowledge.SearchQuery) ([]*knowledge.SearchResult, error) {
	return s.indexRepo.Search(query)
}
