package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mengri/nbcoder/domain/knowledge"
	knowledgeApp "github.com/mengri/nbcoder/application/knowledge"
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
	docs := router.Group("/knowledge")
	{
		docs.POST("/documents", h.UploadDocument)
		docs.POST("/search", h.Search)
		docs.GET("/documents/:id", h.GetDocument)
	}
}

func (h *KnowledgeHandler) UploadDocument(c *gin.Context) {
	var doc knowledge.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.knowledgeService.UploadDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doc)
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
	c.JSON(http.StatusOK, doc)
}
