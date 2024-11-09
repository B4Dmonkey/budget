package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"github.com/B4Dmonkey/ExtendSqliteAndConnect/DatabaseConnect"
)


func Connect(ctx context.Context) (*sql.DB, error) {

	log.Println("Connecting to database...")


	var ok bool

	var DATABASE_LOC string

	if DATABASE_LOC, ok = os.LookupEnv("DATABASE_LOC"); !ok || DATABASE_LOC == "" {
		return nil, errors.New("DATABASE_LOC is not set")
	}

	return database_connect.ConnectToExtendedSqliteDatabase(DATABASE_LOC)
}
