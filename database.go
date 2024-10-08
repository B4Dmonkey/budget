package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var conn *sql.DB

func CreateDatabase() error {
	log.Println("Connecting to database...")
	var err error
	DATABASE_LOC := os.Getenv("DATABASE_LOC")

	if DATABASE_LOC == "" {
		return fmt.Errorf("DATABASE_LOC is not set")
	}

	conn, err = sql.Open("sqlite3_extended", DATABASE_LOC)
	if conn == nil || err != nil {
		return fmt.Errorf("Failed to connect the database: %s", DATABASE_LOC)
	}

	if conn.Ping() != nil {
		return fmt.Errorf("Failed to ping database: %s", DATABASE_LOC)
	}

	log.Println("Database connected successfully")
	return nil
}
