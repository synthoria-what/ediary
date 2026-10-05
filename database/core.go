package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal("sql open: %w", err)
		return nil, fmt.Errorf("sql open: %w", err)
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("close db: %v", err)
			}
		}()
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}
