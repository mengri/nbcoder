package notify

import (
	"fmt"
	"strings"
	"time"
)

type NotificationTemplate struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	EventType       string    `json:"event_type"`
	SubjectTemplate string    `json:"subject_template"`
	BodyTemplate    string    `json:"body_template"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewNotificationTemplate(id, name, eventType, subjectTemplate, bodyTemplate string) *NotificationTemplate {
	now := time.Now().UTC()
	return &NotificationTemplate{
		ID:              id,
		Name:            name,
		EventType:       eventType,
		SubjectTemplate: subjectTemplate,
		BodyTemplate:    bodyTemplate,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (t *NotificationTemplate) Render(data map[string]string) (subject string, body string, err error) {
	subject, err = renderTemplate(t.SubjectTemplate, data)
	if err != nil {
		return "", "", fmt.Errorf("render subject: %w", err)
	}
	body, err = renderTemplate(t.BodyTemplate, data)
	if err != nil {
		return "", "", fmt.Errorf("render body: %w", err)
	}
	return subject, body, nil
}

func renderTemplate(tpl string, data map[string]string) (string, error) {
	var missing []string
	result := tpl
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	for _, open := range findAllBetween(result, "{{", "}}") {
		missing = append(missing, open)
	}
	if len(missing) > 0 {
		return "", fmt.Errorf("unresolved placeholders: %v", missing)
	}
	return result, nil
}

func findAllBetween(s, left, right string) []string {
	var results []string
	for {
		start := strings.Index(s, left)
		if start == -1 {
			break
		}
		end := strings.Index(s, right)
		if end == -1 || end <= start {
			break
		}
		results = append(results, s[start+len(left):end])
		s = s[end+len(right):]
	}
	return results
}
