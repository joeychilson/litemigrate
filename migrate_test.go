package litemigrate_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/joeychilson/litemigrate"
)

const testDBPath = ":memory:"

func TestNew(t *testing.T) {
	migrations := &litemigrate.Migrations{}
	db, err := litemigrate.New(testDBPath, migrations)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if db == nil {
		t.Error("expected non-nil database, got nil")
	}

	err = db.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMigrateUp(t *testing.T) {
	migrations := &litemigrate.Migrations{
		{
			Version:     1,
			Description: "Create test table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY);`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec(`DROP TABLE test;`)
				return err
			},
		},
	}

	db, err := litemigrate.New(testDBPath, migrations)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer db.Close()

	err = db.MigrateUp(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	version, err := db.CurrentVersion(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if version != 1 {
		t.Errorf("expected version 1, got %d", version)
	}
}

func TestMigrateDown(t *testing.T) {
	migrations := &litemigrate.Migrations{
		{
			Version:     1,
			Description: "Create test table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY);`)
				return err
			},
			Down: func(tx *sql.Tx) error {
				_, err := tx.Exec(`DROP TABLE test;`)
				return err
			},
		},
	}

	db, err := litemigrate.New(testDBPath, migrations)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer db.Close()

	err = db.MigrateUp(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = db.MigrateDown(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	version, err := db.CurrentVersion(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if version != 0 {
		t.Errorf("expected version 0, got %d", version)
	}
}

func TestInvalidMigration(t *testing.T) {
	migrations := &litemigrate.Migrations{
		{
			Version:     0,
			Description: "Invalid migration",
			Up:          nil,
			Down:        nil,
		},
	}

	db, err := litemigrate.New(testDBPath, migrations)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer db.Close()

	err = db.MigrateUp(context.Background())
	expectedErr := fmt.Errorf("invalid migration: version and description must be set")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
