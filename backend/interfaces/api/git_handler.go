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
		gitRoutes.POST("/pull-requests/auto", h.CreatePullRequestWithDescription)
		gitRoutes.GET("/pull-requests/:id", h.GetPullRequest)
		gitRoutes.POST("/pull-requests/:id/merge", h.MergePullRequest)
		gitRoutes.POST("/pull-requests/:id/close", h.ClosePullRequest)
		gitRoutes.POST("/branches/validate", h.ValidateBranch)
	}
}

func (h *GitHandler) CreatePullRequest(c *gin.Context) {
	var req struct {
		Title        string `json:"title" binding:"required"`
		SourceBranch string `json:"source_branch" binding:"required"`
		TargetBranch string `json:"target_branch" binding:"required"`
		ProjectID    string `json:"project_id"`
		Author       string `json:"author"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pr, err := h.gitService.CreatePullRequest(req.Title, req.SourceBranch, req.TargetBranch, req.ProjectID, req.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toPullRequestResponse(pr))
}

func (h *GitHandler) CreatePullRequestWithDescription(c *gin.Context) {
	var req struct {
		Title        string `json:"title" binding:"required"`
		SourceBranch string `json:"source_branch" binding:"required"`
		TargetBranch string `json:"target_branch" binding:"required"`
		ProjectID    string `json:"project_id"`
		Author       string `json:"author"`
		Commits      []struct {
			Hash    string `json:"hash"`
			Message string `json:"message"`
			Author  string `json:"author"`
		} `json:"commits"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	commits := make([]*git.Commit, len(req.Commits))
	for i, c := range req.Commits {
		commits[i] = &git.Commit{Hash: c.Hash, Message: c.Message, Author: c.Author}
	}
	pr, err := h.gitService.CreatePullRequestWithDescription(req.Title, req.SourceBranch, req.TargetBranch, req.ProjectID, req.Author, commits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toPullRequestResponse(pr))
}

func (h *GitHandler) GetPullRequest(c *gin.Context) {
	id := c.Param("id")
	pr, err := h.gitService.GetPullRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pull request not found"})
		return
	}
	c.JSON(http.StatusOK, toPullRequestResponse(pr))
}

func (h *GitHandler) MergePullRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.gitService.MergePullRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pull request merged"})
}

func (h *GitHandler) ClosePullRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.gitService.ClosePullRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pull request closed"})
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

func toPullRequestResponse(pr *git.PullRequest) gin.H {
	resp := gin.H{
		"id":             pr.ID,
		"title":          pr.Title,
		"description":    pr.Description,
		"source_branch":  pr.SourceBranch,
		"target_branch":  pr.TargetBranch,
		"status":         string(pr.Status),
		"project_id":     pr.ProjectID,
		"author":         pr.Author,
		"created_at":     pr.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":     pr.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if pr.GeneratedDesc != "" {
		resp["generated_desc"] = pr.GeneratedDesc
	}
	return resp
}
