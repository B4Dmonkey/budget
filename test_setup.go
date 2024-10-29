package main

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	"my-budget/database/orm"
)

type AppTest struct {
	t          *testing.T
	a          *App
	testServer *httptest.Server
}

type Domain struct {
	t    *testing.T
	conn *sql.DB
	ctx  context.Context
	tx   *sql.Tx
	q    *orm.Queries
	txq  *orm.Queries
}

func NewDomainTest(t *testing.T) *Domain {
	ctx := context.Background()
	conn := setupTestDbConnection(t, ctx)
	dt := Domain{t: t, conn: conn, ctx: ctx}

	return &dt
}

func (dt *Domain) setupQueries() *orm.Queries {
	return orm.New(dt.conn)
}

func (dt *Domain) teardown() { teardownTestDbConnection(dt.t, dt.conn) }

func NewAppTest(t *testing.T) *AppTest {
	ctx := context.Background()
	conn := setupTestDbConnection(t, ctx)

	app := New(ctx, WithDbConnection(conn), WithOrmQueries(orm.New(conn)))

	a := AppTest{t: t, a: app, testServer: httptest.NewServer(app.mux)}
	a.setup()

	return &a
}

func (at *AppTest) setup() {}

func (at *AppTest) teardown() {
	teardownTestDbConnection(at.t, at.a.db_conn)
	at.testServer.Close()
}

func setupTestDbConnection(t *testing.T, ctx context.Context) *sql.DB {
	conn, err := ConnectToDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	if _, err := conn.ExecContext(ctx, ddl); err != nil {
		t.Fatalf("Failed to create schema: %s", err)
	}

	return conn
}

func teardownTestDbConnection(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close database connection: %s", err)
	}
}
