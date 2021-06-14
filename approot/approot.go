package approot

import (
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/postgres"
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
