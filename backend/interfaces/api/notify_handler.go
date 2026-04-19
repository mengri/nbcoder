package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	notifyApp "github.com/mengri/nbcoder/application/notify"
	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/pkg/uid"
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
		notifications.POST("", h.Send)
		notifications.POST("/broadcast", h.Broadcast)
		notifications.GET("", h.ListByRecipient)
		notifications.GET("/:id", h.GetByID)
		notifications.PUT("/:id/read", h.MarkRead)
		notifications.POST("/subscribe", h.Subscribe)
		notifications.DELETE("/subscribe/:id", h.Unsubscribe)
		notifications.GET("/subscriptions", h.GetSubscriptions)
		notifications.POST("/channels", h.RegisterChannel)
		notifications.GET("/channels", h.ListChannels)
	}
}

func (h *NotifyHandler) Send(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		EventType string `json:"event_type"`
		Channel   string `json:"channel" binding:"required"`
		Recipient string `json:"recipient" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	n, err := h.notifyService.Send(req.Title, req.Content, req.EventType, notify.ChannelType(req.Channel), req.Recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toNotificationResponse(n))
}

func (h *NotifyHandler) Broadcast(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		EventType string `json:"event_type"`
		Recipient string `json:"recipient" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notifications, err := h.notifyService.SendToAllChannels(req.Title, req.Content, req.EventType, req.Recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]gin.H, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toNotificationResponse(n))
	}
	c.JSON(http.StatusCreated, result)
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
	result := make([]gin.H, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toNotificationResponse(n))
	}
	c.JSON(http.StatusOK, result)
}

func (h *NotifyHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	n, err := h.notifyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if n == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	c.JSON(http.StatusOK, toNotificationResponse(n))
}

func (h *NotifyHandler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.notifyService.MarkRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

func (h *NotifyHandler) Subscribe(c *gin.Context) {
	var req struct {
		Recipient string `json:"recipient" binding:"required"`
		EventType string `json:"event_type" binding:"required"`
		Channel   string `json:"channel" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := h.notifyService.Subscribe(uid.NewID(), req.Recipient, req.EventType, notify.ChannelType(req.Channel))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":         sub.ID,
		"recipient":  sub.Recipient,
		"event_type": sub.EventType,
		"channel":    string(sub.Channel),
		"muted":      sub.Muted,
	})
}

func (h *NotifyHandler) Unsubscribe(c *gin.Context) {
	id := c.Param("id")
	if err := h.notifyService.Unsubscribe(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unsubscribed"})
}

func (h *NotifyHandler) GetSubscriptions(c *gin.Context) {
	recipient := c.Query("recipient")
	if recipient == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipient is required"})
		return
	}
	subs, err := h.notifyService.GetSubscriptions(recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]gin.H, 0, len(subs))
	for _, s := range subs {
		result = append(result, gin.H{
			"id":         s.ID,
			"recipient":  s.Recipient,
			"event_type": s.EventType,
			"channel":    string(s.Channel),
			"muted":      s.Muted,
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *NotifyHandler) RegisterChannel(c *gin.Context) {
	var req struct {
		Type   string `json:"type" binding:"required"`
		Config string `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel := notify.NewChannel(uid.NewID(), notify.ChannelType(req.Type), req.Config)
	if err := h.notifyService.RegisterChannel(channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": channel.ID, "type": string(channel.Type)})
}

func (h *NotifyHandler) ListChannels(c *gin.Context) {
	channelType := c.Query("type")
	channels, err := h.notifyService.ListChannels(notify.ChannelType(channelType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make([]gin.H, 0, len(channels))
	for _, ch := range channels {
		result = append(result, gin.H{"id": ch.ID, "type": string(ch.Type), "config": ch.Config})
	}
	c.JSON(http.StatusOK, result)
}

func toNotificationResponse(n *notify.Notification) gin.H {
	resp := gin.H{
		"id":         n.ID,
		"title":      n.Title,
		"content":    n.Content,
		"event_type": n.EventType,
		"channel":    string(n.Channel),
		"recipient":  n.Recipient,
		"status":     string(n.Status),
		"read":       n.Read,
		"created_at": n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if n.SentAt != nil {
		resp["sent_at"] = n.SentAt.Format("2006-01-02T15:04:05Z")
	}
	return resp
}
