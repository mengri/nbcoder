package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mengri/nbcoder/domain/notify"
	notifyApp "github.com/mengri/nbcoder/application/notify"
)

type NotifyHandler struct {
	notifyService *notifyApp.NotifyService
}

func NewNotifyHandler(notifyService *notifyApp.NotifyService) *NotifyHandler {
	return &NotifyHandler{
		notifyService: notifyService,
	}
}

func (h *NotifyHandler) RegisterRoutes(router *gin.RouterGroup) {
	notifications := router.Group("/notifications")
	{
		notifications.GET("", h.ListByRecipient)
		notifications.POST("/subscribe", h.Subscribe)
	}
}

func (h *NotifyHandler) ListByRecipient(c *gin.Context) {
	recipient := c.Query("recipient")
	if recipient == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipient is required"})
		return
	}
	notifications, err := h.notifyService.GetByRecipient(recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *NotifyHandler) Subscribe(c *gin.Context) {
	var sub notify.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.notifyService.Subscribe(&sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "subscribed"})
}
