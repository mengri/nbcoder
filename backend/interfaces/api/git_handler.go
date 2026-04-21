package api

import (
	"github.com/gin-gonic/gin"
	gitApp "github.com/mengri/nbcoder/application/git"
	"github.com/mengri/nbcoder/domain/git"
	"github.com/mengri/nbcoder/pkg/response"
)

type GitHandler struct {
	gitService    *gitApp.GitService
	reviewService *gitApp.ReviewService
}

func NewGitHandler(gitService *gitApp.GitService) *GitHandler {
	return &GitHandler{
		gitService: gitService,
	}
}

func NewGitHandlerWithReview(gitService *gitApp.GitService, reviewService *gitApp.ReviewService) *GitHandler {
	return &GitHandler{
		gitService:    gitService,
		reviewService: reviewService,
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
		gitRoutes.POST("/pull-requests/:id/squash-merge", h.SquashMergePullRequest)
		gitRoutes.POST("/pull-requests/:id/reviews", h.CreateReview)
		gitRoutes.PUT("/pull-requests/:id/reviews/:reviewId/approve", h.ApproveReview)
		gitRoutes.PUT("/pull-requests/:id/reviews/:reviewId/reject", h.RejectReview)
		gitRoutes.GET("/pull-requests/:id/reviews", h.GetReviews)
		gitRoutes.POST("/branches/validate", h.ValidateBranch)
	}
}

func (h *GitHandler) CreatePullRequest(c *gin.Context) {
	var req struct {
		Title        string `json:"title" binding:"required"`
		SourceBranch string `json:"sourceBranch" binding:"required"`
		TargetBranch string `json:"targetBranch" binding:"required"`
		ProjectID    string `json:"projectId"`
		Author       string `json:"author"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	pr, err := h.gitService.CreatePullRequest(req.Title, req.SourceBranch, req.TargetBranch, req.ProjectID, req.Author)
	if err != nil {
		response.Error(c, "创建 Pull Request 失败："+err.Error())
		return
	}
	response.Created(c, toPullRequestResponse(pr))
}

func (h *GitHandler) CreatePullRequestWithDescription(c *gin.Context) {
	var req struct {
		Title        string `json:"title" binding:"required"`
		SourceBranch string `json:"sourceBranch" binding:"required"`
		TargetBranch string `json:"targetBranch" binding:"required"`
		ProjectID    string `json:"projectId"`
		Author       string `json:"author"`
		Commits      []struct {
			Hash    string `json:"hash"`
			Message string `json:"message"`
			Author  string `json:"author"`
		} `json:"commits"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	commits := make([]*git.Commit, len(req.Commits))
	for i, c := range req.Commits {
		commits[i] = &git.Commit{Hash: c.Hash, Message: c.Message, Author: c.Author}
	}
	pr, err := h.gitService.CreatePullRequestWithDescription(req.Title, req.SourceBranch, req.TargetBranch, req.ProjectID, req.Author, commits)
	if err != nil {
		response.Error(c, "创建 Pull Request 失败："+err.Error())
		return
	}
	response.Created(c, toPullRequestResponse(pr))
}

func (h *GitHandler) GetPullRequest(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	pr, err := h.gitService.GetPullRequest(id, projectName)
	if err != nil {
		response.Error(c, "获取 Pull Request 失败："+err.Error())
		return
	}
	if pr == nil {
		response.NotFound(c, "Pull Request 不存在")
		return
	}
	response.Success(c, toPullRequestResponse(pr))
}

func (h *GitHandler) MergePullRequest(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	if err := h.gitService.MergePullRequest(id, projectName); err != nil {
		response.Error(c, "合并 Pull Request 失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *GitHandler) ClosePullRequest(c *gin.Context) {
	id := c.Param("id")
	projectName := c.Query("projectName")
	if err := h.gitService.ClosePullRequest(id, projectName); err != nil {
		response.Error(c, "关闭 Pull Request 失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *GitHandler) SquashMergePullRequest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CommitMessage string `json:"commitMessage" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	projectName := c.Query("projectName")
	if err := h.gitService.SquashMergePullRequest(id, projectName, req.CommitMessage); err != nil {
		response.Error(c, "Squash 合并 Pull Request 失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *GitHandler) CreateReview(c *gin.Context) {
	prID := c.Param("id")
	var req struct {
		Reviewer string `json:"reviewer" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	projectName := c.Query("projectName")
	review, err := h.reviewService.CreateReview(prID, projectName, req.Reviewer)
	if err != nil {
		response.Error(c, "创建 Review 失败："+err.Error())
		return
	}
	response.Created(c, toReviewResponse(review))
}

func (h *GitHandler) ApproveReview(c *gin.Context) {
	prID := c.Param("id")
	reviewID := c.Param("reviewId")
	var req struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	review, err := h.reviewService.ApproveReview(prID, reviewID, req.Comment)
	if err != nil {
		response.Error(c, "批准 Review 失败："+err.Error())
		return
	}
	response.Success(c, toReviewResponse(review))
}

func (h *GitHandler) RejectReview(c *gin.Context) {
	prID := c.Param("id")
	reviewID := c.Param("reviewId")
	var req struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	review, err := h.reviewService.RejectReview(prID, reviewID, req.Comment)
	if err != nil {
		response.Error(c, "拒绝 Review 失败："+err.Error())
		return
	}
	response.Success(c, toReviewResponse(review))
}

func (h *GitHandler) GetReviews(c *gin.Context) {
	prID := c.Param("id")
	reviews, err := h.reviewService.GetReviews(prID)
	if err != nil {
		response.Error(c, "获取 Review 列表失败："+err.Error())
		return
	}
	resp := make([]gin.H, len(reviews))
	for i, r := range reviews {
		resp[i] = toReviewResponse(r)
	}
	response.Success(c, gin.H{"reviews": resp})
}

func (h *GitHandler) ValidateBranch(c *gin.Context) {
	var req struct {
		Pattern string `json:"pattern"`
		Branch  string `json:"branch"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	policy := &git.BranchPolicy{AllowedPattern: req.Pattern}
	valid := h.gitService.ValidateBranch(policy, req.Branch)
	response.Success(c, gin.H{"valid": valid})
}

func toPullRequestResponse(pr *git.PullRequest) gin.H {
	resp := gin.H{
		"id":             pr.ID,
		"title":          pr.Title,
		"description":    pr.Description,
		"sourceBranch":   pr.SourceBranch,
		"targetBranch":   pr.TargetBranch,
		"status":         string(pr.Status),
		"projectName":    pr.ProjectName,
		"author":         pr.Author,
		"createdAt":      pr.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt":      pr.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if pr.GeneratedDesc != "" {
		resp["generatedDesc"] = pr.GeneratedDesc
	}
	if pr.SquashCommitMsg != "" {
		resp["squashCommitMsg"] = pr.SquashCommitMsg
	}
	return resp
}

func toReviewResponse(r *git.Review) gin.H {
	resp := gin.H{
		"id":            r.ID,
		"pullRequestId": r.PullRequestID,
		"reviewer":      r.Reviewer,
		"status":        string(r.Status),
		"comment":       r.Comment,
		"createdAt":     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if r.ReviewedAt != nil {
		resp["reviewedAt"] = r.ReviewedAt.Format("2006-01-02T15:04:05Z")
	}
	return resp
}
