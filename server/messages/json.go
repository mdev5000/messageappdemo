package messages

import (
	"time"

	"github.com/mdev5000/messageappdemo/messages"
)

const (
	fieldIsPalindrome = "isPalindrome"
)

type modifyMessageJSON struct {
	Message string `json:"message"`
}

func (m *modifyMessageJSON) toModifyMessage() messages.ModifyMessage {
	return messages.ModifyMessage{
		Message: m.Message,
	}
}

type MessageListResponseJSON struct {
	Messages []MessageResponseJSON `json:"messages,omitempty"`
}

type MessageResponseJSON struct {
	Id           messages.MessageId      `json:"id,omitempty"`
	Version      messages.MessageVersion `json:"version,omitempty"`
	CreatedAt    *time.Time              `json:"created_at,omitempty"`
	UpdatedAt    *time.Time              `json:"updated_at,omitempty"`
	Message      string                  `json:"message,omitempty"`
	IsPalindrome *bool                   `json:"isPalindrome,omitempty"`
}

type fieldsMap = map[string]struct{}

// Remove out fields that does not exist in the data layer (or they will cause errors).
func filterDynamicFields(fields fieldsMap) fieldsMap {
	out := fieldsMap{}
	for field := range fields {
		if field != fieldIsPalindrome {
			out[field] = struct{}{}
		}
	}
	return out
}

// Similar to messageToJsonValue, but only specified values.
func queryMessageToJsonValue(message *messages.Message, fields fieldsMap) MessageResponseJSON {
	mr := MessageResponseJSON{
		Id:      message.Id,
		Version: message.Version,
		Message: message.Message,
	}
	if hasField(fields, messages.FieldCreatedAt) {
		mr.CreatedAt = &message.CreatedAt
	}
	if hasField(fields, messages.FieldUpdatedAt) {
		mr.UpdatedAt = &message.UpdatedAt
	}
	if hasField(fields, fieldIsPalindrome) {
		isPalindrome := messages.IsPalindrome(message)
		mr.IsPalindrome = &isPalindrome
	}
	return mr
}

func hasField(fields map[string]struct{}, field string) bool {
	_, found := fields[field]
	return found
}

func messageToJsonValue(message *messages.Message) MessageResponseJSON {
	isPalindrome := messages.IsPalindrome(message)
	return MessageResponseJSON{
		Id:           message.Id,
		Version:      message.Version,
		CreatedAt:    &message.CreatedAt,
		UpdatedAt:    &message.UpdatedAt,
		Message:      message.Message,
		IsPalindrome: &isPalindrome,
	}
}
