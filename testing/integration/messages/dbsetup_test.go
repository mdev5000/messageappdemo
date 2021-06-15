package messages

import (
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/mdev5000/qlik_message/testutil/dbpool"
	"log"
	"os"
	"testing"
)

var runDbTests bool
var pool *dbpool.DbPool

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
// db, closeDb := acquireDb(t)
// defer closeDb()
// // do db stuff...
//
func acquireDb(t *testing.T) (*postgres.DB, func()) {
	if !runDbTests {
		t.SkipNow()
	}
	return pool.AcquireDb()
}

func TestMain(m *testing.M) {
	// Do not run any db tests when NODB environment variable is set to 1
	if os.Getenv("NODB") == "1" {
		runDbTests = false
		os.Exit(m.Run())
	}

	runDbTests = true
	pool = dbpool.NewDbPool()
	pool.SetupSchema = data.SetupSchema
	pool.PurgeDb = data.PurgeDb

	if err := pool.Setup(); err != nil {
		pool.Close(false)
		log.Fatalf("Failed to start pool:\n%s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	pool.Close(true)
	os.Exit(code)
}
