package db

import (
	"context"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sauravsinghs/secure-n-banked/util"
)

var (
	testStore Store
	connPool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		log.Fatal("cannot resolve project root", err)
	}

	migrationURL := config.MigrationURL
	if migrationURL == "" {
		migrationURL = "file://db/migration"
	}
	if len(migrationURL) >= 7 && migrationURL[:7] == "file://" {
		migrationPath := filepath.Join(projectRoot, "db", "migration")
		migrationURL = (&url.URL{Scheme: "file", Path: filepath.ToSlash(migrationPath)}).String()
	}

	migration, err := migrate.New(migrationURL, config.DBSource)
	if err != nil {
		log.Fatal("cannot create migration instance", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot run migrations", err)
	}

	connPool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}