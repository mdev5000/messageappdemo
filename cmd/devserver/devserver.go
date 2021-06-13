package main

import (
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/mdev5000/qlik_message/server"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	db, err := postgres.Open("postgres", "postgres", "postgres")
	if err != nil {
		return err
	}

	// Setup the database schema.
	if _, err := db.Exec(data.Schema); err != nil {
		return err
	}

	log := logging.New()
	messagesRepo := data.NewMessageRepository(db)

	services := server.Services{
		Log:             log,
		MessagesService: messages.NewService(log, messagesRepo),
	}

	handler, err := server.Handler(services, server.Config{
		LogRequest: true,
	})
	if err != nil {
		return err
	}

	return http.ListenAndServe("localhost:8000", handler)
}
