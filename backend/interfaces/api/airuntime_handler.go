package api

import (
	"github.com/gin-gonic/gin"
	airuntimeApp "github.com/mengri/nbcoder/application/airuntime"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/domain/airuntime"
	"github.com/mengri/nbcoder/pkg/response"
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
		ai.POST("/call", h.CallModel)
	}
}

func (h *AIRuntimeHandler) RegisterProvider(c *gin.Context) {
	var req dto.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	provider := &airuntime.Provider{
		ID:        "provider-" + req.Name,
		Name:      req.Name,
		APIKeyRef: req.APIKeyRef,
		Models:    []*airuntime.Model{},
	}

	if err := h.aiRuntimeService.RegisterProvider(provider); err != nil {
		response.Error(c, "注册提供商失败："+err.Error())
		return
	}

	response.Created(c, provider)
}

func (h *AIRuntimeHandler) GetProvider(c *gin.Context) {
	id := c.Param("id")
	provider, ok := h.aiRuntimeService.GetProvider(id)
	if !ok {
		response.NotFound(c, "提供商不存在")
		return
	}
	response.Success(c, provider)
}

func (h *AIRuntimeHandler) CallModel(c *gin.Context) {
	var req dto.CallModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	messages := make([]airuntime.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = airuntime.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := h.aiRuntimeService.CallModel(c.Request.Context(), req.ProviderID, req.ModelID, messages, req.AgentID)
	if err != nil {
		response.Error(c, "调用模型失败："+err.Error())
		return
	}

	response.Success(c, resp)
}
