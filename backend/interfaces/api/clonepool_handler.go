package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clonepoolApp "github.com/mengri/nbcoder/application/clonepool"
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
		RepositoryID string `json:"repository_id"`
		TaskID       string `json:"task_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inst, err := h.clonePoolService.AcquireInstance(req.RepositoryID, req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inst)
}

func (h *ClonePoolHandler) Release(c *gin.Context) {
	instanceID := c.Param("id")
	if err := h.clonePoolService.ReleaseInstance(instanceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "instance released"})
}

func (h *ClonePoolHandler) CreateInstance(c *gin.Context) {
	var req struct {
		RepositoryID string `json:"repository_id" binding:"required"`
		RepoURL      string `json:"repo_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inst, err := h.clonePoolService.CreateCloneInstance(c.Request.Context(), req.RepositoryID, req.RepoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inst)
}

func (h *ClonePoolHandler) CommitChanges(c *gin.Context) {
	instanceID := c.Param("id")
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.clonePoolService.CommitChanges(c.Request.Context(), instanceID, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "changes committed"})
}

func (h *ClonePoolHandler) GetStatus(c *gin.Context) {
	instanceID := c.Param("id")
	status, err := h.clonePoolService.GetStatus(c.Request.Context(), instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
