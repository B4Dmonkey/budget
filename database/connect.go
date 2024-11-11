package database

import (
	"database/sql"
	"log"
	"github.com/B4Dmonkey/ExtendSqliteAndConnect"
)


func Connect(dbLocation string) *sql.DB {
	if connection, err := database_connect.ConnectToExtendedSqliteDatabase(dbLocation); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil
	} else {
		return connection
	}
}
