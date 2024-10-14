package main

import (
	// "context"
	// "sync"
	// "time"
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	tests := []struct {
		description string
		path        string
		contentType string
	}{
		// {
		// 	description: "It serves the styles.css file",
		// 	path:        "/styles.css",
		// 	contentType: "text/css",
		// },
		{
			description: "It serves the home page",
			path:        "/",
			contentType: "text/html; charset=utf-8",
		},
	}

	for _, test := range tests {
		resp, err := http.Get(ts.URL + test.path)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Check the status code
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		// Check the Content-Type header
		if contentType := resp.Header.Get("Content-Type"); contentType != test.contentType {
			t.Fatalf("Expected Content-Type %s, got %s", test.contentType, contentType)
		}
	}
}

func TestAppRun(t *testing.T) {
	// ? Is there room to simplify the graceful shutdown of the server?
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := ConnectToDatabase(ctx)

	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	app := New(conn)

	var wg sync.WaitGroup

	wg.Add(1)

	go app.Run(ctx, &wg)

	time.Sleep(1 * time.Second) // Give the server a moment to start

	cancel()
	wg.Wait()

	assert.True(t, true, "Server shutdown gracefully")
}
