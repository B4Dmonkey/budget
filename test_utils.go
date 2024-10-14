package main

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	"my-budget/database/orm"
)

type DomainTest struct {
	t    *testing.T
	conn *sql.DB
	ctx  context.Context
	tx   *sql.Tx
	q    *orm.Queries
	txq  *orm.Queries
}

func SetupDomainTest(t *testing.T) *DomainTest {
	ctx := context.Background()
	conn := setupTestDbConnection(t, ctx)
	dt := DomainTest{t: t, conn: conn, ctx: ctx}
	dt.setup()

	return &dt
}

func (dt *DomainTest) setup() *orm.Queries {
	dt.q = orm.New(dt.conn)
	tx, err := dt.conn.Begin()

	if err != nil {
		dt.t.Fatalf("Failed to start transaction: %s", err)
	}

	dt.tx = tx
	dt.txq = dt.q.WithTx(dt.tx)

	return dt.q
}

func (dt *DomainTest) teardown() { teardownTestDbConnection(dt.t, dt.conn) }

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

func setupTestServer(t *testing.T) *httptest.Server {
	app := New(context.Background())
	// assert.Equal(t, "test_address", app.server.Addr)
	ts := httptest.NewServer(app.mux)

	return ts
}

type AppTest struct {
	t          *testing.T
	a          *App
	testServer *httptest.Server
}

func NewAppTest(t *testing.T) *AppTest {
	ctx := context.Background()

	app := New(ctx)

	a := AppTest{t: t, a: app, testServer: httptest.NewServer(app.mux)}
	a.setup()

	return &a
}

func (at *AppTest) setup() {}

func (at *AppTest) teardown() { 
	teardownTestDbConnection(at.t, at.a.db_conn)
	at.testServer.Close()
}
