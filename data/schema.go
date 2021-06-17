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

func SetupSchema(db *postgres.DB) error {
	_, err := db.Exec(Schema)
	return err
}

func PurgeDb(db *postgres.DB) error {
	if _, err := db.Exec("delete from messages"); err != nil {
		return err
	}
	return nil
}
