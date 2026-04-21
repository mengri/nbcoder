package api

import (
	"github.com/gin-gonic/gin"
	clonepoolApp "github.com/mengri/nbcoder/application/clonepool"
	"github.com/mengri/nbcoder/pkg/response"
)

type ClonePoolHandler struct {
	clonePoolService *clonepoolApp.ClonePoolService
}

func NewClonePoolHandler(clonePoolService *clonepoolApp.ClonePoolService) *ClonePoolHandler {
	return &ClonePoolHandler{
		clonePoolService: clonePoolService,
	}
}

func (h *ClonePoolHandler) RegisterRoutes(router *gin.RouterGroup) {
	clones := router.Group("/clone-pool")
	{
		clones.POST("/acquire", h.Acquire)
		clones.POST("/:id/release", h.Release)
		clones.POST("/create", h.CreateInstance)
		clones.POST("/:id/commit", h.CommitChanges)
		clones.GET("/:id/status", h.GetStatus)
	}
}

func (h *ClonePoolHandler) Acquire(c *gin.Context) {
	var req struct {
		RepositoryID string `json:"repositoryId"`
		TaskID       string `json:"taskId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	inst, err := h.clonePoolService.AcquireInstance(req.RepositoryID, req.TaskID)
	if err != nil {
		response.Error(c, "获取实例失败："+err.Error())
		return
	}
	response.Success(c, inst)
}

func (h *ClonePoolHandler) Release(c *gin.Context) {
	instanceID := c.Param("id")
	if err := h.clonePoolService.ReleaseInstance(instanceID); err != nil {
		response.Error(c, "释放实例失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *ClonePoolHandler) CreateInstance(c *gin.Context) {
	var req struct {
		RepositoryID string `json:"repositoryId" binding:"required"`
		RepoURL      string `json:"repoUrl" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	inst, err := h.clonePoolService.CreateCloneInstance(c.Request.Context(), req.RepositoryID, req.RepoURL)
	if err != nil {
		response.Error(c, "创建实例失败："+err.Error())
		return
	}

	response.Created(c, inst)
}

func (h *ClonePoolHandler) CommitChanges(c *gin.Context) {
	instanceID := c.Param("id")
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.clonePoolService.CommitChanges(c.Request.Context(), instanceID, req.Message); err != nil {
		response.Error(c, "提交更改失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ClonePoolHandler) GetStatus(c *gin.Context) {
	instanceID := c.Param("id")
	status, err := h.clonePoolService.GetStatus(c.Request.Context(), instanceID)
	if err != nil {
		response.Error(c, "获取状态失败："+err.Error())
		return
	}

	response.Success(c, status)
}
