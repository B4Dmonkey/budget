package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	sqlite "github.com/mattn/go-sqlite3"
)

func currentTimestamp() string { return time.Now().Format("2006-01-02 15:04:05") }

func newUUID() string { return uuid.New().String() }

// init registers the sqlite3 driver with the extended functions
// This MUST be run first before any other database operations
func init() {
	log.Println("Extending sqlite3 driver...")

	var CONNECT_ONCE sync.Once

	CONNECT_ONCE.Do(func() {
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
	log.Println("Extended sqlite3 driver successfully")
}
func ConnectToDatabase(ctx context.Context) (*sql.DB, error) {

	log.Println("Connecting to database...")

	var connection *sql.DB

	var ok bool

	var err error

	var DATABASE_LOC string

	if DATABASE_LOC, ok = os.LookupEnv("DATABASE_LOC"); !ok || DATABASE_LOC == "" {
		return nil, errors.New("DATABASE_LOC is not set")
	}

	connection, err = sql.Open("sqlite3_extended", DATABASE_LOC)
	if connection == nil || err != nil {
		return nil, fmt.Errorf("Failed to connect the database: %s", DATABASE_LOC)
	}

	if connection.Ping() != nil {
		return nil, fmt.Errorf("Failed to ping database: %s", DATABASE_LOC)
	}

	log.Println("Database connected successfully")

	return connection, nil
}
