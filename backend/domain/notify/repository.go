package notify

type NotificationRepo interface {
	Save(notification *Notification) error
	FindByID(id string) (*Notification, error)
	FindByRecipient(recipient string) ([]*Notification, error)
	FindByEventType(eventType string) ([]*Notification, error)
	FindByRecipientAndChannel(recipient string, channel ChannelType) ([]*Notification, error)
	Update(notification *Notification) error
	MarkRead(id string) error
}

type SubscriptionRepo interface {
	Save(sub *Subscription) error
	FindByRecipient(recipient string) ([]*Subscription, error)
	FindByEventType(eventType string) ([]*Subscription, error)
	FindByRecipientAndEventType(recipient, eventType string) ([]*Subscription, error)
	Update(sub *Subscription) error
	Delete(id string) error
}

type ChannelRepo interface {
	Save(channel *Channel) error
	FindByID(id string) (*Channel, error)
	FindByType(channelType ChannelType) ([]*Channel, error)
	Delete(id string) error
}
