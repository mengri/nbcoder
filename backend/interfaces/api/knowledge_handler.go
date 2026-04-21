package api

import (
	"github.com/gin-gonic/gin"
	knowledgeApp "github.com/mengri/nbcoder/application/knowledge"
	"github.com/mengri/nbcoder/domain/knowledge"
	"github.com/mengri/nbcoder/pkg/response"
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
		ProjectName string `json:"projectName" binding:"required"`
		DirectoryID string `json:"directoryId"`
		Content     string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	doc, err := h.knowledgeService.UploadDocument(req.Name, req.Path, req.ProjectName, req.DirectoryID, req.Content)
	if err != nil {
		response.Error(c, "上传文档失败："+err.Error())
		return
	}
	response.Created(c, toDocumentResponse(doc))
}

func (h *KnowledgeHandler) BatchUploadDocuments(c *gin.Context) {
	var req struct {
		Documents []struct {
			Name        string `json:"name" binding:"required"`
			Path        string `json:"path"`
			ProjectName string `json:"projectName" binding:"required"`
			DirectoryID string `json:"directoryId"`
			Content     string `json:"content"`
		} `json:"documents" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	docs := make([]struct{ Name, Path, ProjectName, DirectoryID, Content string }, len(req.Documents))
	for i, d := range req.Documents {
		docs[i] = struct{ Name, Path, ProjectName, DirectoryID, Content string }{
			Name: d.Name, Path: d.Path, ProjectName: d.ProjectName, DirectoryID: d.DirectoryID, Content: d.Content,
		}
	}
	result, err := h.knowledgeService.BatchUploadDocuments(docs)
	if err != nil {
		response.Error(c, "批量上传文档失败："+err.Error())
		return
	}
	resp := make([]gin.H, 0, len(result))
	for _, doc := range result {
		resp = append(resp, toDocumentResponse(doc))
	}
	response.Created(c, resp)
}

func (h *KnowledgeHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	doc, err := h.knowledgeService.GetDocument(id, projectName)
	if err != nil {
		response.Error(c, "获取文档失败："+err.Error())
		return
	}
	if doc == nil {
		response.NotFound(c, "文档不存在")
		return
	}
	response.Success(c, toDocumentResponse(doc))
}

func (h *KnowledgeHandler) ListDocuments(c *gin.Context) {
	projectName := c.Query("projectName")
	docs, err := h.knowledgeService.ListDocuments(projectName)
	if err != nil {
		response.Error(c, "获取文档列表失败："+err.Error())
		return
	}
	resp := make([]gin.H, 0, len(docs))
	for _, doc := range docs {
		resp = append(resp, toDocumentResponse(doc))
	}
	response.Success(c, resp)
}

func (h *KnowledgeHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	if err := h.knowledgeService.DeleteDocument(id, projectName); err != nil {
		response.Error(c, "删除文档失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *KnowledgeHandler) CreateDirectory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		ParentID    string `json:"parentId"`
		ProjectName string `json:"projectName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	dir, err := h.knowledgeService.CreateDirectory(req.Name, req.ParentID, req.ProjectName)
	if err != nil {
		response.Error(c, "创建目录失败："+err.Error())
		return
	}
	response.Created(c, toDirectoryResponse(dir))
}

func (h *KnowledgeHandler) GetDirectory(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	dir, err := h.knowledgeService.GetDirectory(id, projectName)
	if err != nil {
		response.Error(c, "获取目录失败："+err.Error())
		return
	}
	if dir == nil {
		response.NotFound(c, "目录不存在")
		return
	}
	response.Success(c, toDirectoryResponse(dir))
}

func (h *KnowledgeHandler) ListDirectories(c *gin.Context) {
	projectName := c.Query("projectName")
	dirs, err := h.knowledgeService.ListDirectories(projectName)
	if err != nil {
		response.Error(c, "获取目录列表失败："+err.Error())
		return
	}
	resp := make([]gin.H, 0, len(dirs))
	for _, dir := range dirs {
		resp = append(resp, toDirectoryResponse(dir))
	}
	response.Success(c, resp)
}

func (h *KnowledgeHandler) ListDocumentsByDirectory(c *gin.Context) {
	dirID := c.Param("id")
	docs, err := h.knowledgeService.ListDocumentsByDirectory(dirID)
	if err != nil {
		response.Error(c, "获取目录文档失败："+err.Error())
		return
	}
	resp := make([]gin.H, 0, len(docs))
	for _, doc := range docs {
		resp = append(resp, toDocumentResponse(doc))
	}
	response.Success(c, resp)
}

func (h *KnowledgeHandler) ListSubDirectories(c *gin.Context) {
	parentID := c.Param("id")
	dirs, err := h.knowledgeService.ListSubDirectories(parentID)
	if err != nil {
		response.Error(c, "获取子目录失败："+err.Error())
		return
	}
	resp := make([]gin.H, 0, len(dirs))
	for _, dir := range dirs {
		resp = append(resp, toDirectoryResponse(dir))
	}
	response.Success(c, resp)
}

func (h *KnowledgeHandler) DeleteDirectory(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	if err := h.knowledgeService.DeleteDirectory(id, projectName); err != nil {
		response.Error(c, "删除目录失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *KnowledgeHandler) Search(c *gin.Context) {
	var query knowledge.SearchQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	results, err := h.knowledgeService.Search(&query)
	if err != nil {
		response.Error(c, "搜索失败："+err.Error())
		return
	}
	response.Success(c, results)
}

func toDocumentResponse(doc *knowledge.Document) gin.H {
	return gin.H{
		"id":           doc.ID,
		"name":         doc.Name,
		"path":         doc.Path,
		"projectName":  doc.ProjectName,
		"directoryId":  doc.DirectoryID,
		"version":      doc.Version,
		"createdAt":    doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt":    doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toDirectoryResponse(dir *knowledge.Directory) gin.H {
	return gin.H{
		"id":         dir.ID,
		"name":       dir.Name,
		"parentId":   dir.ParentID,
		"projectName":  dir.ProjectName,
		"path":       dir.Path,
		"createdAt":  dir.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt":  dir.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
