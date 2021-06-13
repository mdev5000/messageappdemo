package dbpool

import (
	"fmt"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/ory/dockertest/v3"
	"log"
)

const testDbUser = "pguser"
const testDbPassword = "secret"
const testDbName = "messagesdb"

type DbPool struct {
	SetupSchema      func(db *postgres.DB) error
	PurgeDb          func(db *postgres.DB) error
	sharedDbInstance *postgres.DB
	pool             *dockertest.Pool
	resource         *dockertest.Resource
}

func NewDbPool() *DbPool {
	return &DbPool{}
}

// Setup sets up a PostgreSQL database that is loaded via Docker. This function starts up a container instance of the
// database and ensures the database can be reached.
func (d *DbPool) Setup() error {
	var err error

	// Uses a sensible default on windows (tcp/http) and linux/osx (socket).
	d.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not connect to docker: \n%w", err)
	}

	// Pulls an image, creates a container based on it and runs it.
	d.resource, err = d.pool.Run("postgres", "13.3", []string{
		"POSTGRES_USER=" + testDbUser,
		"POSTGRES_PASSWORD=" + testDbPassword,
		"POSTGRES_DB=" + testDbName,
	})
	if err != nil {
		return fmt.Errorf("could not start resource: \n%w", err)
	}

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet.
	if err := d.pool.Retry(func() error {
		var err error
		d.sharedDbInstance, err = postgres.OpenTest(testDbName, testDbUser, testDbPassword, d.resource.GetPort("5432/tcp"))
		if err != nil {
			return err
		}
		return d.sharedDbInstance.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to sharedDbInstance via docker: 'n%w", err)
	}

	if d.SetupSchema != nil {
		if err := d.SetupSchema(d.sharedDbInstance); err != nil {
			return err
		}
	}

	return nil
}

func (d *DbPool) Close(errIsFatal bool) {
	if err := d.pool.Purge(d.resource); err != nil {
		if errIsFatal {
			log.Fatalf("Could not purge resource: \n%s", err)
		} else {
			log.Printf("Could not purge resource: \n%s", err)
		}
	}
}

// AcquireDb acquires a database instance. You must call close you are finished with the database.  This functions
// currently does 1 thing, but can potentially do 2 at some point.
//
// The first is ensure the database is in a clean state prior to running a test. This means existing database is purged
// from the database.
//
// The second is thing is guarding access to database resources. Currently the database runner only has a single
// database instance, since there is limited testing required. However, at some point it may be required to run multiple
// database instances to improve test performance (and run the db tests in parallel). This function would then act as a
// pool manager, serving database instances as required to test functions.
//
// Ex.
// db, closeDb := acquireDb()
// defer closeDb()
// // do db stuff...
//
func (d *DbPool) AcquireDb() (*postgres.DB, func()) {
	if d.sharedDbInstance == nil {
		panic("dbpool has not been setup, did you run Setup()?")
	}
	return d.sharedDbInstance, func() {
		if d.PurgeDb != nil {
			if err := d.PurgeDb(d.sharedDbInstance); err != nil {
				panic(err)
			}
		}
	}
}
