// Package messages contains the core domain logic for Messages.
package messages

import (
	"time"
)

const (
	noOp = 0

	// MaxMessageCharLength is the maximum number of characters a string can contain. Note this is characters not bytes
	// so max bytes would be MaxMessageCharLength * 4 for UTC strings. Also extended grapheme clusters will count as
	// multiple characters (ex. "🤦🏼‍♂️", paste into https://fsymbols.com/emoticons/maker/ to understand).
	// Also see https://hsivonen.fi/string-length/.
	MaxMessageCharLength = 512
)

type Field = string

const (
	FieldId        = "id"
	FieldVersion   = "version"
	FieldMessage   = "message"
	FieldCreatedAt = "createdAt"
	FieldUpdatedAt = "updatedAt"
)

var AllFields = map[string]struct{}{
	FieldId:        {},
	FieldVersion:   {},
	FieldMessage:   {},
	FieldCreatedAt: {},
	FieldUpdatedAt: {},
}

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

// IsPalindrome determines is a Message is a palindrome.
//
// Implementation decisions
//
// Prior to determining if a message is a palindrome it is encoded into NFC, since the service doesn't strictly require
// posted messages to be in NFC and pre-converting the message to NFC could confuse users when their text does not match
// what was originally saved. The conversion is done at the time of palindrome testing. This allows letters with
// combining characters to be treated as a single letter and allows for a more intuitive notion of what is a palindrome
// (ex. éé).
//
// The implementation assumes extended grapheme clusters (ex. "🤦🏼‍♂️") are not palindromes. And more specifically
// whether emojis are palindromes or not left largely undefined.
//
// There is no special handling for hidden characters. This may confuse users, so might be worth adjusting in the
// future.
//
func IsPalindrome(msg *Message) bool {
	return isPalindrome(msg.Message)
}
