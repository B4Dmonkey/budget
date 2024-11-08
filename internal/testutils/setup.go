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

	"github.com/stretchr/testify/assert"
)

type Domain struct {
	t      *testing.T
	conn   *sql.DB
	ctx    context.Context
	tx     *sql.Tx
	q      *orm.Queries
	txq    *orm.Queries
	Assert *assert.Assertions
}

func NewDomainTest(t *testing.T) *Domain {
	ctx := context.Background()
	conn := setupTestDbConnection(t, ctx)
	dt := Domain{t: t, conn: conn, ctx: ctx, Assert: assert.New(t)}

	return &dt
}

func (dt *Domain) Conn() *sql.DB {
	return dt.conn
}

func (dt *Domain) Ctx() context.Context {
	return dt.ctx
}

func (dt *Domain) ResetTestState() {
	dt.Teardown()
	dt.ctx = context.Background()
	dt.conn = setupTestDbConnection(dt.t, dt.ctx)
}

func (dt *Domain) setupQueries() *orm.Queries {
	return orm.New(dt.conn)
}

func (dt *Domain) Teardown() { teardownTestDbConnection(dt.t, dt.conn) }

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
