package migrations

import (
	"database/sql"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

// note: this embeds the migration scripts into the binary
// this benefits us so that we keep production migration process simpler
// instead of bringing goose cli or the entire repo into the VPS
// and doing the migration there, we can just stop the app and pass an argument
// and the app itself does the migration up, since everything needed is in
// the binary of the app.

func Up(databaseUrl string) error {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}
