package server

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

func acquireDb(t *testing.T) (*postgres.DB, func()) {
	if !runDbTests {
		t.SkipNow()
	}
	return pool.AcquireDb()
}

func TestMain(m *testing.M) {
	// Do not run db tests is NODB environment variable set to 1
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
