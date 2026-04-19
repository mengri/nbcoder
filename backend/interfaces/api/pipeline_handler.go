package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pipelineApp "github.com/mengri/nbcoder/application/pipeline"
	"github.com/mengri/nbcoder/application/dto"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uid.NewID()
	aggregate, err := h.pipelineService.CreatePipeline(id, req.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.PipelineResponse{
		ID:     aggregate.Pipeline.ID,
		CardID: aggregate.Pipeline.CardID,
	})
}

func (h *PipelineHandler) StartNextStage(c *gin.Context) {
	pipelineID := c.Param("id")
	if err := h.pipelineService.StartNextStage(pipelineID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "stage started"})
}

func (h *PipelineHandler) CompleteStage(c *gin.Context) {
	pipelineID := c.Param("id")
	if err := h.pipelineService.CompleteStage(pipelineID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "stage completed"})
}

func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	pipelineID := c.Param("id")
	pl, err := h.pipelineService.GetPipeline(pipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pipeline not found"})
		return
	}
	c.JSON(http.StatusOK, pl)
}
