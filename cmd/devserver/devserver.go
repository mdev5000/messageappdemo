package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/mdev5000/messageappdemo/approot"
	"github.com/mdev5000/messageappdemo/data"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/mdev5000/messageappdemo/messages"
	"github.com/mdev5000/messageappdemo/postgres"
	"github.com/mdev5000/messageappdemo/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var noseed bool
	flag.BoolVar(&noseed, "noseed", false, "When true will not purge and re-seed the database.")
	flag.Parse()

	db, err := postgres.OpenDev("postgres", "postgres", "postgres")
	if err != nil {
		return err
	}

	// Setup the database schema.
	if _, err := db.Exec(data.Schema); err != nil {
		return err
	}

	log := logging.New()
	services := approot.Setup(db, log)

	// Seed the database with dev data.
	if !noseed {
		fmt.Println("Delete existing data...")
		if err := data.PurgeDb(db); err != nil {
			return err
		}
		fmt.Println("Seeding database with dev data...")
		if err := seed(services.MessagesService); err != nil {
			return err
		}
	}

	handler, err := server.Handler(server.Services{
		Log:             services.Log,
		MessagesService: services.MessagesService,
	}, server.Config{
		LogRequest: true,
	})
	if err != nil {
		return err
	}

	fmt.Println("Running at localhost:8000")
	return http.ListenAndServe("localhost:8000", handler)
}

func seed(msgService *messages.Service) error {
	for i := 0; i < 100; i++ {
		if _, err := msgService.Create(messages.ModifyMessage{
			Message: fmt.Sprintf("message %d", i),
		}); err != nil {
			return err
		}
	}
	return nil
}
