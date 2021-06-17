// Package approot acts as the connector between the core domain logic and external dependencies (ex. database).
package approot

import (
	"github.com/mdev5000/messageappdemo/data"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/mdev5000/messageappdemo/messages"
	"github.com/mdev5000/messageappdemo/postgres"
)

type Services struct {
	Log             *logging.Logger
	MessagesService *messages.Service
}

func Setup(db *postgres.DB, log *logging.Logger) *Services {
	messagesRepo := data.NewMessageRepository(db)
	services := Services{
		Log:             log,
		MessagesService: messages.NewService(log, messagesRepo),
	}
	return &services
}
