package messages

import (
	"time"
)

const NoOp = 0

type Field = string

const (
	FieldId        = "id"
	FieldVersion   = "version"
	FieldMessage   = "message"
	FieldCreatedAt = "createdAt"
	FieldUpdatedAt = "updatedAt"
)

type CreateMessage struct {
	Message string `db:"message"`

	// CreatedAt is only used for creation of a message and will be ignored for update operations.
	CreatedAt time.Time `db:"created_at"`
}

type Repository interface {
	Create(cm CreateMessage) (MessageId, error)
	DeleteById(id MessageId) error
	GetAllQuery(query MessageQuery, messages *[]*Message) error
	GetById(id MessageId, m *Message) error
	UpdateById(id MessageId, m ModifyMessage) (MessageVersion, error)
}

type ModifyMessage struct {
	Message string
}

type MessageId = int64
type MessageVersion = int

type Message struct {
	Id        MessageId      `db:"id"`
	Version   MessageVersion `db:"version"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	Message   string         `db:"message"`
}

type MessageQuery struct {
	Fields map[Field]struct{}
	Limit  uint64
	Offset uint64
}

func IsPalindrome(msg *Message) bool {
	return isPalindrome(msg.Message)
}
