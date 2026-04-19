package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gitApp "github.com/mengri/nbcoder/application/git"
	"github.com/mengri/nbcoder/domain/git"
)

type GitHandler struct {
	gitService *gitApp.GitService
}

func NewGitHandler(gitService *gitApp.GitService) *GitHandler {
	return &GitHandler{
		gitService: gitService,
	}
}

func (h *GitHandler) RegisterRoutes(router *gin.RouterGroup) {
	gitRoutes := router.Group("/git")
	{
		gitRoutes.POST("/pull-requests", h.CreatePullRequest)
		gitRoutes.POST("/pull-requests/:id/merge", h.MergePullRequest)
		gitRoutes.POST("/branches/validate", h.ValidateBranch)
	}
}

func (h *GitHandler) CreatePullRequest(c *gin.Context) {
	var req struct {
		Title        string `json:"title"`
		SourceBranch string `json:"source_branch"`
		TargetBranch string `json:"target_branch"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := generateID()
	pr, err := h.gitService.CreatePullRequest(id, req.Title, req.SourceBranch, req.TargetBranch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pr)
}

func (h *GitHandler) MergePullRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.gitService.MergePullRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pull request merged"})
}

func (h *GitHandler) ValidateBranch(c *gin.Context) {
	var req struct {
		Pattern string `json:"pattern"`
		Branch  string `json:"branch"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	policy := &git.BranchPolicy{AllowedPattern: req.Pattern}
	valid := h.gitService.ValidateBranch(policy, req.Branch)
	c.JSON(http.StatusOK, gin.H{"valid": valid})
}
