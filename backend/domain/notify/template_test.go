package notify

import (
	"testing"
	"time"
)

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestNotificationTemplate_Render_Success(t *testing.T) {
	tpl := NewNotificationTemplate(
		"tpl-1",
		"Task Completed",
		"TASK_COMPLETED",
		"Task {{task_name}} completed",
		"Hello {{user}}, task {{task_name}} has been completed successfully.",
	)
	subject, body, err := tpl.Render(map[string]string{
		"task_name": "build",
		"user":      "alice",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject != "Task build completed" {
		t.Errorf("subject = %q, want %q", subject, "Task build completed")
	}
	if body != "Hello alice, task build has been completed successfully." {
		t.Errorf("body = %q, want %q", body, "Hello alice, task build has been completed successfully.")
	}
}

func TestNotificationTemplate_Render_MissingPlaceholder(t *testing.T) {
	tpl := NewNotificationTemplate(
		"tpl-2",
		"Build Failed",
		"BUILD_FAILED",
		"Build {{project}} failed",
		"Build {{project}} failed with error: {{error}}",
	)
	_, _, err := tpl.Render(map[string]string{
		"project": "nbcoder",
	})
	if err == nil {
		t.Fatal("expected error for missing placeholder, got nil")
	}
}

func TestNotificationTemplate_Render_NoPlaceholders(t *testing.T) {
	tpl := NewNotificationTemplate(
		"tpl-3",
		"Static Notice",
		"STATIC",
		"System Notice",
		"The system will undergo maintenance.",
	)
	subject, body, err := tpl.Render(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject != "System Notice" {
		t.Errorf("subject = %q, want %q", subject, "System Notice")
	}
	if body != "The system will undergo maintenance." {
		t.Errorf("body = %q, want %q", body, "The system will undergo maintenance.")
	}
}

func TestNotificationTemplate_Render_DuplicatePlaceholder(t *testing.T) {
	tpl := NewNotificationTemplate(
		"tpl-4",
		"Deploy",
		"DEPLOY",
		"{{env}} deploy",
		"Deploying to {{env}}. Good luck {{env}}!",
	)
	subject, body, err := tpl.Render(map[string]string{
		"env": "production",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject != "production deploy" {
		t.Errorf("subject = %q, want %q", subject, "production deploy")
	}
	if body != "Deploying to production. Good luck production!" {
		t.Errorf("body = %q, want %q", body, "Deploying to production. Good luck production!")
	}
}

func TestNotificationTemplate_Render_ExtraDataIgnored(t *testing.T) {
	tpl := NewNotificationTemplate(
		"tpl-5",
		"Simple",
		"SIMPLE",
		"Hello {{name}}",
		"Welcome, {{name}}",
	)
	subject, body, err := tpl.Render(map[string]string{
		"name":    "bob",
		"ignored": "value",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject != "Hello bob" {
		t.Errorf("subject = %q, want %q", subject, "Hello bob")
	}
	if body != "Welcome, bob" {
		t.Errorf("body = %q, want %q", body, "Welcome, bob")
	}
}

func TestNewNotificationHistory(t *testing.T) {
	now := parseTime("2025-01-01T00:00:00Z")
	h := NewNotificationHistory(
		"h-1",
		"notif-1",
		"user-1",
		ChannelEmail,
		HistoryStatusSent,
		now,
		"",
	)
	if h.ID != "h-1" {
		t.Errorf("ID = %q, want %q", h.ID, "h-1")
	}
	if h.NotificationID != "notif-1" {
		t.Errorf("NotificationID = %q, want %q", h.NotificationID, "notif-1")
	}
	if h.Recipient != "user-1" {
		t.Errorf("Recipient = %q, want %q", h.Recipient, "user-1")
	}
	if h.Channel != ChannelEmail {
		t.Errorf("Channel = %q, want %q", h.Channel, ChannelEmail)
	}
	if h.Status != HistoryStatusSent {
		t.Errorf("Status = %q, want %q", h.Status, HistoryStatusSent)
	}
	if h.Error != "" {
		t.Errorf("Error = %q, want empty", h.Error)
	}
}

func TestNewNotificationHistory_WithErr(t *testing.T) {
	now := parseTime("2025-01-01T00:00:00Z")
	h := NewNotificationHistory(
		"h-2",
		"notif-2",
		"user-2",
		ChannelWebSocket,
		HistoryStatusFailed,
		now,
		"connection refused",
	)
	if h.Status != HistoryStatusFailed {
		t.Errorf("Status = %q, want %q", h.Status, HistoryStatusFailed)
	}
	if h.Error != "connection refused" {
		t.Errorf("Error = %q, want %q", h.Error, "connection refused")
	}
}
