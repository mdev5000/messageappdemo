package uris

import (
	"fmt"

	"github.com/mdev5000/messageappdemo/messages"
)

func Message(messageId messages.MessageId) string {
	return fmt.Sprintf("/messages/%d", messageId)
}
