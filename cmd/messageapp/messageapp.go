package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mdev5000/qlik_message/approot"
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/mdev5000/qlik_message/server"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "When set, migrations will be run prior to starting the application.")
	flag.Usage = func() {
		fmt.Println("Message App")
		fmt.Println("")
		fmt.Println("  REST API server that manages messages.")
		fmt.Println("")
		fmt.Println("Flags:")
		fmt.Println("")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Environment variables:")
		fmt.Println("")
		fmt.Println("  DATABASE_URL       The url to the database. [required]")
		fmt.Println("  MIGRATE            When set to 1, migrations will be run prior to starting the application.")
		fmt.Println("")
	}
	flag.Parse()

	log := logging.New()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Error("DATABASE_URL was empty.")
		return errors.New("environment variable DATABASE_URL cannot be empty")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Error("PORT was empty.")
		return errors.New("environment variable PORT cannot be empty")
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	migrateEnv := os.Getenv("MIGRATE")
	if migrateEnv == "1" {
		migrate = true
	}

	db, err := postgres.OpenUrl(dbUrl)
	if err != nil {
		return err
	}

	// Setup the database schema.
	if migrate {
		fmt.Println("Running migrations...")
		if _, err := db.Exec(data.Schema); err != nil {
			log.Errorf("Faied to open db connection: %s", err)
			return err
		}
		fmt.Println("Migrations run.")
	}

	services := approot.Setup(db, log)

	handler, err := server.Handler(server.Services{
		Log:             services.Log,
		MessagesService: services.MessagesService,
	}, server.Config{
		LogRequest: true,
	})
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Running at %s\n", addr)
	return http.ListenAndServe(addr, handler)
}
