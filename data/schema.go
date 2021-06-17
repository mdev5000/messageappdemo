package data

import "github.com/mdev5000/messageappdemo/postgres"

const Schema = `
create table if not exists messages (
	id serial,
	version integer not null,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
    message text not null
);
`

// SetupSchema sets up the current database schema. It is idempotent and is safe to run multiple times.
func SetupSchema(db *postgres.DB) error {
	_, err := db.Exec(Schema)
	return err
}

// PurgeDb deletes all database form the database this should be used only for testing.
func PurgeDb(db *postgres.DB) error {
	if _, err := db.Exec("delete from messages"); err != nil {
		return err
	}
	return nil
}
