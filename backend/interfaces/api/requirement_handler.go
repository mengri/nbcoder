package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	reqApp "github.com/mengri/nbcoder/application/requirement"
	"github.com/mengri/nbcoder/application/dto"
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
		cards.POST("/:id/confirm", h.ConfirmCard)
		cards.GET("/:id", h.GetCard)
	}
}

func (h *RequirementHandler) CreateCard(c *gin.Context) {
	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := generateID()
	card, err := h.reqService.CreateCard(id, req.Title, req.Description, req.Original, req.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.CardResponse{
		ID:        card.ID,
		Title:     card.Title,
		Status:    string(card.Status),
		ProjectID: card.ProjectID,
		CreatedAt: card.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: card.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *RequirementHandler) ConfirmCard(c *gin.Context) {
	cardID := c.Param("id")
	if err := h.reqService.ConfirmCard(cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card confirmed"})
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
	c.JSON(http.StatusOK, dto.CardResponse{
		ID:        card.ID,
		Title:     card.Title,
		Status:    string(card.Status),
		ProjectID: card.ProjectID,
		CreatedAt: card.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: card.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
