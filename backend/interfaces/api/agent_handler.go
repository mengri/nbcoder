package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	agentApp "github.com/mengri/nbcoder/application/agent"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/pkg/uid"
)

type AgentHandler struct {
	agentService *agentApp.AgentService
}

func NewAgentHandler(agentService *agentApp.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

func (h *AgentHandler) RegisterRoutes(router *gin.RouterGroup) {
	agents := router.Group("/agents")
	{
		agents.POST("/tasks", h.CreateTask)
		agents.POST("/tasks/:id/assign", h.AssignTask)
		agents.POST("/tasks/:id/complete", h.CompleteTask)
		agents.POST("/tasks/:id/fail", h.FailTask)
		agents.GET("/tasks/:id", h.GetTask)
	}
}

func (h *AgentHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uid.NewID()
	aggregate, err := h.agentService.CreateTask(id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.TaskResponse{
		ID:          aggregate.Task.ID,
		Name:        aggregate.Task.Name,
		Description: aggregate.Task.Description,
		Status:      string(aggregate.Task.Status),
		CreatedAt:   aggregate.Task.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   aggregate.Task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *AgentHandler) AssignTask(c *gin.Context) {
	taskID := c.Param("id")
	var req dto.AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.agentService.AssignTask(taskID, req.AgentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task assigned"})
}

func (h *AgentHandler) CompleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if err := h.agentService.CompleteTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task completed"})
}

func (h *AgentHandler) FailTask(c *gin.Context) {
	taskID := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.agentService.FailTask(taskID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task failed"})
}

func (h *AgentHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	task, err := h.agentService.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, dto.TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      string(task.Status),
		AssignedTo:  task.AssignedTo,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
