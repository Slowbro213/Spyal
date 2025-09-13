package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
	"spyal/db"
)

//go:embed seeds/*.sql
var seedsFS embed.FS

const seedsTable = `CREATE TABLE IF NOT EXISTS _seeds_applied (
	name TEXT PRIMARY KEY,
	applied_at TIMESTAMP NOT NULL DEFAULT NOW()
)`

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	database, err := db.Connect(db.Config{DSN: dsn})
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := seed(database.DB); err != nil {
		log.Fatal(err)
	}
	log.Println("Seeding complete")
}

func seed(db *sqlx.DB) error {
	if _,err := db.Exec(seedsTable); err != nil {
		return err
	}

	// already applied
	applied := make(map[string]bool)
	var names []string
	if err := db.Select(&names, `SELECT name FROM _seeds_applied`); err == nil {
		for _, n := range names {
			applied[n] = true
		}
	}

	// list embedded seed files
	entries, err := seedsFS.ReadDir("seeds")
	if err != nil {
		return err
	}
	var toRun []fs.DirEntry
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".sql" {
			toRun = append(toRun, e)
		}
	}
	sort.Slice(toRun, func(i, j int) bool {
		return toRun[i].Name() < toRun[j].Name()
	})

	for _, e := range toRun {
		name := e.Name()
		if applied[name] {
			continue
		}
		log.Printf("applying seed %s", name)

		sqlBytes, err := seedsFS.ReadFile(filepath.Join("seeds", name))
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT INTO _seeds_applied (name) VALUES ($1)`, name); err != nil {
			return err
		}
	}
	return nil
}
