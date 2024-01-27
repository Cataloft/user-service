package storage

import (
	"database/sql"
	"github.com/pressly/goose"
	"log"
)

func UpMigrations(db *sql.DB) {

	migration := goose.

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "internal/storage/migrations"); err != nil {
		log.Fatal(err)
	}
}
