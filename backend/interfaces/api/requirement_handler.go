package api

import (
	"github.com/gin-gonic/gin"
	reqApp "github.com/mengri/nbcoder/application/requirement"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/domain/requirement"
	"github.com/mengri/nbcoder/pkg/response"
	"github.com/mengri/nbcoder/pkg/uid"
)

type RequirementHandler struct {
	reqService *reqApp.RequirementService
}

func NewRequirementHandler(reqService *reqApp.RequirementService) *RequirementHandler {
	return &RequirementHandler{
		reqService: reqService,
	}
}

func (h *RequirementHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/projects/:id/cards", h.CreateCard)
	router.GET("/projects/:id/cards", h.ListCards)
	router.GET("/projects/:id/cards/:cardId", h.GetCard)
	router.PUT("/projects/:id/cards/:cardId", h.UpdateCard)
	router.DELETE("/projects/:id/cards/:cardId", h.DeleteCard)
	router.POST("/projects/:id/cards/:cardId/confirm", h.ConfirmCard)
	router.POST("/projects/:id/cards/:cardId/start", h.StartCard)
	router.POST("/projects/:id/cards/:cardId/complete", h.CompleteCard)
	router.POST("/projects/:id/cards/:cardId/dependencies", h.AddDependency)
	router.DELETE("/projects/:id/cards/:cardId/dependencies/:depId", h.RemoveDependency)
	router.GET("/projects/:id/cards/:cardId/dependencies", h.GetDependencies)
}

func (h *RequirementHandler) CreateCard(c *gin.Context) {
	projectID := c.Param("id")
	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id := uid.NewID()
	priority := requirement.Priority(req.Priority)
	card, err := h.reqService.CreateCard(id, req.Title, req.Description, req.Original, projectID, priority)
	if err != nil {
		response.Error(c, "创建卡片失败："+err.Error())
		return
	}
	response.Created(c, toCardResponse(card))
}

func (h *RequirementHandler) GetCard(c *gin.Context) {
	cardID := c.Param("cardId")
	card, err := h.reqService.GetCard(cardID)
	if err != nil {
		response.Error(c, "获取卡片失败："+err.Error())
		return
	}
	if card == nil {
		response.NotFound(c, "卡片不存在")
		return
	}
	response.Success(c, toCardResponse(card))
}

func (h *RequirementHandler) ListCards(c *gin.Context) {
	projectID := c.Param("id")
	status := c.Query("status")

	if status != "" {
		cards, err := h.reqService.ListCardsByStatus(requirement.CardStatus(status))
		if err != nil {
			response.Error(c, "获取卡片列表失败："+err.Error())
			return
		}
		response.Success(c, toCardResponseList(cards))
		return
	}

	cards, err := h.reqService.ListCards(projectID)
	if err != nil {
		response.Error(c, "获取卡片列表失败："+err.Error())
		return
	}
	response.Success(c, toCardResponseList(cards))
}

func (h *RequirementHandler) UpdateCard(c *gin.Context) {
	cardID := c.Param("cardId")
	var req dto.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var priority requirement.Priority
	if req.Priority != nil {
		priority = requirement.Priority(*req.Priority)
	}
	card, err := h.reqService.UpdateCard(cardID, req.Title, req.Description, priority, req.StructuredOutput, req.PipelineID)
	if err != nil {
		response.Error(c, "更新卡片失败："+err.Error())
		return
	}
	response.Success(c, toCardResponse(card))
}

func (h *RequirementHandler) DeleteCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if err := h.reqService.DeleteCard(cardID); err != nil {
		response.Error(c, "删除卡片失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) ConfirmCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if err := h.reqService.ConfirmCard(cardID); err != nil {
		response.Error(c, "确认卡片失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) StartCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if err := h.reqService.StartCard(cardID); err != nil {
		response.Error(c, "启动卡片失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) CompleteCard(c *gin.Context) {
	cardID := c.Param("cardId")
	if err := h.reqService.CompleteCard(cardID); err != nil {
		response.Error(c, "完成卡片失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) AddDependency(c *gin.Context) {
	cardID := c.Param("cardId")
	var req dto.AddDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id := uid.NewID()
	dep, err := h.reqService.AddDependency(id, cardID, req.DependsOnID, requirement.DependencyType(req.Type))
	if err != nil {
		response.Error(c, "添加依赖失败："+err.Error())
		return
	}
	response.Created(c, gin.H{
		"id":            dep.ID,
		"cardId":       dep.CardID,
		"dependsOnId":  dep.DependsOnID,
		"type":          string(dep.Type),
	})
}

func (h *RequirementHandler) RemoveDependency(c *gin.Context) {
	depID := c.Param("depId")
	if err := h.reqService.RemoveDependency(depID); err != nil {
		response.Error(c, "删除依赖失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RequirementHandler) GetDependencies(c *gin.Context) {
	cardID := c.Param("cardId")
	deps, err := h.reqService.GetDependencies(cardID)
	if err != nil {
		response.Error(c, "获取依赖列表失败："+err.Error())
		return
	}
	result := make([]gin.H, 0, len(deps))
	for _, dep := range deps {
		result = append(result, gin.H{
			"id":           dep.ID,
			"cardId":      dep.CardID,
			"dependsOnId": dep.DependsOnID,
			"type":        string(dep.Type),
		})
	}
	response.Success(c, result)
}

func toCardResponse(card *requirement.Card) dto.CardResponse {
	return dto.CardResponse{
		ID:               card.ID,
		Title:            card.Title,
		Description:      card.Description,
		Original:         card.Original,
		Status:           string(card.Status),
		Priority:         string(card.Priority),
		StructuredOutput: card.StructuredOutput,
		PipelineID:       card.PipelineID,
		ProjectName:      card.ProjectName,
		SupersededBy:     card.SupersededBy,
		CreatedAt:        card.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        card.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toCardResponseList(cards []*requirement.Card) []dto.CardResponse {
	result := make([]dto.CardResponse, 0, len(cards))
	for _, card := range cards {
		result = append(result, toCardResponse(card))
	}
	return result
}
