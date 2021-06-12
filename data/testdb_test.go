package data

import (
	"fmt"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

const testDbUser = "pguser"
const testDbPassword = "secret"
const testDbName = "messagesdb"

// Do not reference this directly! Use acquireDb() instead.
var sharedDbInstance *postgres.DB

// The PostgreSQL database is loaded via Docker. This function starts up a container instance of the database and
// ensures the database can be reached.
func setupDbPool() (*dockertest.Pool, *dockertest.Resource, error) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket).
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: \n%w", err)
	}

	// Pulls an image, creates a container based on it and runs it.
	resource, err := pool.Run("postgres", "13.3", []string{
		"POSTGRES_USER=" + testDbUser,
		"POSTGRES_PASSWORD=" + testDbPassword,
		"POSTGRES_DB=" + testDbName,
	})
	if err != nil {
		return pool, nil, fmt.Errorf("could not start resource: \n%w", err)
	}

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet.
	if err := pool.Retry(func() error {
		var err error
		sharedDbInstance, err = postgres.OpenTest(testDbName, testDbUser, testDbPassword, resource.GetPort("5432/tcp"))
		if err != nil {
			return err
		}
		return sharedDbInstance.Ping()
	}); err != nil {
		return pool, resource, fmt.Errorf("could not connect to sharedDbInstance via docker: 'n%w", err)
	}

	return pool, resource, nil
}

func closeDbPool(pool *dockertest.Pool, resource *dockertest.Resource, errIsFatal bool) {
	if err := pool.Purge(resource); err != nil {
		if errIsFatal {
			log.Fatalf("Could not purge resource: \n%s", err)
		} else {
			log.Printf("Could not purge resource: \n%s", err)
		}
	}
}

func setupDbSchema() error {
	_, err := sharedDbInstance.Exec(Schema)
	return err
}

func TestMain(m *testing.M) {
	// Setup a running PostgreSQL instance.
	pool, resource, err := setupDbPool()
	if err != nil {
		if pool != nil {
			closeDbPool(pool, resource, false)
		}
		log.Fatalf("Failed to start pool:\n%s", err)
	}

	// Apply the required database schema.
	if err := setupDbSchema(); err != nil {
		log.Fatalf("Failed to create schema: \n%s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	closeDbPool(pool, resource, true)

	os.Exit(code)
}

// Acquire a database instance. You must call close you are finished with the database.  This functions currently
// does 1 thing, but can potentially do 2 at some point.
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
func acquireDb() (*postgres.DB, func()) {
	return sharedDbInstance, func() {
		purgeDatabase(sharedDbInstance)
	}
}

// Removes all data from the database.
func purgeDatabase(db *postgres.DB) {
	db.MustExec("delete from messages")
}
