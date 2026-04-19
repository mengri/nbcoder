package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	airuntimeApp "github.com/mengri/nbcoder/application/airuntime"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/pkg/uid"
)

type AIRuntimeHandler struct {
	aiRuntimeService *airuntimeApp.AIRuntimeService
}

func NewAIRuntimeHandler(aiRuntimeService *airuntimeApp.AIRuntimeService) *AIRuntimeHandler {
	return &AIRuntimeHandler{
		aiRuntimeService: aiRuntimeService,
	}
}

func (h *AIRuntimeHandler) RegisterRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai-runtime")
	{
		ai.POST("/providers", h.RegisterProvider)
		ai.GET("/providers/:id", h.GetProvider)
	}
}

func (h *AIRuntimeHandler) RegisterProvider(c *gin.Context) {
	var req dto.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	provider := &airuntime.Provider{
		ID:        uid.NewID(),
		Name:      req.Name,
		APIKeyRef: req.APIKeyRef,
	}
	if err := h.aiRuntimeService.RegisterProvider(provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, provider)
}

func (h *AIRuntimeHandler) GetProvider(c *gin.Context) {
	id := c.Param("id")
	provider, ok := h.aiRuntimeService.GetProvider(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
		return
	}
	c.JSON(http.StatusOK, provider)
}
