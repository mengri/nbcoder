package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	reqApp "github.com/mengri/nbcoder/application/requirement"
	"github.com/mengri/nbcoder/application/dto"
	"github.com/mengri/nbcoder/domain/requirement"
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
	cards := router.Group("/cards")
	{
		cards.POST("", h.CreateCard)
		cards.GET("", h.ListCards)
		cards.GET("/:id", h.GetCard)
		cards.PUT("/:id", h.UpdateCard)
		cards.DELETE("/:id", h.DeleteCard)
		cards.POST("/:id/confirm", h.ConfirmCard)
		cards.POST("/:id/start", h.StartCard)
		cards.POST("/:id/complete", h.CompleteCard)
		cards.POST("/:id/dependencies", h.AddDependency)
		cards.DELETE("/:id/dependencies/:depId", h.RemoveDependency)
		cards.GET("/:id/dependencies", h.GetDependencies)
	}
}

func (h *RequirementHandler) CreateCard(c *gin.Context) {
	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uid.NewID()
	priority := requirement.Priority(req.Priority)
	card, err := h.reqService.CreateCard(id, req.Title, req.Description, req.Original, req.ProjectID, priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toCardResponse(card))
}

func (h *RequirementHandler) GetCard(c *gin.Context) {
	cardID := c.Param("id")
	card, err := h.reqService.GetCard(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}
	c.JSON(http.StatusOK, toCardResponse(card))
}

func (h *RequirementHandler) ListCards(c *gin.Context) {
	projectID := c.Query("project_id")
	status := c.Query("status")

	if status != "" {
		cards, err := h.reqService.ListCardsByStatus(requirement.CardStatus(status))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, toCardResponseList(cards))
		return
	}

	cards, err := h.reqService.ListCards(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toCardResponseList(cards))
}

func (h *RequirementHandler) UpdateCard(c *gin.Context) {
	cardID := c.Param("id")
	var req dto.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var priority requirement.Priority
	if req.Priority != nil {
		priority = requirement.Priority(*req.Priority)
	}
	card, err := h.reqService.UpdateCard(cardID, req.Title, req.Description, priority, req.StructuredOutput, req.PipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toCardResponse(card))
}

func (h *RequirementHandler) DeleteCard(c *gin.Context) {
	cardID := c.Param("id")
	if err := h.reqService.DeleteCard(cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card deleted"})
}

func (h *RequirementHandler) ConfirmCard(c *gin.Context) {
	cardID := c.Param("id")
	if err := h.reqService.ConfirmCard(cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card confirmed"})
}

func (h *RequirementHandler) StartCard(c *gin.Context) {
	cardID := c.Param("id")
	if err := h.reqService.StartCard(cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card started"})
}

func (h *RequirementHandler) CompleteCard(c *gin.Context) {
	cardID := c.Param("id")
	if err := h.reqService.CompleteCard(cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card completed"})
}

func (h *RequirementHandler) AddDependency(c *gin.Context) {
	cardID := c.Param("id")
	var req dto.AddDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uid.NewID()
	dep, err := h.reqService.AddDependency(id, cardID, req.DependsOnID, requirement.DependencyType(req.Type))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":            dep.ID,
		"card_id":       dep.CardID,
		"depends_on_id": dep.DependsOnID,
		"type":          string(dep.Type),
	})
}

func (h *RequirementHandler) RemoveDependency(c *gin.Context) {
	depID := c.Param("depId")
	if err := h.reqService.RemoveDependency(depID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "dependency removed"})
}

func (h *RequirementHandler) GetDependencies(c *gin.Context) {
	cardID := c.Param("id")
	deps, err := h.reqService.GetDependencies(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]gin.H, 0, len(deps))
	for _, dep := range deps {
		result = append(result, gin.H{
			"id":            dep.ID,
			"card_id":       dep.CardID,
			"depends_on_id": dep.DependsOnID,
			"type":          string(dep.Type),
		})
	}
	c.JSON(http.StatusOK, result)
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
		ProjectID:        card.ProjectID,
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
