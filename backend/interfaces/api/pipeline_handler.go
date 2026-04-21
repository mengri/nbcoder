package api

import (
	"github.com/gin-gonic/gin"
	pipelineApp "github.com/mengri/nbcoder/application/pipeline"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/pkg/response"
	"github.com/mengri/nbcoder/pkg/uid"
)

type PipelineHandler struct {
	pipelineService *pipelineApp.PipelineService
}

func NewPipelineHandler(pipelineService *pipelineApp.PipelineService) *PipelineHandler {
	return &PipelineHandler{
		pipelineService: pipelineService,
	}
}

func (h *PipelineHandler) RegisterRoutes(router *gin.RouterGroup) {
	pipelines := router.Group("/pipelines")
	{
		pipelines.POST("", h.CreatePipeline)
		pipelines.POST("/:id/start", h.StartNextStage)
		pipelines.POST("/:id/complete", h.CompleteStage)
		pipelines.GET("/:id", h.GetPipeline)
	}
}

func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req dto.CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id := uid.NewID()
	aggregate, err := h.pipelineService.CreatePipeline(id, req.CardID)
	if err != nil {
		response.Error(c, "创建流水线失败："+err.Error())
		return
	}
	response.Created(c, dto.PipelineResponse{
		ID:     aggregate.Pipeline.ID,
		CardID: aggregate.Pipeline.CardID,
	})
}

func (h *PipelineHandler) StartNextStage(c *gin.Context) {
	pipelineID := c.Param("id")
	if err := h.pipelineService.StartNextStage(pipelineID); err != nil {
		response.Error(c, "启动流水线失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *PipelineHandler) CompleteStage(c *gin.Context) {
	pipelineID := c.Param("id")
	if err := h.pipelineService.CompleteStage(pipelineID); err != nil {
		response.Error(c, "完成流水线阶段失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	pipelineID := c.Param("id")
	pl, err := h.pipelineService.GetPipeline(pipelineID)
	if err != nil {
		response.Error(c, "获取流水线失败："+err.Error())
		return
	}
	if pl == nil {
		response.NotFound(c, "流水线不存在")
		return
	}
	response.Success(c, pl)
}
