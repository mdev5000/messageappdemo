package uris

import (
	"fmt"
	"github.com/mdev5000/qlik_message/messages"
)

func ReadMessage(messageId messages.MessageId) string {
	return fmt.Sprintf("/messages/%d", messageId)
}
