package testutils

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"

	"my-budget/database/orm"
	"my-budget/internal/db"
)

type Domain struct {
	t    *testing.T
	Conn *sql.DB
	Ctx  context.Context
	tx   *sql.Tx
	q    *orm.Queries
	txq  *orm.Queries
}

func NewDomainTest(t *testing.T) *Domain {
	ctx := context.Background()
	conn := setupTestDbConnection(t, ctx)
	dt := Domain{t: t, Conn: conn, Ctx: ctx}

	return &dt
}

func (dt *Domain) setupQueries() *orm.Queries {
	return orm.New(dt.Conn)
}

func (dt *Domain) teardown() { teardownTestDbConnection(dt.t, dt.Conn) }

func setupTestDbConnection(t *testing.T, ctx context.Context) *sql.DB {
	conn, err := db.ConnectToDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Failed to get current working directory: %s", err)
	}

	filePath := filepath.Join(cwd, "../../database/schema.sql")
	schema, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	if _, err := conn.ExecContext(ctx, string(schema)); err != nil {
		t.Fatalf("Failed to create schema: %s", err)
	}

	return conn
}

func teardownTestDbConnection(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close database connection: %s", err)
	}
}
