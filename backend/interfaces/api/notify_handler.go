package api

import (
	"github.com/gin-gonic/gin"
	notifyApp "github.com/mengri/nbcoder/application/notify"
	"github.com/mengri/nbcoder/domain/notify"
	"github.com/mengri/nbcoder/pkg/response"
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
		EventType string `json:"eventType"`
		Channel   string `json:"channel" binding:"required"`
		Recipient string `json:"recipient" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	n, err := h.notifyService.Send(req.Title, req.Content, req.EventType, notify.ChannelType(req.Channel), req.Recipient)
	if err != nil {
		response.Error(c, "发送通知失败："+err.Error())
		return
	}
	response.Created(c, toNotificationResponse(n))
}

func (h *NotifyHandler) Broadcast(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		EventType string `json:"eventType"`
		Recipient string `json:"recipient" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	notifications, err := h.notifyService.SendToAllChannels(req.Title, req.Content, req.EventType, req.Recipient)
	if err != nil {
		response.Error(c, "广播通知失败："+err.Error())
		return
	}
	result := make([]gin.H, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toNotificationResponse(n))
	}
	response.Created(c, result)
}

func (h *NotifyHandler) ListByRecipient(c *gin.Context) {
	recipient := c.Query("recipient")
	if recipient == "" {
		response.BadRequest(c, "recipient is required")
		return
	}
	notifications, err := h.notifyService.GetByRecipient(recipient)
	if err != nil {
		response.Error(c, "获取通知列表失败："+err.Error())
		return
	}
	result := make([]gin.H, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toNotificationResponse(n))
	}
	response.Success(c, result)
}

func (h *NotifyHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	n, err := h.notifyService.GetByID(id)
	if err != nil {
		response.Error(c, "获取通知失败："+err.Error())
		return
	}
	if n == nil {
		response.NotFound(c, "通知不存在")
		return
	}
	response.Success(c, toNotificationResponse(n))
}

func (h *NotifyHandler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.notifyService.MarkRead(id); err != nil {
		response.Error(c, "标记通知已读失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *NotifyHandler) Subscribe(c *gin.Context) {
	var req struct {
		Recipient string `json:"recipient" binding:"required"`
		EventType string `json:"eventType" binding:"required"`
		Channel   string `json:"channel" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	sub, err := h.notifyService.Subscribe(uid.NewID(), req.Recipient, req.EventType, notify.ChannelType(req.Channel))
	if err != nil {
		response.Error(c, "订阅失败："+err.Error())
		return
	}
	response.Created(c, gin.H{
		"id":         sub.ID,
		"recipient":  sub.Recipient,
		"eventType":  sub.EventType,
		"channel":    string(sub.Channel),
		"muted":      sub.Muted,
	})
}

func (h *NotifyHandler) Unsubscribe(c *gin.Context) {
	id := c.Param("id")
	if err := h.notifyService.Unsubscribe(id); err != nil {
		response.Error(c, "取消订阅失败："+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *NotifyHandler) GetSubscriptions(c *gin.Context) {
	recipient := c.Query("recipient")
	if recipient == "" {
		response.BadRequest(c, "recipient is required")
		return
	}
	subs, err := h.notifyService.GetSubscriptions(recipient)
	if err != nil {
		response.Error(c, "获取订阅列表失败："+err.Error())
		return
	}
	result := make([]gin.H, 0, len(subs))
	for _, s := range subs {
		result = append(result, gin.H{
			"id":         s.ID,
			"recipient":  s.Recipient,
			"eventType":  s.EventType,
			"channel":    string(s.Channel),
			"muted":      s.Muted,
		})
	}
	response.Success(c, result)
}

func (h *NotifyHandler) RegisterChannel(c *gin.Context) {
	var req struct {
		Type   string `json:"type" binding:"required"`
		Config string `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	channel := notify.NewChannel(uid.NewID(), notify.ChannelType(req.Type), req.Config)
	if err := h.notifyService.RegisterChannel(channel); err != nil {
		response.Error(c, "注册通道失败："+err.Error())
		return
	}
	response.Created(c, gin.H{"id": channel.ID, "type": string(channel.Type)})
}

func (h *NotifyHandler) ListChannels(c *gin.Context) {
	channelType := c.Query("type")
	channels, err := h.notifyService.ListChannels(notify.ChannelType(channelType))
	if err != nil {
		response.Error(c, "获取通道列表失败："+err.Error())
		return
	}
	result := make([]gin.H, 0, len(channels))
	for _, ch := range channels {
		result = append(result, gin.H{"id": ch.ID, "type": string(ch.Type), "config": ch.Config})
	}
	response.Success(c, result)
}

func toNotificationResponse(n *notify.Notification) gin.H {
	resp := gin.H{
		"id":        n.ID,
		"title":     n.Title,
		"content":   n.Content,
		"eventType": n.EventType,
		"channel":   string(n.Channel),
		"recipient": n.Recipient,
		"status":    string(n.Status),
		"read":      n.Read,
		"createdAt": n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if n.SentAt != nil {
		resp["sentAt"] = n.SentAt.Format("2006-01-02T15:04:05Z")
	}
	return resp
}
