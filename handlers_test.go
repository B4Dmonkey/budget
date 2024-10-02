package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	mockEnv := new(MockEnvConfig)
	mockEnv.On("IsDev").Return(true)
	mockEnv.On("Addr").Return("test_address")
	// ctx := context.Background()
	app, err := New()
	assert.Nil(t, err, "Error creating app")
	ts := httptest.NewServer(app.Mux)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.Nil(t, err, "Error creating request")
	client := &http.Client{}
	_, err = client.Do(req)
	assert.Nil(t, err, "Error performing request")
}
func TestUpload(t *testing.T) {
	mockEnv := new(MockEnvConfig)
	mockEnv.On("IsDev").Return(true)
	mockEnv.On("Addr").Return("test_address")
	// ctx := context.Background()
	app, err := New()
	assert.Nil(t, err, "Error creating app")
	ts := httptest.NewServer(app.Mux)
	defer ts.Close()

	filePath := "/Users/appstack/Downloads/Chase Activity Sept 27.CSV"

	// Open the file
	file, err := os.Open(filePath)
	assert.Nil(t, err, "Error opening file")
	defer file.Close()

	// Create a new multipart form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	assert.Nil(t, err, "Error creating form file")

	// Copy the file content to the multipart form
	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/upload", body)
	assert.Nil(t, err, "Error creating request")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := &http.Client{}
	_, err = client.Do(req)
	assert.Nil(t, err, "Error performing request")
	// todo: should finish this request
}
