package main

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"
)

func setupTestDb(t *testing.T, ctx context.Context) *sql.DB {
	conn, err := ConnectToDatabase(ctx)

	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	return conn
}

func setupTestServer(t *testing.T) *httptest.Server {
	app := New(setupTestDb(t, context.Background()))
	// assert.Equal(t, "test_address", app.server.Addr)
	ts := httptest.NewServer(app.mux)

	return ts
}
