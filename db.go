package main

import (
	"database/sql"
	_ "embed"

	_ "github.com/marcboeker/go-duckdb"
)

var (
	//go:embed assets/setup.sql
	dbSetup string
)

func newDB() (*sql.DB, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(dbSetup)
	if err != nil {
		return nil, err
	}

	return db, err
}
