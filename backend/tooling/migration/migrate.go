package migration

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // pgsql driver
	_ "github.com/golang-migrate/migrate/v4/source/file" // file driver
)

func Migrate(migrationsPath, databaseURL string) error {
	m, err := migrate.New(
		"file://" + migrationsPath,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("database is already up-to-date")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("migrations applied successfully")
	return nil
}

func Rollback(migrationsPath, databaseURL string, steps int) error {
	m, err := migrate.New(
		"file://" + migrationsPath,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Printf("rolled back last %d migration(s) successfully\n", steps)
	return nil
}
