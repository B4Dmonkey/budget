package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	sqlite "github.com/mattn/go-sqlite3"
)

var conn *sql.DB
var connectOnce sync.Once

func init() {
	connectOnce.Do(func() {
		sql.Register("sqlite3_extended", &sqlite.SQLiteDriver{
			ConnectHook: func(conn *sqlite.SQLiteConn) error {

				if err := conn.RegisterFunc("uuid", newUUID, false); err != nil {
					return err
				}

				if err := conn.RegisterFunc("current_timestamp", currentTimestamp, false); err != nil {
					return err
				}

				return nil
			},
		})
	})
}

func newUUID() string {
	return uuid.New().String()
}

func currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

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
