package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mdev5000/messageappdemo/approot"
	"github.com/mdev5000/messageappdemo/data"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/mdev5000/messageappdemo/postgres"
	"github.com/mdev5000/messageappdemo/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var migrate bool
	var tls bool
	flag.BoolVar(&migrate, "migrate", false, "When set, migrations will be run prior to starting the application.")
	flag.BoolVar(&tls, "tls", false, "When set, if the CERT and KEY environment variables are empty server will panic. This ensure the serve cannot be run in non-TLS mode.")

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
		fmt.Println("  CERT            	  TLS certificate file to use.")
		fmt.Println("  KEY            	  TLS key file to use.")
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

	cert := os.Getenv("CERT")
	if cert == "" && tls {
		return fmt.Errorf("CERT environment variable cannot be empty")
	}

	key := os.Getenv("KEY")
	if key == "" && tls {
		return fmt.Errorf("KEY environment variable cannot be empty")
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
	s := http.Server{
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		Handler:           handler,
		Addr:              addr,
	}
	if cert != "" {
		return s.ListenAndServeTLS(cert, key)
	} else {
		return s.ListenAndServe()
	}
}
