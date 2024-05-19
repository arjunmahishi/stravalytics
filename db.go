package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const (
	insertTmpl = "INSERT INTO %s (%s) VALUES (%s)"
	selectTmpl = "SELECT %s FROM %s"
)

var (
	//go:embed assets/setup.sql
	dbSetup string
)

func queryRunners(id int, db driver.Conn, queryChan <-chan string, resultChan chan<- bool) {
	for query := range queryChan {
		err := db.Exec(context.Background(), query)
		resultChan <- true
		if err != nil {
			log.Println(err)
		}
		// log.Printf("[%d] Executed query (err: %v)", id, err)
	}
}

func bulkInsertData(db driver.Conn, queries []string) {
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

func newDB() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", *dbHost, *dbPort)},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "stravalytics", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	fmt.Print("starting setup")
	setupQuries := strings.Split(strings.TrimSpace(dbSetup), ";")
	setupQuries = setupQuries[:len(setupQuries)-1]
	for _, query := range setupQuries {
		if err := conn.Exec(ctx, query); err != nil {
			fmt.Print("...ERROR: \n")
			return nil, err
		}
	}

	fmt.Print("...DONE\n")
	return conn, nil
}
