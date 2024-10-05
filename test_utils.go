package main

import (
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


func setupTestServer(t *testing.T) *httptest.Server {
	mockEnv := new(MockEnvConfig)
	mockEnv.On("IsDev").Return(true)
	mockEnv.On("Addr").Return(":80")
	// ctx := context.Background()
	app, err := New()
	assert.Nil(t, err, "Error creating app")
  // assert.Equal(t, "test_address", app.Server.Addr)
	ts := httptest.NewServer(app.Mux)
	return ts
}