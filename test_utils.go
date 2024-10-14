package main

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	"my-budget/database/orm"
)

type TestQuery struct {
	t    *testing.T
	conn *sql.DB
	ctx  context.Context
	tx   *sql.Tx
	q    *orm.Queries
	txq  *orm.Queries
}

func (tq *TestQuery) setup() *orm.Queries {
	if _, err := tq.conn.ExecContext(tq.ctx, ddl); err != nil {
		tq.t.Fatalf("Failed to create schema: %s", err)
	}
	
	tq.q = orm.New(tq.conn)
	tx, err := tq.conn.Begin()

	if err != nil {
		tq.t.Fatalf("Failed to start transaction: %s", err)
	}

	tq.tx = tx
	tq.txq = tq.q.WithTx(tq.tx)

	return tq.q
}

func (tq *TestQuery) teardown() {
	if err := tq.conn.Close(); err != nil {
		tq.t.Fatalf("Failed to close database connection: %s", err)
	}
}

func setupTestDbConnection(t *testing.T, ctx context.Context) *sql.DB {
	conn, err := ConnectToDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	return conn
}

func teardownTestDbConnection(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close database connection: %s", err)
	}
}

func setupTestServer(t *testing.T) *httptest.Server {
	app := New(setupTestDbConnection(t, context.Background()))
	// assert.Equal(t, "test_address", app.server.Addr)
	ts := httptest.NewServer(app.mux)

	return ts
}
