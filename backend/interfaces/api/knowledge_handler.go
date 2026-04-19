package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	knowledgeApp "github.com/mengri/nbcoder/application/knowledge"
	"github.com/mengri/nbcoder/domain/knowledge"
)

type KnowledgeHandler struct {
	knowledgeService *knowledgeApp.KnowledgeService
}

func NewKnowledgeHandler(knowledgeService *knowledgeApp.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{
		knowledgeService: knowledgeService,
	}
}

func (h *KnowledgeHandler) RegisterRoutes(router *gin.RouterGroup) {
	kb := router.Group("/knowledge")
	{
		kb.POST("/documents", h.UploadDocument)
		kb.POST("/documents/batch", h.BatchUploadDocuments)
		kb.GET("/documents", h.ListDocuments)
		kb.GET("/documents/:id", h.GetDocument)
		kb.DELETE("/documents/:id", h.DeleteDocument)
		kb.POST("/directories", h.CreateDirectory)
		kb.GET("/directories", h.ListDirectories)
		kb.GET("/directories/:id", h.GetDirectory)
		kb.GET("/directories/:id/documents", h.ListDocumentsByDirectory)
		kb.GET("/directories/:id/subdirectories", h.ListSubDirectories)
		kb.DELETE("/directories/:id", h.DeleteDirectory)
		kb.POST("/search", h.Search)
	}
}

func (h *KnowledgeHandler) UploadDocument(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Path        string `json:"path"`
		ProjectID   string `json:"project_id" binding:"required"`
		DirectoryID string `json:"directory_id"`
		Content     string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc, err := h.knowledgeService.UploadDocument(req.Name, req.Path, req.ProjectID, req.DirectoryID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toDocumentResponse(doc))
}

func (h *KnowledgeHandler) BatchUploadDocuments(c *gin.Context) {
	var req struct {
		Documents []struct {
			Name        string `json:"name" binding:"required"`
			Path        string `json:"path"`
			ProjectID   string `json:"project_id" binding:"required"`
			DirectoryID string `json:"directory_id"`
			Content     string `json:"content"`
		} `json:"documents" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	docs := make([]struct{ Name, Path, ProjectID, DirectoryID, Content string }, len(req.Documents))
	for i, d := range req.Documents {
		docs[i] = struct{ Name, Path, ProjectID, DirectoryID, Content string }{
			Name: d.Name, Path: d.Path, ProjectID: d.ProjectID, DirectoryID: d.DirectoryID, Content: d.Content,
		}
	}
	result, err := h.knowledgeService.BatchUploadDocuments(docs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]gin.H, 0, len(result))
	for _, doc := range result {
		resp = append(resp, toDocumentResponse(doc))
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *KnowledgeHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	doc, err := h.knowledgeService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	c.JSON(http.StatusOK, toDocumentResponse(doc))
}

func (h *KnowledgeHandler) ListDocuments(c *gin.Context) {
	projectID := c.Query("project_id")
	docs, err := h.knowledgeService.ListDocuments(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]gin.H, 0, len(docs))
	for _, doc := range docs {
		resp = append(resp, toDocumentResponse(doc))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *KnowledgeHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.knowledgeService.DeleteDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}

func (h *KnowledgeHandler) CreateDirectory(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ParentID  string `json:"parent_id"`
		ProjectID string `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dir, err := h.knowledgeService.CreateDirectory(req.Name, req.ParentID, req.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toDirectoryResponse(dir))
}

func (h *KnowledgeHandler) GetDirectory(c *gin.Context) {
	id := c.Param("id")
	dir, err := h.knowledgeService.GetDirectory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if dir == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "directory not found"})
		return
	}
	c.JSON(http.StatusOK, toDirectoryResponse(dir))
}

func (h *KnowledgeHandler) ListDirectories(c *gin.Context) {
	projectID := c.Query("project_id")
	dirs, err := h.knowledgeService.ListDirectories(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]gin.H, 0, len(dirs))
	for _, dir := range dirs {
		resp = append(resp, toDirectoryResponse(dir))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *KnowledgeHandler) ListDocumentsByDirectory(c *gin.Context) {
	dirID := c.Param("id")
	docs, err := h.knowledgeService.ListDocumentsByDirectory(dirID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]gin.H, 0, len(docs))
	for _, doc := range docs {
		resp = append(resp, toDocumentResponse(doc))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *KnowledgeHandler) ListSubDirectories(c *gin.Context) {
	parentID := c.Param("id")
	dirs, err := h.knowledgeService.ListSubDirectories(parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]gin.H, 0, len(dirs))
	for _, dir := range dirs {
		resp = append(resp, toDirectoryResponse(dir))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *KnowledgeHandler) DeleteDirectory(c *gin.Context) {
	id := c.Param("id")
	if err := h.knowledgeService.DeleteDirectory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "directory deleted"})
}

func (h *KnowledgeHandler) Search(c *gin.Context) {
	var query knowledge.SearchQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	results, err := h.knowledgeService.Search(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func toDocumentResponse(doc *knowledge.Document) gin.H {
	return gin.H{
		"id":           doc.ID,
		"name":         doc.Name,
		"path":         doc.Path,
		"project_id":   doc.ProjectID,
		"directory_id": doc.DirectoryID,
		"version":      doc.Version,
		"created_at":   doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":   doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toDirectoryResponse(dir *knowledge.Directory) gin.H {
	return gin.H{
		"id":         dir.ID,
		"name":       dir.Name,
		"parent_id":  dir.ParentID,
		"project_id": dir.ProjectID,
		"path":       dir.Path,
		"created_at": dir.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at": dir.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
