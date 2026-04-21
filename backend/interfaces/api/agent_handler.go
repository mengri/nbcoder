package api

import (
	"github.com/gin-gonic/gin"
	agentApp "github.com/mengri/nbcoder/application/agent"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/pkg/response"
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
		response.BadRequest(c, err.Error())
		return
	}
	id := uid.NewID()
	aggregate, err := h.agentService.CreateTask(id, req.Name, req.Description)
	if err != nil {
		response.Error(c, "创建任务失败："+err.Error())
		return
	}
	response.Created(c, dto.TaskResponse{
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
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.agentService.AssignTask(taskID, req.AgentID); err != nil {
		response.Error(c, "分配任务失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *AgentHandler) CompleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if err := h.agentService.CompleteTask(taskID); err != nil {
		response.Error(c, "完成任务失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *AgentHandler) FailTask(c *gin.Context) {
	taskID := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.agentService.FailTask(taskID, req.Reason); err != nil {
		response.Error(c, "标记任务失败失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *AgentHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	task, err := h.agentService.GetTask(taskID)
	if err != nil {
		response.Error(c, "获取任务失败："+err.Error())
		return
	}
	if task == nil {
		response.NotFound(c, "任务不存在")
		return
	}
	response.Success(c, dto.TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      string(task.Status),
		AssignedTo:  task.AssignedTo,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
