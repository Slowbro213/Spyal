package migration

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // pgsql driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // file driver
)

func Migrate(migrationsPath, databaseURL string) error {
	// Debug prints to ensure we are pointing to the right path and DB
	wd, _ := os.Getwd()
	log.Printf("DEBUG: current working dir: %s\n", wd)

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Printf("WARN: failed to get abs path for %s: %v\n", migrationsPath, err)
		absPath = migrationsPath
	}
	log.Printf("DEBUG: migrationsPath (abs) = %s\n", absPath)
	log.Printf("DEBUG: databaseURL = %s\n", databaseURL)

	// List files in migrations dir (show what migrate will see)
	if files, err := os.ReadDir(absPath); err == nil {
		log.Println("DEBUG: migration files found:")
		for _, f := range files {
			info, _ := f.Info()
			log.Printf("  - %s  (isDir=%v size=%d)\n", f.Name(), f.IsDir(), info.Size())
		}
	} else {
		log.Printf("DEBUG: failed to read migrations dir '%s': %v\n", absPath, err)
	}

	// Provide a fallback: ensure we pass a proper file:// absolute path to migrate
	sourceURL := "file://" + absPath

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations (source=%s): %w", sourceURL, err)
	}

	// Apply all pending migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("database is already up-to-date (no change)")
		} else {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
	} else {
		log.Println("migrations applied successfully (Up)")
	}

	// Print version & dirty flag if available
	if v, dirty, verr := m.Version(); verr == nil {
		log.Printf("DEBUG: migrate version: %d, dirty=%v\n", v, dirty)
	} else {
		// If there's no recorded version it returns ErrNilVersion
		log.Printf("DEBUG: migrate Version() returned error: %v\n", verr)
	}

	// Also log the files in source dir recursively (useful if migrationsPath is actually a file or empty)
	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(absPath, path)
			log.Printf("DEBUG: found source file: %s\n", rel)
		}
		return nil
	})
	if err != nil {
		log.Printf("DEBUG: walk error: %v\n", err)
	}

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
