package notify

import (
	"testing"
)

func TestNewNotification(t *testing.T) {
	n := NewNotification("id-1", "Title", "Content", "CardCreated", ChannelSystem, "user-1")
	if n.ID != "id-1" {
		t.Errorf("expected ID id-1, got %s", n.ID)
	}
	if n.Status != NotificationPending {
		t.Errorf("expected status PENDING, got %s", n.Status)
	}
	if n.Read {
		t.Error("expected notification to be unread")
	}
}

func TestNotification_MarkSent(t *testing.T) {
	n := NewNotification("id-1", "Title", "Content", "CardCreated", ChannelSystem, "user-1")
	n.MarkSent()
	if n.Status != NotificationSent {
		t.Errorf("expected status SENT, got %s", n.Status)
	}
	if n.SentAt == nil {
		t.Error("expected SentAt to be set")
	}
}

func TestNotification_MarkFailed(t *testing.T) {
	n := NewNotification("id-1", "Title", "Content", "CardCreated", ChannelSystem, "user-1")
	n.MarkFailed()
	if n.Status != NotificationFailed {
		t.Errorf("expected status FAILED, got %s", n.Status)
	}
}

func TestNotification_MarkDelivered(t *testing.T) {
	n := NewNotification("id-1", "Title", "Content", "CardCreated", ChannelSystem, "user-1")
	n.MarkDelivered()
	if n.Status != NotificationDelivered {
		t.Errorf("expected status DELIVERED, got %s", n.Status)
	}
}

func TestNotification_Validate(t *testing.T) {
	n := NewNotification("id-1", "Title", "Content", "CardCreated", ChannelSystem, "user-1")
	if err := n.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	emptyTitle := NewNotification("id-2", "", "Content", "CardCreated", ChannelSystem, "user-1")
	if err := emptyTitle.Validate(); err == nil {
		t.Error("expected error for empty title")
	}

	emptyRecipient := NewNotification("id-3", "Title", "Content", "CardCreated", ChannelSystem, "")
	if err := emptyRecipient.Validate(); err == nil {
		t.Error("expected error for empty recipient")
	}
}

func TestChannelType_IsValid(t *testing.T) {
	types := []ChannelType{ChannelWebSocket, ChannelEmail, ChannelSystem}
	for _, ct := range types {
		if !ct.IsValid() {
			t.Errorf("expected channel type %s to be valid", ct)
		}
	}
	if ChannelType("INVALID").IsValid() {
		t.Error("expected INVALID to be invalid")
	}
}

func TestChannel_Validate(t *testing.T) {
	ch := NewChannel("ch-1", ChannelSystem, "config")
	if err := ch.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	invalid := NewChannel("ch-2", ChannelType("INVALID"), "")
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for invalid channel type")
	}
}

func TestChannelDispatcher(t *testing.T) {
	d := NewChannelDispatcher()

	if d.HasSender(ChannelSystem) {
		t.Error("expected no sender registered initially")
	}

	sender := &mockSender{supported: []ChannelType{ChannelSystem}}
	d.Register(sender)

	if !d.HasSender(ChannelSystem) {
		t.Error("expected system sender to be registered")
	}

	n := NewNotification("id-1", "Title", "Content", "CardCreated", ChannelSystem, "user-1")
	if err := d.Dispatch(n); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	nEmail := NewNotification("id-2", "Title", "Content", "CardCreated", ChannelEmail, "user-1")
	if err := d.Dispatch(nEmail); err == nil {
		t.Error("expected error for unregistered channel")
	}
}

func TestSubscription(t *testing.T) {
	sub := NewSubscription("sub-1", "user-1", "CardCreated", ChannelSystem)
	if sub.Muted {
		t.Error("expected subscription to be unmuted initially")
	}
	sub.Mute()
	if !sub.Muted {
		t.Error("expected subscription to be muted")
	}
	sub.Unmute()
	if sub.Muted {
		t.Error("expected subscription to be unmuted")
	}
}

type mockSender struct {
	supported []ChannelType
	last      *Notification
}

func (s *mockSender) Send(n *Notification) error {
	s.last = n
	return nil
}

func (s *mockSender) Supports(ct ChannelType) bool {
	for _, st := range s.supported {
		if st == ct {
			return true
		}
	}
	return false
}
