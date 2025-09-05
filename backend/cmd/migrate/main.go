package main

import (
	"log"
	"os"

	"spyal/db"
	"spyal/tooling/migration"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	if err := migration.Migrate("./migrations", dsn); err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(db.Config{DSN: dsn})
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	log.Println("DB ready")
}
