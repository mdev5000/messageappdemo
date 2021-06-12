package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB = sqlx.DB

type Message struct {
	Id      int    `db:"id"`
	Version int    `db:"version"`
	Message string `db:"message"`
}

func Open(dbname, user, password string) (*DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbname, password))
	return db, err
}

// OpenTest is a convenience function for setting up the database for testing.
func OpenTest(dbname, user, password, port string) (*DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable", user, dbname, password, port))
	return db, err
}
