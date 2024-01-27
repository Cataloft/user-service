package storage

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
)

func UpMigrations(db *sql.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "internal/storage/migrations"); err != nil {
		log.Fatal(err)
	}
}
