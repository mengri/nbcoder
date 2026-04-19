package notify

type NotificationRepo interface {
	Save(notification *Notification) error
	FindByID(id string) (*Notification, error)
	FindByRecipient(recipient string) ([]*Notification, error)
	FindByEventType(eventType string) ([]*Notification, error)
	MarkRead(id string) error
}

type SubscriptionRepo interface {
	Save(sub *Subscription) error
	FindByRecipient(recipient string) ([]*Subscription, error)
	FindByEventType(eventType string) ([]*Subscription, error)
	Update(sub *Subscription) error
	Delete(id string) error
}
