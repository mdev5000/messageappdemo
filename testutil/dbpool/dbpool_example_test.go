package dbpool_test

import (
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/mdev5000/qlik_message/testutil/dbpool"
	"log"
)

func testWithDb(db *postgres.DB) {
}

func ExampleDbPool() {
	// In your test pre-setup.
	pool := dbpool.NewDbPool()
	pool.SetupSchema = data.SetupSchema
	pool.PurgeDb = data.PurgeDb

	if err := pool.Setup(); err != nil {
		pool.Close(false)
		log.Fatalf("Failed to start pool:\n%s", err)
	}

	// Run you test here and in your test.
	db, closeDb := pool.AcquireDb()
	defer closeDb()
	testWithDb(db)

	// Usually this is used with TestMain, so you can't defer this because os.Exit doesn't care for defer (when using
	// TestMain).
	pool.Close(true)
}
