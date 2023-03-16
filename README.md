# litemigrate

A simple SQLite migration library for Go.

## Installation

```bash
go get github.com/joeychilson/litemigrate
```

## Example

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joeychilson/litemigrate"
)

func main() {
	ctx := context.Background()

	// Create the migrations slice.
	migrations := litemigrate.Migrations{
		{
			Version:     1,
			Description: "create users table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS users (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL
					);
				`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec("DROP TABLE IF EXISTS users;")
				return err
			},
		},
		{
			Version:     2,
			Description: "add email column to users table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec("ALTER TABLE users ADD COLUMN email TEXT;")
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec("ALTER TABLE users DROP COLUMN email;")
				return err
			},
		},
	}

	// Create a new database instance.
	db, err := litemigrate.New("test.db", &migrations)
	if err != nil {
		log.Fatalf("failed to create database instance: %v", err)
	}

	// Migrate up to the latest version.
	err = db.MigrateUp(ctx)
	if err != nil {
		log.Fatalf("failed to migrate up: %v", err)
	}

	// Migrate down to the previous version.
	err = db.MigrateDown(ctx, 1)
	if err != nil {
		log.Fatalf("failed to migrate down: %v", err)
	}

	// Get the current version of the database.
	version, err := db.CurrentVersion(ctx)
	if err != nil {
		log.Fatalf("failed to get current version: %v", err)
	}
	fmt.Printf("current database version: %d\n", version)
}
```
