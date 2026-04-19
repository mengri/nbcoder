package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	projectApp "github.com/mengri/nbcoder/application/project"
	"github.com/mengri/nbcoder/application/dto"
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
		projects.GET("/:id", h.GetProject)
		projects.GET("", h.ListProjects)
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
	c.JSON(http.StatusCreated, dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		RepoURL:     project.RepoURL,
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
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
	c.JSON(http.StatusOK, dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		RepoURL:     project.RepoURL,
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var result []dto.ProjectResponse
	for _, p := range projects {
		result = append(result, dto.ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			RepoURL:     p.RepoURL,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	c.JSON(http.StatusOK, result)
}
