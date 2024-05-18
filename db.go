package main

import (
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

const (
	insertTmpl = "INSERT INTO %s (%s) VALUES (%s)"
	selectTmpl = "SELECT %s FROM %s"
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

func queryRunners(id int, db *sql.DB, queryChan <-chan string, resultChan chan<- bool) {
	for query := range queryChan {
		_, err := db.Exec(query)
		resultChan <- true
		log.Printf("[%d] Executed query (err: %v)", id, err)
	}
}

func bulkInsertData(db *sql.DB, queries []string) {
	queryChan := make(chan string, len(queries))
	resultChan := make(chan bool, len(queries))

	// start workers
	for i := 0; i < *workers; i++ {
		i := i
		go queryRunners(i, db, queryChan, resultChan)
	}

	for _, query := range queries {
		queryChan <- query
	}

	for i := 0; i < len(queries); i++ {
		<-resultChan
	}

	close(queryChan)
	close(resultChan)
}
