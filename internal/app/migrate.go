package app

import (
	"database/sql"
	"fmt"

	"github.com/keremeti/iq-progers/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func applyMigrations(conf config.Goose) error {
	db, err := sql.Open(conf.Driver, conf.Url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("applyMigrations - SetDialect: %v", err)
	}

	if err := goose.Up(db, conf.MigrationDir); err != nil {
		return fmt.Errorf("applyMigrations - Up: %v", err)
	}
	return nil
}
