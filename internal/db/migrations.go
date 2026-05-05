package db

import (
	"database/sql"
	"fmt"
	"os"
)

func Migrate(db *sql.DB, schemaPath string) error {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("execute schema: %w", err)
	}

	db.Exec("ALTER TABLE posts ADD COLUMN league TEXT NOT NULL DEFAULT ''")
	db.Exec("ALTER TABLE comments ADD COLUMN parent_id INTEGER NOT NULL DEFAULT 0")
	db.Exec("ALTER TABLE users ADD COLUMN favorite_team TEXT NOT NULL DEFAULT ''")

	return nil
}

func Seed(db *sql.DB, seedPath string) error {
	seed, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}

	if _, err := db.Exec(string(seed)); err != nil {
		return fmt.Errorf("execute seed: %w", err)
	}

	return nil
}
