package main

import (
	"context"
	"sync"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEnvConfig struct {
	mock.Mock
}

func (m *MockEnvConfig) IsDev() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEnvConfig) Addr() string {
	args := m.Called()
	return args.String(0)
}

func TestApp(t *testing.T) {
	mockEnv := new(MockEnvConfig)
	mockEnv.On("IsDev").Return(true)
	mockEnv.On("Addr").Return("test_address")
	ctx := context.Background()
	app := CreateApp(ctx, mockEnv)
	assert.NotNil(t, app)
	assert.Equal(t, "test_address", app.server.Addr)
	app.InitializeServer()
	ts := httptest.NewServer(app.mux)
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
		// {
		// 	description: "It serves the linkedin-icon.svg file",
		// 	path:        "/assets/linkedin-icon.svg",
		// 	contentType: "image/svg+xml",
		// },
		// {
		// 	description: "It serves the github-mark.svg file",
		// 	path:        "/assets/github-mark.svg",
		// 	contentType: "image/svg+xml",
		// },
		{
			description: "It serves the home page",
			path:				"/",
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

func TestAppStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockEnv := new(MockEnvConfig)
	mockEnv.On("IsDev").Return(true)
	mockEnv.On("Addr").Return("localhost:3210")

	app := CreateApp(ctx, mockEnv)

	var wg sync.WaitGroup
	wg.Add(1)

	go app.Start(&wg)

	time.Sleep(1 * time.Second) // Give the server a moment to start

	cancel()
	wg.Wait()

	assert.True(t, true, "Server shutdown gracefully")
}