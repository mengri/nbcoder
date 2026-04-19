package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mengri/nbcoder/application/dto"
	projectApp "github.com/mengri/nbcoder/application/project"
	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/pkg/uid"
)

type ProjectHandler struct {
	projectService *projectApp.ProjectService
}

func NewProjectHandler(projectService *projectApp.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	projects := router.Group("/projects")
	{
		projects.POST("", h.CreateProject)
		projects.POST("/init", h.InitProject)
		projects.GET("", h.ListProjects)
		projects.GET("/:id", h.GetProject)
		projects.PUT("/:id", h.UpdateProject)
		projects.DELETE("/:id", h.DeleteProject)
		projects.POST("/:id/archive", h.ArchiveProject)
		projects.POST("/:id/activate", h.ActivateProject)
		projects.GET("/:id/configs", h.GetConfigs)
		projects.PUT("/:id/configs", h.SetConfig)
		projects.GET("/:id/configs/history", h.GetConfigHistory)
		projects.GET("/:id/standards", h.GetStandards)
		projects.PUT("/:id/standards", h.UpdateStandards)
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uid.NewID()
	project, err := h.projectService.CreateProject(id, req.Name, req.Description, req.RepoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toProjectResponse(project))
}

func (h *ProjectHandler) InitProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.projectService.InitProject(req.Name, req.Description, req.RepoURL, req.BranchStrategy, req.TechStack, req.CodingConventions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := dto.InitProjectResponse{
		Project:     toProjectResponse(result.Project),
		Directories: result.Directories.Dirs,
	}
	for _, cfg := range result.Configs {
		resp.Configs = append(resp.Configs, dto.ConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
	}
	if result.Standards != nil {
		resp.Standards = &dto.StandardsResponse{
			ID:                result.Standards.ID,
			BranchStrategy:    result.Standards.BranchStrategy,
			TechStack:         result.Standards.TechStack,
			CodingConventions: result.Standards.CodingConventions,
		}
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	project, err := h.projectService.GetProject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, toProjectResponse(project))
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]dto.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		result = append(result, toProjectResponse(p))
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.projectService.UpdateProject(id, req.Name, req.Description, req.RepoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toProjectResponse(project))
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.projectService.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

func (h *ProjectHandler) ArchiveProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.projectService.ArchiveProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project archived"})
}

func (h *ProjectHandler) ActivateProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.projectService.ActivateProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project activated"})
}

func (h *ProjectHandler) GetConfigs(c *gin.Context) {
	id := c.Param("id")
	configs, err := h.projectService.GetConfigs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]dto.ConfigResponse, 0, len(configs))
	for _, cfg := range configs {
		result = append(result, dto.ConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProjectHandler) SetConfig(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cfg, err := h.projectService.SetConfig(id, req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.ConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
}

func (h *ProjectHandler) GetStandards(c *gin.Context) {
	id := c.Param("id")
	std, err := h.projectService.GetStandards(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if std == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "standards not found"})
		return
	}
	c.JSON(http.StatusOK, dto.StandardsResponse{
		ID:                std.ID,
		BranchStrategy:    std.BranchStrategy,
		TechStack:         std.TechStack,
		CodingConventions: std.CodingConventions,
	})
}

func (h *ProjectHandler) GetConfigHistory(c *gin.Context) {
	id := c.Param("id")
	logs, err := h.projectService.GetConfigHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]dto.ConfigChangeLogResponse, 0, len(logs))
	for _, l := range logs {
		result = append(result, dto.ConfigChangeLogResponse{
			ID:        l.ID,
			ProjectID: l.ProjectID,
			ConfigKey: l.ConfigKey,
			OldValue:  l.OldValue,
			NewValue:  l.NewValue,
			ChangedAt: l.ChangedAt.Format("2006-01-02T15:04:05Z"),
			ChangedBy: l.ChangedBy,
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProjectHandler) UpdateStandards(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		BranchStrategy    string `json:"branch_strategy"`
		TechStack         string `json:"tech_stack"`
		CodingConventions string `json:"coding_conventions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	std, err := h.projectService.UpdateStandards(id, req.BranchStrategy, req.TechStack, req.CodingConventions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.StandardsResponse{
		ID:                std.ID,
		BranchStrategy:    std.BranchStrategy,
		TechStack:         std.TechStack,
		CodingConventions: std.CodingConventions,
	})
}

func toProjectResponse(p *project.Project) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		RepoURL:     p.RepoURL,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
