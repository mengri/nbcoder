package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mengri/nbcoder/application/dto"
	projectApp "github.com/mengri/nbcoder/application/project"
	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/pkg/response"
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
		projects.GET("/:name", h.GetProject)
		projects.PUT("/:name", h.UpdateProject)
		projects.DELETE("/:name", h.DeleteProject)
		projects.POST("/:name/archive", h.ArchiveProject)
		projects.POST("/:name/activate", h.ActivateProject)
		projects.GET("/:name/configs", h.GetConfigs)
		projects.PUT("/:name/configs", h.SetConfig)
		projects.GET("/:name/configs/history", h.GetConfigHistory)
		projects.GET("/:name/standards", h.GetStandards)
		projects.PUT("/:name/standards", h.UpdateStandards)
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	project, err := h.projectService.CreateProject(req.Name, req.Description, req.RepoURL)
	if err != nil {
		response.Error(c, "创建项目失败："+err.Error())
		return
	}
	response.Created(c, toProjectResponse(project))
}

func (h *ProjectHandler) InitProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.projectService.InitProject(req.Name, req.Description, req.RepoURL, req.BranchStrategy, req.TechStack, req.CodingConventions)
	if err != nil {
		response.Error(c, "初始化项目失败："+err.Error())
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
	response.Created(c, resp)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	name := c.Param("name")
	project, err := h.projectService.GetProject(name)
	if err != nil {
		response.Error(c, "获取项目失败："+err.Error())
		return
	}
	if project == nil {
		response.NotFound(c, "项目不存在")
		return
	}
	response.Success(c, toProjectResponse(project))
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		response.Error(c, "获取项目列表失败："+err.Error())
		return
	}
	result := make([]dto.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		result = append(result, toProjectResponse(p))
	}
	response.Success(c, result)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	name := c.Param("name")
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	project, err := h.projectService.UpdateProject(name, req.Name, req.Description, req.RepoURL)
	if err != nil {
		response.Error(c, "更新项目失败："+err.Error())
		return
	}
	response.Success(c, toProjectResponse(project))
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	name := c.Param("name")
	if err := h.projectService.DeleteProject(name); err != nil {
		response.Error(c, "删除项目失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *ProjectHandler) ArchiveProject(c *gin.Context) {
	name := c.Param("name")
	if err := h.projectService.ArchiveProject(name); err != nil {
		response.Error(c, "归档项目失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *ProjectHandler) ActivateProject(c *gin.Context) {
	name := c.Param("name")
	if err := h.projectService.ActivateProject(name); err != nil {
		response.Error(c, "激活项目失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *ProjectHandler) GetConfigs(c *gin.Context) {
	name := c.Param("name")
	configs, err := h.projectService.GetConfigs(name)
	if err != nil {
		response.Error(c, "获取配置失败："+err.Error())
		return
	}
	result := make([]dto.ConfigResponse, 0, len(configs))
	for _, cfg := range configs {
		result = append(result, dto.ConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
	}
	response.Success(c, result)
}

func (h *ProjectHandler) SetConfig(c *gin.Context) {
	name := c.Param("name")
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cfg, err := h.projectService.SetConfig(name, req.Key, req.Value)
	if err != nil {
		response.Error(c, "设置配置失败："+err.Error())
		return
	}
	response.Success(c, dto.ConfigResponse{ID: cfg.ID, Key: cfg.Key, Value: cfg.Value})
}

func (h *ProjectHandler) GetStandards(c *gin.Context) {
	name := c.Param("name")
	std, err := h.projectService.GetStandards(name)
	if err != nil {
		response.Error(c, "获取规范失败："+err.Error())
		return
	}
	if std == nil {
		response.NotFound(c, "规范不存在")
		return
	}
	response.Success(c, dto.StandardsResponse{
		ID:                std.ID,
		BranchStrategy:    std.BranchStrategy,
		TechStack:         std.TechStack,
		CodingConventions: std.CodingConventions,
	})
}

func (h *ProjectHandler) GetConfigHistory(c *gin.Context) {
	name := c.Param("name")
	logs, err := h.projectService.GetConfigHistory(name)
	if err != nil {
		response.Error(c, "获取配置历史失败："+err.Error())
		return
	}
	result := make([]dto.ConfigChangeLogResponse, 0, len(logs))
	for _, l := range logs {
		result = append(result, dto.ConfigChangeLogResponse{
			ID:          l.ID,
			ProjectName: l.ProjectName,
			ConfigKey:   l.ConfigKey,
			OldValue:    l.OldValue,
			NewValue:    l.NewValue,
			ChangedAt:   l.ChangedAt.Format("2006-01-02T15:04:05Z"),
			ChangedBy:   l.ChangedBy,
		})
	}
	response.Success(c, result)
}

func (h *ProjectHandler) UpdateStandards(c *gin.Context) {
	name := c.Param("name")
	var req struct {
		BranchStrategy    string `json:"branchStrategy"`
		TechStack         string `json:"techStack"`
		CodingConventions string `json:"codingConventions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	std, err := h.projectService.UpdateStandards(name, req.BranchStrategy, req.TechStack, req.CodingConventions)
	if err != nil {
		response.Error(c, "更新规范失败："+err.Error())
		return
	}
	response.Success(c, dto.StandardsResponse{
		ID:                std.ID,
		BranchStrategy:    std.BranchStrategy,
		TechStack:         std.TechStack,
		CodingConventions: std.CodingConventions,
	})
}

func toProjectResponse(p *project.Project) dto.ProjectResponse {
	return dto.ProjectResponse{
		Name:        p.Name,
		Description: p.Description,
		RepoURL:     p.RepoURL,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
